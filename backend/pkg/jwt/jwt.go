// Package jwt - JWT token generation and validation
package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Manager struct {
	cfg *config.JWT
	rdb *redis.Client
}

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID    string `json:"sub"`
	SessionID string `json:"sid"`
	Role      string `json:"role,omitempty"`
	CompanyID string `json:"company_id,omitempty"`
	TenantID  string `json:"tenant_id,omitempty"`
}

type SessionData struct {
	RefreshToken string `json:"rt"`
	UserAgent    string `json:"ua"`
	ClientIP     string `json:"ip"`
	CreatedAt    int64  `json:"created_at"`
	ExpiresAt    int64  `json:"expires_at"`
}

func New(cfg *config.JWT, rdb *redis.Client) *Manager {
	return &Manager{
		cfg: cfg,
		rdb: rdb,
	}
}

func (m *Manager) Generate(ctx context.Context, userID, userAgent, clientIP string) (accessToken, refreshToken string, err error) {
	sessionID := uuid.New().String()
	refreshToken = uuid.New().String()

	accessToken, err = m.generateAccessToken(userID, sessionID)
	if err != nil {
		return "", "", err
	}

	sessionData := SessionData{
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		CreatedAt:    time.Now().Unix(),
		ExpiresAt:    time.Now().Add(m.cfg.RefreshExpireMinutes).Unix(),
	}

	jsonData, err := json.Marshal(sessionData)
	if err != nil {
		return "", "", fmt.Errorf("json marshal error: %w", err)
	}

	key := m.getSessionKey(userID, sessionID)
	err = m.rdb.Set(ctx, key, jsonData, m.cfg.RefreshExpireMinutes).Err()
	if err != nil {
		return "", "", fmt.Errorf("redis set error: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (m *Manager) ValidateAccessToken(ctx context.Context, tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.SecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New(codes.AuthTokenExpired.String())
		}
		return nil, errors.New(codes.AuthTokenInvalid.String())
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New(codes.AuthTokenInvalid.String())
	}

	isBlocked, err := m.rdb.Exists(ctx, m.getBlockKey(claims.UserID)).Result()
	if err != nil {
		return nil, fmt.Errorf("redis check error: %w", err)
	}
	if isBlocked > 0 {
		return nil, errors.New(codes.UserBlocked.String())
	}

	exists, err := m.rdb.Exists(ctx, m.getSessionKey(claims.UserID, claims.SessionID)).Result()
	if err != nil || exists == 0 {
		return nil, errors.New(codes.SessionRevoked.String())
	}

	return claims, nil
}

func (m *Manager) Refresh(ctx context.Context, userID, sessionID, oldRefreshToken, currentUserAgent string) (newAccess, newRefresh string, err error) {
	key := m.getSessionKey(userID, sessionID)

	val, err := m.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", "", errors.New(codes.SessionRevoked.String())
	}
	if err != nil {
		return "", "", fmt.Errorf("redis error: %w", err)
	}

	var sessionData SessionData
	if err := json.Unmarshal([]byte(val), &sessionData); err != nil {
		return "", "", fmt.Errorf("json unmarshal error: %w", err)
	}

	if sessionData.RefreshToken != oldRefreshToken {
		m.rdb.Del(ctx, key)
		return "", "", errors.New(codes.SessionMismatch.String())
	}

	if sessionData.UserAgent != currentUserAgent {
		return "", "", errors.New(codes.SessionRevoked.String())
	}

	newRefresh = uuid.New().String()
	newAccess, err = m.generateAccessToken(userID, sessionID)
	if err != nil {
		return "", "", err
	}

	sessionData.RefreshToken = newRefresh

	newJSONData, err := json.Marshal(sessionData)
	if err != nil {
		return "", "", err
	}

	err = m.rdb.Set(ctx, key, newJSONData, m.cfg.RefreshExpireMinutes).Err()
	if err != nil {
		return "", "", fmt.Errorf("redis update error: %w", err)
	}

	return newAccess, newRefresh, nil
}

func (m *Manager) Logout(ctx context.Context, userID, sessionID string) error {
	return m.rdb.Del(ctx, m.getSessionKey(userID, sessionID)).Err()
}

func (m *Manager) LogoutAll(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("user:%s:session:*", userID)

	iter := m.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := m.rdb.Del(ctx, iter.Val()).Err(); err != nil {
			fmt.Printf("failed to delete session key %s: %v\n", iter.Val(), err)
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetUserSessions(ctx context.Context, userID string) ([]map[string]interface{}, error) {
	pattern := fmt.Sprintf("user:%s:session:*", userID)
	var sessions []map[string]interface{}

	iter := m.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := m.rdb.Get(ctx, key).Result()
		if err == nil {
			var s SessionData
			if json.Unmarshal([]byte(val), &s) == nil {
				sid := key[len(key)-36:]

				sessions = append(sessions, map[string]interface{}{
					"session_id": sid,
					"user_agent": s.UserAgent,
					"client_ip":  s.ClientIP,
					"created_at": s.CreatedAt,
					"is_current": false, 
				})
			}
		}
	}
	return sessions, nil
}

func (m *Manager) generateAccessToken(userID, sessionID string) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "eirsystem",
			Subject:   userID,
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.cfg.AccessExpireMinutes)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:    userID,
		SessionID: sessionID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.SecretKey))
}

func (m *Manager) getSessionKey(userID, sessionID string) string {
	return fmt.Sprintf("user:%s:session:%s", userID, sessionID)
}

func (m *Manager) getBlockKey(userID string) string {
	return "user:blocked:" + userID
}

func (m *Manager) ParseTokenUnverified(tokenStr string) (*CustomClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, &CustomClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
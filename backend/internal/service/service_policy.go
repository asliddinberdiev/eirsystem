package service

import "github.com/casbin/casbin/v3"

type Policy interface {
	AddRoleToUser(userID string, roleName string, clinicID string) error
	SetupDefaultPolicies(clinicID string) error
}

type policyService struct {
	enforcer *casbin.Enforcer
}

func NewPolicyService(enforcer *casbin.Enforcer) Policy {
	return &policyService{
		enforcer: enforcer,
	}
}

func (s *policyService) AddRoleToUser(userID string, roleName string, clinicID string) error {
	_, err := s.enforcer.AddGroupingPolicy(userID, roleName, clinicID)
	return err
}

func (s *policyService) SetupDefaultPolicies(clinicID string) error {
	// Clinic Owner Permissions
	if _, err := s.enforcer.AddPolicy("role:owner", clinicID, "/api/v1/*", "GET"); err != nil {
		return err
	}
	if _, err := s.enforcer.AddPolicy("role:owner", clinicID, "/api/v1/*", "POST"); err != nil {
		return err
	}
	if _, err := s.enforcer.AddPolicy("role:owner", clinicID, "/api/v1/*", "DELETE"); err != nil {
		return err
	}

	// Doctor Permissions
	if _, err := s.enforcer.AddPolicy("role:doctor", clinicID, "/api/v1/patients", "GET"); err != nil {
		return err
	}
	if _, err := s.enforcer.AddPolicy("role:doctor", clinicID, "/api/v1/appointments", "POST"); err != nil {
		return err
	}

	// Nurse Permissions
	if _, err := s.enforcer.AddPolicy("role:nurse", clinicID, "/api/v1/patients", "GET"); err != nil {
		return err
	}
	if _, err := s.enforcer.AddPolicy("role:nurse", clinicID, "/api/v1/vitals", "POST"); err != nil {
		return err
	}

	// Test Permissions
	if _, err := s.enforcer.AddPolicy("role:owner", clinicID, "/api/v1/test/owner", "GET"); err != nil {
		return err
	}
	if _, err := s.enforcer.AddPolicy("role:doctor", clinicID, "/api/v1/test/doctor", "GET"); err != nil {
		return err
	}
	if _, err := s.enforcer.AddPolicy("role:nurse", clinicID, "/api/v1/test/nurse", "GET"); err != nil {
		return err
	}

	return nil
}

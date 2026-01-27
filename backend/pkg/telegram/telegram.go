package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/asliddinberdiev/eirsystem/pkg/logger"
)

var (
	globalBot *Bot
	once      sync.Once
	wg        sync.WaitGroup
)

type Bot struct {
	token   string
	chatID  string
	msgChan chan string
	client  *http.Client
	log     logger.Logger
}

func Init(log logger.Logger, token, chatID string) {
	once.Do(func() {
		if token == "" || chatID == "" {
			log.Warn("Telegram token or ChatID is empty. Bot disabled.")
			return
		}

		globalBot = &Bot{
			log:    log,
			token:  token,
			chatID: chatID,
			msgChan: make(chan string, 100),
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
		}

		wg.Go(globalBot.worker)
		log.Info("Telegram bot initialized successfully")
	})
}

func Send(msg string) {
	if globalBot == nil {
		return
	}

	defer func() {
		if recover() != nil {
		}
	}()

	select {
	case globalBot.msgChan <- msg:
	default:
	}
}

func Close() {
	if globalBot != nil {
		close(globalBot.msgChan)
		wg.Wait()
	}
}

func (b *Bot) worker() {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", b.token)

	for msg := range b.msgChan {
		payload := map[string]string{
			"chat_id":    b.chatID,
			"text":       msg,
			"parse_mode": "HTML",
		}

		body, _ := json.Marshal(payload)
		resp, err := b.client.Post(url, "application/json", bytes.NewBuffer(body))

		if err != nil {
			b.log.Error("Telegram send failed", logger.Error(err))
		} else {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b.log.Warn("Telegram API returned non-200 status", logger.Int("status", resp.StatusCode))
			}
		}
	}
}
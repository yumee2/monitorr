package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const apiBase = "https://api.telegram.org"

type TelegramNotifier struct {
	botToken   string
	chatID     string
	httpClient *http.Client
}

func NewTelegramNotifier(botToken, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		botToken:   botToken,
		chatID:     chatID,
		httpClient: &http.Client{},
	}
}

func (n *TelegramNotifier) Send(ctx context.Context, text string) error {
	body, err := json.Marshal(map[string]string{
		"chat_id": n.chatID,
		"text":    text,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/bot%s/sendMessage", apiBase, n.botToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: unexpected status %d", resp.StatusCode)
	}
	return nil
}

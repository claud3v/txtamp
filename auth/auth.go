package auth

import (
	"context"
	"fmt"
	"strings"
	"time"
	"txtamp/config"
	"txtamp/navidrome"

	tea "charm.land/bubbletea/v2"
)

const connectTimeout = 10 * time.Second

type connectResultMsg struct {
	connectedTo string
	err         error
}

func connectServer(alias, baseURL, username, password string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
		defer cancel()

		client := navidrome.Client{
			BaseURL:  baseURL,
			Username: username,
			Password: password,
		}

		if err := client.Ping(ctx); err != nil {
			return connectResultMsg{err: err}
		}

		credentials := config.Credentials{
			Alias:    alias,
			Host:     baseURL,
			Username: username,
			Password: password,
		}

		if err := config.SaveCredentials(credentials); err != nil {
			return connectResultMsg{err: err}
		}

		return connectResultMsg{connectedTo: connectedName(credentials)}
	}
}

func loadSavedConnection() (config.Credentials, string, bool, error) {
	credentials, found, err := config.LoadCredentials()
	if err != nil || !found {
		return credentials, "", found, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client := navidrome.Client{
		BaseURL:  credentials.Host,
		Username: credentials.Username,
		Password: credentials.Password,
	}

	if err := client.Ping(ctx); err != nil {
		return credentials, "", true, fmt.Errorf("saved connection failed: %w", err)
	}

	return credentials, connectedName(credentials), true, nil
}

func connectedName(credentials config.Credentials) string {
	if strings.TrimSpace(credentials.Alias) != "" {
		return credentials.Alias
	}

	return credentials.Host
}

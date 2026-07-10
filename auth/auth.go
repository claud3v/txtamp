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

type ConnectResultMsg struct {
	ConnectedTo string
	Client      navidrome.Client
	Err         error
}

type Connection struct {
	Alias       string
	Host        string
	Username    string
	Password    string
	ConnectedTo string
}

func (c Connection) Client() navidrome.Client {
	return navidrome.Client{
		BaseURL:  c.Host,
		Username: c.Username,
		Password: c.Password,
	}
}

func connectServer(alias, baseURL, username, password string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
		defer cancel()

		connection := Connection{
			Alias:    alias,
			Host:     baseURL,
			Username: username,
			Password: password,
		}
		client := connection.Client()

		if err := client.Ping(ctx); err != nil {
			return ConnectResultMsg{Err: err}
		}

		credentials := credentialsFromConnection(connection)

		if err := config.SaveCredentials(credentials); err != nil {
			return ConnectResultMsg{Err: err}
		}

		return ConnectResultMsg{
			ConnectedTo: connectedName(credentials),
			Client:      client,
		}
	}
}

func LoadSavedConnection() (Connection, bool, error) {
	credentials, found, err := config.LoadCredentials()
	if err != nil || !found {
		return connectionFromCredentials(credentials), found, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	connection := connectionFromCredentials(credentials)
	if err := connection.Client().Ping(ctx); err != nil {
		return connection, true, fmt.Errorf("saved connection failed: %w", err)
	}

	connection.ConnectedTo = connectedName(credentials)

	return connection, true, nil
}

func credentialsFromConnection(connection Connection) config.Credentials {
	return config.Credentials{
		Alias:    connection.Alias,
		Host:     connection.Host,
		Username: connection.Username,
		Password: connection.Password,
	}
}

func connectionFromCredentials(credentials config.Credentials) Connection {
	return Connection{
		Alias:    credentials.Alias,
		Host:     credentials.Host,
		Username: credentials.Username,
		Password: credentials.Password,
	}
}

func connectedName(credentials config.Credentials) string {
	if strings.TrimSpace(credentials.Alias) != "" {
		return credentials.Alias
	}

	return credentials.Host
}

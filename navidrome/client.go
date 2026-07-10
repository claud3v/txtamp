package navidrome

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiVersion = "1.16.1"
	clientName = "txtamp"
)

type Client struct {
	BaseURL  string
	Username string
	Password string

	HTTPClient *http.Client
}

type Playlist struct {
	ID        string
	Name      string
	SongCount int
	Duration  int
}

type Song struct {
	ID       string
	Title    string
	Artist   string
	Album    string
	Track    int
	Duration int
}

func (c Client) Ping(ctx context.Context) error {
	var response subsonicResponse
	if err := c.get(ctx, "ping.view", nil, &response); err != nil {
		return err
	}

	if response.Response.Status != "ok" {
		if response.Response.Error.Message != "" {
			return fmt.Errorf("navidrome ping failed: %s", response.Response.Error.Message)
		}

		return fmt.Errorf("navidrome ping failed: status %q", response.Response.Status)
	}

	return nil
}

func (c Client) ListPlaylists(ctx context.Context) ([]Playlist, error) {
	var response playlistsResponse
	if err := c.get(ctx, "getPlaylists.view", nil, &response); err != nil {
		return nil, err
	}

	if err := checkStatus(response.Response.Status, response.Response.Error, "list playlists"); err != nil {
		return nil, err
	}

	return response.Response.Playlists.Playlist, nil
}

func (c Client) GetPlaylist(ctx context.Context, id string) ([]Song, error) {
	params := url.Values{}
	params.Set("id", id)

	var response playlistResponse
	if err := c.get(ctx, "getPlaylist.view", params, &response); err != nil {
		return nil, err
	}

	if err := checkStatus(response.Response.Status, response.Response.Error, "get playlist"); err != nil {
		return nil, err
	}

	return response.Response.Playlist.Entry, nil
}

func (c Client) StreamURL(id string) (string, error) {
	params := url.Values{}
	params.Set("id", id)

	return c.requestURL("stream.view", params)
}

func (c Client) get(ctx context.Context, endpoint string, params url.Values, out any) error {
	requestURL, err := c.requestURL(endpoint, params)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("navidrome returned HTTP %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

func (c Client) requestURL(endpoint string, params url.Values) (string, error) {
	if strings.TrimSpace(c.BaseURL) == "" {
		return "", fmt.Errorf("base URL cannot be blank")
	}

	base, err := url.Parse(strings.TrimRight(c.BaseURL, "/"))
	if err != nil {
		return "", fmt.Errorf("parsing base URL: %w", err)
	}

	if base.Scheme == "" || base.Host == "" {
		return "", fmt.Errorf("base URL must include scheme and host")
	}

	salt, err := randomSalt()
	if err != nil {
		return "", err
	}

	values := url.Values{}
	for key, vals := range params {
		for _, val := range vals {
			values.Add(key, val)
		}
	}

	values.Set("u", c.Username)
	values.Set("t", authToken(c.Password, salt))
	values.Set("s", salt)
	values.Set("v", apiVersion)
	values.Set("c", clientName)
	values.Set("f", "json")

	base.Path = strings.TrimRight(base.Path, "/") + "/rest/" + strings.TrimLeft(endpoint, "/")
	base.RawQuery = values.Encode()

	return base.String(), nil
}

func randomSalt() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generating auth salt: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func authToken(password, salt string) string {
	sum := md5.Sum([]byte(password + salt))
	return hex.EncodeToString(sum[:])
}

type subsonicResponse struct {
	Response struct {
		Status string        `json:"status"`
		Error  subsonicError `json:"error"`
	} `json:"subsonic-response"`
}

type playlistsResponse struct {
	Response struct {
		Status    string        `json:"status"`
		Error     subsonicError `json:"error"`
		Playlists struct {
			Playlist []Playlist `json:"playlist"`
		} `json:"playlists"`
	} `json:"subsonic-response"`
}

type playlistResponse struct {
	Response struct {
		Status   string        `json:"status"`
		Error    subsonicError `json:"error"`
		Playlist struct {
			Entry []Song `json:"entry"`
		} `json:"playlist"`
	} `json:"subsonic-response"`
}

type subsonicError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func checkStatus(status string, responseError subsonicError, action string) error {
	if status == "ok" {
		return nil
	}

	if responseError.Message != "" {
		return fmt.Errorf("navidrome %s failed: %s", action, responseError.Message)
	}

	return fmt.Errorf("navidrome %s failed: status %q", action, status)
}

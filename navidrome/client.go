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

type Artist struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	AlbumCount int    `json:"albumCount"`
}

type Album struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Artist    string `json:"artist"`
	Year      int    `json:"year"`
	SongCount int    `json:"songCount"`
	Duration  int    `json:"duration"`
}

type Song struct {
	ID       string
	Title    string
	Artist   string
	Album    string
	Track    int
	Duration int
}

type SearchResult struct {
	Artists []Artist `json:"artist"`
	Albums  []Album  `json:"album"`
	Songs   []Song   `json:"song"`
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

func (c Client) ListArtists(ctx context.Context) ([]Artist, error) {
	var response artistsResponse
	if err := c.get(ctx, "getArtists.view", nil, &response); err != nil {
		return nil, err
	}

	if err := checkStatus(response.Response.Status, response.Response.Error, "list artists"); err != nil {
		return nil, err
	}

	var artists []Artist
	for _, index := range response.Response.Artists.Index {
		artists = append(artists, index.Artist...)
	}

	return artists, nil
}

func (c Client) GetArtistAlbums(ctx context.Context, id string) ([]Album, error) {
	params := url.Values{}
	params.Set("id", id)

	var response artistResponse
	if err := c.get(ctx, "getArtist.view", params, &response); err != nil {
		return nil, err
	}

	if err := checkStatus(response.Response.Status, response.Response.Error, "get artist"); err != nil {
		return nil, err
	}

	return response.Response.Artist.Album, nil
}

func (c Client) GetAlbumSongs(ctx context.Context, id string) ([]Song, error) {
	params := url.Values{}
	params.Set("id", id)

	var response albumResponse
	if err := c.get(ctx, "getAlbum.view", params, &response); err != nil {
		return nil, err
	}

	if err := checkStatus(response.Response.Status, response.Response.Error, "get album"); err != nil {
		return nil, err
	}

	return response.Response.Album.Song, nil
}

func (c Client) Search(ctx context.Context, query string) (SearchResult, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("artistCount", "20")
	params.Set("albumCount", "20")
	params.Set("songCount", "50")

	var response searchResponse
	if err := c.get(ctx, "search3.view", params, &response); err != nil {
		return SearchResult{}, err
	}

	if err := checkStatus(response.Response.Status, response.Response.Error, "search"); err != nil {
		return SearchResult{}, err
	}

	return response.Response.SearchResult, nil
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

type artistsResponse struct {
	Response struct {
		Status  string        `json:"status"`
		Error   subsonicError `json:"error"`
		Artists struct {
			Index []struct {
				Artist []Artist `json:"artist"`
			} `json:"index"`
		} `json:"artists"`
	} `json:"subsonic-response"`
}

type artistResponse struct {
	Response struct {
		Status string        `json:"status"`
		Error  subsonicError `json:"error"`
		Artist struct {
			Album []Album `json:"album"`
		} `json:"artist"`
	} `json:"subsonic-response"`
}

type albumResponse struct {
	Response struct {
		Status string        `json:"status"`
		Error  subsonicError `json:"error"`
		Album  struct {
			Song []Song `json:"song"`
		} `json:"album"`
	} `json:"subsonic-response"`
}

type searchResponse struct {
	Response struct {
		Status       string        `json:"status"`
		Error        subsonicError `json:"error"`
		SearchResult SearchResult  `json:"searchResult3"`
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

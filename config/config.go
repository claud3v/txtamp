package config

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	configDirName = ".txtamp"
	envFileName   = "config.env"
	keyFileName   = ".key"
)

type Credentials struct {
	Alias    string
	Host     string
	Username string
	Password string
}

func SaveCredentials(credentials Credentials) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("finding home directory: %w", err)
	}

	return saveCredentials(
		filepath.Join(home, configDirName, envFileName),
		filepath.Join(home, configDirName, keyFileName),
		credentials,
	)
}

func LoadCredentials() (Credentials, bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Credentials{}, false, fmt.Errorf("finding home directory: %w", err)
	}

	return loadCredentials(
		filepath.Join(home, configDirName, envFileName),
		filepath.Join(home, configDirName, keyFileName),
	)
}

func saveCredentials(envPath, keyPath string, credentials Credentials) error {
	if err := os.MkdirAll(filepath.Dir(envPath), 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	key, err := loadOrCreateKey(keyPath)
	if err != nil {
		return err
	}

	encryptedPassword, err := encryptPassword(credentials.Password, key)
	if err != nil {
		return err
	}

	var contents strings.Builder
	if strings.TrimSpace(credentials.Alias) != "" {
		fmt.Fprintf(&contents, "alias=%s\n", quoteEnvValue(credentials.Alias))
	}
	fmt.Fprintf(&contents, "host=%s\n", quoteEnvValue(credentials.Host))
	fmt.Fprintf(&contents, "username=%s\n", quoteEnvValue(credentials.Username))
	fmt.Fprintf(&contents, "password=%s\n", quoteEnvValue(encryptedPassword))

	if err := os.WriteFile(envPath, []byte(contents.String()), 0600); err != nil {
		return fmt.Errorf("writing credentials: %w", err)
	}

	return nil
}

func loadCredentials(envPath, keyPath string) (Credentials, bool, error) {
	values, err := readEnvFile(envPath)
	if errors.Is(err, os.ErrNotExist) {
		return Credentials{}, false, nil
	}
	if err != nil {
		return Credentials{}, false, err
	}

	key, err := loadKey(keyPath)
	if err != nil {
		return Credentials{}, true, err
	}

	password, err := decryptPassword(values["password"], key)
	if err != nil {
		return Credentials{}, true, err
	}

	return Credentials{
		Alias:    values["alias"],
		Host:     values["host"],
		Username: values["username"],
		Password: password,
	}, true, nil
}

func readEnvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		return nil, fmt.Errorf("reading credentials: %w", err)
	}
	defer file.Close()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("invalid credentials line %q", line)
		}

		unquoted, err := strconv.Unquote(strings.TrimSpace(value))
		if err != nil {
			return nil, fmt.Errorf("parsing credentials value for %s: %w", strings.TrimSpace(key), err)
		}

		values[strings.TrimSpace(key)] = unquoted
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading credentials: %w", err)
	}

	return values, nil
}

func loadKey(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading encryption key: %w", err)
	}

	key, err := hex.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		return nil, fmt.Errorf("decoding encryption key: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes")
	}

	return key, nil
}

func loadOrCreateKey(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		key, err := hex.DecodeString(strings.TrimSpace(string(data)))
		if err != nil {
			return nil, fmt.Errorf("decoding encryption key: %w", err)
		}
		if len(key) != 32 {
			return nil, fmt.Errorf("encryption key must be 32 bytes")
		}

		return key, nil
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading encryption key: %w", err)
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("generating encryption key: %w", err)
	}

	if err := os.WriteFile(path, []byte(hex.EncodeToString(key)+"\n"), 0600); err != nil {
		return nil, fmt.Errorf("writing encryption key: %w", err)
	}

	return key, nil
}

func encryptPassword(password string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("creating cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("creating cipher mode: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generating password nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(password), nil)
	encrypted := append(nonce, ciphertext...)

	return "v1:" + base64.StdEncoding.EncodeToString(encrypted), nil
}

func decryptPassword(password string, key []byte) (string, error) {
	if !strings.HasPrefix(password, "v1:") {
		return "", fmt.Errorf("unsupported password format")
	}

	encrypted, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(password, "v1:"))
	if err != nil {
		return "", fmt.Errorf("decoding password: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("creating cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("creating cipher mode: %w", err)
	}

	if len(encrypted) < gcm.NonceSize() {
		return "", fmt.Errorf("encrypted password is too short")
	}

	nonce := encrypted[:gcm.NonceSize()]
	ciphertext := encrypted[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypting password: %w", err)
	}

	return string(plaintext), nil
}

func quoteEnvValue(value string) string {
	return strconv.Quote(value)
}

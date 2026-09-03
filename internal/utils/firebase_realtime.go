package utils

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"nalakarsa/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type serviceAccountFile struct {
	ProjectID   string `json:"project_id"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	TokenURI    string `json:"token_uri"`
}

type tokenCache struct {
	mu          sync.RWMutex
	accessToken string
	expiresAt   time.Time
}

var (
	saData     *serviceAccountFile
	saOnce     sync.Once
	parsedKey  *rsa.PrivateKey
	globalAuth tokenCache
)

func loadServiceAccount() {
	var data []byte
	var err error
	if envJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON"); strings.TrimSpace(envJSON) != "" {
		data = []byte(envJSON)
	} else {
		filePath := "serviceAccountKey.json"
		data, err = os.ReadFile(filePath)
		if err != nil {
			return
		}
	}

	var sa serviceAccountFile
	if err := json.Unmarshal(data, &sa); err != nil {
		log.Printf("[FIREBASE AUTH] Failed to unmarshal service account json: %v", err)
		return
	}

	block, _ := pem.Decode([]byte(sa.PrivateKey))
	if block == nil {
		log.Printf("[FIREBASE AUTH] Failed to decode PEM block of private key")
		return
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("[FIREBASE AUTH] Failed to parse PKCS8 private key: %v", err)
		return
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		log.Printf("[FIREBASE AUTH] Parsed key is not RSA private key")
		return
	}

	saData = &sa
	parsedKey = rsaKey
	log.Println("[FIREBASE AUTH] Successfully loaded Firebase Service Account for project:", sa.ProjectID)
}

func getFirebaseAccessToken() (string, error) {
	saOnce.Do(loadServiceAccount)
	if saData == nil || parsedKey == nil {
		return "", nil
	}

	globalAuth.mu.RLock()
	if globalAuth.accessToken != "" && time.Now().Before(globalAuth.expiresAt.Add(-2*time.Minute)) {
		token := globalAuth.accessToken
		globalAuth.mu.RUnlock()
		return token, nil
	}
	globalAuth.mu.RUnlock()

	globalAuth.mu.Lock()
	defer globalAuth.mu.Unlock()
	if globalAuth.accessToken != "" && time.Now().Before(globalAuth.expiresAt.Add(-2*time.Minute)) {
		return globalAuth.accessToken, nil
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   saData.ClientEmail,
		"sub":   saData.ClientEmail,
		"aud":   "https://oauth2.googleapis.com/token",
		"iat":   now.Unix(),
		"exp":   now.Add(1 * time.Hour).Unix(),
		"scope": "https://www.googleapis.com/auth/firebase.database https://www.googleapis.com/auth/userinfo.email",
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedJWT, err := jwtToken.SignedString(parsedKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT for service account: %w", err)
	}

	tokenURI := saData.TokenURI
	if tokenURI == "" {
		tokenURI = "https://oauth2.googleapis.com/token"
	}

	reqBody := fmt.Sprintf("grant_type=urn%%3Aietf%%3Aparams%%3Aoauth%%3Agrant-type%%3Ajwt-bearer&assertion=%s", signedJWT)
	req, err := http.NewRequest("POST", tokenURI, strings.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange token with Google OAuth2: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("google oauth token exchange failed (%d): %s", resp.StatusCode, string(b))
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	globalAuth.accessToken = res.AccessToken
	globalAuth.expiresAt = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)
	return res.AccessToken, nil
}
func PushToFirebaseRealtime(path string, payload interface{}, cfg *config.Config) {
	go func() {
		saOnce.Do(loadServiceAccount)

		baseURL := ""
		if cfg != nil && cfg.FirebaseDatabaseURL != "" {
			baseURL = cfg.FirebaseDatabaseURL
		} else if saData != nil && saData.ProjectID != "" {
			baseURL = fmt.Sprintf("https://%s-default-rtdb.asia-southeast1.firebasedatabase.app", saData.ProjectID)
		}

		if baseURL == "" {
			return
		}

		baseURL = strings.TrimRight(baseURL, "/")
		cleanPath := strings.Trim(path, "/")
		url := fmt.Sprintf("%s/%s.json", baseURL, cleanPath)

		if cfg != nil && cfg.FirebaseDatabaseSecret != "" {
			url = fmt.Sprintf("%s?auth=%s", url, cfg.FirebaseDatabaseSecret)
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Printf("[FIREBASE RTDB] Failed to marshal payload for %s: %v", cleanPath, err)
			return
		}

		client := &http.Client{Timeout: 5 * time.Second}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("[FIREBASE RTDB] Failed to create request for %s: %v", cleanPath, err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		accessToken, tokenErr := getFirebaseAccessToken()
		if tokenErr == nil && accessToken != "" && (cfg == nil || cfg.FirebaseDatabaseSecret == "") {
			req.Header.Set("Authorization", "Bearer "+accessToken)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[FIREBASE RTDB] Network error pushing to %s: %v", cleanPath, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			log.Printf("[FIREBASE RTDB] Successfully broadcasted message to path: %s", cleanPath)
		} else {
			respBody, _ := io.ReadAll(resp.Body)
			log.Printf("[FIREBASE RTDB] Non-2xx response from Firebase (%d) on %s: %s", resp.StatusCode, cleanPath, string(respBody))
		}
	}()
}

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"nalakarsa/internal/config"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type supabaseUploadResponse struct {
	Key string `json:"Key"`
}

type supabaseDeleteResponse struct {
	Message string `json:"message"`
}

type supabaseDeleteRequest struct {
	Prefixes []string `json:"prefixes"`
	Paths    []string `json:"paths,omitempty"`
}

func CreateChatSignedURL(path string, cfg *config.Config) (string, error) {
	if cfg.SupabaseURL == "" || cfg.SupabaseServiceRoleKey == "" || cfg.SupabaseChatBucket == "" {
		return "", errors.New("supabase chat storage is not configured")
	}

	payload, err := json.Marshal(map[string]interface{}{
		"expiresIn": 3600,
		"paths":     []string{path},
	})
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/storage/v1/object/sign/%s", strings.TrimRight(cfg.SupabaseURL, "/"), cfg.SupabaseChatBucket)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseServiceRoleKey)
	req.Header.Set("apikey", cfg.SupabaseServiceRoleKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("supabase signed URL error (%d): %s", resp.StatusCode, bytes.TrimSpace(body))
	}

	var result struct {
		SignedURLs []struct {
			SignedURL string `json:"signedURL"`
		} `json:"signedURLs"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if len(result.SignedURLs) == 0 || result.SignedURLs[0].SignedURL == "" {
		return "", errors.New("signed URL was not returned")
	}
	return result.SignedURLs[0].SignedURL, nil
}

func UploadChatAttachment(conversationID uuid.UUID, fileHeader *multipart.FileHeader, cfg *config.Config) (string, error) {
	if cfg.SupabaseURL == "" || cfg.SupabaseServiceRoleKey == "" || cfg.SupabaseChatBucket == "" {
		return "", errors.New("supabase chat storage is not configured")
	}

	allowedTypes := map[string]bool{
		"image/jpeg": true, "image/jpg": true, "image/png": true,
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}
	mimeType := fileHeader.Header.Get("Content-Type")
	if !allowedTypes[mimeType] {
		return "", errors.New("unsupported file type")
	}
	if cfg.SupabaseChatMaxFileSize > 0 && fileHeader.Size > cfg.SupabaseChatMaxFileSize {
		return "", errors.New("file size exceeds the maximum limit")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	objectPath := fmt.Sprintf("conversations/%s/pending/%s%s", conversationID, uuid.New(), ext)
	baseURL := strings.TrimRight(cfg.SupabaseURL, "/")
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", baseURL, cfg.SupabaseChatBucket, objectPath)
	req, err := http.NewRequest(http.MethodPost, uploadURL, file)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseServiceRoleKey)
	req.Header.Set("apikey", cfg.SupabaseServiceRoleKey)
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("x-upsert", "false")
	req.ContentLength = fileHeader.Size

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload chat attachment: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return "", fmt.Errorf("supabase chat upload error (%d): %s", resp.StatusCode, bytes.TrimSpace(body))
	}
	return objectPath, nil
}

func UploadAvatarToSupabase(userID uuid.UUID, fileHeader *multipart.FileHeader, cfg *config.Config) (string, error) {
	if cfg.SupabaseURL == "" || cfg.SupabaseServiceRoleKey == "" || cfg.SupabaseStorageBucket == "" {
		return "", errors.New("supabase storage configuration (URL, service role key, or storage bucket) is not configured")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == ".jpeg" {
		ext = ".jpg"
	}
	objectPath := fmt.Sprintf("%s%s", userID.String(), ext)

	baseURL := strings.TrimRight(cfg.SupabaseURL, "/")
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", baseURL, cfg.SupabaseStorageBucket, objectPath)

	parsedUploadURL, err := url.Parse(uploadURL)
	if err != nil {
		return "", fmt.Errorf("invalid supabase URL: %w", err)
	}
	query := parsedUploadURL.Query()
	query.Set("upsert", "true")
	parsedUploadURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, parsedUploadURL.String(), file)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseServiceRoleKey)
	req.Header.Set("apikey", cfg.SupabaseServiceRoleKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("x-upsert", "true")
	req.Header.Set("cache-control", "no-cache")
	req.ContentLength = fileHeader.Size

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to Supabase: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("supabase storage error (%d): %s", resp.StatusCode, bytes.TrimSpace(body))
	}

	if len(body) > 0 {
		var parsedResp supabaseUploadResponse
		if err := json.Unmarshal(body, &parsedResp); err != nil {
			return "", fmt.Errorf("failed to parse supabase response: %w", err)
		}
		if parsedResp.Key != "" && parsedResp.Key != objectPath {
			objectPath = parsedResp.Key
		}
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", baseURL, cfg.SupabaseStorageBucket, objectPath)
	return publicURL, nil
}

func CleanupOtherAvatarVariants(userID uuid.UUID, keepURL string, cfg *config.Config) error {
	if cfg.SupabaseURL == "" || cfg.SupabaseServiceRoleKey == "" || cfg.SupabaseStorageBucket == "" {
		return errors.New("supabase storage configuration (URL, service role key, or storage bucket) is not configured")
	}

	baseURL := strings.TrimRight(cfg.SupabaseURL, "/")
	keepPath, err := extractSupabaseObjectPath(keepURL, cfg)
	if err != nil {
		return err
	}

	currentExt := strings.ToLower(filepath.Ext(keepPath))
	var paths []string
	for _, ext := range []string{".jpg", ".jpeg", ".png"} {
		if ext != currentExt {
			paths = append(paths, fmt.Sprintf("%s%s", userID.String(), ext))
		}
	}
	if len(paths) == 0 {
		return nil
	}

	payload, err := json.Marshal(supabaseDeleteRequest{
		Prefixes: paths,
		Paths:    paths,
	})
	if err != nil {
		return fmt.Errorf("failed to encode delete request: %w", err)
	}

	deleteURL := fmt.Sprintf("%s/storage/v1/object/%s", baseURL, cfg.SupabaseStorageBucket)
	req, err := http.NewRequest(http.MethodDelete, deleteURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.SupabaseServiceRoleKey)
	req.Header.Set("apikey", cfg.SupabaseServiceRoleKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete old avatar from Supabase: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("supabase bulk delete error (%d): %s", resp.StatusCode, bytes.TrimSpace(body))
	}

	return nil
}

func extractSupabaseObjectPath(avatarURL string, cfg *config.Config) (string, error) {
	parsedAvatarURL, err := url.Parse(avatarURL)
	if err != nil {
		return "", fmt.Errorf("invalid avatar URL: %w", err)
	}

	parsedBaseURL, err := url.Parse(strings.TrimRight(cfg.SupabaseURL, "/"))
	if err != nil {
		return "", fmt.Errorf("invalid supabase URL: %w", err)
	}

	if !strings.EqualFold(parsedAvatarURL.Host, parsedBaseURL.Host) {
		return "", nil
	}

	publicPrefix := fmt.Sprintf("/storage/v1/object/public/%s/", cfg.SupabaseStorageBucket)
	if !strings.HasPrefix(parsedAvatarURL.Path, publicPrefix) {
		return "", nil
	}

	objectPath := strings.TrimPrefix(parsedAvatarURL.Path, publicPrefix)
	objectPath, err = url.PathUnescape(objectPath)
	if err != nil {
		return "", fmt.Errorf("failed to decode object path: %w", err)
	}

	return objectPath, nil
}

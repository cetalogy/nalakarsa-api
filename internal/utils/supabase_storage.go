package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"nalakarsa/internal/config"
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
	Paths []string `json:"paths"`
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
	objectPath := fmt.Sprintf("avatars/%s%s", userID.String(), ext)

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

	candidates := []string{
		fmt.Sprintf("avatars/%s.jpg", userID.String()),
		fmt.Sprintf("avatars/%s.jpeg", userID.String()),
		fmt.Sprintf("avatars/%s.png", userID.String()),
	}
	paths := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate != keepPath {
			paths = append(paths, candidate)
		}
	}
	if len(paths) == 0 {
		return nil
	}

	payload, err := json.Marshal(supabaseDeleteRequest{Paths: paths})
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

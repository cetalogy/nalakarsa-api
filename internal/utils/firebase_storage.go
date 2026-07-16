package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"nalakarsa/internal/config"
	"path/filepath"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

// UploadAvatarToFirebase uploads a multipart file to Firebase Storage and returns the public URL.
func UploadAvatarToFirebase(fileHeader *multipart.FileHeader, cfg *config.Config) (string, error) {
	if cfg.FirebaseProjectID == "" || cfg.FirebaseStorageBucket == "" {
		return "", errors.New("firebase configurations (project ID or storage bucket) are not configured")
	}

	ctx := context.Background()

	// Initialize Firebase App
	var opt option.ClientOption
	if cfg.FirebaseCredentialJSONPath != "" {
		opt = option.WithCredentialsFile(cfg.FirebaseCredentialJSONPath)
	} else {
		return "", errors.New("firebase credential file path is required")
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{
		StorageBucket: cfg.FirebaseStorageBucket,
		ProjectID:     cfg.FirebaseProjectID,
	}, opt)
	if err != nil {
		return "", fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	// Initialize Storage Client
	storageClient, err := app.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to initialize storage client: %w", err)
	}

	bucket, err := storageClient.Bucket(cfg.FirebaseStorageBucket)
	if err != nil {
		return "", fmt.Errorf("failed to get storage bucket: %w", err)
	}

	// Open the file from the header
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate a unique filename in GCS
	ext := filepath.Ext(fileHeader.Filename)
	uniqueFilename := fmt.Sprintf("avatars/%s%s", uuid.New().String(), ext)

	// Create a bucket object writer
	obj := bucket.Object(uniqueFilename)
	wc := obj.NewWriter(ctx)
	wc.ContentType = fileHeader.Header.Get("Content-Type")

	// Copy file data to GCS
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to copy file to bucket: %w", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	// Firebase storage public URL pattern:
	// https://firebasestorage.googleapis.com/v0/b/<bucket-name>/o/<encoded-file-path>?alt=media
	escapedFilename := strings.ReplaceAll(uniqueFilename, "/", "%2F")
	publicURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", cfg.FirebaseStorageBucket, escapedFilename)

	return publicURL, nil
}

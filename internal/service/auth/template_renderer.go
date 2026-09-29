package authservice

import (
	"html/template"
	"strings"

	commonemail "nalakarsa/internal/common/email"
)

type emailMessage struct {
	subject string
	body    string
}

type WelcomeData struct {
	RecipientName string
	LoginLink     string
}

type VerificationData struct {
	RecipientName string
	Code          string
}

type PasswordResetData struct {
	RecipientName string
	ResetLink     string
}

type PasswordChangedData struct {
	RecipientName string
}

func buildWelcomeEmail(data WelcomeData) (emailMessage, error) {
	body, err := renderTemplate("templates/welcome.html", struct {
		RecipientName string
		LoginLink     string
	}{RecipientName: displayName(data.RecipientName), LoginLink: data.LoginLink})
	if err != nil {
		return emailMessage{}, err
	}
	return emailMessage{subject: "Selamat datang di Nalakarsa", body: body}, nil
}

func buildVerificationEmail(data VerificationData) (emailMessage, error) {
	body, err := renderTemplate("templates/email_verification.html", struct {
		RecipientName string
		Code          string
	}{RecipientName: displayName(data.RecipientName), Code: data.Code})
	if err != nil {
		return emailMessage{}, err
	}
	return emailMessage{subject: "Verifikasi Email Nalakarsa", body: body}, nil
}

func buildPasswordResetEmail(data PasswordResetData) (emailMessage, error) {
	body, err := renderTemplate("templates/password_reset.html", struct {
		RecipientName string
		ResetLink     string
	}{RecipientName: displayName(data.RecipientName), ResetLink: data.ResetLink})
	if err != nil {
		return emailMessage{}, err
	}
	return emailMessage{subject: "Reset Password Nalakarsa", body: body}, nil
}

func buildPasswordChangedEmail(data PasswordChangedData) (emailMessage, error) {
	body, err := renderTemplate("templates/password_changed.html", struct{ RecipientName string }{RecipientName: displayName(data.RecipientName)})
	if err != nil {
		return emailMessage{}, err
	}
	return emailMessage{subject: "Password Nalakarsa berhasil diubah", body: body}, nil
}

func renderTemplate(path string, data any) (string, error) {
	source, err := commonemail.Templates.ReadFile(path)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New(path).Parse(string(source))
	if err != nil {
		return "", err
	}
	var rendered strings.Builder
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", err
	}
	return rendered.String(), nil
}

func displayName(name string) string {
	name = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(name, "\r", ""), "\n", ""))
	if name == "" {
		return "Pengguna Nalakarsa"
	}
	return name
}

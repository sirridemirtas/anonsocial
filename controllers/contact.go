package controllers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/utils"
)

const (
	googleFormURL = "https://docs.google.com/forms/d/e/1FAIpQLSdKSnGiN6BMg9dm53aB0FYhNR2hT-ZsLE1IQUxQehlON5wbYg/formResponse"
)

// ValidSubjects contains the allowed subject options
var ValidSubjects = []string{"Genel", "Destek", "Öneri", "Teknik", "Şikayet"}

// ContactFormData represents the contact form request body
type ContactFormData struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// isValidSubject checks if the subject is one of the valid options
func isValidSubject(subject string) bool {
	for _, validSubject := range ValidSubjects {
		if subject == validSubject {
			return true
		}
	}
	return false
}

// SubmitContactForm handles the contact form submission
func SubmitContactForm(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var formData ContactFormData
	if err := c.ShouldBindJSON(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate form fields
	if strings.TrimSpace(formData.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "İsim alanı zorunludur"})
		return
	}

	if !utils.ValidateEmail(formData.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçerli bir e-posta adresi girmelisiniz"})
		return
	}

	if !isValidSubject(formData.Subject) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Konu, şu seçeneklerden biri olmalıdır: %s", strings.Join(ValidSubjects, ", ")),
		})
		return
	}

	if strings.TrimSpace(formData.Message) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mesaj alanı zorunludur"})
		return
	}

	// Prepare form data for Google Form
	googleFormData := url.Values{}
	googleFormData.Add("entry.1278516670", formData.Name)
	googleFormData.Add("entry.594480513", formData.Email)
	googleFormData.Add("entry.508744151", formData.Subject)
	googleFormData.Add("entry.1252619669", formData.Message)

	// Create HTTP client with context and timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", googleFormURL, strings.NewReader(googleFormData.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Form işlenirken bir hata oluştu"})
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "AnonSocial/1.0")

	// Send request to Google Form
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Form gönderilemedi, lütfen daha sonra tekrar deneyin"})
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Form işlenirken bir hata oluştu, lütfen daha sonra tekrar deneyin"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Mesajınız başarıyla gönderildi. En kısa sürede size dönüş yapacağız."})
}

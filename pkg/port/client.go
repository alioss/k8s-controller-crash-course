package port

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// PortClient handles communication with Port.io API
type PortClient struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	Token        string
	HTTPClient   *http.Client
}

// NewPortClient creates a new Port.io API client
func NewPortClient(baseURL, clientID, clientSecret string) *PortClient {
	return &PortClient{
		BaseURL:      baseURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TokenResponse represents the response from Port's auth endpoint
type TokenResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}

// Authenticate gets an access token from Port.io
func (c *PortClient) Authenticate() error {
	authURL := fmt.Sprintf("%s/v1/auth/access_token", c.BaseURL)
	
	payload := map[string]string{
		"clientId":     c.ClientID,
		"clientSecret": c.ClientSecret,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal auth payload: %w", err)
	}
	
	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Port: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
	}
	
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}
	
	c.Token = tokenResp.AccessToken
	log.Info().Msg("Successfully authenticated with Port.io")
	
	return nil
}

// FrontendPageEntity represents a FrontendPage entity in Port
type FrontendPageEntity struct {
	Identifier string                 `json:"identifier"`
	Title      string                 `json:"title"`
	Blueprint  string                 `json:"blueprint"`
	Properties map[string]interface{} `json:"properties"`
}

// CreateOrUpdateEntity creates or updates a FrontendPage entity in Port
func (c *PortClient) CreateOrUpdateEntity(entity FrontendPageEntity) error {
	if c.Token == "" {
		if err := c.Authenticate(); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}
	
	entityURL := fmt.Sprintf("%s/v1/blueprints/%s/entities", c.BaseURL, entity.Blueprint)
	
	jsonData, err := json.Marshal(entity)
	if err != nil {
		return fmt.Errorf("failed to marshal entity: %w", err)
	}
	
	// Try to update first (PATCH)
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", entityURL, entity.Identifier), bytes.NewBuffer(jsonData))
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		
		resp, err := c.HTTPClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				log.Info().Str("entity", entity.Identifier).Msg("Updated entity in Port")
				return nil
			}
		}
	}
	
	// If update failed, try to create (POST)
	req, err = http.NewRequest("POST", entityURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create entity request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create entity, status: %d", resp.StatusCode)
	}
	
	log.Info().Str("entity", entity.Identifier).Msg("Created entity in Port")
	return nil
}

// DeleteEntity deletes a FrontendPage entity from Port
func (c *PortClient) DeleteEntity(blueprint, identifier string) error {
	if c.Token == "" {
		if err := c.Authenticate(); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}
	
	entityURL := fmt.Sprintf("%s/v1/blueprints/%s/entities/%s", c.BaseURL, blueprint, identifier)
	
	req, err := http.NewRequest("DELETE", entityURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete entity: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete entity, status: %d", resp.StatusCode)
	}
	
	log.Info().Str("entity", identifier).Msg("Deleted entity from Port")
	return nil
}

// ReportActionStatus reports the status of a self-service action back to Port
func (c *PortClient) ReportActionStatus(runID, status, message string) error {
	if c.Token == "" {
		if err := c.Authenticate(); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}
	
	statusURL := fmt.Sprintf("%s/v1/actions/runs/%s", c.BaseURL, runID)
	
	payload := map[string]interface{}{
		"status": map[string]interface{}{
			"status":  status,
			"message": message,
		},
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal status payload: %w", err)
	}
	
	req, err := http.NewRequest("PATCH", statusURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create status request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to report status: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to report status, status code: %d", resp.StatusCode)
	}
	
	log.Info().Str("runID", runID).Str("status", status).Msg("Reported action status to Port")
	return nil
}

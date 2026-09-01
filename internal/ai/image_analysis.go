package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"travel-diary-backend/internal/models"
)

type ImageAnalyzer interface {
	Analyze(ctx context.Context, imageURL, model string) (models.TripImageAnalysis, error)
}
type OpenAIImageAnalyzer struct {
	key    string
	client *http.Client
}

func NewOpenAIImageAnalyzer(key string) *OpenAIImageAnalyzer {
	return &OpenAIImageAnalyzer{key: key, client: &http.Client{Timeout: 90 * time.Second}}
}

type chatRequest struct {
	Model          string         `json:"model"`
	Messages       []chatMessage  `json:"messages"`
	ResponseFormat responseFormat `json:"response_format"`
}
type chatMessage struct {
	Role    string        `json:"role"`
	Content []contentPart `json:"content"`
}
type contentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *imageURL `json:"image_url,omitempty"`
}
type imageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail"`
}
type responseFormat struct {
	Type       string `json:"type"`
	JSONSchema schema `json:"json_schema"`
}
type schema struct {
	Name   string         `json:"name"`
	Strict bool           `json:"strict"`
	Schema map[string]any `json:"schema"`
}
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (o *OpenAIImageAnalyzer) Analyze(ctx context.Context, imageURLValue, model string) (models.TripImageAnalysis, error) {
	if o.key == "" {
		return models.TripImageAnalysis{}, fmt.Errorf("OPENAI_KEY is not configured")
	}
	reqBody := chatRequest{Model: model, Messages: []chatMessage{{Role: "user", Content: []contentPart{{Type: "text", Text: "Analyze this travel photograph. Return only JSON. Do not identify people. Normalize every score from 0 to 1."}, {Type: "image_url", ImageURL: &imageURL{URL: imageURLValue, Detail: "low"}}}}}, ResponseFormat: responseFormat{Type: "json_schema", JSONSchema: schema{Name: "travel_image_analysis", Strict: true, Schema: imageSchema()}}}
	data, _ := json.Marshal(reqBody)
	var last error
	for attempt := 0; attempt < 3; attempt++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(data))
		req.Header.Set("Authorization", "Bearer "+o.key)
		req.Header.Set("Content-Type", "application/json")
		res, err := o.client.Do(req)
		if err == nil && res.StatusCode >= 200 && res.StatusCode < 300 {
			var out chatResponse
			err = json.NewDecoder(res.Body).Decode(&out)
			res.Body.Close()
			if err == nil && len(out.Choices) > 0 {
				var a models.TripImageAnalysis
				if err = json.Unmarshal([]byte(out.Choices[0].Message.Content), &a); err == nil {
					return a, nil
				}
			}
		} else if res != nil {
			last = fmt.Errorf("openai request failed with status %d", res.StatusCode)
			res.Body.Close()
		} else {
			last = err
		}
		time.Sleep(time.Duration(1<<attempt) * time.Second)
	}
	return models.TripImageAnalysis{}, last
}
func imageSchema() map[string]any {
	return map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{
		"caption": map[string]any{"type": "string"}, "scene": map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{"primary": map[string]any{"type": "string"}, "secondary": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, "environment": map[string]any{"type": "string"}}, "required": []string{"primary", "secondary", "environment"}},
		"visual":        map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{"timeOfDay": map[string]any{"type": "string"}, "weather": map[string]any{"type": "string"}, "lighting": map[string]any{"type": "string"}, "dominantMood": map[string]any{"type": "string"}}, "required": []string{"timeOfDay", "weather", "lighting", "dominantMood"}},
		"content":       map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{"subjects": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, "peopleCount": map[string]any{"type": "integer"}}, "required": []string{"subjects", "peopleCount"}},
		"story":         map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{"importance": map[string]any{"type": "number"}, "storyRole": map[string]any{"type": "string"}, "emotionalTone": map[string]any{"type": "string"}, "keywords": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}}, "required": []string{"importance", "storyRole", "emotionalTone", "keywords"}},
		"quality":       map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{"aestheticScore": map[string]any{"type": "number"}, "sharpnessScore": map[string]any{"type": "number"}, "exposureScore": map[string]any{"type": "number"}}, "required": []string{"aestheticScore", "sharpnessScore", "exposureScore"}},
		"composition":   map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{"subjectPosition": map[string]any{"type": "string"}, "negativeSpace": map[string]any{"type": "array", "items": map[string]any{"type": "string"}}, "textSafeAreas": map[string]any{"type": "array", "items": map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{"position": map[string]any{"type": "string"}, "confidence": map[string]any{"type": "number"}}, "required": []string{"position", "confidence"}}}}, "required": []string{"subjectPosition", "negativeSpace", "textSafeAreas"}},
		"locationGuess": map[string]any{"type": "object", "additionalProperties": false, "properties": map[string]any{"name": map[string]any{"type": "string"}, "confidence": map[string]any{"type": "number"}}, "required": []string{"name", "confidence"}},
	}, "required": []string{"caption", "scene", "visual", "content", "story", "quality", "composition", "locationGuess"}}
}

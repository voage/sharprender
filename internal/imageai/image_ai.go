package imageai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/voage/sharprender-api/internal/imagescraper"
)

func GetRecommendations(image imagescraper.Image) (*Recommendation, error) {

	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		log.Fatalf("OPENAI_KEY not found")
	}

	client := openai.NewClient(apiKey)

	prompt := generatePrompt(image)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert in image optimization. When you respond, provide only the JSON object specified, with no additional text, explanations, or formatting. Do not include any narrative, preamble, or code blocks.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return nil, err
	}

	recommendation, err := parseResponse(resp.Choices[0].Message.Content)
	if err != nil {
		return nil, err
	}
	return recommendation, nil
}

func generatePrompt(image imagescraper.Image) string {
	prompt := fmt.Sprintf(
		`Given the following image properties:

Source URL: %s
Alt Text: %s
Format: %s
Width: %d pixels
Height: %d pixels
File Size: %d bytes

As an expert in image optimization, provide recommendations on optimizing this image. Make it as descriptive as possible. Output **only** the following JSON object, with no additional text, explanations, or formatting:

{
    "format_recommendations": "string",
    "resize_recommendations": "string",
    "compression_recommendations": "string",
    "caching_recommendations": "string",
	"other_recommendations": "string", 
}

Respond **only** with this JSON structure, without code blocks, extra text, or markdown.`,
		image.Src,
		image.Alt,
		image.Format,
		image.Width,
		image.Height,
		image.Size,
	)
	return prompt
}

func parseResponse(reply string) (*Recommendation, error) {

	jsonStart := strings.Index(reply, "{")
	jsonEnd := strings.LastIndex(reply, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonStr := reply[jsonStart : jsonEnd+1]

	var recommendation Recommendation
	err := json.Unmarshal([]byte(jsonStr), &recommendation)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %v\nResponse was: %s", err, jsonStr)
	}
	return &recommendation, nil
}

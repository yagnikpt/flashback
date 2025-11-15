package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/yagnikpt/flashback/internal/utils"
	"google.golang.org/genai"
)

func (app *App) GenerateEmbeddingForNote(content, taskType string) ([]float32, error) {
	ctx := context.Background()
	contents := []*genai.Content{
		genai.NewContentFromText(content, genai.RoleUser),
	}
	result, err := app.Gemini.Models.EmbedContent(ctx,
		"gemini-embedding-001",
		contents,
		&genai.EmbedContentConfig{
			TaskType:             taskType,
			OutputDimensionality: genai.Ptr[int32](768),
		},
	)
	if err != nil {
		return nil, err
	}
	return result.Embeddings[0].Values, nil
}

func (app *App) GenerateMetadataForSimpleNote(content string) (map[string]string, error) {
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(string(utils.SimpleTextExtractionPrompt), genai.RoleUser),
		ResponseMIMEType:  "application/json",
		ResponseJsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"tldr": map[string]any{
					"type":        "string",
					"description": "Short summary (1 sentence) of the content.",
				},
				"tags": map[string]any{
					"type":        "string",
					"description": "Tags, topics, or keywords related to the text in string of array with [] format.",
				},
			},
			"additionalProperties": false,
		},
	}

	ctx := context.Background()
	result, err := app.Gemini.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(content),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("error generating simple note metadata: %w", err)
	}

	res := map[string]string{}
	err = json.Unmarshal([]byte(result.Text()), &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling simple note metadata: %w", err)
	}

	return res, nil
}

func (app *App) GenerateMetadataForWebNote(content string) (map[string]string, error) {
	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(string(utils.WebExtractionPrompt), genai.RoleUser),
		ResponseMIMEType:  "application/json",
		ResponseJsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"image": map[string]any{
					"type":        "string",
					"description": "Primary image or Open Graph image for the page.",
				},
				"image_main": map[string]any{
					"type":        "string",
					"description": "true or omitted. Whether the image is the main focus of the page.",
				},
				"tldr": map[string]any{
					"type":        "string",
					"description": "Short summary (1 sentence) of the page or its content.",
				},
				"description": map[string]any{
					"type":        "string",
					"description": "Meta description or short contextual summary of the page.",
				},
				"tags": map[string]any{
					"type":        "string",
					"description": "Tags, topics, or keywords related to the page in string of array with [] format.",
				},
			},
			"additionalProperties": false,
		},
	}

	ctx := context.Background()
	result, err := app.Gemini.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(content),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("error generating web metadata: %w", err)
	}
	if result.Text() == "" {
		return nil, fmt.Errorf("empty response received from metadata generation")
	}

	res := map[string]string{}
	err = json.Unmarshal([]byte(result.Text()), &res)
	if err != nil {
		log.Fatal(err)
	}

	image, ok := res["image"]
	if ok {
		imageMetadata, err := app.GenerateMetadataForImage(image)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		priorityImage, imageMainOk := res["image_main"]
		for k, v := range imageMetadata {
			if k == "tags" && imageMainOk && priorityImage == "true" {
				if existingTags, exists := res["tags"]; exists && existingTags != "" {
					var imageTags []string
					var webPageTags []string
					err := json.Unmarshal([]byte(v), &imageTags)
					if err != nil {
						return nil, fmt.Errorf("error unmarshaling image tags: %w", err)
					}
					err = json.Unmarshal([]byte(existingTags), &webPageTags)
					if err != nil {
						return nil, fmt.Errorf("error unmarshaling webpage tags: %w", err)
					}
					combinedTags := append(webPageTags, imageTags...)
					combinedTags = utils.UniqueStrings(combinedTags)
					combinedTagsJson, err := json.Marshal(combinedTags)
					if err != nil {
						return nil, fmt.Errorf("error marshaling combined tags: %w", err)
					}
					res["tags"] = string(combinedTagsJson)
					continue
				}
			} else {
				continue
			}
			if k == "tldr" {
				if existingTldr, exists := res["tldr"]; exists && existingTldr != "" {
					if imageMainOk && priorityImage == "true" {
						res["tldr"] = v
					}
					continue
				}
			}
			res[k] = v
		}
	}

	return res, nil
}

func (app *App) GenerateMetadataForImage(imageUrl string) (map[string]string, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(string(utils.ImageExtractionPrompt), genai.RoleUser),
		ResponseMIMEType:  "application/json",
		ResponseJsonSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"tldr": map[string]any{
					"type":        "string",
					"description": "Short summary (1 sentence) of the image.",
				},
				"tags": map[string]any{
					"type":        "string",
					"description": "Tags, topics, or keywords related to the image in string of array with [] format.",
				},
			},
			"additionalProperties": false,
		},
	}
	parts := []*genai.Part{
		{InlineData: &genai.Blob{Data: data, MIMEType: "image/jpeg"}},
	}
	contents := []*genai.Content{{Parts: parts}}
	ctx := context.Background()
	result, err := app.Gemini.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, config)
	if err != nil {
		errrr, _ := err.(genai.APIError)
		log.Println(errrr)
		return nil, fmt.Errorf("%s", errrr.Status)
	}
	res := map[string]string{}
	err = json.Unmarshal([]byte(result.Text()), &res)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error unmarshaling image metadata: %w", err)
	}

	return res, nil
}

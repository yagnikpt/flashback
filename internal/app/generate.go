package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yagnikpt/flashback/internal/utils"
	"google.golang.org/genai"
)

func (app *App) GenerateEmbeddingForNote(ctx context.Context, content, taskType string) ([]float32, error) {
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

func (app *App) GenerateMetadataForSimpleNote(ctx context.Context, content string) (map[string]string, error) {
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
					"description": "Tags, topics, or keywords related to the text in string of array with [] format. Omit if none found.",
				},
			},
			"additionalProperties": false,
		},
	}

	result, err := app.Gemini.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(content),
		config,
	)
	if err != nil {
		return nil, err
	}

	res := map[string]string{}
	err = json.Unmarshal([]byte(result.Text()), &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling simple note metadata: %w", err)
	}

	return res, nil
}

func (app *App) GenerateMetadataForWebNote(ctx context.Context, content string) (map[string]string, error) {
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

	result, err := app.Gemini.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(content),
		config,
	)
	if err != nil {
		errs := err.(genai.APIError)
		return nil, genai.APIError{
			Code:    errs.Code,
			Message: errs.Message,
			Status:  errs.Status,
			Details: errs.Details,
		}
	}
	if result.Text() == "" {
		return nil, fmt.Errorf("empty response received from metadata generation")
	}

	res := map[string]string{}
	err = json.Unmarshal([]byte(result.Text()), &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling metadata")
	}

	image, ok := res["image"]
	imageMain, imageMainOk := res["image_main"]
	if ok && imageMainOk && imageMain == "true" {
		imageMetadata, err := app.GenerateMetadataForImage(ctx, image)
		if err != nil {
			errs := err.(genai.APIError)
			return nil, genai.APIError{
				Code:    errs.Code,
				Message: errs.Message,
				Status:  errs.Status,
				Details: errs.Details,
			}
		}
		for k, v := range imageMetadata {
			if k == "tags" {
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
			}
			res[k] = v
		}
	}

	return res, nil
}

func (app *App) GenerateMetadataForImage(ctx context.Context, imageUrl string) (map[string]string, error) {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
	result, err := app.Gemini.Models.GenerateContent(ctx, "gemini-2.0-flash", contents, config)
	if err != nil {
		return nil, err
	}
	res := map[string]string{}
	err = json.Unmarshal([]byte(result.Text()), &res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling image metadata: %w", err)
	}

	return res, nil
}

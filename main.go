package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/dotprompt"
	"github.com/firebase/genkit/go/plugins/googleai"
	"github.com/invopop/jsonschema"
)

type promptInput struct {
	URL string `json:"url"`
}

func main() {
	ctx := context.Background()

	// Initialize the Google AI plugin
	if err := googleai.Init(ctx, nil); err != nil {
		log.Fatalf("Failed to initialize Google AI plugin: %v", err)
	}

	// Define the webLoader tool
	webLoader := ai.DefineTool(
		"webLoader",
		"Loads a webpage and returns the textual content.",
		func(ctx context.Context, input struct {
			URL string `json:"url"`
		}) (string, error) {
			return fetchWebContent(input.URL)
		},
	)

	model := googleai.Model("gemini-1.5-flash")

	summarizePrompt, err := dotprompt.Define("summarizePrompt",
		"First, fetch this link: {{url}}. Then, summarize the content within 20 words.",
		dotprompt.Config{
			Model: model,
			Tools: []ai.Tool{webLoader},
			GenerationConfig: &ai.GenerationCommonConfig{
				Temperature: 1,
			},
			InputSchema:  jsonschema.Reflect(promptInput{}),
			OutputFormat: ai.OutputFormatText,
		},
	)
	if err != nil {
		log.Fatalf("Failed to initialize prompt: %v", err)
	}

	// Define a flow that fetches a webpage and summarizes its content
	genkit.DefineFlow("summarizeFlow", func(ctx context.Context, input string) (string, error) {
		resp, err := summarizePrompt.Generate(ctx,
			&dotprompt.PromptRequest{
				Variables: &promptInput{
					URL: input,
				},
			},
			nil,
		)
		if err != nil {
			return "", fmt.Errorf("failed to generate summary: %w", err)
		}
		return resp.Text(), nil
	})

	// Initialize Genkit
	if err := genkit.Init(ctx, nil); err != nil {
		log.Fatalf("Failed to initialize Genkit: %v", err)
	}
}

// fetchWebContent fetches and processes the content from the provided URL
func fetchWebContent(url string) (string, error) {
	// Fetch the content from the provided URL
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer res.Body.Close()

	// Read the HTML content
	html, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Load the HTML content into goquery for parsing
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Remove unnecessary elements
	doc.Find("script, style, noscript").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	// Prefer 'article' content, fallback to 'body' if not available
	article := doc.Find("article").Text()
	if article != "" {
		return article, nil
	}

	body := doc.Find("body").Text()
	return body, nil
}

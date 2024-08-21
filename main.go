package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googleai"
)

func main() {
	ctx := context.Background()

	// Initialize environment variables
	initEnv()

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

	// Define a flow that fetches a webpage and summarizes its content
	genkit.DefineFlow("summarizeFlow", func(ctx context.Context, input string) (string, error) {
		m := googleai.Model("gemini-1.5-flash")
		resp, err := ai.Generate(
			ctx,
			m,
			ai.WithConfig(&ai.GenerationCommonConfig{Temperature: 1}),
			ai.WithTextPrompt(fmt.Sprintf(`First, fetch this link: %s. Then, summarize the content within 20 words.`, input)),
			ai.WithTools(webLoader),
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

// initEnv initializes environment variables
func initEnv() {
	if os.Getenv("GOOGLE_GENAI_API_KEY") == "" {
		log.Fatal("GOOGLE_GENAI_API_KEY environment variable is not set")
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

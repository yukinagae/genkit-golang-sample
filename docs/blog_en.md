## Introduction

When developing LLM applications, JavaScript/TypeScript and Python are often the primary programming language choices. However, it's now possible to develop LLM applications using Go as well. This article will introduce how to do that. As a simple example, we'll create an API that summarizes web content when given a URL.

The final product is equivalent to the TypeScript implementation below, so reading it will deepen your understanding of Firebase Genkit:

https://medium.com/@yukinagae/your-first-guide-to-getting-started-with-firebase-genkit-6948d88e8a92

Note: As Genkit's Go support is in alpha, specifications may change significantly or contain bugs, so please don't use it in production environments yet. Enjoy it for private or verification purposes only!

## What is Firebase Genkit?

[Firebase Genkit](https://firebase.google.com/docs/genkit) is an open-source framework that helps developers create AI-powered applications. Although developed by the Firebase team, it doesn't depend on Firebase services, so it can be used on any platform including AWS, GCP, or Azure.

There are various LLMs available recently, and using them appropriately for different purposes can yield better results. Typically, to execute an LLM, you either directly call a REST API or use an SDK. For example, you might use [Gemini](https://ai.google.dev/) for input A and [OpenAI](https://openai.com/api/) for input B, depending on the use case.

Firebase Genkit makes it easy to switch between these use cases. It also provides a developer UI, offering an all-in-one solution for developers to adjust LLM parameters, debug, and troubleshoot.

## Using Firebase Genkit with Go

The final product is available in the following repository, so please refer to it:

https://github.com/yukinagae/genkit-golang-sample

### Setup

- Go: Install Go by following [Go - Download and Install](https://go.dev/doc/install).
- Genkit: Install Genkit by following [Firebase Genkit - Getting Started](https://firebase.google.com/docs/genkit-go/get-started-go).

Check the versions of Go and Genkit with the following commands:

```bash
$ go version
go version go1.23.2 darwin/arm64
$ genkit --version
0.9.1
```

You can either create your own Genkit project based on the official documentation [Get started with Genkit using Go (alpha)](https://firebase.google.com/docs/genkit-go/get-started-go) and modify the code, or clone the sample project:

```bash
$ git clone https://github.com/yukinagae/genkit-golang-sample
```

### Running Locally

This example uses [Gemini](https://ai.google.dev/), so obtain an API key in advance.

You can set the API key as an environment variable and start the Genkit server with the following commands:

```bash
$ export GOOGLE_GENAI_API_KEY=your_api_key
$ genkit start -o -- go run main.go
```

The `-o` option automatically opens a browser, allowing you to access the developer UI at http://localhost:4000 to input prompts and debug.

Issues for the developer UI are managed on [firebase/genkit](https://github.com/firebase/genkit), and development seems active, so it may improve significantly in the future.
For more on the developer UI and basic Genkit usage, please refer to my previous article.

https://medium.com/@yukinagae/your-first-guide-to-getting-started-with-firebase-genkit-6948d88e8a92

## Code Explanation

### 1. Initializing Google API Plugin and Getting the Model

Each LLM is provided as a plugin, and can be used by initializing and obtaining the model as follows.
This example obtains the Gemini model:

```go
import "github.com/firebase/genkit/go/plugins/googleai"

// Initialize the Google AI plugin
if err := googleai.Init(ctx, nil); err != nil {
	log.Fatalf("Failed to initialize Google AI plugin: %v", err)
}

// Get the model
model := googleai.Model("gemini-1.5-flash")
// or
model := googleai.Model("gemini-1.5-pro")
```

If you want to use [OpenAI](https://openai.com/api/), import the dedicated plugin and initialize/obtain the model as follows:

```go
import "github.com/yukinagae/genkit-go-plugins/plugins/openai"

// Initialize the OpenAI plugin
if err := openai.Init(ctx, nil); err != nil {
	log.Fatalf("Failed to initialize OpenAI plugin: %v", err)
}

// Get the model
model := openai.Model("gpt-4o-mini")
// or
model := openai.Model("gpt-4o")
```

### 2. Defining the Prompt

Use the `dotprompt.Define()` function to define the prompt. Here you can specify what processing you want to do with the prompt, input/output schemas, and parameters like `Temperature`.

```go
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
```

Genkit also supports Function Calling, so you can define tools using the `ai.DefineTool()` function.

```go
// Define the webLoader tool
webLoader := ai.DefineTool(
	"webLoader",
	"Loads a webpage and returns the textual content.",
	func(ctx context.Context, input promptInput) (string, error) {
		return fetchWebContent(input.URL)
	},
)
```

### 3. Defining the Flow

Use the `genkit.DefineFlow()` function to define a `Flow`. In Genkit, a `Flow` is a single function that combines multiple processes, intended to be called as an API. Here, we define a flow that summarizes the content when given a web content URL, but in actual use cases, it can express more complex processes.

```go
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
```

### 4. Starting the Genkit Server

You can initialize Genkit and start the server using the `genkit.Init()` function.
By default, the server starts on `127.0.0.1:3400`.

You can change settings like the server address by passing `&genkit.Options{}` as an argument.

```go
// Initialize Genkit
if err := genkit.Init(ctx, nil); err != nil {
	log.Fatalf("Failed to initialize Genkit: %v", err)
}
```

### 5. Calling the API

The API server is running on `127.0.0.1:3400`, so let's access it with curl.
If you get a result like this, it's successful:

```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"data": "https://firebase.google.com/docs/genkit-go/get-started-go"}' http://127.0.0.1:3400/summarizeFlow
{"result": "Firebase Genkit for Go is in alpha and helps you build AI features with Go. \n"}
```

## Conclusion

Using Firebase Genkit, we were able to develop an LLM application in Go. If you want to know more, read the [official documentation](https://firebase.google.com/docs/genkit-go/get-started-go). Or check out the [Awesome list](https://github.com/xavidop/awesome-firebase-genkit) which has many articles, videos, and samples for reference.

If you have any questions or feedback, please use [Discord](https://discord.gg/qXt5zzQKpc) or [GitHub](https://github.com/firebase/genkit/issues).

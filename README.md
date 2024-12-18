# genkit-golang-sample

`genkit-golang-sample` is a sample repository to help you learn how to build LLM applications with Firebase Genkit using Golang.

- [Requirements](#requirements)
- [Usage](#usage)
- [License](#license)

## Requirements

- **Go**: Follow the [Go - Download and install](https://go.dev/doc/install) to install Go.
- **Genkit**: Follow the [Firebase Genkit - Get started](https://firebase.google.com/docs/genkit/get-started) to install Genkit.

Verify your installations:

```bash
$ go version
go version go1.23.2 darwin/arm64
$ genkit --version
0.9.1
```

## Usage

Set your API key and start Genkit:

```bash
$ export GOOGLE_GENAI_API_KEY=your_api_key
$ genkit start -o -- go run main.go # Starts Genkit
```

Open your browser and navigate to [http://localhost:4001](http://localhost:4001) to access the Genkit UI.

## License

MIT

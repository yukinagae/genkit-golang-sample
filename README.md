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
v22.4.1
$ genkit --version
0.5.4
```

## Usage

Set your API key and start Genkit:

```bash
$ export GOOGLE_GENAI_API_KEY=your_api_key
$ make genkit # Starts Genkit
```

Open your browser and navigate to [http://localhost:4000](http://localhost:4000) to access the Genkit UI.

## License

MIT

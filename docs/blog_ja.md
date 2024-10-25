# Firebase Genkit を活用して Go で LLM 入門

## はじめに

LLM アプリを開発する際のプログラミング言語の選択としては、主に JavaScript/TypeScript、Python などが挙げられます。しかし、Go でも LLM アプリを開発できるようになったので、本記事ではその方法を紹介します。簡易的な例として、Web コンテンツの URL を渡すとその内容を要約する API を作成します。

最終的な成果物は以下の TypeScript による実装と等価のため、こちらも読むとより Firebase Genkit の理解が深まると思います。

[Zenn - Firebase Genkit 入門](https://zenn.dev/cureapp/articles/ab5382ce510c8c)

※注意事項: 本記事は 2024-10-25 時点の Firebase Genkit のドキュメントに基づいています。また、Genkit の Go サポートは alpha のため、仕様が大幅変更されたり不具合を含む可能性があるため、本番環境ではまだ使用しないでください。あくまでもプライベートや検証用とで楽しんでください！

## Firebase Genkit とは？

[Firebase Genkit](https://firebase.google.com/docs/genkit) は、開発者が AI 搭載のアプリケーションを作成するのを支援するオープンソースフレームワークです。Firebase チームが開発していますが、Firebase のサービスに依存しないため、Firebase 以外の AWS/GCP/Azure いずれのプラットフォームでも利用できます。

最近では様々な LLM があり、用途によって使い分けることでより良い結果を得ることができます。通常、LLM を実行するには、REST API を直接叩くか、SDK を使用するかのいずれかです。例えば、ある入力 A については [Gemini](https://ai.google.dev/) を使用し、別の入力 B については [OpenAI](https://openai.com/api/) を使用するというように、用途によって使い分ける必要があります。

Firebase Genkit は、このような用途による使い分けを容易にすることができます。また、開発者用 UI が提供されているため、開発者が LLM のパラメータを調整したり、デバッグやトラブルシューティングをサポートする仕組みがオールインワンで提供されています。

## Go で Firebase Genkit を使う

最終的な成果物は以下のリポジトリにあるので、ぜひ参考にしてください。

https://github.com/yukinagae/genkit-golang-sample

### セットアップ

- Go: [Go - ダウンロードとインストール](https://go.dev/doc/install) に従って Go をインストールしてください。
- Genkit: [Firebase Genkit - はじめに](https://firebase.google.com/docs/genkit-go/get-started-go) に従って Genkit をインストールしてください。

以下のコマンドで Go と Genkit のバージョンを確認してください。

```bash
$ $ go version
go version go1.23.2 darwin/arm64
$ genkit --version
0.5.10
```

公式ドキュメントの [Get started with Genkit using Go (alpha) ](https://firebase.google.com/docs/genkit-go/get-started-go) を元に自分で Genkit プロジェクトを作成してコードをいじってもいいですし、サンプルプロジェクトをクローンしてもいいです。

```bash
$ git clone https://github.com/yukinagae/genkit-golang-sample
```

### ローカルで起動してみる

今回の例では [Gemini](https://ai.google.dev/) を使用するので、事前に API キーを取得しておいてください。

以下のコマンドで環境変数に API キーを設定し、Genkit サーバーを起動することができます。

```bash
$ export GOOGLE_GENAI_API_KEY=your_api_key
$ genkit start -o # Starts Genkit
```

`-o` オプションをつけることで、自動的にブラウザが開き、http://localhost:4000 にアクセスすることで開発者用 UI でプロンプトを入力したり、デバッグしたりできます。

開発者用 UI の issue は [firebase/genkit](https://github.com/firebase/genkit) で管理されており、開発が活発な様子なので、今後大幅に改善される可能性があります。
開発者用 UI や基礎的な Genkit の使い方については、私の以前の記事を参考にしてください。

[Zenn - Firebase Genkit 入門](https://zenn.dev/cureapp/articles/ab5382ce510c8c)

## コードの説明

### 1. Google API プラグインの初期化

```go
// Initialize the Google AI plugin
if err := googleai.Init(ctx, nil); err != nil {
	log.Fatalf("Failed to initialize Google AI plugin: %v", err)
}
```

### 2. モデル・ツール・プロンプトの定義

```go
// Define the webLoader tool
webLoader := ai.DefineTool(
	"webLoader",
	"Loads a webpage and returns the textual content.",
	func(ctx context.Context, input promptInput) (string, error) {
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
```

### 3. フローの定義

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

### 4. Genkit サーバーの起動

```go
// Initialize Genkit
if err := genkit.Init(ctx, nil); err != nil {
	log.Fatalf("Failed to initialize Genkit: %v", err)
}
```

## まとめ

Firebase Genkit を使うことで、Go でも LLM アプリケーションを開発することができました。もっと詳しく知りたい人は、[公式ドキュメント](https://firebase.google.com/docs/genkit-go/get-started-go) を読んでみてください。または[Awesome リスト](https://github.com/xavidop/awesome-firebase-genkit)にも多くの記事や動画、サンプルがあるので参考にしてください。

もし質問やフィードバックがあれば、[Discord](https://discord.gg/qXt5zzQKpc) や [GitHub](https://github.com/firebase/genkit/issues) を活用してください。

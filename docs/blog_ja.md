## はじめに

LLM アプリを開発する際のプログラミング言語の選択としては、主に JavaScript/TypeScript、Python などが挙げられます。しかし、Go でも LLM アプリを開発できるようになったので、本記事ではその方法を紹介します。簡易的な例として、Web コンテンツの URL を渡すとその内容を要約する API を作成します。

最終的な成果物は以下の TypeScript による実装と等価のため、こちらも読むとより Firebase Genkit の理解が深まると思います。

https://zenn.dev/cureapp/articles/ab5382ce510c8c

※注意事項: Genkit の Go サポートは alpha のため、仕様が大幅変更されたり不具合を含む可能性があるため、本番環境ではまだ使用しないでください。あくまでもプライベートや検証用として楽しんでください！

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
0.9.1
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
$ genkit start -o -- go run main.go
```

`-o` オプションをつけることで、自動的にブラウザが開き、http://localhost:4000 にアクセスすることで開発者用 UI でプロンプトを入力したり、デバッグしたりできます。

開発者用 UI の issue は [firebase/genkit](https://github.com/firebase/genkit) で管理されており、開発が活発な様子なので、今後大幅に改善される可能性があります。
開発者用 UI や基礎的な Genkit の使い方については、私の以前の記事を参考にしてください。

https://zenn.dev/cureapp/articles/ab5382ce510c8c

## コードの説明

### 1. Google API プラグインの初期化・モデル取得

各 LLM はプラグインとして提供されており、以下のように初期化・モデル取得により使用することができます。
ここの例では Gemini のモデルを取得しています。

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

もし [OpenAI](https://openai.com/api/) を使いたい場合には、以下のように専用プラグインを import し、初期化・モデル取得を行います。

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

### 2. プロンプトの定義

`dotprompt.Define()` という関数を使ってプロンプトを定義します。ここでどのような処理をしたいかプロンプトで指示したり、入出力のスキーマや Temperature などのパラメータを指定します。

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

また Genkit は Function Calling をサポートしているため、`ai.DefineTool()` という関数を使ってツールを定義することができます。

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

### 3. フローの定義

`genkit.DefineFlow()` という関数を使って `フロー` を定義します。Genkit における`フロー` は複数の処理をまとめた単一の関数で、API として呼び出す際にはこのフローを呼び出すことを想定しています。ここでは、Web コンテンツの URL を渡すとその内容を要約するフローを定義していますが、実際のユースケースではより複雑な処理を表現することができます。

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

`genkit.Init()` という関数を使って Genkit を初期化およびサーバー起動することができます。
初期値では `127.0.0.1:3400` でサーバーが起動します。

引数で `&genkit.Options{})` を渡すことで、サーバーのアドレスなどの設定を変更することができます。

```go
// Initialize Genkit
if err := genkit.Init(ctx, nil); err != nil {
	log.Fatalf("Failed to initialize Genkit: %v", err)
}
```

### 5. API の呼び出し

API サーバーは `127.0.0.1:3400` で起動しているので、curl でアクセスしてみます。
以下のような結果が返ってくれば成功です。

```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"data": "https://firebase.google.com/docs/genkit-go/get-started-go"}' http://127.0.0.1:3400/summarizeFlow
{"result": "Firebase Genkit for Go is in alpha and helps you build AI features with Go. \n"}
```

## まとめ

Firebase Genkit を使うことで、Go でも LLM アプリケーションを開発することができました。もっと詳しく知りたい人は、[公式ドキュメント](https://firebase.google.com/docs/genkit-go/get-started-go) を読んでみてください。または [Awesome リスト](https://github.com/xavidop/awesome-firebase-genkit)にも多くの記事や動画、サンプルがあるので参考にしてください。

もし質問やフィードバックがあれば、[Discord](https://discord.gg/qXt5zzQKpc) や [GitHub](https://github.com/firebase/genkit/issues) を活用してください。

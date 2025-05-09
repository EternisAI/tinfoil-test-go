package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/tinfoilsh/tinfoil-go"
)

func main() {
	// Create a client for a specific enclave and model repository
	client, err := tinfoil.NewClientWithParams("llama3-3-70b.model.tinfoil.sh", "tinfoilsh/confidential-llama3-3-70b",
		option.WithAPIKey("xxx"),
	)
	if err != nil {
		panic(err)
	}

	// Create chat completion
	completion, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: openai.F("llama3-3-70b"),
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Hello!"),
			}),
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(completion.Choices[0].Message.Content)
}

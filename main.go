package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
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

	tools := []openai.ChatCompletionToolParam{
		{
			Type: "function",
			Function: openai.FunctionDefinitionParam{
				Name:        "image_tool",
				Description: param.NewOpt("This tool creates an image based on the user prompt"),
				Parameters: openai.FunctionParameters{
					"type": "object",
					"properties": map[string]any{
						"prompt": map[string]string{
							"type": "string",
						},
					},
					"required": []string{"prompt"},
				},
			},
		},
	}

	// Create chat completion
	completion, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: "llama3-3-70b",
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Hello!"),
			},
			Tools: tools,
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("completion", completion.Choices[0].Message.Content)

	stream := client.Chat.Completions.NewStreaming(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: "llama3-3-70b",
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Hello!"),
			},
			Tools: tools,
		},
	)

	for stream.Next() {
		chunk := stream.Current()
		fmt.Println("streamchunk", chunk)
	}

}

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
				openai.UserMessage("What is your name?"),
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
				openai.UserMessage("What is your name?"),
			},
			Tools: tools,
		},
	)
	acc := openai.ChatCompletionAccumulator{}

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if content, ok := acc.JustFinishedContent(); ok {
			println("Content stream finished:", content)
		}

		// if using tool calls
		if tool, ok := acc.JustFinishedToolCall(); ok {
			println("Tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			println("Refusal stream finished:", refusal)
		}

		if len(chunk.Choices) > 0 {
			println(chunk.Choices[0].Delta.Content)
		}
	}

	if stream.Err() != nil {
		panic(stream.Err())
	}

	_ = acc.Choices[0].Message.Content

}

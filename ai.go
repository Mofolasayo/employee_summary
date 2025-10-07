package main

import (
    "context"
    openai "github.com/sashabaranov/go-openai"
    "os"
)

func SummarizeText(text string) string {
    apiKey := os.Getenv("OPENAI_API_KEY")
    client := openai.NewClient(apiKey)

    resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
        Model: "gpt-4o-mini",
        Messages: []openai.ChatCompletionMessage{
            {Role: "system", Content: "You are an assistant that summarizes weekly employee reports."},
            {Role: "user", Content: text},
        },
    })

    if err != nil || len(resp.Choices) == 0 {
        return "Error generating summary: " + err.Error()
    }

    return resp.Choices[0].Message.Content
}

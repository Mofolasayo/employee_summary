package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func SummarizeText(text string) string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "Missing GEMINI_API_KEY in environment"
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=%s", apiKey)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": "Summarize the following weekly employee report clearly, detailed, and concisely:\n\n" + text},
				},
			},
		},
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "Request error: " + err.Error()
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("API error: %s\n%s", resp.Status, string(body))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "JSON parse error: " + err.Error()
	}

	candidates, ok := data["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "No candidates returned from Gemini API"
	}

	first := candidates[0].(map[string]interface{})
	content := first["content"].(map[string]interface{})
	parts := content["parts"].([]interface{})
	textOut := parts[0].(map[string]interface{})["text"].(string)
	return textOut
}

package main

import (
    "fmt"
    "os"
    "github.com/robfig/cron/v3"
)

func StartScheduler() {
    c := cron.New()
    c.AddFunc("@weekly", func() {
        files, _ := os.ReadDir("uploads")
        combined := ""
        for _, f := range files {
            content, _ := os.ReadFile("uploads/" + f.Name())
            combined += string(content) + "\n\n"
        }
        if combined == "" {
            fmt.Println("No reports found this week.")
            return
        }
        summary := SummarizeText(combined)
        fmt.Println("Generated weekly summary:\n", summary)
        SendEmail(summary)
    })
    c.Start()
}

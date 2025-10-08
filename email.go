package main

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(summary string) {
	from := os.Getenv("EMAIL_FROM")
	pass := os.Getenv("EMAIL_PASSWORD")
	to := os.Getenv("EMAIL_TO")

	if from == "" || pass == "" || to == "" {
		fmt.Println("Missing email environment variables")
		return
	}

	msg := "Subject: Weekly Summary Report\n\n" + summary

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from,
		[]string{to},
		[]byte(msg),
	)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return
	}

	fmt.Println("Email sent to", to)
}

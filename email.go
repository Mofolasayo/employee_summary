package main

import (
    "gopkg.in/gomail.v2"
    "os"
)

func SendEmail(summary string) {
    from := os.Getenv("EMAIL_FROM")
    to := os.Getenv("EMAIL_TO")
    password := os.Getenv("EMAIL_PASSWORD")

    m := gomail.NewMessage()
    m.SetHeader("From", from)
    m.SetHeader("To", to)
    m.SetHeader("Subject", "Weekly Work Summary")
    m.SetBody("text/plain", summary)

    d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

    if err := d.DialAndSend(m); err != nil {
        panic(err)
    }
}

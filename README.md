# Employee Weekly Summary Server

This Go-based microservice allows employees to upload weekly work summaries.  
It automatically summarizes their reports using **Gemini AI** and sends a simulated email summary (or real email if SMTP credentials are configured).

---

## Features
- Upload employee reports through `/upload`
- Generates summaries using **Gemini AI (via API)**
- Displays results instantly in the terminal
- Includes a weekly scheduler to summarize all reports
- Configurable `.env` file for API keys and credentials

---

## Setup

### 1. Clone the repository
```bash
git clone https://github.com/YOUR_USERNAME/employee-summary.git
cd employee-summary

---

## Swagger

https://deonna-maximal-kyleigh.ngrok-free.dev/swagger/index.html#/default/get_summaries

## Upload an employee report
curl -X POST -F "employee=John" -F "summary=@test.txt" \
https://deonna-maximal-kyleigh.ngrok-free.dev/upload


##
The app uses a built-in cron scheduler that automatically generates and sends all summaries every week (you can adjust the schedule in scheduler.go).
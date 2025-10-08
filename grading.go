package main

import (
	"fmt"
	"time"
)

func EvaluateEmployees() string {
	records := GetAllRecords()
	if len(records) == 0 {
		return "No employee submissions found this week."
	}

	report := "Employee Performance Evaluation:\n\n"

	for _, r := range records {
		daysSinceLast := int(time.Since(r.LastSubmitted).Hours() / 24)
		status := "Serious"
		reason := fmt.Sprintf("Submitted %d reports so far", r.TotalReports)

		if daysSinceLast > 10 {
			status = "Unserious"
			reason = fmt.Sprintf("No report for %d days", daysSinceLast)
		}

		report += fmt.Sprintf("👤 %s — %s (%s)\n", r.Name, status, reason)
	}

	return report
}

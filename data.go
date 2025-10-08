package main

import (
	"encoding/json"
	"os"
	"time"
)

type EmployeeRecord struct {
	Name          string    `json:"name"`
	Submissions   []string  `json:"submissions"`
	LastSubmitted time.Time `json:"last_submitted"`
	TotalReports  int       `json:"total_reports"`
}

func SaveSubmission(name, filename string) {
	var records []EmployeeRecord
	file := "employees.json"

	data, _ := os.ReadFile(file)
	json.Unmarshal(data, &records)

	found := false
	for i, r := range records {
		if r.Name == name {
			r.Submissions = append(r.Submissions, filename)
			r.LastSubmitted = time.Now()
			r.TotalReports++
			records[i] = r
			found = true
			break
		}
	}

	if !found {
		records = append(records, EmployeeRecord{
			Name:          name,
			Submissions:   []string{filename},
			LastSubmitted: time.Now(),
			TotalReports:  1,
		})
	}

	out, _ := json.MarshalIndent(records, "", "  ")
	os.WriteFile(file, out, 0644)
}

func GetAllRecords() []EmployeeRecord {
	var records []EmployeeRecord
	data, _ := os.ReadFile("employees.json")
	json.Unmarshal(data, &records)
	return records
}

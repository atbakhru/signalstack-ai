package models

import "time"

type EvaluationQuestion struct {
	ID               int64     `json:"id"`
	Question         string    `json:"question"`
	ExpectedSource   string    `json:"expected_source"`
	ExpectedKeywords []string  `json:"expected_keywords"`
	CreatedAt        time.Time `json:"created_at"`
}

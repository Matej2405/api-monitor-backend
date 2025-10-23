package models

import "time"

type APIRequest struct {
	ID           int       `json:"id"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	ResponseCode int       `json:"response_code"`
	ResponseTime int       `json:"response_time"`
	ResponseBody string    `json:"response_body,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type Problem struct {
	ID          int       `json:"id"`
	RequestID   int       `json:"request_id"`
	ProblemType string    `json:"problem_type"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

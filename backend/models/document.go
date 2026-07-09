package models

import "time"

type Document struct {
	ID          int64      `json:"id"`
	SourceID    int64      `json:"source_id"`
	ExternalID  string     `json:"external_id"`
	Title       string     `json:"title"`
	Summary     string     `json:"summary"`
	URL         string     `json:"url"`
	Author      string     `json:"author"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	RawS3Key    string     `json:"raw_s3_key"`
	ContentHash string     `json:"content_hash"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type NormalizedDocument struct {
	Source      string            `json:"source"`
	ExternalID  string            `json:"external_id"`
	Title       string            `json:"title"`
	Summary     string            `json:"summary"`
	URL         string            `json:"url"`
	Author      string            `json:"author"`
	PublishedAt *time.Time        `json:"published_at,omitempty"`
	Content     string            `json:"content"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

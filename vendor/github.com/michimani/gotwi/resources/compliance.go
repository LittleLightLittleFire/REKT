package resources

import "time"

type ComplianceType string

const (
	ComplianceTypeTweets ComplianceType = "tweets"
	ComplianceTypeUsers  ComplianceType = "users"
)

type Compliance struct {
	ID                string         `json:"id"`
	Resumable         bool           `json:"resumable"`
	Status            string         `json:"status"`
	CreatedAt         *time.Time     `json:"created_at"`
	Type              ComplianceType `json:"type"`
	Name              string         `json:"name"`
	UploadURL         string         `json:"upload_url"`
	UploadExpiresAt   *time.Time     `json:"upload_expires_at"`
	DownloadURL       string         `json:"download_url"`
	DownloadExpiresAt *time.Time     `json:"download_expires_at"`
}

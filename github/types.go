package github

import "time"

// Package represents a GitHub package
type Package struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	PackageType string    `json:"package_type"`
	Visibility  string    `json:"visibility"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Repository  *struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository,omitempty"`
	VersionCount int `json:"version_count"`
}

// PackageVersion represents a version of a GitHub package
type PackageVersion struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	PackageHTMLURL string    `json:"package_html_url"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Metadata       struct {
		PackageType string `json:"package_type"`
		Container   struct {
			Tags []string `json:"tags"`
		} `json:"container"`
	} `json:"metadata"`
}

// Tags returns the tags for container images
func (v *PackageVersion) Tags() []string {
	return v.Metadata.Container.Tags
}

// TagsString returns tags as a comma-separated string
func (v *PackageVersion) TagsString() string {
	tags := v.Tags()
	if len(tags) == 0 {
		return "<untagged>"
	}
	result := tags[0]
	for i := 1; i < len(tags) && i < 3; i++ {
		result += ", " + tags[i]
	}
	if len(tags) > 3 {
		result += "..."
	}
	return result
}

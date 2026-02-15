package github

import "testing"

func TestPackageVersion_Tags(t *testing.T) {
	tests := []struct {
		name     string
		version  PackageVersion
		expected []string
	}{
		{
			name:     "returns tags from metadata",
			version:  PackageVersion{Metadata: struct{ PackageType string `json:"package_type"`; Container struct{ Tags []string `json:"tags"` } `json:"container"` }{Container: struct{ Tags []string `json:"tags"` }{Tags: []string{"v1.0", "latest"}}}},
			expected: []string{"v1.0", "latest"},
		},
		{
			name:     "returns empty slice when no tags",
			version:  PackageVersion{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.version.Tags()
			if len(result) != len(tt.expected) {
				t.Errorf("Tags() = %v, want %v", result, tt.expected)
				return
			}
			for i, tag := range result {
				if tag != tt.expected[i] {
					t.Errorf("Tags()[%d] = %v, want %v", i, tag, tt.expected[i])
				}
			}
		})
	}
}

func TestPackageVersion_TagsString(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		expected string
	}{
		{
			name:     "empty tags returns untagged",
			tags:     nil,
			expected: "<untagged>",
		},
		{
			name:     "single tag",
			tags:     []string{"v1.0"},
			expected: "v1.0",
		},
		{
			name:     "two tags",
			tags:     []string{"v1.0", "latest"},
			expected: "v1.0, latest",
		},
		{
			name:     "three tags",
			tags:     []string{"v1.0", "v1", "latest"},
			expected: "v1.0, v1, latest",
		},
		{
			name:     "more than three tags truncates",
			tags:     []string{"v1.0", "v1", "latest", "stable"},
			expected: "v1.0, v1, latest...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &PackageVersion{}
			v.Metadata.Container.Tags = tt.tags
			result := v.TagsString()
			if result != tt.expected {
				t.Errorf("TagsString() = %q, want %q", result, tt.expected)
			}
		})
	}
}



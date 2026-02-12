package updater

import (
	"reflect"
	"testing"

	"github.com/blang/semver"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    semver.Version
		wantErr bool
	}{
		{
			name:  "clean version",
			input: "1.2.3",
			want:  semver.MustParse("1.2.3"),
		},
		{
			name:  "version with v prefix",
			input: "v1.2.3",
			want:  semver.MustParse("1.2.3"),
		},
		{
			name:  "version with pre-release",
			input: "1.2.3-beta.1",
			want:  semver.MustParse("1.2.3-beta.1"),
		},
		{
			name:  "version with v prefix and pre-release",
			input: "v1.2.3-beta.1",
			want:  semver.MustParse("1.2.3-beta.1"),
		},
		{
			name:    "invalid version",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

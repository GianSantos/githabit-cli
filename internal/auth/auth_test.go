package auth

import (
	"testing"
)

func TestCheckRequiredScopes(t *testing.T) {
	tests := []struct {
		name    string
		scopes  string
		wantErr bool
	}{
		{
			name:    "all scopes present",
			scopes:  "repo, read:user, read:org",
			wantErr: false,
		},
		{
			name:    "scopes with extra",
			scopes:  "repo, read:user, read:org, workflow, write:discussion",
			wantErr: false,
		},
		{
			name:    "missing repo",
			scopes:  "read:user, read:org",
			wantErr: true,
		},
		{
			name:    "missing read:user",
			scopes:  "repo, read:org",
			wantErr: true,
		},
		{
			name:    "missing read:org",
			scopes:  "repo, read:user",
			wantErr: true,
		},
		{
			name:    "empty scopes",
			scopes:  "",
			wantErr: true,
		},
		{
			name:    "single scope only",
			scopes:  "repo",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckRequiredScopes(tt.scopes)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckRequiredScopes(%q) error = %v, wantErr %v", tt.scopes, err, tt.wantErr)
			}
		})
	}
}

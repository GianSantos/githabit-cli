package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v60/github"
	"github.com/zalando/go-keyring"
)

const (
	serviceName = "githabit"
	accountName = "github_token"
)

// SaveToken securely stores the GitHub PAT in the system keyring.
func SaveToken(token string) error {
	return keyring.Set(serviceName, accountName, token)
}

// GetToken retrieves the GitHub PAT from the system keyring.
func GetToken() (string, error) {
	return keyring.Get(serviceName, accountName)
}

// DeleteToken removes the token from the system keyring.
func DeleteToken() error {
	return keyring.Delete(serviceName, accountName)
}

// ValidateToken checks if the token is valid and has the required scopes.
func ValidateToken(ctx context.Context, token string) error {
	ts := github.NewClient(nil).WithAuthToken(token)

	// We call the 'user' endpoint to verify the token and get scope headers
	user, resp, err := ts.Users.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to verify token: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token: received status %s", resp.Status)
	}

	// GitHub returns scopes in the 'X-OAuth-Scopes' header
	scopes := resp.Header.Get("X-OAuth-Scopes")
	if err := CheckRequiredScopes(scopes); err != nil {
		return err
	}

	fmt.Printf("✓ Authenticated as %s\n", user.GetLogin())
	return nil
}

// CheckRequiredScopes validates that the given comma-separated scopes include repo, read:user, and read:org.
func CheckRequiredScopes(scopes string) error {
	required := []string{"repo", "read:user", "read:org"}
	found := strings.Split(scopes, ", ")

	for _, req := range required {
		isFound := false
		for _, f := range found {
			if f == req {
				isFound = true
				break
			}
		}
		if !isFound {
			return fmt.Errorf("missing required scope: %s", req)
		}
	}
	return nil
}

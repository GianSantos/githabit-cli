package cmd

import (
	"context"
	"fmt"

	"github.com/GianSantos/githabit-cli/internal/auth"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure GitHub PAT and validate scopes",
	Long:  "Prompts for a Personal Access Token, validates repo and read:user scopes, and saves it securely to the system keyring.",
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	var token string
	prompt := &survey.Password{
		Message: "Enter your GitHub Personal Access Token:",
	}
	if err := survey.AskOne(prompt, &token); err != nil {
		return fmt.Errorf("prompt canceled or failed: %w", err)
	}
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	ctx := context.Background()
	if err := auth.ValidateToken(ctx, token); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := auth.SaveToken(token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	fmt.Println("✓ Token saved securely to keyring.")
	return nil
}

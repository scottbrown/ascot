package main

import (
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/charmbracelet/lipgloss"
	"github.com/scottbrown/ascot"
	"github.com/spf13/cobra"

	"context"
	"fmt"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "ascot",
	Short: "Ascot is an AWS toolkit for SecOps and GRC",
	Long: `A suite of various tools to inspect AWS resources that
          security operations and GRC (compliance) teams often
          need to use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ShowRequiredPermissions {
			printRequiredPermissions([]string{"none"})
			return nil
		}

		if HowItWorks {
			fmt.Println(headingStyle.Render("Logic:"))
			fmt.Println("- Call sts:GetCallerIdentity")
			fmt.Println("- Return the Arn.")
			fmt.Println("")
			fmt.Println(headingStyle.Render("Notes:"))
			fmt.Println("- It requires no permissions to use, making it safe to check whether the authentication to AWS was successful.")
			return nil
		}

		cfg, err := ascot.GetAWSConfig(ascot.DEFAULT_REGION, Profile)
		if err != nil {
			return err
		}

		client := sts.NewFromConfig(cfg)

		resp, err := client.GetCallerIdentity(context.TODO(),
			&sts.GetCallerIdentityInput{},
		)

		if err != nil {
			return err
		}

		highlightStyle := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "12", Dark: "86"})

		fmt.Println("AWS login was successful.")
		fmt.Printf("You are currently logged in as %s\n", highlightStyle.Render(*resp.Arn))

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printRequiredPermissions(privs []string) {
	for _, priv := range privs {
		fmt.Println(priv)
	}
}

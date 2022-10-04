package cmd

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/charmbracelet/lipgloss"
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
	Run: func(cmd *cobra.Command, args []string) {
		var cfg aws.Config
		var err error

		if ShowRequiredPermissions {
			fmt.Println("none")
			return
		}

		if HowItWorks {
			fmt.Println(headingStyle.Render("Logic:"))
			fmt.Println("- Call sts:GetCallerIdentity")
			fmt.Println("- Return the Arn.")
			fmt.Println("")
			fmt.Println(headingStyle.Render("Notes:"))
			fmt.Println("- It requires no permissions to use, making it safe to check whether the authentication to AWS was successful.")
			return
		}

		if Profile != "" {
			cfg, err = config.LoadDefaultConfig(context.TODO(),
				config.WithRegion(DEFAULT_REGION),
				config.WithSharedConfigProfile(Profile),
			)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			// use the default profile
			cfg, err = config.LoadDefaultConfig(context.TODO(),
				config.WithRegion(DEFAULT_REGION),
			)
			if err != nil {
				fmt.Println("AWS login failed")
				fmt.Println(err)
				os.Exit(1)
			}
		}

		client := sts.NewFromConfig(cfg)

		resp, err := client.GetCallerIdentity(context.TODO(),
			&sts.GetCallerIdentityInput{},
		)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		highlightStyle := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "12", Dark: "86"})

		fmt.Println("AWS login was successful.")
		fmt.Printf("You are currently logged in as %s\n", highlightStyle.Render(*resp.Arn))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

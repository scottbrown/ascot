package cmd

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"context"
	"errors"
	"fmt"
	"os"
)

var headingStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.AdaptiveColor{Light: "12", Dark: "86"})

func init() {
	rootCmd.AddCommand(accessKeyOwnerCmd)
}

var alertStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#FF0000"))

var accessKeyOwnerCmd = &cobra.Command{
	Use:   "access-key-owner [access key id]",
	Short: "Finds the owner of a given AWS access key id",
	Long:  `Given an AWS access key, prints the details of the key or nothing if no match`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !ShowRequiredPermissions && len(args) < 1 {
			return errors.New("Missing required argument: access key id")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var cfg aws.Config
		var err error

		if ShowRequiredPermissions {
			fmt.Println("iam:ListAccessKeys")
			fmt.Println("iam:ListUsers")
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
				fmt.Println(err)
				os.Exit(1)
			}
		}

		accessKeyId := args[0]

		client := iam.NewFromConfig(cfg)

		var users []string
		var marker *string
		for {
			resp, err := client.ListUsers(context.TODO(), &iam.ListUsersInput{
				Marker: marker,
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, u := range resp.Users {
				users = append(users, *u.UserName)
			}

			if !resp.IsTruncated {
				break
			}

			marker = resp.Marker
		}

		for i := range users {
			resp, err := client.ListAccessKeys(context.TODO(),
				&iam.ListAccessKeysInput{
					UserName: &users[i],
				},
			)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for _, key := range resp.AccessKeyMetadata {
				if *key.AccessKeyId == accessKeyId {
					if key.Status == "Active" {
						fmt.Println(alertStyle.Render("This key is active"))
					}
					fmt.Printf("%s %s\n", headingStyle.Render("Username:"), *key.UserName)
					fmt.Printf("%s %v\n", headingStyle.Render("Create Date:"), key.CreateDate)
					fmt.Printf("%s %s\n", headingStyle.Render("Status:"), key.Status)
				}
			}
		}
	},
}

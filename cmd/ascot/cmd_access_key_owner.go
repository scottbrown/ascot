package main

import (
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/scottbrown/ascot"
	"github.com/spf13/cobra"

	"errors"
	"fmt"
)

func init() {
	rootCmd.AddCommand(accessKeyOwnerCmd)
}

var accessKeyOwnerCmd = &cobra.Command{
	Use:   "access-key-owner [access key id]",
	Short: "Finds the owner of a given AWS access key id",
	Long:  `Given an AWS access key, prints the details of the key or nothing if no match`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !ShowRequiredPermissions && !HowItWorks && len(args) < 1 {
			return errors.New("Missing required argument: access key id")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var runner ascot.AccessKeyOwnerRunner

		if ShowRequiredPermissions {
			printRequiredPermissions(runner.RequiredPermissions())
			return nil
		}

		if HowItWorks {
			printHowItWorks(runner.HowItWorks())
			return nil
		}

		cfg, err := ascot.GetAWSConfig(ascot.DEFAULT_REGION, Profile)

		if err != nil {
			return err
		}

		accessKeyId := args[0]

		client := iam.NewFromConfig(cfg)

		runner.ListUsersClient = client
		runner.ListAccessKeysClient = client
		runner.AccessKeyId = accessKeyId

		key, err := runner.Run()
		if err != nil {
			return err
		}
		if key.Status == "Active" {
			fmt.Println(alertStyle.Render("This key is active"))
		}
		fmt.Printf("%s %s\n", headingStyle.Render("Username:"), *key.UserName)
		fmt.Printf("%s %v\n", headingStyle.Render("Create Date:"), key.CreateDate)
		fmt.Printf("%s %s\n", headingStyle.Render("Status:"), key.Status)

		return nil
	},
}

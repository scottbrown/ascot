package main

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/scottbrown/ascot"
	"github.com/spf13/cobra"

	"fmt"
	"strings"
)

func init() {
	rootCmd.AddCommand(missingImagesCmd)
}

var missingImagesCmdPrivs []string

var missingImagesCmd = &cobra.Command{
	Use:   "missing-images",
	Short: "Lists any EC2 instances with missing AMIs",
	Long:  `Finds if any EC2 instances are using an AMI that no longer exists, then lists them along with the missing AMI`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var regionRunner ascot.ActiveRegionsRunner
		var runner ascot.MissingImagesRunner

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

		client := ec2.NewFromConfig(cfg)
		regionRunner.Client = *client

		regions, err := regionRunner.Run()
		if err != nil {
			return err
		}

		for _, region := range regions {
			regionalCfg, err := ascot.GetAWSConfig(*region.RegionName, Profile)
			if err != nil {
				return err
			}

			regionalClient := ec2.NewFromConfig(regionalCfg)
			runner.DescribeInstancesClient = regionalClient
			runner.DescribeImagesClient = regionalClient

			missingImages, err := runner.Run()
			if err != nil {
				return err
			}

			// only print if missing AMIs exist
			if len(missingImages) > 0 {
				fmt.Printf("Region: %s\n", headingStyle.Render(*region.RegionName))
				fmt.Println(headingStyle.Render("AMIs Missing:"))
				for imageId, instanceIds := range missingImages {
					fmt.Printf("%s: %s\n", imageId, strings.Join(instanceIds, ", "))
				}
			}
		}

		return nil
	},
}

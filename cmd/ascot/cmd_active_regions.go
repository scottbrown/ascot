package main

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/scottbrown/ascot"
	"github.com/spf13/cobra"

	"fmt"
	"sort"
)

var activeRegionsCmdPrivs []string

var activeRegionsCmd = &cobra.Command{
	Use:   "active-regions",
	Short: "Lists the regions active in the AWS account.",
	Long:  `Reports each region that is listed as active in IAM in the given AWS account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var runner ascot.ActiveRegionsRunner

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
		runner.Client = *client

		regions, err := runner.Run()
		if err != nil {
			return err
		}

		var regionNames []string

		for _, region := range regions {
			regionNames = append(regionNames, *region.RegionName)
		}

		sort.Strings(regionNames)

		for _, name := range regionNames {
			fmt.Println(name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(activeRegionsCmd)
}

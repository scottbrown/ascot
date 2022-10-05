package main

import (
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
		if ShowRequiredPermissions {
			printRequiredPermissions(activeRegionsCmdPrivs)
			return nil
		}

		if HowItWorks {
			fmt.Println(headingStyle.Render("Logic:"))
			fmt.Println("- Call ec2:DescribeRegions")
			fmt.Println("- Loop through each region")
			fmt.Println("- Print the region name")
			return nil
		}

		cfg, err := ascot.GetAWSConfig(ascot.DEFAULT_REGION, Profile)
		if err != nil {
			return err
		}

		regions, err := ascot.GetAllRegions(cfg)
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

	activeRegionsCmdPrivs = []string{
		"ec2:DescribeRegions",
	}
}

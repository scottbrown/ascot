package main

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/scottbrown/ascot"
	"github.com/spf13/cobra"

	"fmt"
)

func init() {
	rootCmd.AddCommand(auditDefaultVpcsCmd)
}

var auditDefaultVpcsCmd = &cobra.Command{
	Use:   "audit-default-vpcs",
	Short: "Validates whether default VPCs exist in all regions",
	Long:  `Default VPCs are not intended to be used, and should not exist in any AWS region.  This command verifies whether they exist in a region (FAIL) or have been removed (PASS)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var regionRunner ascot.ActiveRegionsRunner
		var runner ascot.AuditDefaultVpcsRunner

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
			runner.Client = regionalClient

			vpcs, err := runner.Run()
			if err != nil {
				return err
			}

			if len(vpcs) > 0 {
				fmt.Printf("[%s] %s\n", failStyle.Render("FAIL"), *region.RegionName)
			} else {
				fmt.Printf("[%s] %s\n", passStyle.Render("PASS"), *region.RegionName)
			}
		}

		return nil
	},
}

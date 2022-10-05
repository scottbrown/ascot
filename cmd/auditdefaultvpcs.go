package cmd

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"context"
	"fmt"
)

var auditDefaultVpcsCmdPrivs []string

func init() {
	rootCmd.AddCommand(auditDefaultVpcsCmd)

	auditDefaultVpcsCmdPrivs = []string{
		"ec2:DescribeRegions",
		"ec2:DescribeVpcs",
	}
}

var auditDefaultVpcsCmd = &cobra.Command{
	Use:   "audit-default-vpcs",
	Short: "Validates whether default VPCs exist in all regions",
	Long:  `Default VPCs are not intended to be used, and should not exist in any AWS region.  This command verifies whether they exist in a region (FAIL) or have been removed (PASS)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ShowRequiredPermissions {
			printRequiredPermissions(auditDefaultVpcsCmdPrivs)
			return nil
		}

		if HowItWorks {
			fmt.Println(headingStyle.Render("Logic:"))
			fmt.Println("- Call ec2:DescribeRegions")
			fmt.Println("- Loop through each region")
			fmt.Println("- Call ec2:DescribeVpcs, filtering by is-default")
			fmt.Println("- Print FAIL if any VPCs were returned")
			fmt.Println("- Otherwise PASS")
			return nil
		}

		cfg, err := getAWSConfig(DEFAULT_REGION, Profile)
		if err != nil {
			return err
		}

		regions, err := getAllRegions(cfg)
		if err != nil {
			return err
		}

		for _, region := range regions {
			regionalCfg, err := getAWSConfig(*region.RegionName, Profile)
			if err != nil {
				return err
			}

			client := ec2.NewFromConfig(regionalCfg)

			resp, err := client.DescribeVpcs(context.TODO(),
				&ec2.DescribeVpcsInput{
					Filters: []types.Filter{
						types.Filter{
							Name: aws.String("is-default"),
							Values: []string{
								"true",
							},
						},
					},
				},
			)
			if err != nil {
				return err
			}

			passStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
			failStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))

			if len(resp.Vpcs) > 0 {
				fmt.Printf("[%s] %s\n", failStyle.Render("FAIL"), *region.RegionName)
			} else {
				fmt.Printf("[%s] %s\n", passStyle.Render("PASS"), *region.RegionName)
			}
		}

		return nil
	},
}

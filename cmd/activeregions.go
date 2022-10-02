package cmd

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"

	"context"
	"fmt"
	"os"
	"sort"
)

var activeRegionsCmd = &cobra.Command{
	Use:   "active-regions",
	Short: "Lists the regions active in the AWS account.",
	Long:  `Reports each region that is listed as active in IAM in the given AWS account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg aws.Config
		var err error

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

		client := ec2.NewFromConfig(cfg)

		resp, err := client.DescribeRegions(context.TODO(),
			&ec2.DescribeRegionsInput{},
		)

		if err != nil {
			return err
		}

		var regionNames []string

		for _, region := range resp.Regions {
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

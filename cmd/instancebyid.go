package cmd

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"

	"context"
	"errors"
	"fmt"
	"os"
)

var instanceByIdCmd = &cobra.Command{
	Use:   "instance-by-id [instance-id]",
	Short: "Finds the instance in any region by its ID",
	Long:  `Searches every region for an EC2 instance that matches a given ID`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !ShowRequiredPermissions && len(args) < 1 {
			return errors.New("Missing required argument: instance id")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg aws.Config
		var err error

		if ShowRequiredPermissions {
			fmt.Println("ec2:DescribeRegions")
			fmt.Println("ec2:DescribeInstances")
			return nil
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

		instanceId := args[0]

		client := ec2.NewFromConfig(cfg)

		resp, err := client.DescribeRegions(context.TODO(),
			&ec2.DescribeRegionsInput{},
		)

		if err != nil {
			return err
		}

		for _, region := range resp.Regions {
			// connect to another region
			if Profile != "" {
				cfg, err = config.LoadDefaultConfig(context.TODO(),
					config.WithRegion(*region.RegionName),
					config.WithSharedConfigProfile(Profile),
				)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				// use the default profile
				cfg, err = config.LoadDefaultConfig(context.TODO(),
					config.WithRegion(*region.RegionName),
				)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			regional_client := ec2.NewFromConfig(cfg)
			resp, err := regional_client.DescribeInstances(context.TODO(),
				&ec2.DescribeInstancesInput{
					Filters: []types.Filter{
						types.Filter{
							Name: aws.String("instance-id"),
							Values: []string{
								instanceId,
							},
						},
					},
				},
			)

			if err != nil {
				return err
			}

			// print out the instance details if a match was found
			for _, reservation := range resp.Reservations {
				for _, instance := range reservation.Instances {
					fmt.Printf("%s: %s\n", headingStyle.Render("Region"), *region.RegionName)
					fmt.Printf("%s: %s\n", headingStyle.Render("Instance ID"), *instance.InstanceId)
					fmt.Printf("%s: %s\n", headingStyle.Render("Public IP Address"), *instance.PublicIpAddress)
					fmt.Printf("%s: %s\n", headingStyle.Render("Private IP Address"), *instance.PrivateIpAddress)
					fmt.Printf("%s: %s\n", headingStyle.Render("Image ID"), *instance.ImageId)
					fmt.Printf("%s: %s\n", headingStyle.Render("Instance Type"), instance.InstanceType)
					fmt.Printf("%s: %v\n", headingStyle.Render("Launch Time"), *instance.LaunchTime)
					fmt.Printf("%s: %s\n", headingStyle.Render("State"), instance.State.Name)
					fmt.Printf("%s: %s\n", headingStyle.Render("VPC"), *instance.VpcId)
					fmt.Printf("%s: %s\n", headingStyle.Render("Subnet"), *instance.SubnetId)
					fmt.Printf("%s:\n", headingStyle.Render("Tags"))
					for _, tag := range instance.Tags {
						fmt.Printf("%s: %s\n", headingStyle.Render(*tag.Key), *tag.Value)
					}
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(instanceByIdCmd)
}

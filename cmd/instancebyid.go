package cmd

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/cobra"

	"context"
	"errors"
	"fmt"
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
		if ShowRequiredPermissions {
			fmt.Println("ec2:DescribeRegions")
			fmt.Println("ec2:DescribeInstances")
			return nil
		}

		cfg, err := getAWSConfig(DEFAULT_REGION, Profile)
		if err != nil {
			return err
		}

		instanceId := args[0]

		regions, err := getAllRegions(cfg)
		if err != nil {
			return err
		}

		for _, region := range regions {
			// connect to another region
			regional_cfg, err := getAWSConfig(*region.RegionName, Profile)
			if err != nil {
				return err
			}

			regional_client := ec2.NewFromConfig(regional_cfg)
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

package main

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/scottbrown/ascot"
	"github.com/spf13/cobra"

	"errors"
	"fmt"
)

var instanceByIdCmdPrivs []string

const fancyPrint string = "%s: %s\n"

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
		var regionRunner ascot.ActiveRegionsRunner
		var runner ascot.InstanceByIdRunner

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

		instanceId := args[0]

		client := ec2.NewFromConfig(cfg)
		regionRunner.Client = *client

		regions, err := regionRunner.Run()
		if err != nil {
			return err
		}

		for _, region := range regions {
			// connect to another region
			regionalCfg, err := ascot.GetAWSConfig(*region.RegionName, Profile)
			if err != nil {
				return err
			}

			regionalClient := ec2.NewFromConfig(regionalCfg)
			runner.Client = regionalClient

			instance, err := runner.Run(instanceId)
			if err != nil {
				return err
			}

			if instance.InstanceId != nil {
				fmt.Printf(fancyPrint, headingStyle.Render("Region"), *region.RegionName)
				fmt.Printf(fancyPrint, headingStyle.Render("Instance ID"), *instance.InstanceId)
				fmt.Printf(fancyPrint, headingStyle.Render("Public IP Address"), *instance.PublicIpAddress)
				fmt.Printf(fancyPrint, headingStyle.Render("Private IP Address"), *instance.PrivateIpAddress)
				fmt.Printf(fancyPrint, headingStyle.Render("Image ID"), *instance.ImageId)
				fmt.Printf(fancyPrint, headingStyle.Render("Instance Type"), instance.InstanceType)
				fmt.Printf(fancyPrint, headingStyle.Render("Launch Time"), *instance.LaunchTime)
				fmt.Printf(fancyPrint, headingStyle.Render("State"), instance.State.Name)
				fmt.Printf(fancyPrint, headingStyle.Render("VPC"), *instance.VpcId)
				fmt.Printf(fancyPrint, headingStyle.Render("Subnet"), *instance.SubnetId)
				fmt.Printf("%s:\n", headingStyle.Render("Tags"))
				for _, tag := range instance.Tags {
					fmt.Printf(fancyPrint, headingStyle.Render(*tag.Key), *tag.Value)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(instanceByIdCmd)

	instanceByIdCmdPrivs = []string{
		"ec2:DescribeRegions",
		"ec2:DescribeInstances",
	}
}

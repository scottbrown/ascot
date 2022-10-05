package ascot

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"context"
)

func GetAWSConfig(region string, profile string) (aws.Config, error) {
	var cfg aws.Config
	var err error

	if profile != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
			config.WithSharedConfigProfile(profile),
		)
		if err != nil {
			return aws.Config{}, err
		}
	} else {
		// use the default profile
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(region),
		)
		if err != nil {
			return aws.Config{}, err
		}
	}

	return cfg, nil
}

func GetAllRegions(cfg aws.Config) ([]types.Region, error) {
	client := ec2.NewFromConfig(cfg)

	resp, err := client.DescribeRegions(context.TODO(),
		&ec2.DescribeRegionsInput{},
	)

	if err != nil {
		return []types.Region{}, err
	}

	return resp.Regions, nil
}

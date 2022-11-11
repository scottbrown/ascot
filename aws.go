package ascot

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2_types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iam_types "github.com/aws/aws-sdk-go-v2/service/iam/types"

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

type ActiveRegionsRunner struct {
	Client ec2.Client
}

func (s ActiveRegionsRunner) Run() ([]ec2_types.Region, error) {
	resp, err := s.Client.DescribeRegions(context.TODO(),
		&ec2.DescribeRegionsInput{},
	)

	if err != nil {
		return []ec2_types.Region{}, err
	}

	return resp.Regions, nil
}

func (s ActiveRegionsRunner) RequiredPermissions() []string {
	return []string{"ec2:DescribeRegions"}
}

func (s ActiveRegionsRunner) HowItWorks() ([]string, []string) {
	return []string{
		"Call ec2:DescribeRegions",
		"Loop through each region",
		"Print the region name",
	}, []string{}
}

type AccessKeyOwnerRunner struct {
	ListUsersClient      iam.ListUsersAPIClient
	ListAccessKeysClient iam.ListAccessKeysAPIClient
	AccessKeyId          string
}

func (s AccessKeyOwnerRunner) Run() (iam_types.AccessKeyMetadata, error) {
	var users []string
	var marker *string
	for {
		resp, err := s.ListUsersClient.ListUsers(context.TODO(), &iam.ListUsersInput{
			Marker: marker,
		})
		if err != nil {
			return iam_types.AccessKeyMetadata{}, err
		}

		for _, u := range resp.Users {
			users = append(users, *u.UserName)
		}

		if !resp.IsTruncated {
			break
		}

		marker = resp.Marker
	}

	var foundKey iam_types.AccessKeyMetadata
	for i := range users {
		resp, err := s.ListAccessKeysClient.ListAccessKeys(context.TODO(),
			&iam.ListAccessKeysInput{
				UserName: &users[i],
			},
		)

		if err != nil {
			return iam_types.AccessKeyMetadata{}, err
		}

		for _, key := range resp.AccessKeyMetadata {
			if *key.AccessKeyId == s.AccessKeyId {
				foundKey = key
				break
			}
		}
	}

	return foundKey, nil
}

func (AccessKeyOwnerRunner) RequiredPermissions() []string {
	return []string{"iam:ListAccessKeys", "iam:ListUsers"}
}

func (AccessKeyOwnerRunner) HowItWorks() ([]string, []string) {
	return []string{
		"Call iam:ListUsers",
		"Loop through each user",
		"Call iam:ListAccessKeys for the user",
		"Find a match with the given key",
	}, []string{}
}

type MissingImagesRunner struct {
	DescribeInstancesClient ec2.DescribeInstancesAPIClient
	DescribeImagesClient    ec2.DescribeImagesAPIClient
}

func (s MissingImagesRunner) Run() (map[string][]string, error) {
	var instances []ec2_types.Instance
	var missingImages map[string][]string

	var marker *string
	for {
		resp, err := s.DescribeInstancesClient.DescribeInstances(context.TODO(),
			&ec2.DescribeInstancesInput{
				NextToken: marker,
			},
		)
		if err != nil {
			return missingImages, err
		}

		for _, reservation := range resp.Reservations {
			for _, instance := range reservation.Instances {
				instances = append(instances, instance)
			}
		}

		marker = resp.NextToken

		if resp.NextToken == nil {
			break
		}
	}

	// create a Set
	var imageIds []string
	for _, instance := range instances {
		idArr := append(missingImages[*instance.ImageId], *instance.InstanceId)
		missingImages[*instance.ImageId] = idArr
	}

	for k := range missingImages {
		imageIds = append(imageIds, k)
	}

	resp, err := s.DescribeImagesClient.DescribeImages(context.TODO(),
		&ec2.DescribeImagesInput{
			ExecutableUsers: []string{
				"self",
				"all",
			},
			ImageIds: imageIds,
		},
	)

	if err != nil {
		return missingImages, err
	}

	// match up the resp.Images with imageIds to find what is missing
	for _, images := range resp.Images {
		delete(missingImages, *images.ImageId)
	}

	return missingImages, nil
}

func (s MissingImagesRunner) RequiredPermissions() []string {
	return []string{
		"ec2:DescribeRegions",
		"ec2:DescribeInstances",
		"ec2:DescribeImages",
	}
}

func (s MissingImagesRunner) HowItWorks() ([]string, []string) {
	return []string{
		"Call ec2:DescribeRegions",
		"Loop through each region",
		"Call ec2:DescribeInstances for each region",
		"Get the instance ID and the image ID",
		"Call ec2:DescribeImages with the image ID",
		"Print if the image ID is non-existent",
	}, []string{}
}

type AuditDefaultVpcsRunner struct {
	Client ec2.DescribeVpcsAPIClient
}

func (s AuditDefaultVpcsRunner) Run() ([]ec2_types.Vpc, error) {
	resp, err := s.Client.DescribeVpcs(context.TODO(),
		&ec2.DescribeVpcsInput{
			Filters: []ec2_types.Filter{
				ec2_types.Filter{
					Name: aws.String("is-default"),
					Values: []string{
						"true",
					},
				},
			},
		},
	)

	if err != nil {
		return []ec2_types.Vpc{}, err
	}

	return resp.Vpcs, nil
}

func (s AuditDefaultVpcsRunner) RequiredPermissions() []string {
	return []string{"ec2:DescribeRegions", "ec2:DescribeVpcs"}
}

func (s AuditDefaultVpcsRunner) HowItWorks() ([]string, []string) {
	return []string{
		"- Call ec2:DescribeRegions",
		"- Loop through each region",
		"- Call ec2:DescribeVpcs, filtering by is-default",
		"- Print FAIL if any VPCs were returned",
		"- Otherwise PASS",
	}, []string{}
}

type InstanceByIdRunner struct {
	Client ec2.DescribeInstancesAPIClient
}

func (s InstanceByIdRunner) Run(instanceId string) (ec2_types.Instance, error) {
	resp, err := s.Client.DescribeInstances(context.TODO(),
		&ec2.DescribeInstancesInput{
			Filters: []ec2_types.Filter{
				ec2_types.Filter{
					Name: aws.String("instance-id"),
					Values: []string{
						instanceId,
					},
				},
			},
		},
	)

	if err != nil {
		return ec2_types.Instance{}, err
	}

	// print out the instance details if a match was found
	var instance ec2_types.Instance
	for _, r := range resp.Reservations {
		for _, i := range r.Instances {
			instance = i
		}
	}

	return instance, err
}

func (s InstanceByIdRunner) RequiredPermissions() []string {
	return []string{"ec2:DescribeRegions", "ec2:DescribeInstances"}
}

func (s InstanceByIdRunner) HowItWorks() ([]string, []string) {
	return []string{
		"Call ec2:DescribeRegions",
		"Loop through each region",
		"Call ec2:DescribeInstances, filtering by instance-id",
		"Print instance details if matched with given id",
	}, []string{}
}

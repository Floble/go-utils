package ec2

import (
	"context"
	"time"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Instance struct {
}

func NewEC2Instance() *EC2Instance {
	return new(EC2Instance)
}

type EC2CreateInstanceAPI interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	CreateTags(ctx context.Context,
		params *ec2.CreateTagsInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
}

type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

func (instance *EC2Instance) Create() (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)
	waiter := NewInstanceRunningWaiter(client)

	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-05f7491af5eef733a"),
		InstanceType: types.InstanceTypeT2Micro,
		KeyName: 	  aws.String("Floble"),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	}

	runResult, err := makeInstance(context.TODO(), client, runInput)
	if err != nil {
		return "", err
	}

	tagInput := &ec2.CreateTagsInput{
		Resources: []string{*runResult.Instances[0].InstanceId},
		Tags: []types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("YUMA"),
			},
		},
	}

	_, err = makeTags(context.TODO(), client, tagInput)
	if err != nil {
		return "", err
	}

	describeInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{*runResult.Instances[0].InstanceId},
	}

	err = waiter.Wait(context.TODO(), describeInput, 5 * time.Minute)
	if err != nil {
		return "", err
	}
	
	describeResult, err := getInstances(context.TODO(), client, describeInput)
	if err != nil {
		return "", err
	}

	return *describeResult.Reservations[0].Instances[0].PublicIpAddress, nil
}

func makeInstance(c context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}

func makeTags(c context.Context, api EC2CreateInstanceAPI, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return api.CreateTags(c, input)
}

func getInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}
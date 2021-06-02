package ec2

// Parts of this code were obtain from the AWS documentation
// https://docs.aws.amazon.com/code-samples/latest/catalog/gov2-ec2-CreateInstance-CreateInstancev2.go.html
// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/ec2/describeinstances/
// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/ec2/stopinstances/

import (
	"context"
	"time"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Instance struct {
	id, publicIP string
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

type EC2StopInstancesAPI interface {
	StopInstances(ctx context.Context,
		params *ec2.StopInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
}

func (instance *EC2Instance) Create() error {
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

	runResult, err := makeInstances(context.TODO(), client, runInput)
	if err != nil {
		return err
	}

	tagInput := &ec2.CreateTagsInput{
		Resources: []string{*runResult.Instances[0].InstanceId},
		Tags: []types.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("Managed by YUMA"),
			},
		},
	}

	_, err = makeTags(context.TODO(), client, tagInput)
	if err != nil {
		return err
	}

	describeInput := &ec2.DescribeInstancesInput{
		InstanceIds: []string{*runResult.Instances[0].InstanceId},
	}

	err = waiter.Wait(context.TODO(), describeInput, 5 * time.Minute)
	if err != nil {
		return err
	}
	
	describeResult, err := getInstances(context.TODO(), client, describeInput)
	if err != nil {
		return err
	}

	instance.SetID(*runResult.Instances[0].InstanceId)
	instance.SetPublicIP(*describeResult.Reservations[0].Instances[0].PublicIpAddress)

	return nil
}

func (instance *EC2Instance) Stop() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)
	
	stopInput := &ec2.StopInstancesInput{
		InstanceIds: []string{
			instance.id,
		},
	}

	_, err = stopInstances(context.TODO(), client, stopInput)
	if err != nil {
		return err
	}

	return nil
}

func makeInstances(c context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}

func makeTags(c context.Context, api EC2CreateInstanceAPI, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return api.CreateTags(c, input)
}

func getInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func stopInstances(c context.Context, api EC2StopInstancesAPI, input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
    return api.StopInstances(c, input)
}

func (instance *EC2Instance) GetID() string {
	return instance.id
}

func (instance *EC2Instance) GetPublicIP() string {
	return instance.publicIP
}

func (instance *EC2Instance) SetID(id string) {
	instance.id = id
}

func (instance *EC2Instance) SetPublicIP(publicIP string) {
	instance.publicIP = publicIP
}
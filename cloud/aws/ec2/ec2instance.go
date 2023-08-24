package ec2

// Parts of this code were obtain from the AWS documentation
// https://docs.aws.amazon.com/code-samples/latest/catalog/gov2-ec2-CreateInstance-CreateInstancev2.go.html
// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/ec2/describeinstances/
// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/ec2/stopinstances/

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Instance struct {
	id, publicIP, privateIP string
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

type EC2TerminateInstancesAPI interface {
	TerminateInstances(ctx context.Context,
		params *ec2.TerminateInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

func (instance *EC2Instance) Create() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)
	waiter := ec2.NewInstanceRunningWaiter(client)

	runInput := &ec2.RunInstancesInput{
		ImageId:      aws.String("ami-03cbad7144aeda3eb"),
		//ImageId:      aws.String("ami-04e601abe3e1a910f"),
		InstanceType: types.InstanceTypeT2Xlarge,
		KeyName: 	  aws.String("Floble"),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
	}

	runResult, err := makeInstances(context.TODO(), client, runInput)
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

	instance.SetID(*runResult.Instances[0].InstanceId)
	instance.SetPublicIP(*describeResult.Reservations[0].Instances[0].PublicIpAddress)
	instance.SetPrivateIP(*describeResult.Reservations[0].Instances[0].PrivateIpAddress)

	return nil
}

func (instance *EC2Instance) Delete() error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ec2.NewFromConfig(cfg)

	terminateInput := &ec2.TerminateInstancesInput{
		InstanceIds: []string{
			instance.id,
		},
	}

	_, err = terminateInstances(context.TODO(), client, terminateInput)
	if err != nil {
		return err
	}

	return nil
}

func (instance *EC2Instance) AddToKnownHosts() error {
	cmd := exec.Command("ssh-keyscan", "-t", "rsa", instance.GetPublicIP())
	export := ""
	out, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("ADD TO KNOWN HOSTS ERROR: STDOUT PIPE")
		return err
	}
  
	scanner := bufio.NewScanner(out)
	go func() {
		for scanner.Scan() {
			export += scanner.Text() + "\n"
		}
	}()
  
	err = cmd.Start()
	if err != nil {
		fmt.Println("ADD TO KNOWN HOSTS ERROR: COMMAND START")
		return err  
	}
  
	err = cmd.Wait()
	if err != nil {
		fmt.Println("ADD TO KNOWN HOSTS ERROR: COMMAND WAIT")
		fmt.Println(err.Error())
		return err
	}

	file, err := os.OpenFile("/Users/floble/.ssh/known_hosts", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("ADD TO KNOWN HOSTS ERROR: OPEN FILE")
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(export); err != nil {
		fmt.Println("ADD TO KNOWN HOSTS ERROR: WRITE STRING")
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

func terminateInstances(c context.Context, api EC2TerminateInstancesAPI, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
    return api.TerminateInstances(c, input)
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

func (instance *EC2Instance) GetPrivateIP() string {
	return instance.privateIP
}

func (instance *EC2Instance) SetPrivateIP(privateIP string) {
	instance.privateIP = privateIP
}
package k8sml

import (
	"gopkg.in/yaml.v3"
	"errors"
	"reflect"
	"strings"
	terraform "KubeArch/kubearch/proletarian/terraform"
)

type AmazonWebServices struct {
	ID, Profile, Region string
	RuntimeVariables map[string]string
	Cloud []Cloud
	Policy []*IAMPolicy
}

type tmpAmazonWebServices struct {
	Profile string `yaml:"profile"`
	Region string `yaml:"region"`
	VirtualPrivateCloud map[string][]yaml.Node `yaml:",inline"`
}

func (aws *AmazonWebServices) GetID() string {
	return aws.ID
}

func (aws *AmazonWebServices) GetVariableValue(variable string) interface{} {
	e := reflect.ValueOf(aws).Elem()
	var value interface{}

	for i := 0; i < e.NumField(); i++ {
		if strings.ToLower(e.Type().Field(i).Name) == strings.ToLower(variable) {
			value = e.Field(i).Interface()
			return value
		}
	}

	return nil
}

func (aws *AmazonWebServices) GetCloud() []Cloud {
	return aws.Cloud
}

func (aws *AmazonWebServices) GetType() string {
	providerType := reflect.TypeOf(aws).String()
	return providerType
}

func (aws *AmazonWebServices) GetPolicy() []*IAMPolicy {
	return aws.Policy
}

func (aws *AmazonWebServices) GetRuntimeVariables() map[string]string {
    return aws.RuntimeVariables
}

func (aws *AmazonWebServices) AddRuntimeVariable(key, value string) {
    aws.RuntimeVariables[key] = value
}

func (aws *AmazonWebServices) ExportModule() error {
	provider := terraform.NewProvider("aws")

	e := reflect.ValueOf(aws).Elem()
	
	for i := 0; i < e.NumField(); i++ {
		key := e.Type().Field(i).Name
		value := e.Field(i).Interface()
		if key != "ID" {
			provider.AddVariable(key, value)
		}
	}

	err := provider.Export()

	return err
}

func (aws *AmazonWebServices) UnmarshalYAML(value *yaml.Node) error {
	var tmpAmazonWebServices tmpAmazonWebServices
	
    if err := value.Decode(&tmpAmazonWebServices); err != nil {
        return err
	}
	
	aws.Profile = tmpAmazonWebServices.Profile
	aws.Region = tmpAmazonWebServices.Region

	var cloud []Cloud

	for tag, nodes := range tmpAmazonWebServices.VirtualPrivateCloud {
		switch tag {
		case "VirtualPrivateCloud":
			for _, node := range nodes {
				vpc := &VirtualPrivateCloud{}
				vpc.CloudProvider = aws
			
				if err := node.Decode(vpc); err != nil {
					return err
				}
				
				cloud = append(cloud, vpc)
			}
		default:
			return errors.New("Failed to interpret the cloud of type: \"" + tag + "\"")
		}
	}

	aws.Cloud = cloud

	aws.Policy = make([]*IAMPolicy, 0)
	mPolicy := NewPolicy("AWSCloudProviderMasterPolicy", "master")
	mPolicy.CloudProvider = aws
	mPolicy.AddRuntimeVariable("arn", "")
	mRole := NewIAMRole("aws-cloudprovider-master")
	mRole.Policy = mPolicy
	mPolicy.Role = append(mPolicy.Role, mRole)
	mRole.AddRuntimeVariable("arn", "")

	wPolicy := NewPolicy("AWSCloudProviderWorkerPolicy", "worker")
	wPolicy.CloudProvider = aws
	wPolicy.AddRuntimeVariable("arn", "")
	aws.Policy = append(aws.Policy, mPolicy)
	aws.Policy = append(aws.Policy, wPolicy)
	wRole := NewIAMRole("aws-cloudprovider-worker")
	wRole.Policy = wPolicy
	wPolicy.Role = append(wPolicy.Role, wRole)
	wRole.AddRuntimeVariable("arn", "")
	
    return nil
}
package k8sml

import (
	"errors"
	"gopkg.in/yaml.v3"
)

type CloudEnvironment struct {
	CloudProvider []CloudProvider
}

type tmpCloudEnvironment map[string]yaml.Node

func NewCloudEnvironment() *CloudEnvironment {
	cloudEnv := new(CloudEnvironment)

	return cloudEnv
}

func (cloudEnvironment *CloudEnvironment) UnmarshalYAML(value *yaml.Node) error {
	var tmpCloudEnvironment tmpCloudEnvironment
	
    if err := value.Decode(&tmpCloudEnvironment); err != nil {
        return err
	}
	
	cloudProvider := make([]CloudProvider, 0, len(tmpCloudEnvironment))
	
	for tag, node := range tmpCloudEnvironment {
		switch tag {
		case "AmazonWebServices":
			aws := &AmazonWebServices{}
			if err := node.Decode(aws); err != nil {
				return err
			}
			aws.ID = "aws"
			
			cloudProvider = append(cloudProvider, aws)
		default:
			return errors.New("Failed to interpret the cloud provider of type: \"" + tag + "\"")
		}
	}
	
	cloudEnvironment.CloudProvider = cloudProvider
	
    return nil
}
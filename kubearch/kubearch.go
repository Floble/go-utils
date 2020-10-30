package main

import (
    "log"
    "io/ioutil"
    "os"
    "flag"
    "strings"
    "path/filepath"
    k8sml "KubeArch/kubearch/k8sml"
    ansible "KubeArch/kubearch/proletarian/ansible"
    proletarian "KubeArch/kubearch/proletarian"
)

func main() {
  buildArg := flag.NewFlagSet("build", flag.ExitOnError)
  arch := buildArg.String("a", "", "The declarative architecture model that must be build (Required)")
  ssh := buildArg.String("ssh", "", "The path to the ssh private key required by Ansible (Required)")

  switch os.Args[1] {
  case "build":
    buildArg.Parse(os.Args[2:])
    
    data, err := ioutil.ReadFile(*arch)
    if err != nil {
      log.Fatalf("error: %v", err)
    }

    productionPlan := new(proletarian.ProductionPlan)
    productionPlan = productionPlan.Create(data)

    if err := productionPlan.BuildInfrastructure(); err != nil {
      log.Fatalf("error: %v", err)
    }

    if err := productionPlan.Update("EC2Instance"); err != nil {
      log.Fatalf("error: %v", err)
    }
    if err := productionPlan.Update("VirtualPrivateCloud"); err != nil {
      log.Fatalf("error: %v", err)
    }
    if err := productionPlan.Update("Subnet"); err != nil {
      log.Fatalf("error: %v", err)
    }

    productionPlan.ProductionOrder = nil
    for _, aws := range productionPlan.Products["AmazonWebServices"] {
      productionPlan.ProductionOrder = append(productionPlan.ProductionOrder, aws.(k8sml.Infrastructure))
    }
    for _, tg := range productionPlan.Products["TargetGroup"] {
      productionPlan.ProductionOrder = append(productionPlan.ProductionOrder, tg.(k8sml.Infrastructure))
    }
    for _, nlb := range productionPlan.Products["NetworkLoadBalancer"] {
      productionPlan.ProductionOrder = append(productionPlan.ProductionOrder, nlb.(k8sml.Infrastructure))
    }

    if err := productionPlan.BuildInfrastructure(); err != nil {
      log.Fatalf("error: %v", err)
    }

    if err := productionPlan.Update("NetworkLoadBalancer"); err != nil {
      log.Fatalf("error: %v", err)
    }

    err, defaults := scanForAnsibleDefaults()
    if err != nil {
      log.Fatalf("error: %v", err)
    }

    for _, d := range defaults {
      for entry, _ := range d.Variables {
        element, id, variable := d.ParseVariable(entry)
        elementTitle := strings.Title(element)
        variableTitle := strings.Title(variable)

        product := productionPlan.GetProduct(elementTitle, id)
        value := product.GetVariableValue(variableTitle)

        d.SetVariableValue(element + "." + id + "." + variable, value.(string))
      }

      d.ReplaceVariables()
    }

    if err := productionPlan.DeploySoftware(*ssh); err != nil {
      log.Fatalf("error: %v", err)
    }
  }
}

func scanForAnsibleDefaults() (error, []*ansible.Default) {
	ds := make([]*ansible.Default, 0)
	
	err := filepath.Walk("./factors_of_production/software/",
		func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if strings.Contains(path, "/defaults/main.yml") {
			d := ansible.NewDefault(path)
			ds = append(ds, d)
		}

		return nil
	})

	return err, ds
}
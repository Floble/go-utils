package proletarian

import (
	"bytes"
	"strings"
	"reflect"
	"fmt"
	"bufio"
	"time"
	"log"
	"os"
	"os/exec"
	"gopkg.in/yaml.v3"
	k8sml "KubeArch/kubearch/k8sml"
    ansible "KubeArch/kubearch/proletarian/ansible"
)

type ProductionPlan struct {
	ProductionOrder []k8sml.Infrastructure
	Products map[string][]k8sml.K8sML
  }

  func (productionPlan *ProductionPlan) Create(data []byte) *ProductionPlan {
	cloudEnvironment := k8sml.NewCloudEnvironment()
    err := yaml.Unmarshal(data, &cloudEnvironment)
    if err != nil {
        log.Fatalf("error: %v", err)
    }
	
	productionOrder := make([]k8sml.Infrastructure, 0)
	products := make(map[string][]k8sml.K8sML, 0)
  
	for _, cloudProvider := range cloudEnvironment.CloudProvider {
	  productionOrder = append(productionOrder, cloudProvider)
	  cloudProviderType := strings.Split(reflect.TypeOf(cloudProvider).String(), "*k8sml.")[1]
	  products[cloudProviderType] = append(products[cloudProviderType], cloudProvider)
  
	  for _, policy := range cloudProvider.GetPolicy() {
		productionOrder = append(productionOrder, policy)
		policyType := strings.Split(reflect.TypeOf(policy).String(), "*k8sml.")[1]
		products[policyType] = append(products[policyType], policy)
  
		for _, role := range policy.Role {
		  productionOrder = append(productionOrder, role)
		  roleType := strings.Split(reflect.TypeOf(role).String(), "*k8sml.")[1]
		  products[roleType] = append(products[roleType], role)
		}
	  }
	  
	  for _, cloud := range cloudProvider.GetCloud() {
		productionOrder = append(productionOrder, cloud)
		cloudType := strings.Split(reflect.TypeOf(cloud).String(), "*k8sml.")[1]
		products[cloudType] = append(products[cloudType], cloud)
  
		for _, cidr := range cloud.GetIPv4Cidr() {
		  productionOrder = append(productionOrder, cidr)
		  cidrType := strings.Split(reflect.TypeOf(cidr).String(), "*k8sml.")[1]
		  products[cidrType] = append(products[cidrType], cidr)
		}
  
		igw := cloud.GetInternetGateway()
		if igw != nil {
		  productionOrder = append(productionOrder, igw)
		  igwType := strings.Split(reflect.TypeOf(igw).String(), "*k8sml.")[1]
		  products[igwType] = append(products[igwType], igw)
		}
  
		k8s := cloud.GetKubernetes()
		k8sType := strings.Split(reflect.TypeOf(k8s).String(), "*k8sml.")[1]
		products[k8sType] = append(products[k8sType], k8s)
  
		for _, subnet := range k8s.Subnet {
		  productionOrder = append(productionOrder, subnet)
		  subnetType := strings.Split(reflect.TypeOf(subnet).String(), "*k8sml.")[1]
		  products[subnetType] = append(products[subnetType], subnet)
  
		  for _, ngw := range subnet.NATGateway {
			if ngw != nil {
			  productionOrder = append(productionOrder, ngw)
			  ngwType := strings.Split(reflect.TypeOf(ngw).String(), "*k8sml.")[1]
			  products[ngwType] = append(products[ngwType], ngw)
			}
		  }
  
		  for _, routeTable := range subnet.RouteTable {
			productionOrder = append(productionOrder, routeTable)
			routeTableType := strings.Split(reflect.TypeOf(routeTable).String(), "*k8sml.")[1]
			products[routeTableType] = append(products[routeTableType], routeTable)
  
			for _, route := range routeTable.Route {
			  productionOrder = append(productionOrder, route)
			  routeType := strings.Split(reflect.TypeOf(route).String(), "*k8sml.")[1]
			  products[routeType] = append(products[routeType], route)
			}
		  }
  
		  for _, nlb := range subnet.NetworkLoadBalancer {
			nlbType := strings.Split(reflect.TypeOf(nlb).String(), "*k8sml.")[1]
			products[nlbType] = append(products[nlbType], nlb)
		  }
  
		  for _, virtualFirewall := range subnet.VirtualFirewall {
			productionOrder = append(productionOrder, virtualFirewall)
			virtualFirewallType := strings.Split(reflect.TypeOf(virtualFirewall).String(), "*k8sml.")[1]
			products[virtualFirewallType] = append(products[virtualFirewallType], virtualFirewall)
	
			for _, ingress := range virtualFirewall.GetIngress() {
			  productionOrder = append(productionOrder, ingress)
			  ingressType := strings.Split(reflect.TypeOf(ingress).String(), "*k8sml.")[1]
			  products[ingressType] = append(products[ingressType], ingress)
			}
  
			for _, egress := range virtualFirewall.GetEgress() {
			  productionOrder = append(productionOrder, egress)
			  egressType := strings.Split(reflect.TypeOf(egress).String(), "*k8sml.")[1]
			  products[egressType] = append(products[egressType], egress)
			}
  
			for _, role := range virtualFirewall.GetRoles() {
			  roleType := strings.Split(reflect.TypeOf(role).String(), "*k8sml.")[1]
			  products[roleType] = append(products[roleType], role)
			  for _, vm := range role.VirtualMachines {
				productionOrder = append(productionOrder, vm)
				vmType := strings.Split(reflect.TypeOf(vm).String(), "*k8sml.")[1]
				products[vmType] = append(products[vmType], vm)
				
				keyType := strings.Split(reflect.TypeOf(vm.GetKey()).String(), "*k8sml.")[1]
				products[keyType] = append(products[keyType], vm.GetKey())
  
				imageType := strings.Split(reflect.TypeOf(vm.GetImage()).String(), "*k8sml.")[1]
				products[imageType] = append(products[imageType], vm.GetImage())
			  }
			}
  
			for _, tg := range virtualFirewall.GetTargetGroups() {
			  tgType := strings.Split(reflect.TypeOf(tg).String(), "*k8sml.")[1]
			  products[tgType] = append(products[tgType], tg)
			  role := tg.Target
			  roleType := strings.Split(reflect.TypeOf(role).String(), "*k8sml.")[1]
			  products[roleType] = append(products[roleType], role)
			  for _, vm := range role.VirtualMachines {
				productionOrder = append(productionOrder, vm)
				vmType := strings.Split(reflect.TypeOf(vm).String(), "*k8sml.")[1]
				products[vmType] = append(products[vmType], vm)
				
				keyType := strings.Split(reflect.TypeOf(vm.GetKey()).String(), "*k8sml.")[1]
				products[keyType] = append(products[keyType], vm.GetKey())
  
				imageType := strings.Split(reflect.TypeOf(vm.GetImage()).String(), "*k8sml.")[1]
				products[imageType] = append(products[imageType], vm.GetImage())
			  }
			}
		  }
		}
	  }
	}
  
	productionPlan.ProductionOrder = productionOrder
	productionPlan.Products = products
	
	return productionPlan
  }
  
  func (productionPlan *ProductionPlan) BuildInfrastructure() error {
	if err := productionPlan.CleanupInfrastructure(); err != nil {
	  log.Fatalf("error: %v", err)
	}
	
	if err := productionPlan.Export(); err != nil {
	  log.Fatalf("error: %v", err)
	}
  
	if err := productionPlan.ProduceInfrastructure(); err != nil {
	  log.Fatalf("error: %v", err)
	}
  
	return nil
  }
  
  func (productionPlan *ProductionPlan) Export() error {
	file, err := os.OpenFile("variables.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	  if err != nil {
		  return err
	  }
  
	  defer file.Close()
	
	for i := 0; i < len(productionPlan.ProductionOrder); i++ {
	  productionPlan.ProductionOrder[i].ExportModule()
	}
  
	return nil
  }
  
  func (productionPlan *ProductionPlan) DeploySoftware(ssh string) error {
	for _, product := range productionPlan.Products["Role"] {
	  role := product.(*k8sml.Role)
	  if err := role.ExportInventory(); err != nil {
		return err
	  }
	}
  
	inventory := ansible.NewInventory("all:vars")
	inventory.AddVirtualMachine("")
	inventory.AddVariable("", "ansible_python_interpreter", "/usr/bin/python3")
	inventory.AddVariable("", "ansible_ssh_private_key_file", ssh)
  
	if err := inventory.Export(); err != nil {
	  return err
	}
  
	playbook := ansible.NewPlaybook()
  
	localhost := new(k8sml.Role)
	localhost.ID = "localhost"
	localhost.Software = make([]string, 0)
	localhost.Software = append(localhost.Software, "kubearch_localhost")
	productionPlan.Products["Role"] = append(productionPlan.Products["Role"], localhost)
  
	jumphost := new(k8sml.Role)
	jumphost.ID = "Jumphost"
	jumphost.Software = make([]string, 0)
	jumphost.Software = append(jumphost.Software, "kubearch_public_hosts")
	productionPlan.Products["Role"] = append(productionPlan.Products["Role"], jumphost)
  
	all := new(k8sml.Role)
	all.ID = "all"
	all.Software = make([]string, 0)
	all.Software = append(all.Software, "kubearch_all")
	productionPlan.Products["Role"] = append(productionPlan.Products["Role"], all)
  
	for _, product := range productionPlan.Products["Role"] {
	  role := product.(*k8sml.Role)
	  playbook.AddSoftwareStack(role.ID, role.Software)
	}
  
	if err := playbook.Export(); err != nil {
	  return err
	}
  
	time.Sleep(5 * time.Second)
	cmd := exec.Command("ansible-playbook", "-i", "hosts", "build.yml", "-vvv")
  
	out, err := cmd.StdoutPipe()
	if err != nil {
	  return err
	}
  
	scanner := bufio.NewScanner(out)
	go func() {
	  for scanner.Scan() {
		fmt.Println(scanner.Text())
	  }
	}()
  
	err = cmd.Start()
	if err != nil {
	  return err
	}
  
	err = cmd.Wait()
	if err != nil {
	  return err
	}
  
	return nil
  }
  
  func (productionPlan *ProductionPlan) Correction() {
	jumpHosts := make([]k8sml.VirtualMachine, 0)
  
	iamroles := productionPlan.Products["IAMRole"]
	var mRole *k8sml.IAMRole
	var wRole *k8sml.IAMRole
	for _, iamrole := range iamroles {
	  if iamrole.GetID() == "aws-cloudprovider-master" {
		mRole = iamrole.(*k8sml.IAMRole)
	  } else {
		wRole = iamrole.(*k8sml.IAMRole)
	  }
	}
	
	for _, route := range productionPlan.Products["Route"] {
	  route := route.(*k8sml.Route)
	  targetID := route.Target.GetTargetID()
  
	  switch route.Target.(type) {
	  case *k8sml.InternetGateway:
		for _, igw := range productionPlan.Products["InternetGateway"] {
		  igw := igw.(*k8sml.InternetGateway)
  
		  if targetID == igw.ID {
			route.Target = igw
		  }
		}
	  case *k8sml.NATGateway:
		for _, ngw := range productionPlan.Products["NATGateway"] {
		  ngw := ngw.(*k8sml.NATGateway)
  
		  if targetID == ngw.ID {
			route.Target = ngw
		  }
		}
	  }
	}
  
	for _, product := range productionPlan.Products["Role"] {
	  role := product.(*k8sml.Role)
  
	  if role.ID == "Jumphost" {
		jumpHosts = role.VirtualMachines
	  }
	}
  
	for _, product := range productionPlan.Products["Role"] {
	  role := product.(*k8sml.Role)
  
	  for _, vm := range role.VirtualMachines {
		vm.SetJumpHosts(jumpHosts)
	  }
	}
  
	for _, product := range productionPlan.Products["TargetGroup"] {
	  tg := product.(*k8sml.TargetGroup)
	  lbID := tg.LoadBalancer.ID
  
	  for _, product := range productionPlan.Products["NetworkLoadBalancer"] {
		nlb := product.(*k8sml.NetworkLoadBalancer)
  
		if nlb.ID == lbID { 
		  tg.LoadBalancer = nlb
		  nlb.TargetGroup = tg
		}
	  }
	}
  
	for _, vm := range productionPlan.Products["EC2Instance"] {
	  instance := vm.(*k8sml.EC2Instance)
	  if instance.IAMRole != nil && instance.IAMRole.ID == "aws-cloudprovider-master" {
		instance.IAMRole = mRole
		mRole.VirtualMachine = append(mRole.VirtualMachine, instance)
	  } else {
		instance.IAMRole = wRole
		wRole.VirtualMachine = append(wRole.VirtualMachine, instance)
	  }
	}
  }
  
  func (productionPlan *ProductionPlan) Init() error {
	cmd := exec.Command("terraform", "init")
  
	out, err := cmd.StdoutPipe()
	if err != nil {
	  return err
	}
  
	scanner := bufio.NewScanner(out)
	go func() {
	  for scanner.Scan() {
		fmt.Println(scanner.Text())
	  }
	}()
  
	err = cmd.Start()
	if err != nil {
	  return err
	}
  
	err = cmd.Wait()
	if err != nil {
	  return err
	}
  
	return nil
  }
  
  func (productionPlan *ProductionPlan) ProduceInfrastructure() error {
	productionPlan.Correction()
	productionPlan.Init()
	
	cmd := exec.Command("terraform", "apply", "-auto-approve")
  
	out, err := cmd.StdoutPipe()
	if err != nil {
	  return err
	}
  
	scanner := bufio.NewScanner(out)
	go func() {
	  for scanner.Scan() {
		fmt.Println(scanner.Text())
	  }
	}()
  
	err = cmd.Start()
	if err != nil {
	  return err
	}
  
	err = cmd.Wait()
	if err != nil {
	  return err
	}
  
	return nil
  }
  
  func (productionPlan *ProductionPlan) Update(productType string) error {
	for _, product := range productionPlan.Products[productType] {
	  runtimeVariables := product.(k8sml.Infrastructure).GetRuntimeVariables()
  
	  for key, _ := range runtimeVariables {
		productType = strings.Split(reflect.TypeOf(product).String(), "*k8sml.")[1]
		variable := product.GetID() + "_" + strings.ToLower(productType) + "_" + key
  
		cmd := exec.Command("terraform", "output", variable)
		var out bytes.Buffer
		cmd.Stdout = &out
  
		if err := cmd.Run(); err != nil {
		  return err
		}
  
		runtimeVariables[key] = strings.TrimSuffix(out.String(), "\n")
	  }
	}
  
	return nil
  }
  
  func (productionPlan *ProductionPlan) CleanupInfrastructure() error {
	files := make([]string, 0)
	files = append(files, "main.tf", "variables.tf", "outputs.tf", "terraform.tfstate")
  
	for _, file := range files {
	  if _, err := os.Stat(file); err == nil {
		if err := os.Remove(file); err != nil {
		  return err
		}
	  }
	}
	
	return nil
  }
  
  func (productionPlan *ProductionPlan) GetProduct(element, id string) k8sml.K8sML {
	for _ , product := range productionPlan.Products[element] {
	  if product.GetID() == id {
		return product
	  }
	}
  
	return nil
  }
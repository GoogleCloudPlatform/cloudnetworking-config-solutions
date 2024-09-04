// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package integrationtest

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
)

const (
	terraformDirectoryPath   = "../../../02-networking"
	region                   = "us-west2"
	peerASN                  = 64513
	psaRangeName             = "testpsarange"
	psaRange                 = "10.0.64.0/20"
	tunnel1BGPPeerASNAddress = "169.254.1.1"
	tunnel1SharedSecret      = "secret1"
	tunnel2BGPPeerASNAddress = "169.254.2.1"
	tunnel2SharedSecret      = "secret2"
)

var (
	projectID          = os.Getenv("TF_VAR_project_id")
	uniqueID           = rand.Int() //included as a suffix to the VPC and subnet names.
	networkName        = fmt.Sprintf("test-vpc-existing-%d", uniqueID)
	subnetworkName     = fmt.Sprintf("test-subnet-existing-%d", uniqueID)
	subnetworkIPCIDR   = "10.0.0.0/24"
	createInterconnect = true
)

// Name of the deployed dedicated interconnect received after deploying the resource in the test lab
// e.g. dedicated-ix-vpn-client-0
var deployedInterconnectName = os.Getenv("deployed_interconnect_name")

// Variables for Interconnect configuration.
var interconnectProjectID = os.Getenv("TF_VAR_interconnect_project_id")

var zone = "us-west2-a"
var subnetworkIPCidr = "10.0.0.0/24"
var deletionProtection = false

// Variables for Interconnect configuration.
var firstInterconnectName = "cso-lab-interconnect-1"
var secondInterconnectName = "cso-lab-interconnect-2"
var userSpecifiedIPRange = []string{"0.0.0.0/0", "199.36.154.8/30"}

// First vlan attachment configuration values.
var firstVaAsn = "65418"
var firstVlanAttachmentName = "vlan-attachment-a"
var firstVaBandwidth = "BPS_1G"

// Second vlan attachment configuration values.
var secondVaAsn = "65418"
var secondVlanAttachmentName = "vlan-attachment-b"
var secondVaBandwidth = "BPS_1G"

/*
This test creates all the resources including the vpc network, subnetwork along with a PSA range.

It then validates if
1. VPC network is created
2. Subnetwork is created
3. PSA range is created
*/
func TestCreateVPCNetworkModule(t *testing.T) {
	//wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)

	var (
		networkName    = fmt.Sprintf("test-vpc-new-%d", uniqueID)
		subnetworkName = fmt.Sprintf("test-subnet-new-%d", uniqueID)
		tfVars         = map[string]any{
			"project_id":             projectID,
			"region":                 region,
			"create_network":         true,
			"create_subnetwork":      true,
			"create_nat":             true,
			"create_havpn":           false,
			"create_scp_policy":      true,
			"subnets_for_scp_policy": []interface{}{subnetworkName},
			"subnets": []any{
				map[string]any{
					"ip_cidr_range": subnetworkIPCIDR,
					"name":          subnetworkName,
					"region":        region,
				},
			},
			"network_name":                 networkName,
			"tunnel_1_bgp_peer_asn":        peerASN,
			"tunnel_2_bgp_peer_asn":        peerASN,
			"tunnel_1_bgp_peer_ip_address": tunnel1BGPPeerASNAddress,
			"tunnel_1_shared_secret":       tunnel1SharedSecret,
			"tunnel_2_bgp_peer_ip_address": tunnel2BGPPeerASNAddress,
			"tunnel_2_shared_secret":       tunnel2SharedSecret,
			"psa_range_name":               psaRangeName,
			"psa_range":                    psaRange,
		}
	)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars:                 tfVars,
		TerraformDir:         terraformDirectoryPath,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 60 seconds to let resource acheive stable state
	time.Sleep(60 * time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	want := networkName
	got := terraform.Output(t, terraformOptions, "name")

	t.Log(" ========= Verify Network Name ========= ")
	if got != want {
		t.Errorf("Network with invalid name created = %v, want = %v", got, want)
	}

	t.Log(" ========= Verify Subnetwork Id ========= ")
	got = terraform.Output(t, terraformOptions, "subnet_ids")
	subnetworkID := fmt.Sprintf("[projects/%s/regions/%s/subnetworks/%s]", projectID, region, subnetworkName)
	wantSubnetworkID := subnetworkID
	if got != wantSubnetworkID {
		t.Errorf("Subnetwork with invalid subnetwork ID is created = %v, want = %v", got, wantSubnetworkID)
	}

	// Verify Service Connection Policy from Terraform Output
	t.Logf("======= Verify Service Connection Policy (Terraform Output) =======")
	output := terraform.OutputJson(t, terraformOptions, "service_connection_policy_details") // Assuming this is your output

	defaultServiceClass := "gcp-memorystore-redis"
	policyName := fmt.Sprintf("SCP-%s-%s", networkName, defaultServiceClass)

	if !gjson.Get(output, "0.name").Exists() { // Assuming "0" is your key, change if needed
		t.Errorf("Service Connection Policy '%s' not found in Terraform output", policyName)
	}

	// Check if policy details are as expected (customize as needed)
	if gjson.Get(output, "0.name").String() != policyName { // Changed key to "0"
		t.Errorf("Service Connection Policy name mismatch: got %s, want %s",
			gjson.Get(output, "0.name").String(), policyName)
	}

	t.Logf("Service Connection Policy '%s' verified successfully in Terraform output.", policyName)

	t.Log(" ========= Verify PSA Range ========= ")
	vpcOutputValue := terraform.OutputJson(t, terraformOptions, "vpc_networks")
	if !gjson.Valid(vpcOutputValue) {
		t.Errorf("Error parsing output, invalid json: %s", vpcOutputValue)
	}
	result := gjson.Parse(vpcOutputValue)
	psaRangeNamePath := fmt.Sprintf("subnets_psa.%s.name", psaRangeName)
	got = gjson.Get(result.String(), psaRangeNamePath).String()
	want = psaRangeName
	if got != want {
		t.Errorf("Invalid PSA range created = %v, want = %v", got, want)
	}
}

/*
This test utilizes existing VPC and subnet and then creates all the resources
along with a PSA range.

It then validates if
1. Existing VPC network exists and can be used.
2. Existing subnetwork exists and can be used.
3. PSA range is created.
*/
func TestExistingVPCNetworkModule(t *testing.T) {
	// wait for 60 seconds to allow resources to be available.
	time.Sleep(60 * time.Second)
	var (
		tfVars = map[string]any{
			"project_id":             projectID,
			"region":                 region,
			"create_network":         false,
			"create_subnetwork":      false,
			"create_nat":             true,
			"create_havpn":           false,
			"create_scp_policy":      false,
			"subnets_for_scp_policy": []string{""},
			"subnets": []any{
				map[string]any{
					"ip_cidr_range": subnetworkIPCIDR,
					"name":          subnetworkName,
					"region":        region,
				},
			},
			"network_name":                 networkName,
			"tunnel_1_bgp_peer_asn":        peerASN,
			"tunnel_2_bgp_peer_asn":        peerASN,
			"tunnel_1_bgp_peer_ip_address": tunnel1BGPPeerASNAddress,
			"tunnel_1_shared_secret":       tunnel1SharedSecret,
			"tunnel_2_bgp_peer_ip_address": tunnel2BGPPeerASNAddress,
			"tunnel_2_shared_secret":       tunnel2SharedSecret,
			"psa_range_name":               psaRangeName,
			"psa_range":                    psaRange,
		}
	)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars:                 tfVars,
		TerraformDir:         terraformDirectoryPath,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})

	// Create VPC and subnet outside of the terraform module.
	createVPCSubnets(t, projectID, networkName, subnetworkName, region)

	// Delete VPC and subnet created outside of the terraform module.
	defer deleteVPCSubnets(t, projectID, networkName, subnetworkName, region)

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// wait for 60 seconds to let resource acheive stable state.
	time.Sleep(60 * time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	want := networkName
	got := terraform.Output(t, terraformOptions, "name")

	t.Log(" ========= Verify Network Name ========= ")
	if got != want {
		t.Errorf("Network with invalid name created = %v, want = %v", got, want)
	}

	t.Logf(" ========= Verify Subnetwork Id ========= ")
	got = terraform.Output(t, terraformOptions, "subnet_ids")
	subnetworkID := fmt.Sprintf("[projects/%s/regions/%s/subnetworks/%s]", projectID, region, subnetworkName)
	wantSubnetworkID := subnetworkID
	if got != wantSubnetworkID {
		t.Errorf("Subnetwork with invalid sub network id created = %v, want = %v", got, wantSubnetworkID)
	}

	// Create SCP outside of terraform
	defaultServiceClass := "gcp-memorystore-redis"
	policyName := fmt.Sprintf("SCP-%s-%s", networkName, defaultServiceClass)
	createServiceConnectionPolicy(t, projectID, region, networkName, policyName, subnetworkName, defaultServiceClass, 5) //Pass the correct parameters

	t.Logf("======= Verify Service Connection Policy (Terraform Output) =======")
	output := gjson.Parse(terraform.OutputJson(t, terraformOptions, "service_connection_policy_details"))

	if !output.Get(defaultServiceClass + ".name").Exists() {
		t.Logf("Service Connection Policy '%s' was correctly not created by terraform", policyName)
	}
	t.Logf("======= Verify Service Connection Policy using gcloud =======")

	// Check if policy exists using gcloud describe
	scpOutput, err := shell.RunCommandAndGetOutputE(t, shell.Command{
		Command: "gcloud",
		Args: []string{
			"network-connectivity", "service-connection-policies", "describe", policyName,
			"--project", projectID,
			"--region", region,
			"--verbosity=none",
			"--format", "json", // Format output as JSON for easy parsing
		},
	})

	if err != nil {
		t.Errorf("Error: Service Connection Policy '%s' not found or could not be described: %s", policyName, err)
	}

	// Validate if valid JSON received
	if !gjson.Valid(scpOutput) {
		t.Errorf("Error parsing output, invalid json: %s", scpOutput)
	}
	// Parse gcloud output using gjson
	policyDetails := gjson.Parse(scpOutput)

	expectedPolicyName := fmt.Sprintf("projects/%s/locations/%s/serviceConnectionPolicies/%s", projectID, region, policyName)
	gotPolicyName := gjson.Get(policyDetails.String(), "name").String()
	// Check if policy details are as expected
	if gotPolicyName != expectedPolicyName {
		t.Errorf("Service Connection Policy name mismatch: got %s, want %s", gotPolicyName, expectedPolicyName)
	} else {
		t.Logf("=============== Service Connection Policy '%s' verified successfully.================", gotPolicyName)
	}
	gotNetworkName := gjson.Get(policyDetails.String(), "network").String()
	expectedNetworkName := fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkName)
	if gotNetworkName != expectedNetworkName {
		t.Errorf("Service Connection Policy network mismatch: got %s, want %s", gotNetworkName, expectedNetworkName)
	} else {
		t.Logf("=============== Service Connection Policy Network '%s' verified successfully.================", gotNetworkName)
	}

	//Get self link of the existing subnet you created manually
	existingSubnetSelfLink := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", projectID, region, subnetworkName)

	subnets := policyDetails.Get("pscConfig.subnetworks")
	if len(subnets.Array()) > 0 {
		if subnets.Array()[0].String() != existingSubnetSelfLink {
			t.Errorf("Service Connection Policy subnetwork mismatch: got %s, want %s", subnets.Array()[0].String(), existingSubnetSelfLink)
		}
	} else {
		t.Log("No subnets specified in Service Connection Policy, which is acceptable in this test scenario.")
	}

	t.Log(" ========= Verify PSA Range ========= ")
	vpcOutputValue := terraform.OutputJson(t, terraformOptions, "vpc_networks")
	if !gjson.Valid(vpcOutputValue) {
		t.Errorf("Error parsing output, invalid json: %s", vpcOutputValue)
	}
	result := gjson.Parse(vpcOutputValue)
	psaRangeNamePath := fmt.Sprintf("subnets_psa.%s.name", psaRangeName)
	got = gjson.Get(result.String(), psaRangeNamePath).String()
	want = psaRangeName
	if got != want {
		t.Errorf("Invalid PSA range created = %v, want = %v", got, want)
	}

}

/*
deleteVPCSubnets is a helper function which deletes the VPC and subnets after
completion of the test expecting to use existing VPC and subnets.
*/
func deleteVPCSubnets(t *testing.T, projectID string, networkName string, subnetworkName string, region string) {
	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "subnets", "delete", subnetworkName, "--region=" + region, "--project=" + projectID, "--quiet"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}

	// Sleep for 60 seconds to ensure the deleted subnets is reliably reflected.
	time.Sleep(60 * time.Second)

	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "delete", networkName, "--project=" + projectID, "--quiet"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
}

/*
createVPCSubnets is a helper function which creates the VPC and subnets before
execution of the test expecting to use existing VPC and subnets.
*/

func createVPCSubnets(t *testing.T, projectID string, networkName string, subnetworkName string, region string) {
	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", networkName, "--project=" + projectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom", "--verbosity=none"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
	time.Sleep(60 * time.Second)
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "subnets", "create", subnetworkName, "--network=" + networkName, "--project=" + projectID, "--range=" + subnetworkIPCIDR, "--region=" + region, "--format=json", "--enable-private-ip-google-access", "--enable-flow-logs", "--verbosity=none"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===Error %s Encountered while executing %s", err, text)
	}
}

// Function to create Service Connection Policy
func createServiceConnectionPolicy(t *testing.T, projectID string, region string, networkName string, policyName string, subnetworkName string, serviceClass string, connectionLimit int) {
	// Get subnet self link from subnet ID using gcloud command
	subnetSelfLink := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/regions/%s/subnetworks/%s", projectID, region, subnetworkName)

	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"network-connectivity", "service-connection-policies", "create",
			policyName, // Add the policyName here as the first argument after "create"
			"--project", projectID,
			"--region", region,
			"--network", networkName,
			"--service-class", serviceClass,
			"--subnets", subnetSelfLink,
			"--psc-connection-limit", fmt.Sprintf("%d", connectionLimit),
			"--quiet",
		},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("Error creating Service Connection Policy: %s", err)
	}
}

/*
	TestInterconnectWithVPCCreation tests the creation of

interconnect.tf example by creating a new vpc and a new subnet.
*/
func TestInterconnectWithVPCCreation(t *testing.T) {
	if deployedInterconnectName == "" {
		t.Skip("Skipping Interconnect testing.")
	}
	deploymentNumber, err := strconv.Atoi(deployedInterconnectName[len(deployedInterconnectName)-1:])
	if err != nil {
		t.Errorf("Deployment number is not an int, using default value for deployment number.")
		deploymentNumber = 1
	}
	var icRouterBgpAsn = 65000 + deploymentNumber
	var firstVaBgpRange = fmt.Sprintf("169.254.6%d.0/29", deploymentNumber)
	var firstVlanTag = 600 + deploymentNumber
	var secondVaBgpRange = fmt.Sprintf("169.254.6%d.8/29", deploymentNumber)
	var secondVlanTag = 600 + deploymentNumber
	var tfVars = map[string]any{
		"project_id":          projectID,
		"region":              region,
		"create_network":      true,
		"create_subnetwork":   true,
		"create_nat":          true,
		"create_havpn":        false,
		"create_scp_policy":   false,
		"create_interconnect": createInterconnect,
		"subnets": []any{
			map[string]any{
				"ip_cidr_range": subnetworkIPCIDR,
				"name":          subnetworkName,
				"region":        region,
			},
		},
		"network_name":                 networkName,
		"tunnel_1_bgp_peer_asn":        peerASN,
		"tunnel_2_bgp_peer_asn":        peerASN,
		"tunnel_1_bgp_peer_ip_address": tunnel1BGPPeerASNAddress,
		"tunnel_1_shared_secret":       tunnel1SharedSecret,
		"tunnel_2_bgp_peer_ip_address": tunnel2BGPPeerASNAddress,
		"tunnel_2_shared_secret":       tunnel2SharedSecret,
		"psa_range_name":               psaRangeName,
		"psa_range":                    psaRange,
		"user_specified_ip_range":      userSpecifiedIPRange,
		"interconnect_project_id":      interconnectProjectID,
		"first_interconnect_name":      firstInterconnectName,
		"second_interconnect_name":     secondInterconnectName,
		"first_va_name":                firstVlanAttachmentName,
		"ic_router_bgp_asn":            icRouterBgpAsn,
		"first_va_asn":                 firstVaAsn,
		"first_va_bandwidth":           firstVaBandwidth,
		"first_va_bgp_range":           firstVaBgpRange,
		"first_vlan_tag":               firstVlanTag,
		"second_va_asn":                secondVaAsn,
		"second_va_name":               secondVlanAttachmentName,
		"second_va_bandwidth":          secondVaBandwidth,
		"second_va_bgp_range":          secondVaBgpRange,
		"second_vlan_tag":              secondVlanTag,
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir:         terraformDirectoryPath,
		Vars:                 tfVars,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})
	initiateTestForNetworkResource(t, terraformOptions, firstVlanTag)
}

/*
TestInterconnectWithoutVPCCreation tests the creation of example by using the existing vpc and  subnet.
*/
func TestInterconnectWithoutVPCCreation(t *testing.T) {
	if deployedInterconnectName == "" {
		t.Skip("Skipping Interconnect testing.")
	}
	deploymentNumber, err := strconv.Atoi(deployedInterconnectName[len(deployedInterconnectName)-1:])
	if err != nil {
		t.Errorf("Deployment number is not an int, using default value for deployment number.")
		deploymentNumber = 1
	}
	var icRouterBgpAsn = 65000 + deploymentNumber
	var firstVaBgpRange = fmt.Sprintf("169.254.6%d.0/29", deploymentNumber)
	var firstVlanTag = 600 + deploymentNumber
	var secondVaBgpRange = fmt.Sprintf("169.254.6%d.8/29", deploymentNumber)
	var secondVlanTag = 600 + deploymentNumber
	ProjectID := projectID
	var tfVars = map[string]any{
		"project_id":          projectID,
		"region":              region,
		"create_network":      false,
		"create_subnetwork":   false,
		"create_nat":          true,
		"create_havpn":        false,
		"create_scp_policy":   false,
		"create_interconnect": createInterconnect,
		"subnets": []any{
			map[string]any{
				"ip_cidr_range": subnetworkIPCIDR,
				"name":          subnetworkName,
				"region":        region,
			},
		},
		"network_name":                 networkName,
		"tunnel_1_bgp_peer_asn":        peerASN,
		"tunnel_2_bgp_peer_asn":        peerASN,
		"tunnel_1_bgp_peer_ip_address": tunnel1BGPPeerASNAddress,
		"tunnel_1_shared_secret":       tunnel1SharedSecret,
		"tunnel_2_bgp_peer_ip_address": tunnel2BGPPeerASNAddress,
		"tunnel_2_shared_secret":       tunnel2SharedSecret,
		"psa_range_name":               psaRangeName,
		"psa_range":                    psaRange,
		"user_specified_ip_range":      userSpecifiedIPRange,
		"interconnect_project_id":      interconnectProjectID,
		"first_interconnect_name":      firstInterconnectName,
		"second_interconnect_name":     secondInterconnectName,
		"first_va_name":                firstVlanAttachmentName,
		"ic_router_bgp_asn":            icRouterBgpAsn,
		"first_va_asn":                 firstVaAsn,
		"first_va_bandwidth":           firstVaBandwidth,
		"first_va_bgp_range":           firstVaBgpRange,
		"first_vlan_tag":               firstVlanTag,
		"second_va_asn":                secondVaAsn,
		"second_va_name":               secondVlanAttachmentName,
		"second_va_bandwidth":          secondVaBandwidth,
		"second_va_bgp_range":          secondVaBgpRange,
		"second_vlan_tag":              secondVlanTag,
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir:         terraformDirectoryPath,
		Vars:                 tfVars,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})
	// Create VPC and subnet outside of the terraform module

	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", networkName, "--project=" + ProjectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom", "--verbosity=none"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "subnets", "create", subnetworkName, "--network=" + networkName, "--project=" + ProjectID, "--range=10.0.0.0/24", "--region=" + region, "--format=json", "--enable-private-ip-google-access", "--enable-flow-logs", "--verbosity=none"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}
	initiateTestForNetworkResource(t, terraformOptions, firstVlanTag)
}

/*
	initiateTestForNetworkResource is a helper function that helps in verification

of the resources being created as part of test.
*/
func initiateTestForNetworkResource(t *testing.T, terraformOptions *terraform.Options, firstVlanTag int) {
	t.Helper()

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 60 seconds to let resource achieve stable state.
	time.Sleep(60 * time.Second)

	log.Println(" ========= Verify Subnet Name ========= ")
	want := networkName
	got := terraform.Output(t, terraformOptions, "name")
	if !cmp.Equal(got, want) {
		t.Errorf("Test Network Name = %v, want = %v", got, want)
	}
	// Validate if interconnects vlans attachments are up & running with Established Connection.
	log.Println(" ====================================================== ")
	log.Println(" ========= Verify Interconnect/VLAN Tunnels ========= ")

	var interconnectAttachmentNameList = []string{firstVlanAttachmentName, secondVlanAttachmentName}

	for _, vlanAttachmentName := range interconnectAttachmentNameList {
		ProjectID := projectID
		text := "compute"
		cmd := shell.Command{
			Command: "gcloud",
			Args:    []string{"compute", "interconnects", "attachments", "describe", vlanAttachmentName, "--region", region, "--project", ProjectID, "--format=json", "--verbosity=none"},
		}
		op, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err != nil {
			log.Printf("===Error %s Encountered while executing %s", err, text)
		}
		if !gjson.Valid(op) {
			t.Fatalf("Error parsing output, invalid json: %s", op)
		}
		result := gjson.Parse(op)
		if err != nil {
			log.Printf("=== Error %s Encountered while executing %s", err, text)
		}
		log.Printf(" \n========= Validating attachment %s ============\n", vlanAttachmentName)
		log.Println(" ========= Check if attach Operation Status is active ========= ")
		want = "OS_ACTIVE"
		got = gjson.Get(result.String(), "operationalStatus").String()
		if !cmp.Equal(got, want) {
			t.Errorf("Test VLAN Operational State = %v, want = %v", got, want)
		}
		log.Println(" ========= Check if state is Active ========= ")
		want = "ACTIVE"
		got = gjson.Get(result.String(), "state").String()
		if !cmp.Equal(got, want) {
			t.Errorf("Test VLAN State = %v, want = %v", got, want)
		}
		log.Println(" ========= Check if type is Dedicated ========= ")
		want = "DEDICATED"
		got = gjson.Get(result.String(), "type").String()
		if !cmp.Equal(got, want) {
			t.Errorf("Test Interconnect type = %v, want = %v", got, want)
		}

		log.Println(" ========= Check if vlan tag is Same as Configured ========= ")
		want = strconv.Itoa(firstVlanTag)
		got = gjson.Get(result.String(), "vlanTag8021q").String()
		if !cmp.Equal(got, want) {
			t.Errorf("Test VLAN tag = %v, want = %v", got, want)
		}

	}
}

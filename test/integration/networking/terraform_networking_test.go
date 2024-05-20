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
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/tidwall/gjson"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	terraformDirectoryPath   = "../../../02-networking"
	region                   = "us-central1"
	peerASN                  = 64513
	psaRangeName             = "testpsarange"
	psaRange                 = "10.0.64.0/20"
	tunnel1BGPPeerASNAddress = "169.254.1.1"
	tunnel1SharedSecret      = "secret1"
	tunnel2BGPPeerASNAddress = "169.254.2.1"
	tunnel2SharedSecret      = "secret2"
)

var (
	projectID        = os.Getenv("TF_VAR_project_id")
	uniqueID         = rand.Int() //included as a suffix to the VPC and subnet names.
	networkName      = fmt.Sprintf("test-vpc-existing-%d", uniqueID)
	subnetworkName   = fmt.Sprintf("test-subnet-existing-%d", uniqueID)
	subnetworkIPCIDR = "10.0.0.0/24"
)

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
			"project_id":        projectID,
			"region":            region,
			"create_network":    true,
			"create_subnetwork": true,
			"create_nat":        true,
			"create_havpn":      false,
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
			"project_id":        projectID,
			"region":            region,
			"create_network":    false,
			"create_subnetwork": false,
			"create_nat":        true,
			"create_havpn":      false,
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

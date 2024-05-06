// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package unittest

import (
	compare "cmp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"golang.org/x/exp/slices"
	"testing"
)

const (
	terraformDirectoryPath   = "../../../02-networking"
	region                   = "us-central1"
	networkName              = "unit-test-vpc-1"
	peerASN                  = 64513
	tunnel1BGPPeerASNAddress = "169.254.1.1"
	tunnel1SharedSecret      = "secret1"
	tunnel2BGPPeerASNAddress = "169.254.2.1"
	tunnel2SharedSecret      = "secret2"
)

// Unit tests for VPC network, subnet, Cloud NAT, and HA VPN creation.
var (
	projectID = "dummy-project-id"
	tfVars    = map[string]any{
		"project_id":        projectID,
		"region":            region,
		"create_network":    true,
		"create_subnetwork": true,
		"create_nat":        true,
		"create_havpn":      true,
		"subnets": []any{
			map[string]any{
				"ip_cidr_range": "10.0.0.0/24",
				"name":          "unit-test-subnet-1",
				"region":        region,
			},
			map[string]any{
				"ip_cidr_range": "10.0.16.0/24",
				"name":          "unit-test-subnet-2",
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
	}
)

func TestInitAndPlanRunWithTfVars(t *testing.T) {
	/*
	 0 = Succeeded with empty diff (no changes)
	 1 = Error
	 2 = Succeeded with non-empty diff (changes present)
	*/
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	want := 2
	got := planExitCode
	if got != want {
		t.Errorf("Test Plan Exit Code = %v, want = %v", got, want)
	}
}
func TestInitAndPlanRunWithoutTfVarsExpectFailureScenario(t *testing.T) {
	/*
	 0 = Succeeded with empty diff (no changes)
	 1 = Error
	 2 = Succeeded with non-empty diff (changes present)
	*/
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	want := 1
	got := planExitCode
	if !cmp.Equal(got, want) {
		t.Errorf("Test Plan Exit Code = %v, want = %v", got, want)
	}
}

func TestResourcesCount(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planStruct := terraform.InitAndPlan(t, terraformOptions)
	resourceCount := terraform.GetResourceCount(t, planStruct)
	if got, want := resourceCount.Add, 21; got != want {
		t.Errorf("Test Resource Count Add = %v, want = %v", got, want)
	}
	if got, want := resourceCount.Change, 0; got != want {
		t.Errorf("Test Resource Count Change = %v, want = %v", got, want)
	}
	if got, want := resourceCount.Destroy, 0; got != want {
		t.Errorf("Test Resource Count Destroy = %v, want = %v", got, want)
	}
}

func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	expectedModuleAddresses := []string{"module.vpc_network", "module.nat[0]", "module.havpn[0]"}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planStruct := terraform.InitAndPlanAndShow(t, terraformOptions)
	content, err := terraform.ParsePlanJSON(planStruct)
	if err != nil {
		t.Fatal(err.Error())
	}
	actualModuleAddresses := make([]string, 0)
	for _, element := range content.ResourceChangesMap {
		if element.ModuleAddress != "" && !slices.Contains(actualModuleAddresses, element.ModuleAddress) {
			actualModuleAddresses = append(actualModuleAddresses, element.ModuleAddress)
		}
	}
	want := expectedModuleAddresses
	got := actualModuleAddresses
	if !cmp.Equal(got, want, cmpopts.SortSlices(compare.Less[string])) {
		t.Errorf("Test Element Mismatch = %v, want = %v", got, want)
	}
}

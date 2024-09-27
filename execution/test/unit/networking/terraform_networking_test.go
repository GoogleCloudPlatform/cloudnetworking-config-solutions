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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"golang.org/x/exp/slices"
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
		"project_id":             projectID,
		"region":                 region,
		"create_network":         true,
		"create_subnetwork":      true,
		"create_nat":             true,
		"create_havpn":           true,
		"create_interconnect":    true,
		"create_scp_policy":      true,
		"subnets_for_scp_policy": []interface{}{"unit-test-subnet-1"}, // Changed to []interface{}
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
		"interconnect_project_id":      interconnectProjectID,
		"first_interconnect_name":      firstInterconnectName,
		"second_interconnect_name":     secondInterconnectName,
		"ic_router_bgp_asn":            icRouterBgpAsn,
		"first_va_asn":                 firstVaAsn,
		"first_va_bandwidth":           firstVaBandwidth,
		"first_va_bgp_range":           firstVaBgpRange,
		"first_vlan_tag":               firstVlanTag,
		"second_va_asn":                secondVaAsn,
		"second_va_bandwidth":          secondVaBandwidth,
		"second_va_bgp_range":          secondVaBgpRange,
		"second_vlan_tag":              secondVlanTag,
	}
)

// variables for Interconnect Configuration
var zone = "us-west2-a"
var subnetworkName = "cloudsql-easy-subnet"
var subnetworkIPCidr = "10.2.0.0/16"
var deletionProtection = false

// variables for Interconnect configuration
var interconnectProjectID = "dummy-interconnect-project-id"
var firstInterconnectName = "cso-lab-interconnect-1"
var secondInterconnectName = "cso-lab-interconnect-2"
var icRouterBgpAsn = 65004

// first vlan attachment configuration values
var firstVaAsn = "65418"
var firstVaBandwidth = "BPS_1G"
var firstVaBgpRange = "169.254.61.0/29"
var firstVlanTag = 601

// second vlan attachment configuration values
var secondVaAsn = "65418"
var secondVaBandwidth = "BPS_1G"
var secondVaBgpRange = "169.254.61.8/29"
var secondVlanTag = 601

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
	if got, want := resourceCount.Add, 29; got != want {
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
	expectedModuleAddresses := []string{"module.vpc_network", "module.vlan_attachment_a[0]", "module.vlan_attachment_b[0]", "module.havpn[0]", "module.nat[0]"}
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
		t.Errorf("Error parsing Terraform plan: %v", err) // Detailed error message
		return                                            // Exit early on parsing error
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

func TestTerraformResourceAddressListMatch(t *testing.T) {
	expectedResourceAddresses := []string{
		"data.google_compute_network.vpc_network",
		"google_compute_route.default[0]",
		"google_network_connectivity_service_connection_policy.policy[0]",
		"module.havpn[0].google_compute_ha_vpn_gateway.ha_gateway[0]",
		"module.havpn[0].google_compute_router.router[0]",
		"module.havpn[0].google_compute_router_interface.router_interface[\"remote-0\"]",
		"module.havpn[0].google_compute_router_interface.router_interface[\"remote-1\"]",
		"module.havpn[0].google_compute_router_peer.bgp_peer[\"remote-0\"]",
		"module.havpn[0].google_compute_router_peer.bgp_peer[\"remote-1\"]",
		"module.havpn[0].google_compute_vpn_tunnel.tunnels[\"remote-0\"]",
		"module.havpn[0].google_compute_vpn_tunnel.tunnels[\"remote-1\"]",
		"module.havpn[0].random_id.secret",
		"module.nat[0].google_compute_router.router[0]",
		"module.nat[0].google_compute_router_nat.nat",
		"module.vlan_attachment_a[0].google_compute_interconnect_attachment.default",
		"module.vlan_attachment_a[0].google_compute_router_interface.default[0]",
		"module.vlan_attachment_a[0].google_compute_router_peer.default[0]",
		"module.vlan_attachment_b[0].google_compute_interconnect_attachment.default",
		"module.vlan_attachment_b[0].google_compute_router_interface.default[0]",
		"module.vlan_attachment_b[0].google_compute_router_peer.default[0]",
		"module.vpc_network.google_compute_global_address.psa_ranges[\"psarange\"]",
		"module.vpc_network.google_compute_network.network[0]",
		"module.vpc_network.google_compute_network_peering_routes_config.psa_routes[0]",
		"module.vpc_network.google_compute_route.gateway[\"private-googleapis\"]",
		"module.vpc_network.google_compute_route.gateway[\"restricted-googleapis\"]",
		"module.vpc_network.google_compute_subnetwork.subnetwork[\"us-central1/unit-test-subnet-1\"]",
		"module.vpc_network.google_compute_subnetwork.subnetwork[\"us-central1/unit-test-subnet-2\"]",
		"module.vpc_network.google_service_networking_connection.psa_connection[0]",
		"google_compute_router.interconnect-router[0]",
		"module.vpc_network.google_compute_shared_vpc_host_project.shared_vpc_host[0]",
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
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
		t.Errorf("Error parsing Terraform plan: %v", err) // Detailed error
		return                                            // Exit early on parsing error
	}

	actualResourceAddresses := make([]string, 0)
	resourcePolicyFound := false

	for _, element := range content.ResourceChangesMap {
		if element.Address != "" { // Check the Address field directly
			if !slices.Contains(actualResourceAddresses, element.Address) {
				actualResourceAddresses = append(actualResourceAddresses, element.Address)
			}

			// Check for your resource based on its complete address
			if element.Address == "google_network_connectivity_service_connection_policy.policy[0]" {
				resourcePolicyFound = true
				// Optionally add further checks on element.Change here (actions, values, etc.)
			}
		}
	}

	if !resourcePolicyFound {
		t.Errorf("Resource 'google_network_connectivity_service_connection_policy.policy[\"gcp-memorystore-redis\"]' not found in the plan.")
	}

	want := expectedResourceAddresses
	got := actualResourceAddresses

	if !cmp.Equal(got, want, cmpopts.SortSlices(compare.Less[string])) {
		t.Errorf("Test Element Mismatch = %v, want = %v", got, want)
	}
}

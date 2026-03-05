// Copyright 2025 Google LLC
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

package common_utils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	stdlib_strconv "strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
CreateVPCSubnets is a helper function which creates the VPC and subnets before
execution of the test expecting to use existing VPC and subnets.
*/

func CreateVPCSubnets(t *testing.T, projectID string, networkName string, subnetworkName string, region string) {
	subnetworkIPCIDR := "10.0.1.0/24"
	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", networkName, "--project=" + projectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom", "--verbosity=none"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===error %s encountered while executing %s", err, text)
	}
	time.Sleep(60 * time.Second)
	if subnetworkName != "" {
		if region == "" {
			region = "us-central1"
		}
		cmd = shell.Command{
			Command: "gcloud",
			Args:    []string{text, "networks", "subnets", "create", subnetworkName, "--network=" + networkName, "--project=" + projectID, "--range=" + subnetworkIPCIDR, "--region=" + region, "--format=json", "--enable-private-ip-google-access", "--enable-flow-logs", "--verbosity=none"},
		}
		_, err = shell.RunCommandAndGetOutputE(t, cmd)
		if err != nil {
			t.Errorf("===error %s encountered while executing %s", err, text)
		}
	} else {
		t.Log("VPC will be created & Subnet will not be created.")
	}

}

/*
DeleteVPCSubnets is a helper function which deletes the VPC and subnets after
completion of the test expecting to use existing VPC and subnets.
*/
func DeleteVPCSubnets(t *testing.T, projectID string, networkName string, subnetworkName string, region string) {
	text := "compute"
	if subnetworkName != "" {
		cmd := shell.Command{
			Command: "gcloud",
			Args:    []string{text, "networks", "subnets", "delete", subnetworkName, "--region=" + region, "--project=" + projectID, "--quiet"},
		}
		_, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err != nil {
			t.Errorf("===error %s encountered while executing %s", err, text)
		}
	}

	// Sleep for 60 seconds to ensure the deleted subnets is reliably reflected.
	time.Sleep(60 * time.Second)

	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "delete", networkName, "--project=" + projectID, "--quiet"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===error %s encountered while executing %s", err, text)
	}

}

/*
CreateServiceConnectionPolicy is a helped function that creates the service
connection policy.
*/
func CreateServiceConnectionPolicy(t *testing.T, projectID string, region string, networkName string, policyName string, subnetworkName string, serviceClass string, connectionLimit int) {
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
		t.Errorf("error creating Service Connection Policy: %s", err)
	}
}

/*
CreateGCEInstance creates a GCE VM with a startup script.
It uses --metadata-from-file for robustness.
*/
func CreateGCEInstance(t *testing.T, projectID, vmName, zone, subnetName, startupScript string, scopes string, hasExternalIP bool, imageProject string, imageFamily string) {
	// Create a temporary file with a predictable name based on the unique vmName.
	scriptFileName := fmt.Sprintf("startup-script-%s.sh", vmName)
	scriptFile, err := os.Create(scriptFileName)
	if err != nil {
		t.Fatalf("Failed to create temp file for startup script: %v", err)
	}
	defer os.Remove(scriptFile.Name()) // Clean up the file afterwards

	if _, err := scriptFile.WriteString(startupScript); err != nil {
		t.Fatalf("Failed to write to temp startup script file: %v", err)
	}
	if err := scriptFile.Close(); err != nil {
		t.Fatalf("Failed to close temp startup script file: %v", err)
	}

	if scopes == "" {
		scopes = "https://www.googleapis.com/auth/cloud-platform"
	}

	if imageProject == "" && imageFamily == "" {
		imageProject = "ubuntu-os-cloud"
		imageFamily = "ubuntu-2204-lts"
	}

	args := []string{
		"compute", "instances", "create", vmName,
		"--project", projectID,
		"--zone", zone,
		"--subnet", subnetName,
		"--scopes=" + scopes,
		"--metadata-from-file", fmt.Sprintf("startup-script=%s", scriptFile.Name()),
	}
	if !hasExternalIP {
		args = append(args, "--no-address")
	}

	if imageProject != "" {
		args = append(args, "--image-project", imageProject)
	}
	if imageFamily != "" {
		args = append(args, "--image-family", imageFamily)
	}

	// Use RunCommandAndGetOutputE for consistent and explicit error handling.
	cmd := shell.Command{
		Command: "gcloud",
		Args:    args,
	}
	if _, err := shell.RunCommandAndGetOutputE(t, cmd); err != nil {
		t.Fatalf("Failed to create GCE instance with command 'gcloud %s': %v", strings.Join(args, " "), err)
	}
}

/*
DeleteGCEInstance cleans up the GCE VM.
*/
func DeleteGCEInstance(t *testing.T, projectID, vmName, zone string) {
	// Use RunCommandAndGetOutputE for consistent and explicit error handling.
	cmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"compute", "instances", "delete", vmName,
			"--project", projectID,
			"--zone", zone,
			"--quiet",
		},
	}
	if _, err := shell.RunCommandAndGetOutputE(t, cmd); err != nil {
		t.Errorf("warning: failed to delete GCE instance '%s', this may require manual cleanup. error: %v", vmName, err)
	}
}

/*
CleanupConfigDir removes all .yaml files from a specified directory path.
*/
func CleanupConfigDir(t *testing.T, configPath string) {
	t.Logf("Cleaning up config directory: %s", configPath)
	files, err := filepath.Glob(filepath.Join(configPath, "*.yaml"))
	if err != nil {
		t.Fatalf("Error finding yaml files to clean up: %v", err)
	}
	for _, file := range files {
		t.Logf("Removing stale config file: %s", file)
		if err := os.Remove(file); err != nil {
			t.Logf("Warning: Could not remove stale config file %s: %v", file, err)
		}
	}
}

/*
GetSerialPortOutput gets the serial port output from a GCE instance.
*/
func GetSerialPortOutput(t *testing.T, projectID, vmName, zone string, port int) (string, error) {
	args := []string{
		"compute", "instances", "get-serial-port-output", vmName,
		"--project=" + projectID,
		"--zone=" + zone,
		fmt.Sprintf("--port=%d", port),
	}

	// Use Go's native exec package to run the command and capture its output
	// without printing it directly to the test log.
	cmd := exec.Command("gcloud", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// If the command fails, return the output (which contains the error message)
		// and the error itself for proper handling.
		return string(output), fmt.Errorf("gcloud get-serial-port-output failed: %w", err)
	}

	return string(output), nil
}

/*
DeletePSA is a helper function which deletes the PSA range after the
execution of the test.
*/

func DeletePSA(t *testing.T, projectID string, networkName string, rangeName string) {
	// Delete PSA IP range
	time.Sleep(60 * time.Second)
	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "addresses", "delete", rangeName, "--project=" + projectID, "--global", "--verbosity=info", "--format=json", "--quiet"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Logf("===Error %s Encountered while executing %s", err, text)
	}
	time.Sleep(60 * time.Second)
	// Delete PSA range
	text = "services"
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "vpc-peerings", "delete", "--service=servicenetworking.googleapis.com", "--project=" + projectID, "--network=" + networkName, "--verbosity=info", "--format=json", "--quiet"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Logf("===Error %s Encountered while executing %s", err, text)
	}
}

/*
CreatePSA is a helper function which creates the PSA range before the
execution of the test.
*/
func CreatePSA(t *testing.T, projectID string, networkName string, rangeName string) {
	// Create an IP range
	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "addresses", "create", rangeName, "--purpose=VPC_PEERING", "--addresses=10.0.64.0", "--prefix-length=20", "--project=" + projectID, "--network=" + networkName, "--global", "--verbosity=info", "--format=json"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===error %s encountered while executing %s", err, text)
	}
	// Create PSA range
	text = "services"
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "vpc-peerings", "connect", "--service=servicenetworking.googleapis.com", "--ranges=" + rangeName, "--project=" + projectID, "--network=" + networkName, "--verbosity=info", "--format=json"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Errorf("===error %s encountered while executing %s", err, text)
	}
	time.Sleep(60 * time.Second)
}

// getProjectNumber retrieves the project number for a given project ID.
func GetProjectNumber(t *testing.T, projectID string) (string, error) {
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"projects", "describe", projectID, "--format=value(projectNumber)"},
	}
	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		return "", fmt.Errorf("error getting project number: %v", err)
	}

	// The gcloud command might output warnings before the actual project number,
	// especially with impersonation. The project number is expected to be the last
	// non-empty line of the output.
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("no output from gcloud projects describe for project ID %s", projectID)
	}

	// Get the last line and trim any surrounding single quotes
	projectNumber := strings.Trim(lines[len(lines)-1], "'")

	// Basic validation that it looks like a number
	if _, err := stdlib_strconv.ParseInt(projectNumber, 10, 64); err != nil {
		return "", fmt.Errorf("extracted project number '%s' is not a valid number. Full output: %s. Error: %v", projectNumber, output, err)
	}

	return projectNumber, nil
}

/*
CreateFirewallPolicy is a helper function that creates a global network
firewall policy for an integration test.
*/
func CreateFirewallPolicy(t *testing.T, projectID, policyName string) {
	t.Logf("Creating Firewall Policy '%s'...", policyName)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"compute", "network-firewall-policies", "create", policyName, "--project=" + projectID, "--description=integration-test-policy", "--global"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if !assert.NoError(t, err, "Firewall Policy creation has Failed") {
		t.FailNow()
	}
}

/*
DeleteFirewallPolicy is a helper function that deletes a global network
firewall policy after a test completes.
*/
func DeleteFirewallPolicy(t *testing.T, projectID, policyName string) {
	t.Logf("--- Deleting Firewall Policy: %s ---", policyName)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"compute", "network-firewall-policies", "delete", policyName, "--project=" + projectID, "--global", "--quiet"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if !assert.NoError(t, err, "Firewall Policy deletion has failed.") {
		t.FailNow()
	}
}

/*
DescribeFirewallPolicyRule is a helper function that describes a specific
firewall policy rule for validation purposes.
*/
func DescribeFirewallPolicyRule(t *testing.T, projectID, policyName, rulePriority string) string {
	t.Logf("Describing Firewall Rule with priority '%s' in policy '%s'...", rulePriority, policyName)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"compute", "network-firewall-policies", "mirroring-rules", "describe", rulePriority, "--project=" + projectID, "--firewall-policy=" + policyName, "--global-firewall-policy"},
	}
	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	if !assert.NoError(t, err, "An Attempt to describe Firewall Policy has failed.") {
		t.FailNow()
	}
	return output
}

/*
CreateMirroringDeploymentGroup is a helper function that creates a packet
mirroring deployment group as a prerequisite for a test.
*/
func CreateMirroringDeploymentGroup(t *testing.T, projectID, dgName, vpcName string) {
	t.Logf("Creating Deployment Group '%s'...", dgName)
	networkURL := fmt.Sprintf("projects/%s/global/networks/%s", projectID, vpcName)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"network-security", "mirroring-deployment-groups", "create", dgName, "--project=" + projectID, "--location=global", "--network=" + networkURL},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if !assert.NoError(t, err, "Mirroring Deployment Group creation has failed.") {
		t.FailNow()
	}
}

/*
DeleteMirroringDeploymentGroup is a helper function that deletes a packet
mirroring deployment group after a test completes.
*/
func DeleteMirroringDeploymentGroup(t *testing.T, projectID, dgName string) {
	t.Logf("--- Deleting Deployment Group: %s ---", dgName)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"network-security", "mirroring-deployment-groups", "delete", dgName, "--project=" + projectID, "--location=global", "--quiet", "--no-async"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if !assert.NoError(t, err, "Mirroring Deployment group deletion has failed.") {
		t.FailNow()
	}
}

/*
CreateMirroringEndpointGroup is a helper function that creates a packet
mirroring endpoint group as a prerequisite for a test.
*/
func CreateMirroringEndpointGroup(t *testing.T, projectID, egName, dgName string) string {
	t.Logf("Creating Endpoint Group '%s'...", egName)
	deploymentGroupURL := fmt.Sprintf("projects/%s/locations/global/mirroringDeploymentGroups/%s", projectID, dgName)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"network-security", "mirroring-endpoint-groups", "create", egName, "--project=" + projectID, "--location=global", "--mirroring-deployment-group=" + deploymentGroupURL},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if !assert.NoError(t, err, "Mirroring Endpoint Group creation failed.") {
		t.FailNow()
	}
	return fmt.Sprintf("projects/%s/locations/global/mirroringEndpointGroups/%s", projectID, egName)
}

/*
DeleteMirroringEndpointGroup is a helper function that deletes a packet
mirroring endpoint group after a test completes.
*/
func DeleteMirroringEndpointGroup(t *testing.T, projectID, egName string) {
	t.Logf("--- Deleting Endpoint Group: %s ---", egName)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"network-security", "mirroring-endpoint-groups", "delete", egName, "--project=" + projectID, "--location=global", "--quiet", "--no-async"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if !assert.NoError(t, err, "Mirroring Endpoint Group deletion failed.") {
		t.FailNow()
	}
}

/*
CreateSecurityProfileAndGroup is a helper function that creates a security
profile and a security profile group for packet mirroring.
*/
func CreateSecurityProfileAndGroup(t *testing.T, orgID, projectID, spName, spgName, endpointGroupID string) {
	t.Logf("Creating Security Profile '%s' and Group '%s'...", spName, spgName)
	// Create the Security Profile
	cmdProfile := shell.Command{
		Command: "gcloud",
		Args:    []string{"network-security", "security-profiles", "custom-mirroring", "create", spName, "--organization=" + orgID, "--location=global", "--billing-project=" + projectID, "--mirroring-endpoint-group=" + endpointGroupID},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmdProfile)
	if !assert.NoError(t, err, "Security Profile creation failed.") {
		t.FailNow()
	}

	// Create the Security Profile Group and add the profile to it
	spPath := fmt.Sprintf("organizations/%s/locations/global/securityProfiles/%s", orgID, spName)
	cmdGroup := shell.Command{
		Command: "gcloud",
		Args:    []string{"network-security", "security-profile-groups", "create", spgName, "--organization=" + orgID, "--location=global", "--billing-project=" + projectID, "--custom-mirroring-profile=" + spPath},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmdGroup)
	if !assert.NoError(t, err, "Security Profile Group creation failed.") {
		t.FailNow()
	}
}

/*
DeleteSecurityProfileAndGroup is a helper function that deletes a security
profile and its associated group after a test completes.
*/
func DeleteSecurityProfileAndGroup(t *testing.T, orgID, spName, spgName string) {
	t.Logf("--- Deleting security profile group '%s' and profile '%s' ---", spgName, spName)
	// We delete the group first, then the profile
	cmdGroup := shell.Command{
		Command: "gcloud",
		Args:    []string{"network-security", "security-profile-groups", "delete", spgName, "--organization=" + orgID, "--location=global", "--quiet"},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmdGroup)
	if !assert.NoError(t, err, "Security Profile Group deletion failed.") {
		t.FailNow()
	}

	cmdProfile := shell.Command{
		Command: "gcloud",
		Args:    []string{"network-security", "security-profiles", "custom-mirroring", "delete", spName, "--organization=" + orgID, "--location=global", "--quiet"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmdProfile)
	if !assert.NoError(t, err, "Security Profile deletion failed.") {
		t.FailNow()
	}
}

// getAttachmentProjectNumber retrieves the project number for the attachment project.
// If TF_VAR_ATTACHMENT_PROJECT_ID is not set, it defaults to the primary project ID.
func GetAttachmentProjectNumber(t *testing.T, projectID string, attachmentProjectID string) (string, error) {
	// If attachmentProjectID is not set, use the primary projectID as fallback.
	if attachmentProjectID == "" {
		attachmentProjectID = projectID
		t.Logf("TF_VAR_ATTACHMENT_PROJECT_ID not set. Defaulting to primary project ID: %s", projectID)
		return GetProjectNumber(t, projectID) // Use the global projectID as the fallback
	}

	return GetProjectNumber(t, attachmentProjectID)
}

/*
GetEnv is a generic utility to read an environment variable or return a fallback.
*/
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

/*
EnsureAppEngineApplicationExists checks if an App Engine app exists and creates one
in the specified region if it does not.
*/
func EnsureAppEngineApplicationExists(t *testing.T, projectID string, region string) {
	const appCreatePropagationWait = 15 * time.Second // Wait after App Engine app creation

	t.Helper()
	t.Logf("Checking if App Engine application exists in project '%s'...", projectID)

	describeCmd := shell.Command{
		Command: "gcloud",
		Args: []string{
			"app", "describe",
			"--project=" + projectID,
		},
		Env: map[string]string{
			"CLOUDSDK_CORE_PROJECT": projectID,
		},
	}

	output, err := shell.RunCommandAndGetOutputE(t, describeCmd)

	if err == nil {
		t.Logf("App Engine application already exists in project '%s'. Description output (may include location):\n%s", projectID, output)
		return
	}
	// Check if error output indicates the app is missing
	if strings.Contains(output, "does not contain an App Engine application") ||
		(err != nil && strings.Contains(err.Error(), "does not contain an App Engine application")) {

		t.Logf("App Engine application does not exist in project '%s'. Attempting to create it in region '%s'...", projectID, region)

		createCmd := shell.Command{
			Command: "gcloud",
			Args: []string{
				"app", "create",
				"--region=" + region,
				"--project=" + projectID,
				"--quiet",
			},
		}
		createOutput, createErr := shell.RunCommandAndGetOutputE(t, createCmd)
		if createErr != nil {
			t.Logf("Failed to create App Engine application in project '%s', region '%s'. Output:\n%s, Error: %s", projectID, region, createOutput, createErr)
		}

		t.Logf("Successfully created App Engine application in project '%s', region '%s'. Output:\n%s", projectID, region, createOutput)

		t.Logf("Waiting %v for App Engine application creation to propagate...", appCreatePropagationWait)
		time.Sleep(appCreatePropagationWait)
	} else {
		// App exists but another error occurred
		if err != nil {
			t.Logf("Failed to describe App Engine application for an unexpected reason. Output:\n%s. Error:%s", output, err)
		}
	}
}

/*
CreateGcsBucket creates a GCS bucket using the 'gcloud storage' command.
*/
func CreateGcsBucket(t *testing.T, projectID string, bucketName string, location string) {
	t.Helper()
	t.Logf("Creating GCS bucket: gs://%s", bucketName)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"storage", "buckets", "create", fmt.Sprintf("gs://%s", bucketName), "--project=" + projectID, "--location=" + location, "--uniform-bucket-level-access"},
	}
	// Run and capture output only on error
	if _, err := shell.RunCommandAndGetOutputE(t, cmd); err != nil {
		t.Logf("Failed to create GCS bucket %s. Error: %s", bucketName, err)
	}
	t.Logf("GCS bucket gs://%s created.", bucketName)
}

/*
DeleteGcsBucket deletes a GCS bucket.
*/
func DeleteGcsBucket(t *testing.T, bucketName string) {
	t.Helper()
	t.Logf("Deleting GCS bucket: gs://%s", bucketName)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"storage", "buckets", "delete", fmt.Sprintf("gs://%s", bucketName), "--quiet"},
	}
	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		if !strings.Contains(output, "BucketNotFoundException: 404") && !strings.Contains(err.Error(), "NotFoundException") {
			t.Logf("Error deleting GCS bucket %s: %v. Output:\n%s", bucketName, err, output)
		} else {
			t.Logf("GCS bucket %s already deleted or not found.", bucketName)
		}
	} else {
		t.Logf("GCS bucket %s deleted.", bucketName)
	}
}

/*
DeleteGcsObjects deletes objects from a GCS path prefix.
*/
func DeleteGcsObjects(t *testing.T, bucketName string, objectPathPrefix string) {
	t.Helper()

	var gcsPath string
	if objectPathPrefix == "" {
		// If the prefix is empty, we want to delete everything at the root of the bucket.
		// The path must be gs://bucket-name/*
		gcsPath = fmt.Sprintf("gs://%s/*", bucketName)
	} else {
		// If a prefix is provided, trim any trailing slash and add the wildcard.
		// The path must be gs://bucket-name/prefix/*
		trimmedPrefix := strings.TrimSuffix(objectPathPrefix, "/")
		gcsPath = fmt.Sprintf("gs://%s/%s/*", bucketName, trimmedPrefix)
	}

	t.Logf("Deleting objects in GCS path: %s", gcsPath)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"storage", "rm", gcsPath, "--recursive"},
	}
	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	// Suppress errors if no objects were found to delete
	if err != nil && !strings.Contains(output, "One or more URLs matched no objects") && !strings.Contains(err.Error(), "One or more URLs matched no objects") {
		t.Logf("Note: Error deleting objects from %s (may be benign if already gone): %v. Output:\n%s", gcsPath, err, output)
	} else {
		t.Logf("Attempted deletion of objects in %s (any matching objects removed or none found).", gcsPath)
	}
}

/*
UploadGcsObjectFromString uploads content from a string to a GCS object by
creating a temporary local file.
*/
func UploadGcsObjectFromString(t *testing.T, projectID string, bucketName string, objectPath string, content string) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "gcs-upload-*.tmp")
	if err != nil {
		t.Logf("Failed to create temp file for GCS upload:Error %s", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up temp file

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Logf("Failed to write content to temp file: Error:%s", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Logf("Failed to close temp file: Error:%s", err)
	}

	gcsDest := fmt.Sprintf("gs://%s/%s", bucketName, objectPath)
	t.Logf("Uploading temp file %s to %s", tmpFile.Name(), gcsDest)
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"storage", "cp", tmpFile.Name(), gcsDest, "--project=" + projectID},
	}
	if _, err := shell.RunCommandAndGetOutputE(t, cmd); err != nil {
		t.Logf("Failed to upload object %s to bucket %s. Error:%s", objectPath, bucketName, err)
	}
	t.Logf("Uploaded object %s successfully.", objectPath)
}

/*
UploadGCSObjectFromFile uploads a local file to a GCS bucket.
This complements UploadGcsObjectFromString.
*/
func UploadGCSObjectFromFile(t *testing.T, projectID string, localFilePath, bucketName, objectName string) (string, error) {
	t.Helper()
	t.Logf("Uploading %s to gs://%s/%s", localFilePath, bucketName, objectName)
	objectUploadPath := fmt.Sprintf("gs://%s/%s", bucketName, objectName)

	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{"storage", "cp", localFilePath, objectUploadPath, "--project=" + projectID},
	}
	_, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		return "", fmt.Errorf("failed to upload %s to %s: %w", localFilePath, objectUploadPath, err)
	}

	// This is the HTTPS URL for App Engine source, gs:// paths are not valid here.
	httpsURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	t.Logf("File %s uploaded successfully. App Engine source URL: %s", filepath.Base(localFilePath), httpsURL)
	return httpsURL, nil
}

/*
DownloadFile downloads a file from a URL to a local destination directory.
*/
func DownloadFile(t *testing.T, url string, destDir string, fileName string) (string, error) {
	filePath := filepath.Join(destDir, fileName)
	t.Logf("Downloading %s to %s", url, filePath)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download %s: status %s", url, resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

/*
CreateZipArchive creates a zip file from a map of source files.
The map key is the path to the file on disk (relative to sourceDir),
and the map value is the path it should have inside the zip file.
*/
func CreateZipArchive(t *testing.T, sourceDir string, targetZipPath string, filesToZip map[string]string) error {
	t.Logf("Creating zip archive %s from contents of %s", targetZipPath, sourceDir)
	zipFile, err := os.Create(targetZipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file %s: %w", targetZipPath, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for sourcePath, pathInZip := range filesToZip {
		fullSourcePath := filepath.Join(sourceDir, sourcePath)
		t.Logf("Adding to zip: %s as %s", fullSourcePath, pathInZip)
		fileToZip, err := os.Open(fullSourcePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s for zipping: %w", fullSourcePath, err)
		}
		defer fileToZip.Close()

		info, err := fileToZip.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %w", fullSourcePath, err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create zip header for %s: %w", fullSourcePath, err)
		}
		header.Name = pathInZip // Use the desired path inside the zip
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip writer for %s: %w", pathInZip, err)
		}
		if _, err = io.Copy(writer, fileToZip); err != nil {
			return fmt.Errorf("failed to write file %s to zip: %w", pathInZip, err)
		}
	}
	t.Logf("Zip archive %s created successfully.", targetZipPath)
	return nil
}

/*
CreateFirewallRules creates the standard firewall rules required
*/
func CreateFirewallRules(t *testing.T, projectID string, networkName string, ruleSuffix string) bool {
	// healthCheckRange := "130.211.0.0/22,35.191.0.0/16"
	allowSourceRanges := "130.211.0.0/22,35.191.0.0/16,10.0.1.0/24"

	rulesToCreate := map[string]string{
		fmt.Sprintf("fw-allow-http-%s", ruleSuffix):  "tcp:80",
		fmt.Sprintf("fw-allow-https-%s", ruleSuffix): "tcp:443",
	}

	allSucceeded := true
	for ruleName, ruleProtoPort := range rulesToCreate {
		t.Logf("Creating firewall rule: %s for %s from source %s", ruleName, ruleProtoPort, allowSourceRanges)
		cmd := shell.Command{
			Command: "gcloud",
			Args: []string{
				"compute", "firewall-rules", "create", ruleName,
				"--project=" + projectID,
				"--network=" + networkName,
				"--direction=INGRESS",
				"--priority=1000",
				"--action=ALLOW",
				"--rules=" + ruleProtoPort,
				"--source-ranges=" + allowSourceRanges,
				"--format=json",
			},
		}
		_, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				t.Logf("Firewall rule %s already exists. Proceeding.", ruleName)
			} else {
				t.Errorf("Error creating firewall rule %s: %v", ruleName, err)
				allSucceeded = false
			}
		} else {
			t.Logf("Firewall rule %s created successfully.", ruleName)
		}
	}
	if !allSucceeded {
		t.Error("One or more firewall rules failed to create properly.")
	}
	return allSucceeded
}

/*
DeleteFirewallRules cleans up the standard firewall rules.
*/
func DeleteFirewallRules(t *testing.T, projectID string, ruleSuffix string) {
	t.Logf("--- Starting Firewall Rule Cleanup for instance suffix: %s ---", ruleSuffix)
	rulesToDelete := []string{
		fmt.Sprintf("fw-allow-http-%s", ruleSuffix),
		fmt.Sprintf("fw-allow-https-%s", ruleSuffix),
	}

	for _, ruleName := range rulesToDelete {
		t.Logf("Attempting to delete firewall rule: %s", ruleName)
		cmd := shell.Command{
			Command: "gcloud",
			Args: []string{
				"compute", "firewall-rules", "delete", ruleName,
				"--project=" + projectID,
				"--quiet",
			},
		}
		_, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err != nil {
			t.Logf("Note: Error deleting firewall rule %s (may be benign if already gone): %v", ruleName, err)
		} else {
			t.Logf("Firewall rule %s deleted successfully or did not exist.", ruleName)
		}
	}
	t.Logf("--- Firewall Rule Cleanup for instance suffix %s finished ---", ruleSuffix)
}

/*
RemoveTempDir is a simple helper to clean up a temporary directory.
*/
func RemoveTempDir(t *testing.T, dirPath string) {
	t.Logf("Cleaning up temporary source directory: %s", dirPath)
	if err := os.RemoveAll(dirPath); err != nil {
		t.Logf("WARN: Failed to remove temp source directory %s: %v", dirPath, err)
	}
}

/*
GetCurrentGcloudUser retrieves the email address of the currently authenticated
gcloud user.
*/
func GetCurrentGcloudUser(t *testing.T) string {
	cmd := shell.Command{Command: "gcloud", Args: []string{"auth", "list", "--filter=status:ACTIVE", "--format=value(account)"}}
	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	require.NoError(t, err, "Failed to get current gcloud user. Ensure gcloud is authenticated.")
	currentUser := strings.TrimSpace(output)
	require.NotEmpty(t, currentUser, "gcloud config get-value account returned empty string.")
	t.Logf("Current gcloud principal identified as: %s", currentUser)
	return currentUser
}

/*
CreateVPCPeering establishes a VPC peering connection between two networks.
*/
func CreateVPCPeering(t *testing.T, projectID, network, peerNetworkURI, peeringName string) {
	t.Logf("Creating peering '%s' from network '%s' to '%s'", peeringName, network, peerNetworkURI)
	cmd := shell.Command{Command: "gcloud", Args: []string{
		"compute", "networks", "peerings", "create", peeringName,
		"--network=" + network,
		"--peer-network=" + peerNetworkURI,
		"--project=" + projectID,
		"--export-custom-routes",
		"--import-custom-routes",
	}}
	shell.RunCommand(t, cmd)
}

/*
DeleteVPCPeering removes a VPC peering connection after a test completes.
*/
func DeleteVPCPeering(t *testing.T, projectID, network, peeringName string) {
	t.Logf("--- Deleting peering '%s' from network '%s' ---", peeringName, network)
	cmd := shell.Command{Command: "gcloud", Args: []string{
		"compute", "networks", "peerings", "delete", peeringName,
		"--network=" + network,
		"--project=" + projectID,
		"--quiet",
	}}
	shell.RunCommand(t, cmd)
}

/*
AddSecurityProfileRuleAndAssociatePolicy adds a rule to an organization-level
firewall policy to apply a security profile group, and then associates the
policy with a VPC.
*/
func AddSecurityProfileRuleAndAssociatePolicy(t *testing.T, orgID, policyName, vpcName, projectID, profileGroupName, srcIPRanges string) {
	profileGroupPath := fmt.Sprintf("organizations/%s/locations/global/securityProfileGroups/%s", orgID, profileGroupName)
	vpcPath := fmt.Sprintf("projects/%s/global/networks/%s", projectID, vpcName)

	t.Logf("Adding rule to policy '%s' to apply security profile group '%s'", policyName, profileGroupPath)
	shell.RunCommand(t, shell.Command{Command: "gcloud", Args: []string{"compute", "firewall-policies", "rules", "create", "1000", "--firewall-policy=" + policyName, "--organization=" + orgID, "--action=apply_security_profile_group", "--security-profile-group=" + profileGroupPath, "--src-ip-ranges=" + srcIPRanges, "--layer4-configs=all", "--enable-logging", "--description=test-rule"}})

	t.Logf("Associating policy '%s' with VPC '%s'", policyName, vpcPath)
	shell.RunCommand(t, shell.Command{Command: "gcloud", Args: []string{"compute", "firewall-policies", "associations", "create", "--firewall-policy=" + policyName, "--organization=" + orgID, fmt.Sprintf("--name=%s-association", policyName), "--replace-association-on-target"}})
}

/*
DeleteSecurityProfileRuleAndPolicyAssociation removes a firewall policy
association from a VPC and deletes the corresponding security profile group
rule from the policy.
*/
func DeleteSecurityProfileRuleAndPolicyAssociation(t *testing.T, orgID, policyName string) {
	if policyName == "" {
		return
	}
	t.Logf("--- Deleting Firewall Policy Association: %s-association ---", policyName)
	shell.RunCommand(t, shell.Command{Command: "gcloud", Args: []string{"compute", "firewall-policies", "associations", "delete", fmt.Sprintf("%s-association", policyName), "--firewall-policy=" + policyName, "--organization=" + orgID}})

	t.Logf("--- Deleting Firewall Policy Rule '1000' from policy '%s' ---", policyName)
	shell.RunCommand(t, shell.Command{Command: "gcloud", Args: []string{"compute", "firewall-policies", "rules", "delete", "1000", "--firewall-policy=" + policyName, "--organization=" + orgID}})
}

/*
CreateOrgFirewallPolicy creates a firewall policy at the organization level.
*/
func CreateOrgFirewallPolicy(t *testing.T, orgID, policyName string) {
	t.Logf("Creating Firewall Policy '%s' in Org '%s'", policyName, orgID)
	shell.RunCommand(t, shell.Command{Command: "gcloud", Args: []string{"compute", "firewall-policies", "create", "--short-name=" + policyName, "--organization=" + orgID, "--description=integ-test-policy"}})
}

/*
DeleteOrgFirewallPolicy deletes an organization-level firewall policy after a
test completes.
*/
func DeleteOrgFirewallPolicy(t *testing.T, orgID, policyName string) {
	if policyName == "" {
		return
	}
	t.Logf("--- Deleting Firewall Policy: %s ---", policyName)
	shell.RunCommand(t, shell.Command{Command: "gcloud", Args: []string{"compute", "firewall-policies", "delete", policyName, "--organization=" + orgID, "--quiet"}})
}

/*
GetOrgIDFromProject retrieves the parent Organization ID for a given Project ID.
*/
func GetOrgIDFromProject(t *testing.T, projectID string) string {
	t.Logf("Attempting to find Organization ID for project '%s'...", projectID)
	args := []string{"projects", "describe", projectID, "--format=value(parent.id)"}
	cmd := shell.Command{Command: "gcloud", Args: args}
	output, err := shell.RunCommandAndGetOutputE(t, cmd)
	require.NoError(t, err, "Failed to run gcloud command to get organization ID for project %s: %v", projectID, err)
	orgID := strings.TrimSpace(output)
	require.NotEmpty(t, orgID, "Organization ID was not found for project %s. Ensure the project is directly under an organization.", projectID)
	t.Logf("Found Organization ID: %s", orgID)
	return orgID
}

/*
GetRegionFromZone is a utility to extract the region from a zone name, for
example, "us-central1-a" becomes "us-central1".
*/
func GetRegionFromZone(t *testing.T, zone string) string {
	lastHyphen := strings.LastIndex(zone, "-")
	if lastHyphen == -1 {
		t.Fatalf("Invalid zone format: %s. Expected format like 'us-central1-a'", zone)
	}
	return zone[:lastHyphen]
}

/*
CreateBiDirectionalVPCPeering establishes a two-way VPC peering connection
between two networks.
*/
func CreateBiDirectionalVPCPeering(t *testing.T, projectID, networkA, networkB string) {
	networkAURI := fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkA)
	networkBURI := fmt.Sprintf("projects/%s/global/networks/%s", projectID, networkB)
	peeringAToB := fmt.Sprintf("peering-to-%s", networkB)
	peeringBToA := fmt.Sprintf("peering-to-%s", networkA)

	t.Logf("Creating peering '%s' from %s to %s", peeringAToB, networkA, networkB)
	cmdAToB := shell.Command{Command: "gcloud", Args: []string{
		"compute", "networks", "peerings", "create", peeringAToB,
		"--network=" + networkA,
		"--peer-network=" + networkBURI,
		"--project=" + projectID,
		"--export-custom-routes",
		"--import-custom-routes",
	}}
	_, err := shell.RunCommandAndGetOutputE(t, cmdAToB)
	require.NoError(t, err, "Failed to create peering from %s to %s", networkA, networkB)

	t.Logf("Creating peering '%s' from %s to %s", peeringBToA, networkB, networkA)
	cmdBToA := shell.Command{Command: "gcloud", Args: []string{
		"compute", "networks", "peerings", "create", peeringBToA,
		"--network=" + networkB,
		"--peer-network=" + networkAURI,
		"--project=" + projectID,
		"--export-custom-routes",
		"--import-custom-routes",
	}}
	_, err = shell.RunCommandAndGetOutputE(t, cmdBToA)
	require.NoError(t, err, "Failed to create peering from %s to %s", networkB, networkA)
}

/*
DeleteBiDirectionalVPCPeering removes a two-way VPC peering connection
between two networks.
*/
func DeleteBiDirectionalVPCPeering(t *testing.T, projectID, networkA, networkB string) {
	peeringAToB := fmt.Sprintf("peering-to-%s", networkB)
	peeringBToA := fmt.Sprintf("peering-to-%s", networkA)

	t.Logf("--- Deleting peering '%s' from network '%s' ---", peeringAToB, networkA)
	cmdAToB := shell.Command{Command: "gcloud", Args: []string{
		"compute", "networks", "peerings", "delete", peeringAToB,
		"--network=" + networkA,
		"--project=" + projectID,
		"--quiet",
	}}
	shell.RunCommand(t, cmdAToB)

	t.Logf("--- Deleting peering '%s' from network '%s' ---", peeringBToA, networkB)
	cmdBToA := shell.Command{Command: "gcloud", Args: []string{
		"compute", "networks", "peerings", "delete", peeringBToA,
		"--network=" + networkB,
		"--project=" + projectID,
		"--quiet",
	}}
	shell.RunCommand(t, cmdBToA)
}

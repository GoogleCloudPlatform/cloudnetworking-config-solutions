# Execution

This directory contains the scripts and Terraform configurations to deploy the CloudNet Config Solutions infrastructure on Google Cloud.

## Overview

The deployment is divided into seven logically isolated stages, each handled by a separate subfolder:

1. **00-bootstrap (Optional):**
   - **Bootstrap stage:** Organization or folder admins with elevated permissions are responsible for running the bootstrap stage once (or as needed when IAM permissions change) to establish the prerequisite resources for the subsequent stages.
   - This stage is optional but recommended.
   - Creates essential GCP resources like service accounts, IAM permissions, and a GCS bucket for Terraform state files.
   - If you skip this stage, ensure you have the required IAM permissions to execute the other stages.
   - This stage creates service accounts that can be impersonated in the provider.tf file for subsequent stages.
   - After successful execution of the bootstrap stage, the output from this stage would look like

    **Bootstrap Output (Example):**
    ```
    consumer_stage_email          = "serviceAccount:consumer-stage-sa@<project-id>.iam.gserviceaccount.com"
    networking_manual_stage_email = "serviceAccount:networking-manual-stage-sa@<project-id>.iam.gserviceaccount.com"
    networking_stage_email        = "serviceAccount:networking-stage-sa@<project-id>.iam.gserviceaccount.com"
    organization_stage_email      = "serviceAccount:organization-stage-sa@<project-id>.iam.gserviceaccount.com"
    producer_stage_email          = "serviceAccount:producer-stage-sa@<project-id>.iam.gserviceaccount.com"
    security_stage_email          = "serviceAccount:security-stage-sa@<project-id>.iam.gserviceaccount.com"
    storage_bucket_name           = "terraform-state"
    ```

    **Generating the provider.tf file**

    1. **Export Variables:**

          ```
          export TF_SERVICE_ACCOUNT="<ENTER THE SERVICE ACCOUNT HERE>"

          export TF_BUCKET_NAME="<ENTER THE GCS BUCKET NAME HERE>"

          export TF_ORGANIZATION_PREFIX="<ENTER THE GCS PREFIX NAME HERE>"
          ```

    2. **Populate Templates:**

        Use `sed` or a similar tool to replace placeholders in your `provider.tf.template` file with the environment variables:
          ```
          sed -e "s|ENTER_TF_SERVICE_ACCOUNT|$TF_SERVICE_ACCOUNT|" \
              -e "s|ENTER_TF_BUCKET_NAME|$TF_BUCKET_NAME|" \
              -e "s|ENTER_TF_ORGANIZATION_PREFIX|$TF_ORGANIZATION_PREFIX|" \
              <stage-folder-name>/provider.tf.template > <stage-folder-name>/provider.tf
          ```

      **Example (Organization Stage):**

      Here is a more specific example  describing in more detail on how to generate these provider.tf file using the output received from the 00-bootstrap stage.
      ```
      export TF_SERVICE_ACCOUNT="organization-stage-sa@<project-id>.iam.gserviceaccount.com"

      export TF_BUCKET_NAME="terraform-state"

      export TF_ORGANIZATION_PREFIX="organization"

      sed \
      -e "s|ENTER_TF_SERVICE_ACCOUNT|$TF_SERVICE_ACCOUNT|" \
      -e "s|ENTER_TF_BUCKET_NAME|$TF_BUCKET_NAME|" \
      -e "s|ENTER_TF_ORGANIZATION_PREFIX|$TF_ORGANIZATION_PREFIX|" \
      01-organization/provider.tf.template > 01-organization/provider.tf

      ```

2. **01-organization:**
   - Manages Google Cloud Project APIs and services within your GCP projects by selectively enabling or disabling them as needed.

3. **02-networking:**
   - Handles network administration tasks.
   - Configures network connectivity (Shared VPC, Private Service Access, Service Connectivity Automation).

4. **03-security:**
   - Defines and manages firewall rules.
   - Secures communication between producer and consumer services.

5. **04-producer:**
   - Deploys GCP-managed producer services.
   - Includes AlloyDB, Cloud SQL, Memorystore Redis clusters, GKE, Vertex AI Vector Search, Vertex AI Online Prediction Endpoint.

6. **05-networking-manual:**
   - This stage establishes Private Service Connect (PSC) for secure, private communication between your consumer project and the producer services created in the "04-producer" stage.
   - **Internal IP Addresses:** Reserved within specific subnets in your consumer project. These addresses act as private endpoints for the PSC connection, providing a secure and stable way to access your producer services.
   - **Forwarding Rules:** Created to direct traffic destined for the reserved internal IP addresses to your producer services through the PSC connection. This ensures seamless communication without exposing your services to the public internet.
   - **Note :** This stage currently supports establishing PSC for Cloud SQL instances only.

7. **06-consumer:**
   - Deploys GCP-managed consumer services.
   - Includes GCE instances and Cloud Run.

## Logical Isolation and Permissions

Each stage is designed to be managed independently by different teams of administrators. This provides granular control and security.  The specific IAM permissions required for each stage are detailed in the `README.md` files within their respective folders.

If you skip the `00-bootstrap` stage, ensure the user executing the specific stage has the required IAM permissions to create the resources defined in that stage.

## Additional Components

- **modules:** Contains custom Terraform modules for reusable infrastructure patterns.
- **run.sh:** A bash script that automates the entire deployment process. Use flags to control which stages are executed (e.g., `--all`, `--networking`).
- **configuration:** Holds the `*.tfvars` files for each stage, providing flexibility for configuration management.
- **provider.tf.template:**  A template file used by Terraform to connect to GCP. Update this file with the appropriate service account details.

## Getting Started

1. **Prerequisites:**
   - Ensure you have Terraform installed and configured.
   - Set up authentication to your GCP project with appropriate permissions.

2. **Configuration:**
   - Review and customize the Terraform variables in the `.tfvars` files located in the `configuration` folder.
   - Each stage has its own `.tfvars` file (e.g., `organization.tfvars`).
   - Update `provider.tf.template` with the necessary service account details if you choose to run the `00-bootstrap` stage.

3. **Deployment:**
   - **Option 1:** Execute individual stages manually:
      1. Navigate to the stage's subfolder.
      2. Run `terraform init` to initialize the environment.
      3. Run `terraform plan -var-file="../../configuration/<stagename>.tfvars"` to preview the changes.
      4. Run `terraform apply -var-file="../../configuration/<stagename>.tfvars"` to apply the configuration.
   - **Option 2:** Use the `run.sh` script for a streamlined deployment:
      - **Example:** To initialize the "networking" stage, run:
        ```bash
        ./run.sh --stage networking --tfcommand init
        ```
      - For complete list of options and usage details, run:
        ```bash
        ./run.sh -h
        ```
        or
        ```bash
        ./run.sh --help
        ```

## Important Notes:

- Refer to the `README.md` files in each stage subfolder for detailed instructions and information specific to that stage's deployment.
- Carefully verify the terraform plan to ensure the desired changes in infrastructure matches the expectations.
- Review and understand the security implications before deploying any of the infrastructure components.

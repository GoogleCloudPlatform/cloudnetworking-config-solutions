## Introduction

The bootstrap stage is the first and most crucial step in setting up your Google Cloud infrastructure using Terraform. It lays the groundwork for subsequent stages (01-organization, 02-networking, 03-security, 04-producer, 05-networking-manual, 06-consumer) by provisioning essential resources and establishing security best practices. This stage focuses on creating the following:
  * **Impersonating Service Accounts:** This stage generates service accounts for each subsequent stage, allowing them to impersonate roles with the necessary permissions for their respective tasks. This approach enhances security by granting only the required privileges to each stage.
  * **Terraform State Bucket:** A Google Cloud Storage bucket is created to store the Terraform state files. This centralizes state management, making it easier to track changes and collaborate on infrastructure updates.

## Pre-Requisites

* IAM Permissions: The user or service account executing Terraform must have the following IAM roles or equivalent permissions:
  * `roles/iam.serviceAccountAdmin` : To create and manage service accounts for the project for which service account needs to be created.
  * `roles/resourcemanager.projectIamAdmin` : Provides permissions to administer allow policies on projects.
  * `roles/storage.admin` : To create and manage Google Cloud Storage buckets.

## Execution Steps:

1. Create `terraform.tfvars`:
    * Make a copy of the provided terraform.tfvars.example file and rename it to terraform.tfvars.
    * Fill in the values for [input variable](#inputs) and other required variables according to your specific requirements.
2. Initialize Terraform:
    `terraform init`
3. Review Execution Plan:
    `terraform plan`
4. Apply Configuration:
    `terraform apply -var-file="../../configuration/bootstrap.tfvars`

### Example

To help you get started, we've provided examples of tfvars files that you can use :

* **Minimal tfvars (Mandatory Fields Only):**
This minimal example includes only the essential fields required to execute the bootstrap stage.
  ```
  bootstrap_project_id                  = "<your-bootstrap-project-id>"
  network_hostproject_id                = "<your-host-project-id>"
  network_serviceproject_id             = "<your-service(producer/consumer)-project-id>"
  organization_stage_administrator      = ["user:user-example@example.com"]
  networking_stage_administrator        = ["user:user-example@example.com"]
  security_stage_administrator          = ["user:user-example@example.com"]
  producer_stage_administrator          = ["user:user-example@example.com"]
  networking_manual_stage_administrator = ["user:user-example@example.com"]
  consumer_stage_administrator          = ["user:user-example@example.com"]
  ```

## Important Considerations:

  * **Security**: Pay close attention to the permissions granted to the service accounts. Follow the principle of least privilege to minimize security risks.
  * **State Management:** The Terraform state bucket is critical for maintaining the state of your infrastructure. Ensure its security and accessibility.
  * **Dependencies:** This bootstrap stage is a prerequisite for all subsequent stages. Make sure it is executed successfully before proceeding with other stages.
  **Note:** You can skip the bootstrap stage if you choose, but you must ensure the following:
  * **Permissions:** The user or service account executing Terraform for each individual stage (01-organization, 02-networking, etc.) must have the necessary IAM permissions outlined in the respective stage's README file.
 * **State File Management:** You are responsible for setting up and maintaining a secure location for Terraform state files for each stage. This could involve using a Google Cloud Storage bucket, a local backend, or another suitable storage mechanism.

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_consumer"></a> [consumer](#module\_consumer) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account | v31.1.0 |
| <a name="module_google_storage_bucket"></a> [google\_storage\_bucket](#module\_google\_storage\_bucket) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/gcs | n/a |
| <a name="module_networking"></a> [networking](#module\_networking) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account | v31.1.0 |
| <a name="module_networking_manual"></a> [networking\_manual](#module\_networking\_manual) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account | v31.1.0 |
| <a name="module_organization"></a> [organization](#module\_organization) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account | v31.1.0 |
| <a name="module_producer"></a> [producer](#module\_producer) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account | v31.1.0 |
| <a name="module_security"></a> [security](#module\_security) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account | v31.1.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_bootstrap_project_id"></a> [bootstrap\_project\_id](#input\_bootstrap\_project\_id) | Google Cloud Project ID which will be used to create the service account and Google Cloud storage buckets. | `string` | n/a | yes |
| <a name="input_consumer_stage_administrator"></a> [consumer\_stage\_administrator](#input\_consumer\_stage\_administrator) | List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com) | `list(string)` | n/a | yes |
| <a name="input_network_hostproject_id"></a> [network\_hostproject\_id](#input\_network\_hostproject\_id) | Google Cloud Project ID for the networking host project to be used to create networking and security resources. | `string` | n/a | yes |
| <a name="input_network_serviceproject_id"></a> [network\_serviceproject\_id](#input\_network\_serviceproject\_id) | Google Cloud Project ID to be used to create Google Cloud resources like consumer and producer services. | `string` | n/a | yes |
| <a name="input_networking_manual_stage_administrator"></a> [networking\_manual\_stage\_administrator](#input\_networking\_manual\_stage\_administrator) | List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com) | `list(string)` | n/a | yes |
| <a name="input_networking_stage_administrator"></a> [networking\_stage\_administrator](#input\_networking\_stage\_administrator) | List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com) | `list(string)` | n/a | yes |
| <a name="input_organization_stage_administrator"></a> [organization\_stage\_administrator](#input\_organization\_stage\_administrator) | List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com) | `list(string)` | n/a | yes |
| <a name="input_producer_stage_administrator"></a> [producer\_stage\_administrator](#input\_producer\_stage\_administrator) | List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com) | `list(string)` | n/a | yes |
| <a name="input_security_stage_administrator"></a> [security\_stage\_administrator](#input\_security\_stage\_administrator) | List of Members to be granted an IAM role. e.g. (group:my-group@example.com),(user:my-user@example.com) | `list(string)` | n/a | yes |
| <a name="input_consumer_sa_name"></a> [consumer\_sa\_name](#input\_consumer\_sa\_name) | Name of the service account to create for consumer stage. | `string` | `"consumer-stage-sa"` | no |
| <a name="input_gcs_bucket_location"></a> [gcs\_bucket\_location](#input\_gcs\_bucket\_location) | Location of the Google Cloud storage bucket. | `string` | `"EU"` | no |
| <a name="input_gcs_bucket_name"></a> [gcs\_bucket\_name](#input\_gcs\_bucket\_name) | Name of the Google Cloud storage bucket. | `string` | `"terraform-state"` | no |
| <a name="input_networking_manual_sa_name"></a> [networking\_manual\_sa\_name](#input\_networking\_manual\_sa\_name) | Name of the service account to create for networking manual stage. | `string` | `"networking-manual-stage-sa"` | no |
| <a name="input_networking_sa_name"></a> [networking\_sa\_name](#input\_networking\_sa\_name) | Name of the service account to create for networking stage. | `string` | `"networking-stage-sa"` | no |
| <a name="input_organization_sa_name"></a> [organization\_sa\_name](#input\_organization\_sa\_name) | Name of the service account to create for organization stage. | `string` | `"organization-stage-sa"` | no |
| <a name="input_producer_sa_name"></a> [producer\_sa\_name](#input\_producer\_sa\_name) | Name of the service account to create for producer stage. | `string` | `"producer-stage-sa"` | no |
| <a name="input_security_sa_name"></a> [security\_sa\_name](#input\_security\_sa\_name) | Name of the service account to create for security stage. | `string` | `"security-stage-sa"` | no |
| <a name="input_versioning"></a> [versioning](#input\_versioning) | The Goocle Cloud storage bucket versioning. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_consumer_stage_email"></a> [consumer\_stage\_email](#output\_consumer\_stage\_email) | Consumer stage service account iam email. |
| <a name="output_networking_manual_stage_email"></a> [networking\_manual\_stage\_email](#output\_networking\_manual\_stage\_email) | Networking manual stage service account iam email. |
| <a name="output_networking_stage_email"></a> [networking\_stage\_email](#output\_networking\_stage\_email) | Networking stage service account iam email. |
| <a name="output_organization_stage_email"></a> [organization\_stage\_email](#output\_organization\_stage\_email) | Organization stage service account iam email. |
| <a name="output_producer_stage_email"></a> [producer\_stage\_email](#output\_producer\_stage\_email) | Producer stage service account iam email. |
| <a name="output_security_stage_email"></a> [security\_stage\_email](#output\_security\_stage\_email) | Security stage service account iam email. |
| <a name="output_storage_bucket_name"></a> [storage\_bucket\_name](#output\_storage\_bucket\_name) | Google Cloud storage bucket name. |
<!-- END_TF_DOCS -->

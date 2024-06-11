## Introduction

This stage is a foundational step in your Terraform-based infrastructure setup. It focuses on enabling and disabling Google Cloud APIs for specific projects within your organization. By proactively managing APIs at the project level, you ensure that only necessary services are activated and reducing potential security risks. This stage also sets the groundwork for subsequent stages by ensuring the required APIs are enabled for use.

## Prerequisites

Before moving to next stages, ensure you have the completed the following pre-requsites:

1. **IAM Permissions:** Grant yourself (or the appropriate users/service accounts) the following IAM roles at the project level (or higher):
    * Service Usage Administrator (`roles/serviceusage.serviceUsageAdmin`)

## Let's Get Started! 🚀
With the prerequisites in place, you can now leverage Terraform to automate the management of your Google Cloud resources in this stage. Here's the workflow:

### Execution Steps

1. **Create your .tfvars files:**
    * Refer the provided `terraform.tfvars.example` file or the [example section](#example) and create a new `terraform.tfvars`.
    * Fill in the values for project_id, activate_apis and other [input variables](#inputs)  according to your specific requirements. Ensure you are using the correct project_id for the project you wish to enable or disable APIs for.

2. **Initialize Terraform:**

    * Open your terminal and navigate to the `organization` directory containing the Terraform configuration.

    * Run the following command to initialize Terraform:

    ```
    terraform init
    ```
3. **Review the Execution Plan:**

    * Use the terraform plan command to generate an execution plan. This will show you the changes Terraform will make to your Google Cloud infrastructure:

    ```
    terraform plan
    ```

Carefully review the plan to ensure it aligns with your intended configuration.

4. **Apply the Configuration:**

    Once you're satisfied with the plan, execute the terraform apply command to enabling/disabling APIs as defined in your `terraform.tfvars` file for the specified project :

    ```
    terraform apply
    ```


### Example

To get you started, we've included sample `terraform.tfvars` configuration files that you can adapt for your organization's specific needs. Remember, the **.tfvars** file requires a map of objects, each with a unique key string. If you use the same key string more than once, the last value assigned to that key will overwrite any previous values.

  ```
  activate_api_identities = {
    "project-01" = {
      project_id = "your-project-id",
      activate_apis = [
        "servicenetworking.googleapis.com",
        "alloydb.googleapis.com",
        "sqladmin.googleapis.com",
        "iam.googleapis.com",
        "compute.googleapis.com",
        "redis.googleapis.com",
      ],
    },
  }
  ```

## Important Considerations:

* Impact: Carefully consider the impact of enabling or disabling APIs, as it will affect the services available within the specified project.

* Dependencies: Ensure that you have enabled any APIs required by subsequent stages (networking, security, etc.) in this stage.

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_activate_project_apis"></a> [activate\_project\_apis](#module\_activate\_project\_apis) | terraform-google-modules/project-factory/google//modules/project_services | 15.0.1 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_activate_api_identities"></a> [activate\_api\_identities](#input\_activate\_api\_identities) | Map of objects containing information required to enable API's for the Google Cloud project. | <pre>map(object({<br>    project_id                  = string,<br>    activate_apis               = list(string),<br>    disable_dependent_services  = optional(bool, false)<br>    disable_services_on_destroy = optional(bool, false)<br>  }))</pre> | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_activated_api_identities"></a> [activated\_api\_identities](#output\_activated\_api\_identities) | Map of objects containing project ID, enabled apis and api\_identities. |
<!-- END_TF_DOCS -->

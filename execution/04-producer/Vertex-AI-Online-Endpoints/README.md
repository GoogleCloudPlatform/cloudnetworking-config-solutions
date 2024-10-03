## Introduction

This folder aims to streamline the deployment of Vertex AI Endpoints for Online Predictions on Google Cloud Platform. It leverages a modular approach, allowing you to define multiple endpoints with distinct configurations through YAML files. This folder utilizes the  `vertex-ai-online-endpoints` module, simplifying the process of creating and managing your Vertex AI infrastructure for Online predictions.

## Pre-requisites

Before you get started, ensure you have the following:

1. **Google Cloud Project:** A GCP project with billing enabled to house your Vertex AI resources.
2. **Terraform Installed:** The Terraform CLI installed on your system. Download the appropriate version from: [https://www.terraform.io/downloads.html](https://www.terraform.io/downloads.html)
3. **Service Account:**  A service account with permissions to manage Vertex AI resources within your GCP project. Ensure it has the required roles, such as `roles/aiplatform.admin`.
4. **Networking Configuration:** If you're enabling VPC peering (`enable_service_networking_to_vpc`), configure a global address range for the peering connection. 
5. **gcloud CLI:**  The Google Cloud SDK (`gcloud`) command-line tool installed. You'll use this to authenticate and manage your GCP resources. Download and install from here: [https://cloud.google.com/sdk/docs/install](https://cloud.google.com/sdk/docs/install)


## Let's Get Started! 🚀

This project uses YAML files to manage your endpoint configurations, enhancing organization and scalability. Follow these steps to deploy:

### Execution Steps

1. **Define Endpoint Configurations:**
   * Navigate to the `configuration/producer/Vertext-AI-Online-Endpoints/config` directory.
   * Create YAML files (e.g., `endpoint.yaml`) containing the settings for your Vertex AI Endpoints.
   * Use the [Example](#example) section as a template or refer to the provided example file (`_endpoint.yaml.example`).
   * Customize these files with your project ID, endpoint names, descriptions, locations, networking details, and any other desired parameters. 

2. **Initialize Terraform:**
   * Open your terminal and navigate to the root directory of this Terraform project.
   * Run `terraform init` to initialize Terraform and download the necessary providers and modules.

3. **Preview Changes with Execution Plan:**
   * Execute `terraform plan` to generate a detailed execution plan. This will display the proposed changes to your infrastructure based on your YAML configurations. 

    ```
    terraform plan -var-file=../../../configuration/producer/Vertex-AI-Online-Endpoints/vertex-ai-online-endpoints.tfvars
    ```

   * Review the output thoroughly to ensure it matches your intended endpoint setup.

4. **Apply and Deploy:**

   * Once satisfied with the plan, run `terraform apply` to deploy your Vertex AI Endpoints. Terraform will create and configure the endpoints as defined in your YAML files. 

    ```
    terraform apply -var-file=../../../configuration/producer/Vertex-AI-Online-Endpoints/vertex-ai-online-endpoints.tfvars
    ```

5. **Manage and Monitor:**
   * After deployment, you can manage and monitor your Vertex AI Endpoints through:
       * **Google Cloud Console:**  Provides a user-friendly interface to interact with your endpoints.
       * **`gcloud` CLI:** Offers command-line tools to manage your Vertex AI resources.

### Example

Below is a sample configuration file (`endpoint1.yaml`) showcasing some of the key settings:

```yaml
display_name: my-first-endpoint
project: your-gcp-project-id
name: endpoint1  
description: "Endpoint for image classification model"
location: us-central1
region: us-central1 
labels:
  environment: production
  team: ml-engineering 
network: projects/your-gcp-project-number/global/networks/your-vpc-network
```

**NOTE** : Network mentioned here should have a service networking peering connection.

## Important Notes: 

* **Customization is Key:**  Tailor the provided example YAML files or create new ones to align with the specific requirements of your Vertex AI Endpoints.
* **Global Address Range:**  If `enable_service_networking_to_vpc` and `create_range` are set to `true`, ensure the `range_name` is globally unique.
* **Security Best Practices:** For production environments, carefully consider and implement appropriate security measures, including IAM roles and network security configurations. 
* **Documentation:**  For comprehensive information, consult the official Google Cloud documentation on Vertex AI Endpoints:  [https://cloud.google.com/vertex-ai/docs/endpoints/](https://cloud.google.com/vertex-ai/docs/endpoints/) 
* **Resource Cleanup:** To remove the endpoints and related resources created by this Terraform project, execute `terraform destroy`. 

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_vertex_endpoints"></a> [vertex\_endpoints](#module\_vertex\_endpoints) | ../../../modules/vertex-ai-online-endpoints/ | n/a |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding Online Endpoints configuration values. | `string` | `"../../../configuration/producer/Vertex-AI-Online-Endpoints/config"` | no |
| <a name="input_description"></a> [description](#input\_description) | The description of the Vertex AI endpoint. | `string` | `"Sample CNCS vertex AI endpoint deployment"` | no |
| <a name="input_display_name"></a> [display\_name](#input\_display\_name) | The display name of the Vertex AI endpoint. | `string` | `"cncs-vertex-ai-display-name"` | no |
| <a name="input_labels"></a> [labels](#input\_labels) | The labels to associate with the Vertex AI endpoint. | `map(string)` | `{}` | no |
| <a name="input_location"></a> [location](#input\_location) | The location of the Vertex AI endpoint. | `string` | `"us-central1"` | no |
| <a name="input_name"></a> [name](#input\_name) | The name of the Vertex AI endpoint. | `string` | `"cncs-vertex-ai-endpoint-name"` | no |
| <a name="input_region"></a> [region](#input\_region) | The region of the Vertex AI endpoint. | `string` | `"us-central1"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_endpoint_configurations"></a> [endpoint\_configurations](#output\_endpoint\_configurations) | Configuration details for all created Vertex AI endpoints. |
| <a name="output_endpoint_configurations_from_yaml"></a> [endpoint\_configurations\_from\_yaml](#output\_endpoint\_configurations\_from\_yaml) | Endpoint configurations read from YAML files. |
<!-- END_TF_DOCS -->

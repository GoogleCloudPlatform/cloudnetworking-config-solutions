# Vector Search - Vertex AI

## Introduction

This Terraform module simplifies the creation and management of vector search instances by automating the creation of indexes and index endpoints, and the deployment of those indexes to their respective endpoints. Vector search can search from billions of semantically similar or semantically related items. With vector search you can leverage the same infrastructure that provides a foundation for Google products such as Google Search, YouTube, and Play.

## Pre-Requisites

Before creating your first vector search instance, ensure you have completed the following prerequisites:

1. **Completed Prior Stages:** Successful deployment of vector search resources depends on the completion of the following stages:
    * **01-organization:** This stage handles the activation of required Google Cloud APIs for vector search.
    * **02-networking:** This stage sets up the necessary network infrastructure, such as VPCs and subnets, to support vector search connectivity.

2. **API Enablement**: Ensure the following Google Cloud APIs have been enabled:

    * Compute Engine API (compute.googleapis.com)
    * AI Platform API (aiplatform.googleapis.com)
    * Cloud Resource Manager API (cloudresourcemanager.googleapis.com)

3. **IAM Permissions**:  Grant yourself (or the appropriate users/service accounts) the following IAM roles at the project level (or higher):

    * AI Platform Admin (roles/aiplatform.admin)

## Let's Get Started! 🚀

With the prerequisites in place and your vector search configuration files ready, you can now leverage Terraform to automate the creation of your vector search instances. Here's the workflow:

### Execution Steps

1. Create your configuration files:

    Create YAML files defining the properties of each vector search instance you want to create. Ensure these files are stored in the `configuration/producer/VectorSearch/config` folder within this vector search folder.

    Each YAML file should map to a single vector search instance, providing details such as  name, project ID, region etc. Each field and its structure are described in the input section below.

    For reference on how to structure your vector search configuration YAML files, see the example section below.

    **Important Note :** Due to limitations in the Terraform Google provider [see issue #19490](https://github.com/hashicorp/terraform-provider-google/issues/19490), creating a deployed index endpoint requires declaring the `region` in your `providers.tf` file. Please ensure your `providers.tf` file includes the `region` field, like so:

    ```
    provider "google" {
      region = "your-google-cloud-region" //e.g. us-central1
    }

    provider "google-beta" {
      region = "your-google-cloud-region" //e.g. us-central1
    }
    ```

2. Initialize Terraform:

    Open your terminal and navigate to the directory containing the Terraform configuration.

    Run the following command to initialize Terraform:

    ```
    terraform init
    ```

3. Review the Execution Plan:

    Use the terraform plan command to generate an execution plan. This will show you the changes Terraform will make to your Google Cloud infrastructure:

    ```
    terraform plan -var-file=../../../configuration/producer/VectorSearch/vectorsearch.tfvars
    ```

4. Apply the Configuration:

    Once you're satisfied with the plan, execute the terraform apply command to provision your vector search instances:

    ```
    terraform apply -var-file=../../../configuration/producer/VectorSearch/vectorsearch.tfvars
    ```

5. Monitor and Manage:

    After the instances are created, you can monitor their status, performance, and logs through the Google Cloud Console or using the Google Cloud CLI. Use Terraform to manage updates and changes to your vector search instances as needed.

## Example YAML Configuration

  ```
  project_id: <test-project>
  region: us-central1
  index_display_name : demo-index-1
  index_update_method : BATCH_UPDATE
  dimension: 2
  approximate_neighbors_count: 150
  shard_size: SHARD_SIZE_SMALL
  distance_measure_type: DOT_PRODUCT_DISTANCE
  index_endpoint_display_name : demo-index-endpoint-1
  index_endpoint_network : projects/<project-number>/global/networks/<network-name>
  tree_ah_config:
    leaf_node_embedding_count: 500
    leaf_nodes_to_search_percent: 7
  brute_force_config: null
  deployed_index_id: deploy_index_id_1
  ```

## Important Notes

- Refer to the official Google Cloud Vertex AI vector search [documentation](https://cloud.google.com/vertex-ai/docs/vector-search/overview) for the most up-to-date information and best practices.
- Order of Execution: Make sure to complete the necessary networking stages before attempting to create vector search instances. Terraform will leverage the resources and configurations established in these prior stages.
- Troubleshooting: If you encounter errors during the vector search instance creation process, verify that all prerequisites are satisfied and that the dependencies between stages are correctly configured.

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_vector_search"></a> [vector\_search](#module\_vector\_search) | ../../../modules/vector-search | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_approximate_neighbors_count"></a> [approximate\_neighbors\_count](#input\_approximate\_neighbors\_count) | The default number of neighbors to find via approximate search before exact reordering is performed. Exact reordering is a procedure where results returned by an approximate search algorithm are reordered via a more expensive distance computation. Required if tree-AH algorithm is used. | `number` | `150` | no |
| <a name="input_automatic_resources"></a> [automatic\_resources](#input\_automatic\_resources) | The minimum number and maximum number of replicas this DeployedModel will be always deployed on. | <pre>object({<br>    min_replica_count = optional(number)<br>    max_replica_count = optional(number)<br>  })</pre> | `null` | no |
| <a name="input_brute_force_config"></a> [brute\_force\_config](#input\_brute\_force\_config) | Configuration options for using brute force search, which simply implements the standard linear search in the database for each query. | `string` | `null` | no |
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding Vector Search configuration values. | `string` | `"../../../configuration/producer/VectorSearch/config"` | no |
| <a name="input_contents_delta_uri"></a> [contents\_delta\_uri](#input\_contents\_delta\_uri) | Allows inserting, updating or deleting the contents of the Matching Engine Index. The string must be a valid Cloud Storage directory path. If this field is set when calling IndexService.UpdateIndex, then no other Index field can be also updated as part of the same call. | `string` | `null` | no |
| <a name="input_dedicated_resources"></a> [dedicated\_resources](#input\_dedicated\_resources) | The type of the machine, minimum number and maximum number of replicas this DeployedModel will be always deployed on. | <pre>object({<br>    machine_spec = object({<br>      machine_type = string<br>    })<br>    min_replica_count = optional(number)<br>    max_replica_count = optional(number)<br>  })</pre> | `null` | no |
| <a name="input_deployed_display_name"></a> [deployed\_display\_name](#input\_deployed\_display\_name) | The display name of the deployment. | `string` | `null` | no |
| <a name="input_deployed_index_auth_config"></a> [deployed\_index\_auth\_config](#input\_deployed\_index\_auth\_config) | The authentication provider that the DeployedIndex uses.A list of allowed JWT issuers. Each entry must be a valid Google service account, in the following format: service-account-name@project-id.iam.gserviceaccount.com | <pre>object({<br>    auth_provider = object({<br>      audiences       = optional(string)<br>      allowed_issuers = optional(list(string))<br>    })<br>  })</pre> | `null` | no |
| <a name="input_deployment_group"></a> [deployment\_group](#input\_deployment\_group) | The deployment group can be no longer than 64 characters (eg: 'test', 'prod'). If not set, we will use the 'default' deployment group. | `string` | `null` | no |
| <a name="input_dimensions"></a> [dimensions](#input\_dimensions) | The number of dimensions of the input vectors. | `number` | `2` | no |
| <a name="input_distance_measure_type"></a> [distance\_measure\_type](#input\_distance\_measure\_type) | The distance measure used in nearest neighbor search. | `string` | `"DOT_PRODUCT_DISTANCE"` | no |
| <a name="input_enable_access_logging"></a> [enable\_access\_logging](#input\_enable\_access\_logging) | If true, private endpoint's access logs are sent to Cloud Logging. | `bool` | `true` | no |
| <a name="input_index_description"></a> [index\_description](#input\_index\_description) | The description of the Index. | `string` | `null` | no |
| <a name="input_index_endpoint_description"></a> [index\_endpoint\_description](#input\_index\_endpoint\_description) | The description of the Index Endpoint. | `string` | `null` | no |
| <a name="input_index_endpoint_labels"></a> [index\_endpoint\_labels](#input\_index\_endpoint\_labels) | Labels to be attached to index endpoint instances. | `map(string)` | `null` | no |
| <a name="input_index_endpoint_network"></a> [index\_endpoint\_network](#input\_index\_endpoint\_network) | The full name of the Google Compute Engine network to which the index endpoint should be peered. Private services access must already be configured for the network. If left unspecified, the index endpoint is not peered with any network. Format: projects/{project}/global/networks/{network}. Where {project} is a project number, as in 12345, and {network} is network name. | `string` | `null` | no |
| <a name="input_index_labels"></a> [index\_labels](#input\_index\_labels) | Labels to be attached to index instances. | `map(string)` | `null` | no |
| <a name="input_index_update_method"></a> [index\_update\_method](#input\_index\_update\_method) | The update method to use with this Index. The value must be the followings. If not set, BATCH\_UPDATE will be used by default. | `string` | `"BATCH_UPDATE"` | no |
| <a name="input_private_service_connect_config"></a> [private\_service\_connect\_config](#input\_private\_service\_connect\_config) | Optional) Optional. Configuration for private service connect. network and privateServiceConnectConfig are mutually exclusive. | <pre>object({<br>    enable_private_service_connect = bool,<br>    project_allowlist              = list(string),<br>  })</pre> | `null` | no |
| <a name="input_public_endpoint_enabled"></a> [public\_endpoint\_enabled](#input\_public\_endpoint\_enabled) | If true, the deployed index will be accessible through public endpoint. | `bool` | `false` | no |
| <a name="input_reserved_ip_ranges"></a> [reserved\_ip\_ranges](#input\_reserved\_ip\_ranges) | A list of reserved ip ranges under the VPC network that can be used for this DeployedIndex. If set, we will deploy the index within the provided ip ranges. | `list(string)` | `[]` | no |
| <a name="input_shard_size"></a> [shard\_size](#input\_shard\_size) | Index data is split into equal parts to be processed. These are called 'shards'. The shard size must be specified when creating an index. | `string` | `"SHARD_SIZE_SMALL"` | no |
| <a name="input_tree_ah_config"></a> [tree\_ah\_config](#input\_tree\_ah\_config) | Configuration options for using the tree-AH algorithm (Shallow tree + Asymmetric Hashing). | <pre>object({<br>    leaf_node_embedding_count    = optional(number)<br>    leaf_nodes_to_search_percent = optional(number)<br>  })</pre> | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_vector_search_instance_details"></a> [vector\_search\_instance\_details](#output\_vector\_search\_instance\_details) | Display Vector search instance attributes including project ID, region, index name, index endpoint name. |
<!-- END_TF_DOCS -->

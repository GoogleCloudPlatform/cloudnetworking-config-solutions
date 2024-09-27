<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_vertex_ai_index.index](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vertex_ai_index) | resource |
| [google_vertex_ai_index_endpoint.index_endpoint](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vertex_ai_index_endpoint) | resource |
| [google_vertex_ai_index_endpoint_deployed_index.basic_deployed_index](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vertex_ai_index_endpoint_deployed_index) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_deployed_index_id"></a> [deployed\_index\_id](#input\_deployed\_index\_id) | The user specified ID of the DeployedIndex. | `string` | n/a | yes |
| <a name="input_index_description"></a> [index\_description](#input\_index\_description) | The description of the Index. | `string` | n/a | yes |
| <a name="input_index_display_name"></a> [index\_display\_name](#input\_index\_display\_name) | The display name of the Index. The name can be up to 128 characters long and can consist of any UTF-8 characters. | `string` | n/a | yes |
| <a name="input_index_endpoint_description"></a> [index\_endpoint\_description](#input\_index\_endpoint\_description) | The description of the Index Endpoint. | `string` | n/a | yes |
| <a name="input_index_endpoint_display_name"></a> [index\_endpoint\_display\_name](#input\_index\_endpoint\_display\_name) | The display name of the Index Endpoint. The name can be up to 128 characters long and can consist of any UTF-8 characters. | `string` | n/a | yes |
| <a name="input_index_endpoint_network"></a> [index\_endpoint\_network](#input\_index\_endpoint\_network) | The full name of the Google Compute Engine network to which the index endpoint should be peered. Private services access must already be configured for the network. If left unspecified, the index endpoint is not peered with any network. Format: projects/{project}/global/networks/{network}. Where {project} is a project number, as in 12345, and {network} is network name. | `string` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The ID of the project in which the resource belongs. If it is not provided, the provider project is used. | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The region of the index endpoint. | `string` | n/a | yes |
| <a name="input_approximate_neighbors_count"></a> [approximate\_neighbors\_count](#input\_approximate\_neighbors\_count) | The default number of neighbors to find via approximate search before exact reordering is performed. Exact reordering is a procedure where results returned by an approximate search algorithm are reordered via a more expensive distance computation. Required if tree-AH algorithm is used. | `number` | `150` | no |
| <a name="input_automatic_resources"></a> [automatic\_resources](#input\_automatic\_resources) | The minimum number and maximum number of replicas this DeployedModel will be always deployed on. | <pre>object({<br>    min_replica_count = optional(number)<br>    max_replica_count = optional(number)<br>  })</pre> | `null` | no |
| <a name="input_brute_force_config"></a> [brute\_force\_config](#input\_brute\_force\_config) | Configuration options for using brute force search, which simply implements the standard linear search in the database for each query. | `string` | `null` | no |
| <a name="input_contents_delta_uri"></a> [contents\_delta\_uri](#input\_contents\_delta\_uri) | Allows inserting, updating or deleting the contents of the Matching Engine Index. The string must be a valid Cloud Storage directory path. If this field is set when calling IndexService.UpdateIndex, then no other Index field can be also updated as part of the same call. | `string` | `""` | no |
| <a name="input_dedicated_resources"></a> [dedicated\_resources](#input\_dedicated\_resources) | The type of the machine, minimum number and maximum number of replicas this DeployedModel will be always deployed on. | <pre>object({<br>    machine_spec = object({<br>      machine_type = string<br>    })<br>    min_replica_count = optional(number)<br>    max_replica_count = optional(number)<br>  })</pre> | `null` | no |
| <a name="input_deployed_display_name"></a> [deployed\_display\_name](#input\_deployed\_display\_name) | The display name of the deployment. | `string` | `null` | no |
| <a name="input_deployed_index_auth_config"></a> [deployed\_index\_auth\_config](#input\_deployed\_index\_auth\_config) | The authentication provider that the DeployedIndex uses.A list of allowed JWT issuers. Each entry must be a valid Google service account, in the following format: service-account-name@project-id.iam.gserviceaccount.com | <pre>object({<br>    auth_provider = object({<br>      audiences       = optional(string)<br>      allowed_issuers = optional(list(string))<br>    })<br>  })</pre> | `null` | no |
| <a name="input_deployment_group"></a> [deployment\_group](#input\_deployment\_group) | The deployment group can be no longer than 64 characters (eg: 'test', 'prod'). If not set, we will use the 'default' deployment group. | `string` | `null` | no |
| <a name="input_dimensions"></a> [dimensions](#input\_dimensions) | The number of dimensions of the input vectors. | `number` | `2` | no |
| <a name="input_distance_measure_type"></a> [distance\_measure\_type](#input\_distance\_measure\_type) | The distance measure used in nearest neighbor search. | `string` | `"DOT_PRODUCT_DISTANCE"` | no |
| <a name="input_enable_access_logging"></a> [enable\_access\_logging](#input\_enable\_access\_logging) | If true, private endpoint's access logs are sent to Cloud Logging. | `bool` | `true` | no |
| <a name="input_index_endpoint_labels"></a> [index\_endpoint\_labels](#input\_index\_endpoint\_labels) | Labels to be attached to index endpoint instances. | `map(string)` | `null` | no |
| <a name="input_index_labels"></a> [index\_labels](#input\_index\_labels) | Labels to be attached to index instances. | `map(string)` | `null` | no |
| <a name="input_index_update_method"></a> [index\_update\_method](#input\_index\_update\_method) | The update method to use with this Index. The value must be the followings. If not set, BATCH\_UPDATE will be used by default. | `string` | `"BATCH_UPDATE"` | no |
| <a name="input_private_service_connect_config"></a> [private\_service\_connect\_config](#input\_private\_service\_connect\_config) | (Optional) Configuration for private service connect. network and privateServiceConnectConfig are mutually exclusive. | <pre>object({<br>    enable_private_service_connect = bool,<br>    project_allowlist              = list(string),<br>  })</pre> | `null` | no |
| <a name="input_public_endpoint_enabled"></a> [public\_endpoint\_enabled](#input\_public\_endpoint\_enabled) | (Optional) If true, the deployed index will be accessible through public endpoint. | `bool` | `false` | no |
| <a name="input_reserved_ip_ranges"></a> [reserved\_ip\_ranges](#input\_reserved\_ip\_ranges) | A list of reserved ip ranges under the VPC network that can be used for this DeployedIndex. If set, we will deploy the index within the provided ip ranges. | `list(string)` | `null` | no |
| <a name="input_shard_size"></a> [shard\_size](#input\_shard\_size) | Index data is split into equal parts to be processed. These are called 'shards'. The shard size must be specified when creating an index. | `string` | `"SHARD_SIZE_SMALL"` | no |
| <a name="input_tree_ah_config"></a> [tree\_ah\_config](#input\_tree\_ah\_config) | Configuration options for using the tree-AH algorithm (Shallow tree + Asymmetric Hashing). | <pre>object({<br>    leaf_node_embedding_count    = optional(number)<br>    leaf_nodes_to_search_percent = optional(number)<br>  })</pre> | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_deploy_id"></a> [deploy\_id](#output\_deploy\_id) | An identifier for the deployed index resource. |
| <a name="output_deploy_index_name"></a> [deploy\_index\_name](#output\_deploy\_index\_name) | The name of the DeployedIndex resource. |
| <a name="output_deployed_indexes"></a> [deployed\_indexes](#output\_deployed\_indexes) | The pointers to DeployedIndexes created from this Index. An Index can be only deleted if all its DeployedIndexes had been undeployed first. |
| <a name="output_index_endpoint_id"></a> [index\_endpoint\_id](#output\_index\_endpoint\_id) | An identifier for the index endpoint resource with format. |
| <a name="output_index_endpoint_name"></a> [index\_endpoint\_name](#output\_index\_endpoint\_name) | The resource name of the Index. |
| <a name="output_index_id"></a> [index\_id](#output\_index\_id) | An identifier for the index resource. |
| <a name="output_index_name"></a> [index\_name](#output\_index\_name) | The resource name of the Index. |
| <a name="output_private_endpoints"></a> [private\_endpoints](#output\_private\_endpoints) | Provides paths for users to send requests directly to the deployed index services running on Cloud via private services access. |
| <a name="output_public_endpoint_domain_name"></a> [public\_endpoint\_domain\_name](#output\_public\_endpoint\_domain\_name) | If publicEndpointEnabled is true, this field will be populated with the domain name to use for this index endpoint. |
<!-- END_TF_DOCS -->

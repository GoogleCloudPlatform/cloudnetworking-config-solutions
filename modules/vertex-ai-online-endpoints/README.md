## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_vertex_ai_endpoint.endpoint](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vertex_ai_endpoint) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding Online Endpoints configuration values. | `string` | n/a | yes |
| <a name="input_description"></a> [description](#input\_description) | The description of the Vertex AI endpoint. | `string` | `"Sample CNCS vertex AI endpoint deployment"` | no |
| <a name="input_display_name"></a> [display\_name](#input\_display\_name) | The display name of the Vertex AI endpoint. | `string` | `"cncs-vertex-ai-display-name"` | no |
| <a name="input_labels"></a> [labels](#input\_labels) | The labels to associate with the Vertex AI endpoint. | `map(string)` | `{}` | no |
| <a name="input_location"></a> [location](#input\_location) | The location of the Vertex AI endpoint. | `string` | n/a | yes |
| <a name="input_name"></a> [name](#input\_name) | The name of the Vertex AI endpoint. | `string` | `"cncs-vertex-ai-endpoint-name"` | no |
| <a name="input_network"></a> [network](#input\_network) | The full name of the Google Compute Engine network to which the Endpoint should be peered. | `string` | n/a | yes |
| <a name="input_project"></a> [project](#input\_project) | The project of the Vertex AI endpoint. | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The region of the Vertex AI endpoint. | `string` | `"us-central1"` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_create_time"></a> [create\_time](#output\_create\_time) | Timestamp when this Endpoint was created. |
| <a name="output_deployed_models"></a> [deployed\_models](#output\_deployed\_models) | The models deployed in this Endpoint. |
| <a name="output_effective_labels"></a> [effective\_labels](#output\_effective\_labels) | All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services. |
| <a name="output_endpoint_configuration"></a> [endpoint\_configuration](#output\_endpoint\_configuration) | Configuration details provided for the Vertex AI endpoint. |
| <a name="output_etag"></a> [etag](#output\_etag) | Used to perform consistent read-modify-write updates. If not set, a blind 'overwrite' update happens. |
| <a name="output_id"></a> [id](#output\_id) | Identifier for the resource with format projects/{{project}}/locations/{{location}}/endpoints/{{name}} |
| <a name="output_model_deployment_monitoring_job"></a> [model\_deployment\_monitoring\_job](#output\_model\_deployment\_monitoring\_job) | Resource name of the Model Monitoring job associated with this Endpoint if monitoring is enabled by CreateModelDeploymentMonitoringJob. Format: projects/{project}/locations/{location}/modelDeploymentMonitoringJobs/{model\_deployment\_monitoring\_job} |
| <a name="output_terraform_labels"></a> [terraform\_labels](#output\_terraform\_labels) | The combination of labels configured directly on the resource and default labels configured on the provider. |
| <a name="output_update_time"></a> [update\_time](#output\_update\_time) | Timestamp when this Endpoint was last updated. |
# Copyright 2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

variable "config_folder_path" {
  description = "Location of YAML files holding Vector Search configuration values."
  type        = string
  default     = "../../../configuration/producer/VectorSearch/config"
}

variable "index_labels" {
  description = "Labels to be attached to index instances."
  type        = map(string)
  default     = null
}

variable "index_description" {
  type        = string
  description = "The description of the Index."
  default     = null
}

variable "index_update_method" {
  default     = "BATCH_UPDATE"
  type        = string
  description = "The update method to use with this Index. The value must be the followings. If not set, BATCH_UPDATE will be used by default."
}

variable "index_endpoint_description" {
  type        = string
  default     = null
  description = "The description of the Index Endpoint."
}

variable "index_endpoint_labels" {
  description = "Labels to be attached to index endpoint instances."
  type        = map(string)
  default     = null
}

variable "contents_delta_uri" {
  type        = string
  default     = null
  description = "Allows inserting, updating or deleting the contents of the Matching Engine Index. The string must be a valid Cloud Storage directory path. If this field is set when calling IndexService.UpdateIndex, then no other Index field can be also updated as part of the same call."
}

variable "dimensions" {
  type        = number
  default     = 2
  description = "The number of dimensions of the input vectors."
}

variable "approximate_neighbors_count" {
  type        = number
  default     = 150
  description = "The default number of neighbors to find via approximate search before exact reordering is performed. Exact reordering is a procedure where results returned by an approximate search algorithm are reordered via a more expensive distance computation. Required if tree-AH algorithm is used."
}

variable "shard_size" {
  type        = string
  default     = "SHARD_SIZE_SMALL"
  description = "Index data is split into equal parts to be processed. These are called 'shards'. The shard size must be specified when creating an index."
}

variable "distance_measure_type" {
  type        = string
  default     = "DOT_PRODUCT_DISTANCE"
  description = "The distance measure used in nearest neighbor search."
}

variable "tree_ah_config" {
  type = object({
    leaf_node_embedding_count    = optional(number)
    leaf_nodes_to_search_percent = optional(number)
  })
  default     = null
  description = "Configuration options for using the tree-AH algorithm (Shallow tree + Asymmetric Hashing). "
}

variable "brute_force_config" {
  type        = string
  default     = null
  description = "Configuration options for using brute force search, which simply implements the standard linear search in the database for each query."
}

variable "deployed_display_name" {
  type        = string
  default     = null
  description = "The display name of the deployment."
}

variable "reserved_ip_ranges" {
  type        = list(string)
  default     = []
  description = "A list of reserved ip ranges under the VPC network that can be used for this DeployedIndex. If set, we will deploy the index within the provided ip ranges."
}

variable "enable_access_logging" {
  type        = bool
  default     = true
  description = "If true, private endpoint's access logs are sent to Cloud Logging."
}

variable "deployment_group" {
  type        = string
  default     = null
  description = "The deployment group can be no longer than 64 characters (eg: 'test', 'prod'). If not set, we will use the 'default' deployment group."
}

variable "automatic_resources" {
  type = object({
    min_replica_count = optional(number)
    max_replica_count = optional(number)
  })
  default     = null
  description = "The minimum number and maximum number of replicas this DeployedModel will be always deployed on."
}

variable "dedicated_resources" {
  type = object({
    machine_spec = object({
      machine_type = string
    })
    min_replica_count = optional(number)
    max_replica_count = optional(number)
  })
  default     = null
  description = "The type of the machine, minimum number and maximum number of replicas this DeployedModel will be always deployed on."
}

variable "deployed_index_auth_config" {
  type = object({
    auth_provider = object({
      audiences       = optional(string)
      allowed_issuers = optional(list(string))
    })
  })
  default     = null
  description = "The authentication provider that the DeployedIndex uses.A list of allowed JWT issuers. Each entry must be a valid Google service account, in the following format: service-account-name@project-id.iam.gserviceaccount.com"
}

variable "public_endpoint_enabled" {
  type        = bool
  default     = false
  description = "If true, the deployed index will be accessible through public endpoint."
}

variable "private_service_connect_config" {
  type = object({
    enable_private_service_connect = bool,
    project_allowlist              = list(string),
  })
  default     = null
  description = "Optional) Optional. Configuration for private service connect. network and privateServiceConnectConfig are mutually exclusive."
}

variable "index_endpoint_network" {
  type        = string
  description = "The full name of the Google Compute Engine network to which the index endpoint should be peered. Private services access must already be configured for the network. If left unspecified, the index endpoint is not peered with any network. Format: projects/{project}/global/networks/{network}. Where {project} is a project number, as in 12345, and {network} is network name."
  default     = null
}

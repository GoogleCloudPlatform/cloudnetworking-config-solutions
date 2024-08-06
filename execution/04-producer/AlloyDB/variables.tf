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

variable "allocated_ip_range" {
  type        = string
  description = "The name of the allocated IP range for the private IP AlloyDB cluster. For example: google-managed-services-default. If set, the instance IPs for this cluster will be created in the allocated range."
  default     = null
}

variable "automated_backup_policy" {
  description = "The automated backup policy for this cluster. If no policy is provided then the default policy will be used. The default policy takes one backup a day, has a backup window of 1 hour, and retains backups for 14 days."
  type = object({
    location      = optional(string)
    backup_window = optional(string)
    enabled       = optional(bool)

    weekly_schedule = optional(object({
      days_of_week = optional(list(string))
      start_times  = list(string)
    })),

    quantity_based_retention_count = optional(number)
    time_based_retention_count     = optional(string)
    labels                         = optional(map(string))
    backup_encryption_key_name     = optional(string)
  })
  default = null
}

variable "read_pool_instance" {
  description = "List of Read Pool Instances to be created."
  type = list(object({
    instance_id        = string
    display_name       = string
    node_count         = optional(number, 1)
    database_flags     = optional(map(string))
    availability_type  = optional(string)
    gce_zone           = optional(string)
    machine_cpu_count  = optional(number, 2)
    ssl_mode           = optional(string)
    require_connectors = optional(bool)
    query_insights_config = optional(object({
      query_string_length     = optional(number)
      record_application_tags = optional(bool)
      record_client_address   = optional(bool)
      query_plans_per_minute  = optional(number)
    }))
  }))
  default = []
  validation {
    condition     = try(alltrue([for rp in var.read_pool_instance : contains(["2", "4", "8", "16", "32", "64", "96", "128"], tostring(rp.machine_cpu_count))]), false) || var.read_pool_instance == null
    error_message = "machine_cpu_count must be one of [2, 4, 8, 16, 32, 64, 96, 128]"
  }
}

variable "cluster_labels" {
  description = "User-defined labels for the alloydb cluster."
  type        = map(string)
  default     = {}
}

variable "cluster_initial_user" {
  description = "Alloy DB Cluster Initial User Credentials."
  type = object({
    user     = optional(string),
    password = string
  })
  default = null
}

variable "database_version" {
  type        = string
  description = "The database engine major version. This is an optional field and it's populated at the Cluster creation time. This field cannot be changed after cluster creation. Possible valus: POSTGRES_14, POSTGRES_15."
  default     = "POSTGRES_15"
}

variable "cluster_encryption_key_name" {
  description = "The fully-qualified resource name of the KMS key for cluster encryption. Each Cloud KMS key is regionalized and has the following format: projects/[PROJECT]/locations/[REGION]/keyRings/[RING]/cryptoKeys/[KEY_NAME]."
  type        = string
  default     = null
}

variable "config_folder_path" {
  description = "Location of YAML files holding AlloyDB configuration values."
  type        = string
  default     = "../../../configuration/producer/AlloyDB/config"
}

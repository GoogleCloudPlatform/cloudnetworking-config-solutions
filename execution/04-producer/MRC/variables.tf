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

variable "shard_count" {
  type        = number
  description = "Number of shards (replicas) in the Redis cluster."
  default     = 3 # You can adjust based on your requirements
}

variable "region" {
  type        = string
  description = "The region in which to create the Redis cluster."
  default     = "us-central1"
}

variable "replica_count" {
  type        = number
  description = "Number of replicas per shard in the Redis cluster."
  default     = 1
}

variable "config_folder_path" {
  description = "Location of YAML files holding MRC configuration values."
  type        = string
  default     = "../../../configuration/producer/MRC/config"
}
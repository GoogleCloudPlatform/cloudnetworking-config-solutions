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

variable "name" {
  type        = string
  description = "The name of the Vertex AI endpoint."
  default     = "cncs-vertex-ai-endpoint-name"
}

variable "display_name" {
  type        = string
  description = "The display name of the Vertex AI endpoint."
  default     = "cncs-vertex-ai-display-name"
}

variable "description" {
  type        = string
  description = "The description of the Vertex AI endpoint."
  default     = "Sample CNCS vertex AI endpoint deployment"
}

variable "location" {
  type        = string
  description = "The location of the Vertex AI endpoint."
  default     = "us-central1"
}

variable "labels" {
  type        = map(string)
  description = "The labels to associate with the Vertex AI endpoint."
  default     = {}
}

variable "config_folder_path" {
  description = "Location of YAML files holding Online Endpoints configuration values."
  type        = string
  default     = "../../../configuration/producer/Vertex-AI-Online-Endpoints/config"
}

variable "region" {
  type        = string
  description = "The region of the Vertex AI endpoint."
  default     = "us-central1"
}

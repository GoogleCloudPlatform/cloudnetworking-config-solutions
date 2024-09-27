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

locals {
  config_folder_path = var.config_folder_path
  instances          = [for file in fileset(local.config_folder_path, "[^_]*.yaml") : yamldecode(file("${local.config_folder_path}/${file}"))]
  instance_list = flatten([
    for instance in try(local.instances, []) : {
      project_id             = instance.project_id
      name                   = instance.name
      region                 = instance.region
      containers             = try(instance.containers, var.containers)
      create_job             = try(instance.create_job, var.create_job)
      custom_audiences       = try(instance.custom_audiences, var.custom_audiences)
      encryption_key         = try(instance.encryption_key, var.encryption_key)
      eventarc_triggers      = try(instance.eventarc_triggers, var.eventarc_triggers)
      iam                    = try(instance.iam, var.iam)
      ingress                = try(instance.ingress, var.ingress)
      labels                 = try(instance.labels, var.labels)
      launch_stage           = try(instance.launch_stage, var.launch_stage)
      prefix                 = try(instance.prefix, var.prefix)
      revision               = try(instance.revision, var.revision)
      service_account        = try(instance.service_account, var.service_account)
      service_account_create = try(instance.service_account_create, var.service_account_create)
      tag_bindings           = try(instance.tag_bindings, var.tag_bindings)
      volumes                = try(instance.volumes, var.volumes)
      vpc_connector_create   = try(instance.vpc_connector_create, var.vpc_connector_create)
    }
  ])
}

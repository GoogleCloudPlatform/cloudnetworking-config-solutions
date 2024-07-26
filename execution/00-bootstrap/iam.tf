
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

/********************************************
 Service Account used to run Organization Stage
*********************************************/
module "organization" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v32.0.1"
  project_id = var.bootstrap_project_id
  name       = var.organization_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.organization_stage_administrator
  }
}
/********************************************
 Service Account used to run Networking Stage
*********************************************/
module "networking" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v32.0.1"
  project_id = var.bootstrap_project_id
  name       = var.networking_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.networking_stage_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.networkAdmin"
    ]
  }
}
/********************************************
 Service Account used to run Security Stage
*********************************************/
module "security" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v32.0.1"
  project_id = var.bootstrap_project_id
  name       = var.security_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.security_stage_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.securityAdmin"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}
/********************************************
 Service Account used to run Producer Stage
*********************************************/
module "producer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v32.0.1"
  project_id = var.bootstrap_project_id
  name       = var.producer_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.producer_stage_administrator
  }
  iam_project_roles = {
    (var.network_serviceproject_id) = [
      "roles/cloudsql.admin",
      "roles/alloydb.admin",
      "roles/redis.admin"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}
/****************************************************
 Service Account used to run Networking Manual Stage
*****************************************************/
module "networking_manual" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v32.0.1"
  project_id = var.bootstrap_project_id
  name       = var.networking_manual_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.networking_manual_stage_administrator
  }
  iam_project_roles = {
    (var.network_serviceproject_id) = [
      "roles/cloudsql.admin",
      "roles/alloydb.admin"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}
/********************************************
 Service Account used to run Consumer Stage
*********************************************/
module "consumer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v32.0.1"
  project_id = var.bootstrap_project_id
  name       = var.consumer_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.consumer_stage_administrator
  }
  iam_project_roles = {
    (var.network_serviceproject_id) = [
      "roles/compute.instanceAdmin.v1",
      "roles/iam.serviceAccountUser"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}


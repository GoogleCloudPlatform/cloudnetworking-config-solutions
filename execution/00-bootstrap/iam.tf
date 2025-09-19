# Copyright 2024-2025 Google LLC
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
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.organization_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.organization_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/iam.serviceAccountUser",
      "roles/serviceusage.serviceUsageAdmin",
    ]
    (var.network_serviceproject_id) = [
      "roles/iam.serviceAccountUser",
      "roles/serviceusage.serviceUsageAdmin",
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run Networking Stage
*********************************************/

module "networking" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.networking_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.networking_administrator
  }
  iam_folder_roles = {
    (var.folder_id) = [
      "roles/compute.xpnAdmin",
    ]
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.networkAdmin",
    ]
    (var.network_serviceproject_id) = [
      "roles/cloudsql.viewer"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run DNS Managed Zones Stage
*********************************************/

module "dns_managed_zones" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.dns_managed_zones_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.dns_managed_zones_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/dns.admin",
      "roles/compute.networkUser"
    ]
    (var.network_serviceproject_id) = [
      "roles/dns.admin",
      "roles/compute.networkUser"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run DNS Response Policy Stage
*********************************************/

module "dns_response_policy" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.dns_response_policy_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.dns_response_policy_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/dns.admin",
      "roles/compute.networkUser"
    ]
    (var.network_serviceproject_id) = [
      "roles/dns.admin",
      "roles/compute.networkUser"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run Security Stage
*********************************************/

module "security" {
  source          = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id      = var.bootstrap_project_id
  organization_id = var.organization_id
  name            = var.security_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.security_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.securityAdmin",
      "roles/compute.orgFirewallPolicyAdmin",
      "roles/networksecurity.securityProfileAdmin",
    ]
  }
  iam_organization_roles = {
    (var.organization_id) = [
      "roles/compute.orgFirewallPolicyAdmin",
      "roles/resourcemanager.organizationViewer",
      "roles/networksecurity.securityProfileAdmin",
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run CloudSQL Producer Stage
*********************************************/

module "cloudsql_producer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.producer_cloudsql_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.producer_cloudsql_administrator
  }
  iam_project_roles = {
    (var.network_serviceproject_id) = [
      "roles/cloudsql.admin"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run AlloyDB Producer Stage
*********************************************/

module "alloydb_producer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.producer_alloydb_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.producer_alloydb_administrator
  }
  iam_project_roles = {
    (var.network_serviceproject_id) = [
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
 Service Account used to run MRC Producer Stage
*********************************************/

module "mrc_producer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.producer_mrc_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.producer_mrc_administrator
  }
  iam_project_roles = {
    (var.network_serviceproject_id) = [
      "roles/redis.admin"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run Vertex AI Producer Stages
*********************************************/

module "vertex_producer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.producer_vertex_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.producer_vertex_administrator
  }
  iam_project_roles = {
    (var.network_serviceproject_id) = [
      "roles/aiplatform.admin"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run GKE Producer Stage
*********************************************/

module "gke_producer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.producer_gke_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.producer_gke_administrator
  }
  iam_project_roles = {
    (var.network_serviceproject_id) = [
      "roles/container.admin",
      "roles/compute.instanceAdmin",
      "roles/iam.serviceAccountAdmin",
      "roles/iam.serviceAccountUser",
      "roles/resourcemanager.projectIamAdmin",
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/************************************************
 Service Account used to run BigQuery Producer Stage
*************************************************/

module "bigquery_producer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.project_id
  name       = var.producer_bq_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.producer_bq_administrator
  }

  iam_project_roles = {
    (var.project_id) = [
      "roles/bigquery.admin"
    ]
  }

  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/****************************************************
 Service Account used to run Producer Connectivity Stage
*****************************************************/

module "producer_connectivity" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.producer_connectivity_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.producer_connectivity_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.networkAdmin",
    ]
    (var.network_serviceproject_id) = [
      "roles/cloudsql.viewer",
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run GCE Consumer Stage
*********************************************/

module "gce_consumer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.consumer_gce_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.consumer_gce_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.networkUser",
    ]
    (var.network_serviceproject_id) = [
      "roles/compute.instanceAdmin.v1",
      "roles/iam.serviceAccountUser",
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run Cloud Run Consumer Stage
*********************************************/

module "cloudrun_consumer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.consumer_cloudrun_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.consumer_cloudrun_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.networkUser",
    ]
    (var.network_serviceproject_id) = [
      "roles/iam.serviceAccountUser",
      "roles/run.admin"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run MIG Consumer Stage
*********************************************/

module "mig_consumer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.consumer_mig_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.consumer_mig_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.networkUser",
    ]
    (var.network_serviceproject_id) = [
      "roles/compute.instanceAdmin.v1",
      "roles/iam.serviceAccountUser",
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run Workbench Consumer Stage
*********************************************/

module "workbench_consumer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.consumer_workbench_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.consumer_workbench_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.networkUser",
    ]
    (var.network_serviceproject_id) = [
      "roles/iam.serviceAccountUser", // Allow impersonation of the service account
      "roles/notebooks.admin",        // Grant access to Notebooks resources
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run Consumer Load Balancing Stage
*********************************************/

module "consumer_load_balancing" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.consumer_lb_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.consumer_lb_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.loadBalancerAdmin"
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run Consumer VPC Access Connector Stage
*********************************************/

module "consumer_vpc_access_connector" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.consumer_vpc_connector_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.consumer_vpc_connector_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/vpcaccess.admin",
    ]
    (var.network_serviceproject_id) = [
      "roles/compute.networkViewer",

      # *Important*: If using Shared VPC and the connector needs to attach to
      # a subnet in the network_project_id (Host Project), the SA *might*
      # also need 'roles/compute.networkUser' on the Host Project or specific subnets.
      # Test if networkViewer is sufficient first.
      # "roles/compute.networkUser",
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run App Engine Consumer Stage
*********************************************/

module "appengine_consumer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.consumer_appengine_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.consumer_appengine_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.networkUser",
    ]
    (var.network_serviceproject_id) = [
      "roles/compute.instanceAdmin.v1",
      "roles/iam.serviceAccountUser",
      "roles/appengine.appAdmin",       // App Engine Admin
      "roles/cloudbuild.builds.editor", // Cloud Build Editor
      "roles/artifactregistry.writer",  // Artifact Registry Writer
      "roles/compute.networkViewer",    // Compute Engine Network Viewer
      "roles/storage.objectViewer",     // Storage Object Viewer
      "roles/vpcaccess.user",           // VPC access connector User
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run UMIG Consumer Stage
*********************************************/

module "umig_consumer" {
  source     = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id = var.bootstrap_project_id
  name       = var.consumer_umig_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.consumer_umig_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/compute.networkUser",
    ]
    (var.network_serviceproject_id) = [
      "roles/compute.instanceAdmin.v1",
      "roles/iam.serviceAccountUser",
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}

/********************************************
 Service Account used to run Network Security Integration
*********************************************/

module "network_security_integration" {
  source          = "github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/iam-service-account?ref=v44.1.0"
  project_id      = var.bootstrap_project_id
  organization_id = var.organization_id
  name            = var.nsi_sa_name
  iam = {
    "roles/iam.serviceAccountTokenCreator" = var.nsi_administrator
  }
  iam_project_roles = {
    (var.network_hostproject_id) = [
      "roles/networksecurity.packetMirroringAdmin",
      "roles/networksecurity.securityProfileAdmin",
    ]
  }
  iam_organization_roles = {
    (var.organization_id) = [
      "roles/resourcemanager.organizationViewer",
      "roles/networksecurity.securityProfileAdmin",
    ]
  }
  iam_storage_roles = {
    (module.google_storage_bucket.name) = [
      "roles/storage.objectAdmin"
    ]
  }
}
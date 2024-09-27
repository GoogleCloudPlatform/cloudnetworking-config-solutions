## Introduction

AlloyDB is designed to supercharge your PostgreSQL experience, making it faster, more scalable and highly available for your most demanding applications. It's a fully managed database service that combines the familiarity of PostgreSQL with Google's innovative database engine.

AlloyDB is your ideal database solution if you're looking for:

* High Performance
* Scalability
* PostgreSQL Compatibility
* Managed Service

## Pre-Requisites

Before creating your first AlloyDB instance, ensure you have the completed the following prerequsites:

1. **Completed Prior Stages:** Successful deployment of AlloyDB resources depends on the completion of the following stages:
    * **01-organization:** This stage handles the activation of required Google Cloud APIs for AlloyDB.
    * **02-networking:** This stage sets up the necessary network infrastructure, such as VPCs and subnets, to support AlloyDB connectivity.
    * **03-security/AlloyDB:** This stage configures firewall rules to allow access to AlloyDB instances on the appropriate ports and IP ranges.

2. **API Enablement:** Ensure the following Google Cloud APIs have been enabled:
   * IAM API (`iam.googleapis.com`)
   * Compute API (`compute.googleapis.com`)
   * AlloyDB API (`alloydb.googleapis.com`)
   * Service Networking API (`servicenetworking.googleapis.com`)
   * Cloud Resource Manager API (`cloudresourcemanager.googleapis.com`)

3. **IAM Permissions:**  Grant yourself (or the appropriate users/service accounts) the following IAM roles at the project level (or higher):
   * AlloyDB Administrator (`roles/alloydb.admin`)
   * Cloud Resource Manager Project IAM Admin (`roles/resourcemanager.projectIamAdmin`) – This is needed to create the AlloyDB service account.

## Let's Get Started! 🚀
With the prerequisites in place and your AlloyDB configuration files ready, you can now leverage Terraform to automate the creation of your AlloyDB instances. Here's the workflow:

### Execution Steps

1. **Create your configuration files:**

    * Create YAML files defining the properties of each AlloyDB instance you want to create. Ensure these files are stored in the configuration/producer/AlloyDB/config folder within this repository.

    * Each YAML file should map to a single AlloyDB instance, providing details such as instance name, region, instance type, storage, and networking configuration. Each field and its structure are described in the [input section](#inputs) below.

    * For reference on how to structure your AlloyDB configuration YAML files, see the [example](#example) section below or refer to sample YAML file at folder location `configuration/producer/AlloyDB/config/instance.yaml.example`. These examples provide templates that you can adapt to your specific needs.


2. **Initialize Terraform:**

    * Open your terminal and navigate to the AlloyDB directory containing the Terraform configuration.

    * Run the following command to initialize Terraform:

    ```
    terraform init
    ```
3. **Review the Execution Plan:**

    * Use the terraform plan command to generate an execution plan. This will show you the changes Terraform will make to your Google Cloud infrastructure:

    ```
    terraform plan -var-file=../../../configuration/producer/AlloyDB/alloydb.tfvars
    ```

Carefully review the plan to ensure it aligns with your intended configuration.

4. **Apply the Configuration:**

    Once you're satisfied with the plan, execute the terraform apply command to provision your AlloyDB instances:

    ```
    terraform apply -var-file=../../../configuration/producer/AlloyDB/alloydb.tfvars
    ```

Terraform will read the YAML files from the `configuration/producer/AlloyDB/config` folder and create the corresponding AlloyDB instances in your Google Cloud project.

5. **Monitor and Manage:**
    * After the instances are created, you can monitor their status, performance, and logs through the Google Cloud Console or using the Google Cloud CLI.

    * Use Terraform to manage updates and changes to your AlloyDB instances as needed.


### Example

To help you get started, we've provided examples of YAML configuration files that you can use as templates for your AlloyDB instances.

* **Minimal YAML (Mandatory Fields Only):**
This minimal example includes only the essential fields required to create a basic AlloyDB instance.

  ```
  project_id: <Project-ID>
  cluster_id: cn-alloydb-cid
  cluster_display_name: cn-alloydb-cid
  region: us-central1
  network_id: projects/<Project-ID>/global/networks/<Network-Name>
  primary_instance:
    instance_id : cn-alloydb-id-12
  ```

* **Comprehensive YAML (All Available Fields):**
This comprehensive example showcases all available fields, allowing you to customize your AlloyDB instance with advanced settings for performance, availability and network configuration.

  ```
  project_id: <Project-ID>
  cluster_id: cn-alloydb-cid
  cluster_display_name: cn-alloydb-cid
  region: us-central1
  network_id: projects/<Project-ID>/global/networks/<Network-Name>
  allocated_ip_range: psarange
  database_version: POSTGRES_15
  primary_instance:
    instance_id : cn-alloydb-id
    display_name : cn-alloydb-id
    instance_type : PRIMARY
    machine_cpu_count : 2
    database_flags : null
  cluster_labels:
    environment: development
  cluster_initial_user:
    user: admin
    password: admin
  read_pool_instance :
    - instance_id: read-instance-1
      display_name: read-instance-1
      node_count: 1
      database_flags: null
      availability_type: ZONAL
      gce_zone: us-central1-a
      machine_cpu_count: 2
      ssl_mode: ALLOW_UNENCRYPTED_AND_ENCRYPTED
      require_connectors: false
  automated_backup_policy :
    location: us-central1
    backup_window:  1800s
    enabled: true
    weekly_schedule:
      days_of_week:
        - MONDAY
      start_times:
        - 2:00:00:00
    quantity_based_retention_count: 1
    time_based_retention_count: null
    labels:
      environment: development
    backup_encryption_key_name: null
  cluster_encryption_key_name: null
  ```

## Important Notes:

* This README is a starting point. Customize it to include specific details about your AlloyDB projects and resources.
Refer to the official Google Cloud AlloyDB documentation for the most up-to-date information and best practices: https://cloud.google.com/alloydb/docs

* Order of Execution: Make sure to complete the 01-organization, 02-networking, and 03-security stages before attempting to create AlloyDB instances. Terraform will leverage the resources and configurations established in these prior stages.

* Troubleshooting: If you encounter errors during the AlloyDB creation process, verify that all prerequisites are satisfied and that the dependencies between stages are correctly configured.

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_alloy_db"></a> [alloy\_db](#module\_alloy\_db) | GoogleCloudPlatform/alloy-db/google | ~> 2.2.0 |


## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
|project_id | The ID of the Google Cloud project where you want to create your AlloyDB instance. | `string` | n/a | yes |
|cluster_id | A unique identifier for your AlloyDB cluster. It must contain only lowercase letters, numbers, and hyphens. | `string` | n/a | yes |
|cluster_display_name_id | A human-readable name for your cluster that will be displayed in the Google Cloud Console. | `string` | n/a | yes |
|region | The Google Cloud region where your AlloyDB cluster will be located. | `string` | n/a | yes |
|network\_id | The Network ID of the VPC network where your AlloyDB instance will be deployed. | `string` | n/a | yes |
|primary_instance| This section configures the primary instance of your AlloyDB cluster, responsible for handling read and write operations. | <pre>object({<br>    instance_id        = string,<br>    display_name       = optional(string),<br>    database_flags     = optional(map(string))<br>    labels             = optional(map(string))<br>    annotations        = optional(map(string))<br>    gce_zone           = optional(string)<br>    availability_type  = optional(string)<br>    machine_cpu_count  = optional(number, 2)<br>    ssl_mode           = optional(string)<br>    require_connectors = optional(bool)<br>    query_insights_config = optional(object({<br>      query_string_length     = optional(number)<br>      record_application_tags = optional(bool)<br>      record_client_address   = optional(bool)<br>      query_plans_per_minute  = optional(number)<br>    }))<br>  })</pre>| n/a | yes
| <a name="input_allocated_ip_range"></a> [allocated\_ip\_range](#input\_allocated\_ip\_range) | The name of the allocated IP range for the private IP AlloyDB cluster. For example: google-managed-services-default. If set, the instance IPs for this cluster will be created in the allocated range. | `string` | `null` | no |
| <a name="input_automated_backup_policy"></a> [automated\_backup\_policy](#input\_automated\_backup\_policy) | The automated backup policy for this cluster. If no policy is provided then the default policy will be used. The default policy takes one backup a day, has a backup window of 1 hour, and retains backups for 14 days. | <pre>object({<br>    location      = optional(string)<br>    backup_window = optional(string)<br>    enabled       = optional(bool)<br><br>    weekly_schedule = optional(object({<br>      days_of_week = optional(list(string))<br>      start_times  = list(string)<br>    })),<br><br>    quantity_based_retention_count = optional(number)<br>    time_based_retention_count     = optional(string)<br>    labels                         = optional(map(string))<br>    backup_encryption_key_name     = optional(string)<br>  })</pre> | `null` | no |
| <a name="input_cluster_encryption_key_name"></a> [cluster\_encryption\_key\_name](#input\_cluster\_encryption\_key\_name) | The fully-qualified resource name of the KMS key for cluster encryption. Each Cloud KMS key is regionalized and has the following format: projects/[PROJECT]/locations/[REGION]/keyRings/[RING]/cryptoKeys/[KEY\_NAME] | `string` | `null` | no |
| <a name="input_cluster_initial_user"></a> [cluster\_initial\_user](#input\_cluster\_initial\_user) | Alloy DB Cluster Initial User Credentials. | <pre>object({<br>    user     = optional(string),<br>    password = string<br>  })</pre> | `null` | no |
| <a name="input_cluster_labels"></a> [cluster\_labels](#input\_cluster\_labels) | User-defined labels for the alloydb cluster. | `map(string)` | `{}` | no |
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding AlloyDB configuration values. | `string` | `"config"` | no |
| <a name="input_database_version"></a> [database\_version](#input\_database\_version) | The database engine major version. This is an optional field and it's populated at the Cluster creation time. This field cannot be changed after cluster creation. Possible valus: POSTGRES\_14, POSTGRES\_15 | `string` | `"POSTGRES_15"` | no |
| <a name="input_read_pool_instance"></a> [read\_pool\_instance](#input\_read\_pool\_instance) | List of Read Pool Instances to be created. | <pre>list(object({<br>    instance_id        = string<br>    display_name       = string<br>    node_count         = optional(number, 1)<br>    database_flags     = optional(map(string))<br>    availability_type  = optional(string)<br>    gce_zone           = optional(string)<br>    machine_cpu_count  = optional(number, 2)<br>    ssl_mode           = optional(string)<br>    require_connectors = optional(bool)<br>    query_insights_config = optional(object({<br>      query_string_length     = optional(number)<br>      record_application_tags = optional(bool)<br>      record_client_address   = optional(bool)<br>      query_plans_per_minute  = optional(number)<br>    }))<br>  }))</pre> | `[]` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cluster_details"></a> [cluster\_details](#output\_cluster\_details) | Display cluster name and details like cluster id, network configuration and state of the AlloyDB cluster created. |
<!-- END_TF_DOCS -->

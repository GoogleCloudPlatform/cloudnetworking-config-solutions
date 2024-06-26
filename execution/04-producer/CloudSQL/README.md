## Introduction

Cloud SQL is Google Cloud's fully managed relational database service, designed to simplify the deployment, maintenance and administration of your databases. It supports popular database engines like MySQL, PostgreSQL and SQL Server freeing you from the complexities of infrastructure management so you can focus on your applications and data.

## Pre-Requisites

Before creating your first Cloud SQL instance, ensure you have completed the following prerequisites:

1. **Completed Prior Stages:** Successful deployment of Cloud SQL resources depends on the completion of the following stages:
    * **01-organization:** This stage handles the activation of required Google Cloud APIs for Cloud SQL.
    * **02-networking:** This stage sets up the necessary network infrastructure, such as VPCs and subnets, to support Cloud SQL connectivity.
    * **03-security/CloudSQL:** This stage configures firewall rules to allow access to Cloud SQL instances on the appropriate ports and IP ranges.

2. **API Enablement:** Ensure the following Google Cloud APIs have been enabled:
   * IAM API (`iam.googleapis.com`)
   * Compute Enginer API (`compute.googleapis.com`)
   * Cloud SQL API (`sqladmin.googleapis.com`)
   * Service Networking API (`servicenetworking.googleapis.com`)
   * Cloud Resource Manager API (`cloudresourcemanager.googleapis.com`)

3. **IAM Permissions:**  Grant yourself (or the appropriate users/service accounts) the following IAM roles at the project level (or higher):
   * Cloud SQL Admin (`roles/cloudsql.admin`)

## Let's Get Started! 🚀
With the prerequisites in place and your Cloud SQL configuration files ready, you can now leverage Terraform to automate the creation of your Cloud SQL instances. Here's the workflow:

### Execution Steps

1. **Create your configuration files:**

    * Create YAML files defining the properties of each Cloud SQL instance you want to create. Ensure these files are stored in the `CloudSQL/config` folder within this repository.

    * Each YAML file should map to a single Cloud SQL instance providing details such as instance name, region, database version and networking configuration. Each field and its structure are described in the [input section](#inputs) below.

    * For reference on how to structure your Cloud SQL configuration YAML files, see the [example](#example) section below or refer to sample YAML file at folder location `CloudSQL/config/instance.yaml.example`. These examples provide templates that you can adapt to your specific needs.


2. **Initialize Terraform:**

    * Open your terminal and navigate to the Cloud SQL directory containing the Terraform configuration.

    * Run the following command to initialize Terraform:

    ```
    terraform init
    ```
3. **Review the Execution Plan:**

    * Use the terraform plan command to generate an execution plan. This will show you the changes Terraform will make to your Google Cloud infrastructure:

    ```
    terraform plan
    ```

Carefully review the plan to ensure it aligns with your intended configuration.

4. **Apply the Configuration:**

    Once you're satisfied with the plan, execute the terraform apply command to provision your Cloud SQL instances:

    ```
    terraform apply
    ```

Terraform will read the YAML files from the `04-producer/CloudSQL/config` folder and create the corresponding Cloud SQL instances in your Google Cloud project.

5. **Monitor and Manage:**
    * After the instances are created, you can monitor their status, performance, and logs through the Google Cloud Console or using the Google Cloud CLI.

    * Use Terraform to manage updates and changes to your Cloud SQL instances as needed.


### Example

To help you get started, we've provided examples of YAML configuration files that you can use as templates for your CloudSQL instances.

* **Minimal YAML (Mandatory Fields Only for PSA-based Network Configuration):**
This minimal example includes only the essential fields required to create a basic Cloud SQL instance.

  ```
  name: cloudsql
  project_id: <your-project-id>
  region: us-central1
  database_version: POSTGRES_15
  network_config:
    connectivity:
      psa_config:
        private_network: projects/<Project-ID>/global/networks/<Network-Name>
  ```

* **Minimal YAML (Mandatory Fields Only for PSC-based Network Configuration):**
This minimal example includes only the essential fields required to create a basic Cloud SQL instance.

  ```
  name: cloudsql
  project_id: <your-project-id>
  region: us-central1
  database_version: POSTGRES_15
  network_config:
    connectivity:
      psc_allowed_consumer_projects : ["your-allowed-consumer-project-id"]
  ```

## Important Notes:

* This README is a starting point. Customize it to include specific details about your Cloud SQL projects and resources.
Refer to the official Google Cloud Cloud SQL documentation for the most up-to-date information and best practices: https://cloud.google.com/cloudsql/docs

* Order of Execution: Make sure to complete the 01-organization, 02-networking, and 03-security stages before attempting to create Cloud SQL instances. Terraform will leverage the resources and configurations established in these prior stages.

* Troubleshooting: If you encounter errors during the Cloud SQL creation process, verify that all prerequisites are satisfied and that the dependencies between stages are correctly configured.

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cloudsql"></a> [cloudsql](#module\_cloudsql) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/cloudsql-instance | v31.0.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_name"></a> [name](#input\_name) | Name of primary instance. | `string` | n/a | yes |
| <a name="input_network_config"></a> [network\_config](#input\_network\_config) | Network configuration for the instance. Only one between private\_network and psc\_config can be used. | <pre>object({<br>    authorized_networks = optional(map(string))<br>    connectivity = object({<br>      public_ipv4 = optional(bool, false)<br>      psa_config = optional(object({<br>        private_network = string<br>        allocated_ip_ranges = optional(object({<br>          primary = optional(string)<br>          replica = optional(string)<br>        }))<br>      }))<br>      psc_allowed_consumer_projects = optional(list(string))<br>    })<br>  })</pre> | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The ID of the project where this instances will be created. | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | Region of the primary instance. | `string` | n/a | yes |
| <a name="input_activation_policy"></a> [activation\_policy](#input\_activation\_policy) | This variable specifies when the instance should be active. Can be either ALWAYS, NEVER or ON\_DEMAND. Default is ALWAYS. | `string` | `"ALWAYS"` | no |
| <a name="input_availability_type"></a> [availability\_type](#input\_availability\_type) | Availability type for the primary replica. Either `ZONAL` or `REGIONAL`. | `string` | `"ZONAL"` | no |
| <a name="input_backup_configuration"></a> [backup\_configuration](#input\_backup\_configuration) | Backup settings for primary instance. Will be automatically enabled if using MySQL with one or more replicas. | <pre>object({<br>    enabled                        = optional(bool, false)<br>    binary_log_enabled             = optional(bool, false)<br>    start_time                     = optional(string, "23:00")<br>    location                       = optional(string)<br>    log_retention_days             = optional(number, 7)<br>    point_in_time_recovery_enabled = optional(bool)<br>    retention_count                = optional(number, 7)<br>  })</pre> | <pre>{<br>  "binary_log_enabled": false,<br>  "enabled": false,<br>  "location": null,<br>  "log_retention_days": 7,<br>  "point_in_time_recovery_enabled": null,<br>  "retention_count": 7,<br>  "start_time": "23:00"<br>}</pre> | no |
| <a name="input_collation"></a> [collation](#input\_collation) | The name of server instance collation. | `string` | `null` | no |
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding Cloud SQL configuration values. | `string` | `"config"` | no |
| <a name="input_connector_enforcement"></a> [connector\_enforcement](#input\_connector\_enforcement) | Specifies if connections must use Cloud SQL connectors. | `string` | `null` | no |
| <a name="input_data_cache"></a> [data\_cache](#input\_data\_cache) | Enable data cache. Only used for Enterprise MYSQL and PostgreSQL. | `bool` | `false` | no |
| <a name="input_database_version"></a> [database\_version](#input\_database\_version) | Database type and version to create. e.g. 'MYSQL\_8\_0','SQLSERVER\_2017\_ENTERPRISE','POSTGRES\_15', | `string` | `"MYSQL_8_0"` | no |
| <a name="input_databases"></a> [databases](#input\_databases) | Databases to create once the primary instance is created. | `list(string)` | `null` | no |
| <a name="input_disk_autoresize_limit"></a> [disk\_autoresize\_limit](#input\_disk\_autoresize\_limit) | The maximum size to which storage capacity can be automatically increased. The default value is 0, which specifies that there is no limit. | `number` | `0` | no |
| <a name="input_disk_size"></a> [disk\_size](#input\_disk\_size) | Disk size in GB. Set to null to enable autoresize. | `number` | `null` | no |
| <a name="input_disk_type"></a> [disk\_type](#input\_disk\_type) | The type of data disk: `PD_SSD` or `PD_HDD`. | `string` | `"PD_SSD"` | no |
| <a name="input_edition"></a> [edition](#input\_edition) | The edition of the instance, can be ENTERPRISE or ENTERPRISE\_PLUS. | `string` | `"ENTERPRISE"` | no |
| <a name="input_encryption_key_name"></a> [encryption\_key\_name](#input\_encryption\_key\_name) | The full path to the encryption key used for the CMEK disk encryption of the primary instance. | `string` | `null` | no |
| <a name="input_flags"></a> [flags](#input\_flags) | Map FLAG\_NAME=>VALUE for database-specific tuning. | `map(string)` | `null` | no |
| <a name="input_gcp_deletion_protection"></a> [gcp\_deletion\_protection](#input\_gcp\_deletion\_protection) | Set Google's deletion protection attribute which applies across all surfaces (UI, API, & Terraform). | `bool` | `true` | no |
| <a name="input_insights_config"></a> [insights\_config](#input\_insights\_config) | Query Insights configuration. Defaults to null which disables Query Insights. | <pre>object({<br>    query_string_length     = optional(number, 1024)<br>    record_application_tags = optional(bool, false)<br>    record_client_address   = optional(bool, false)<br>    query_plans_per_minute  = optional(number, 5)<br>  })</pre> | `null` | no |
| <a name="input_labels"></a> [labels](#input\_labels) | Labels to be attached to all instances. | `map(string)` | `null` | no |
| <a name="input_maintenance_config"></a> [maintenance\_config](#input\_maintenance\_config) | Set maintenance window configuration and maintenance deny period (up to 90 days). Date format: 'yyyy-mm-dd'. | <pre>object({<br>    maintenance_window = optional(object({<br>      day          = number<br>      hour         = number<br>      update_track = optional(string, null)<br>    }), null)<br>    deny_maintenance_period = optional(object({<br>      start_date = string<br>      end_date   = string<br>      start_time = optional(string, "00:00:00")<br>    }), null)<br>  })</pre> | `{}` | no |
| <a name="input_prefix"></a> [prefix](#input\_prefix) | Optional prefix used to generate instance names. | `string` | `null` | no |
| <a name="input_replicas"></a> [replicas](#input\_replicas) | Map of NAME=> {REGION, KMS\_KEY} for additional read replicas. Set to null to disable replica creation. | <pre>map(object({<br>    region              = string<br>    encryption_key_name = optional(string)<br>  }))</pre> | `{}` | no |
| <a name="input_root_password"></a> [root\_password](#input\_root\_password) | Root password of the Cloud SQL instance. Required for MS SQL Server. | `string` | `null` | no |
| <a name="input_ssl"></a> [ssl](#input\_ssl) | Setting to enable SSL, set config and certificates. | <pre>object({<br>    client_certificates = optional(list(string))<br>    require_ssl         = optional(bool)<br>    # More details @ https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/sql_database_instance#ssl_mode<br>    ssl_mode = optional(string)<br>  })</pre> | `{}` | no |
| <a name="input_terraform_deletion_protection"></a> [terraform\_deletion\_protection](#input\_terraform\_deletion\_protection) | Prevent terraform from deleting instances. | `bool` | `true` | no |
| <a name="input_tier"></a> [tier](#input\_tier) | The machine type to use for the instances. | `string` | `"db-g1-small"` | no |
| <a name="input_time_zone"></a> [time\_zone](#input\_time\_zone) | The time\_zone to be used by the database engine (supported only for SQL Server), in SQL Server timezone format. | `string` | `null` | no |
| <a name="input_users"></a> [users](#input\_users) | Map of users to create in the primary instance (and replicated to other replicas). For MySQL, anything after the first `@` (if present) will be used as the user's host. Set PASSWORD to null if you want to get an autogenerated password. The user types available are: 'BUILT\_IN', 'CLOUD\_IAM\_USER' or 'CLOUD\_IAM\_SERVICE\_ACCOUNT'. | <pre>map(object({<br>    password = optional(string)<br>    type     = optional(string)<br>  }))</pre> | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cloudsql_instance_details"></a> [cloudsql\_instance\_details](#output\_cloudsql\_instance\_details) | Display Cloud SQL instance attributes, including name, project ID, region, connection name, IP address (public or private), and database version. |
<!-- END_TF_DOCS -->

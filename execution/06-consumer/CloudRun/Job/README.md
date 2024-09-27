## Introduction

[Cloud Run job](https://cloud.google.com/run/docs/overview/what-is-cloud-run#cloud-run-jobs) is a managed compute platform that lets you run containers directly on top of Google's scalable infrastructure.

If your code executes instructions and then stops (a script is a good example), you can use a Cloud Run job to run your code. You can execute a job from the command line using the gcloud CLI, schedule a recurring job, or run it as part of a workflow.

## Pre-Requisities

Before creating your first Cloud Run job instance, ensure you have completed the following prerequisites:

1. **Completed Prior Stages:** Successful creation of Cloud Run job resources depends on the completion of the following stages:
    * **01-organization:** This stage handles the activation of required Google Cloud APIs for Cloud Run.
    * **02-networking:** This stage sets up the necessary network infrastructure, such as VPCs and subnets, to support Cloud Run outgoing traffic connectivity.

2. **API Enablement:** Ensure the following Google Cloud APIs have been enabled:
   * IAM API (`iam.googleapis.com`)
   * Compute Enginer API (`compute.googleapis.com`)
   * Cloud Run API (`run.googleapis.com`)
   * Service Networking API (`servicenetworking.googleapis.com`)
   * Cloud Resource Manager API (`cloudresourcemanager.googleapis.com`)

3. **IAM Permissions:**  Grant yourself (or the appropriate users/service accounts) the following IAM roles at the project level (or higher):
   * Cloud Run Admin (`roles/run.admin`)

## Let's Get Started! 🚀

With the prerequisites in place and your Cloud Run configuration files ready, you can now leverage Terraform to automate the creation of your Cloud Run instances. Here's the workflow:

### Execution Steps

1. **Create your configuration files:**

    * Create YAML files defining the properties of each Cloud Run job instance you want to create. Ensure these files are stored in the `configuration/consumer/CloudRun/Job/config` folder within this repository.

    * Each YAML file should map to a single Cloud Run job instance providing details such as job name, region, and others configuration. Each field and its structure are described in the [input section](#inputs) below.

    * For reference on how to structure your Cloud Run job configuration YAML files, see the [example](#example) section below or refer to sample YAML file at folder location `configuration/consumer/CloudRun/Job/config/instance.yaml.example`. These examples provide templates that you can adapt to your specific needs.


2. **Initialize Terraform:**

    * Open your terminal and navigate to the Cloud Run job directory containing the Terraform configuration.

    * Run the following command to initialize Terraform:

    ```
    terraform init
    ```
3. **Review the Execution Plan:**

    * Use the terraform plan command to generate an execution plan. This will show you the changes Terraform will make to your Google Cloud infrastructure:

    ```
    terraform plan -var-file=../../../../configuration/consumer/CloudRun/Job/cloudrunjob.tfvars
    ```

Carefully review the plan to ensure it aligns with your intended configuration.

4. **Apply the Configuration:**

    Once you're satisfied with the plan, execute the terraform apply command to provision your Cloud Run job instances:

    ```
    terraform apply -var-file=../../../../configuration/consumer/CloudRun/Job/cloudrunjob.tfvars
    ```

Terraform will read the YAML files from the `configuration/consumer/CloudRun/Job/config/` folder and create the corresponding Cloud Run instances in your Google Cloud project.

5. **Monitor and Manage:**
    * After the instances are created, you can monitor their status, performance, and logs through the Google Cloud Console or using the Google Cloud CLI.

    * Use Terraform to manage updates and changes to your Cloud Run job instances as needed.

### Example

To help you get started, we've provided examples of YAML configuration files that you can use as templates for your Cloud Run job instances.

1. Minimal YAML for Serverless connector (Mandatory Feilds Only) :
    ```
    project_id: <your-project-id>
    region: us-central1
    name: job1
    containers:
      hello-container:
        image: us-docker.pkg.dev/cloudrun/container/hello
    ```

2. Using an existing VPC Access Connector to connect to a VPC from Cloud Run.
    ```
    project_id: <your-project-id>
    region: us-central1
    name: job2
    containers:
      hello-container:
        image: us-docker.pkg.dev/cloudrun/container/hello
    revision:
      vpc_access:
        connector : projects/<your-project-id>/locations/us-central1/connectors/<your-connector-id>
        egress : ALL_TRAFFIC
    ```

  3. Creating a new VPC Access Connector to connect to a VPC from Cloud Run.
      ```
      project_id: <your-project-id>
      region: us-central1
      name: job3
      containers:
        hello-container:
          image: us-docker.pkg.dev/cloudrun/container/hello
      vpc_connector_create:
        ip_cidr_range : 10.10.60.0/28
        network: <your-network-id>
      ```


**Beta Feature :** To use beta features like **Direct VPC Egress**, set the launch stage to a preview stage.

4. Creating Cloud Run job using Direct VPC Egress:
    ```
    project_id: <your-project-id>
    region: us-central1
    name: job4
    launch_stage: BETA
    containers:
      container-name:
        image: us-docker.pkg.dev/cloudrun/container/hello
    revision:
      vpc_access:
        egress: ALL_TRAFFIC
        network: <your-network-id>
        subnetwork: subnet1
    ```

## Important Notes:

* This README is a starting point. Customize it to include specific details about your Cloud Run jobs and relevant configuration.
Refer to the official Google Cloud Cloud Run [documentation](https://cloud.google.com/run/docs/overview/what-is-cloud-run) for the most up-to-date information and best practices.

* Order of Execution: Make sure to complete the 00-bootstrap, 01-organization, 02-networking stages before attempting to create Cloud Run instances. Terraform will leverage the resources and configurations established in these prior stages.

* Troubleshooting: If you encounter errors during the Cloud Run jobs creation process, verify that all prerequisites are satisfied and that the dependencies between stages are correctly configured.

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cloud_run_job"></a> [cloud\_run\_job](#module\_cloud\_run\_job) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/cloud-run-v2 | v34.1.0 |


## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| project\_id | The project ID to deploy to | `string` | n/a | yes |
| region | Cloud Run job deployment region | `string` | n/a | yes |
| name | The name of the Cloud Run job to create | `string` | n/a | yes |
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding Cloud Run job configuration values. | `string` | `"../../../../configuration/consumer/CloudRun/Job/config"` | no |
| <a name="input_containers"></a> [containers](#input\_containers) | Containers in name => attributes format. | <pre>map(object({<br>    image   = string<br>    command = optional(list(string))<br>    args    = optional(list(string))<br>    env     = optional(map(string))<br>    env_from_key = optional(map(object({<br>      secret  = string<br>      version = string<br>    })))<br>    liveness_probe = optional(object({<br>      grpc = optional(object({<br>        port    = optional(number)<br>        service = optional(string)<br>      }))<br>      http_get = optional(object({<br>        http_headers = optional(map(string))<br>        path         = optional(string)<br>      }))<br>      failure_threshold     = optional(number)<br>      initial_delay_seconds = optional(number)<br>      period_seconds        = optional(number)<br>      timeout_seconds       = optional(number)<br>    }))<br>    ports = optional(map(object({<br>      container_port = optional(number)<br>      name           = optional(string)<br>    })))<br>    resources = optional(object({<br>      limits = optional(object({<br>        cpu    = string<br>        memory = string<br>      }))<br>      cpu_idle          = optional(bool)<br>      startup_cpu_boost = optional(bool)<br>    }))<br>    startup_probe = optional(object({<br>      grpc = optional(object({<br>        port    = optional(number)<br>        service = optional(string)<br>      }))<br>      http_get = optional(object({<br>        http_headers = optional(map(string))<br>        path         = optional(string)<br>      }))<br>      tcp_socket = optional(object({<br>        port = optional(number)<br>      }))<br>      failure_threshold     = optional(number)<br>      initial_delay_seconds = optional(number)<br>      period_seconds        = optional(number)<br>      timeout_seconds       = optional(number)<br>    }))<br>    volume_mounts = optional(map(string))<br>  }))</pre> | `n/a` | yes |
| <a name="input_create_job"></a> [create\_job](#input\_create\_job) | Create Cloud Run job instead of Service. | `bool` | `true` | no |
| <a name="input_custom_audiences"></a> [custom\_audiences](#input\_custom\_audiences) | Custom audiences for service. | `list(string)` | `null` | no |
| <a name="input_encryption_key"></a> [encryption\_key](#input\_encryption\_key) | The full resource name of the Cloud KMS CryptoKey. | `string` | `null` | no |
| <a name="input_eventarc_triggers"></a> [eventarc\_triggers](#input\_eventarc\_triggers) | Event arc triggers for different sources. | <pre>object({<br>    audit_log = optional(map(object({<br>      method  = string<br>      service = string<br>    })))<br>    pubsub                 = optional(map(string))<br>    service_account_email  = optional(string)<br>    service_account_create = optional(bool, false)<br>  })</pre> | `{}` | no |
| <a name="input_iam"></a> [iam](#input\_iam) | IAM bindings for Cloud Run service in {ROLE => [MEMBERS]} format. | `map(list(string))` | `{}` | no |
| <a name="input_ingress"></a> [ingress](#input\_ingress) | Ingress settings. | `string` | `null` | no |
| <a name="input_labels"></a> [labels](#input\_labels) | Resource labels. | `map(string)` | `{}` | no |
| <a name="input_launch_stage"></a> [launch\_stage](#input\_launch\_stage) | The launch stage as defined by Google Cloud Platform Launch Stages. | `string` | `null` | no |
| <a name="input_prefix"></a> [prefix](#input\_prefix) | Optional prefix used for resource names. | `string` | `null` | no |
| <a name="input_revision"></a> [revision](#input\_revision) | Revision template configurations. | <pre>object({<br>    name                       = optional(string)<br>    gen2_execution_environment = optional(bool)<br>    max_concurrency            = optional(number)<br>    max_instance_count         = optional(number)<br>    min_instance_count         = optional(number)<br>    vpc_access = optional(object({<br>      connector = optional(string)<br>      egress    = optional(string)<br>      subnet    = optional(string)<br>      tags      = optional(list(string))<br>    }))<br>    timeout = optional(string)<br>  })</pre> | `{}` | no |
| <a name="input_service_account"></a> [service\_account](#input\_service\_account) | Service account email. Unused if service account is auto-created. | `string` | `null` | no |
| <a name="input_service_account_create"></a> [service\_account\_create](#input\_service\_account\_create) | Auto-create service account. | `bool` | `false` | no |
| <a name="input_tag_bindings"></a> [tag\_bindings](#input\_tag\_bindings) | Tag bindings for this service, in key => tag value id format. | `map(string)` | `{}` | no |
| <a name="input_volumes"></a> [volumes](#input\_volumes) | Named volumes in containers in name => attributes format. | <pre>map(object({<br>    secret = optional(object({<br>      name         = string<br>      default_mode = optional(string)<br>      path         = optional(string)<br>      version      = optional(string)<br>      mode         = optional(string)<br>    }))<br>    cloud_sql_instances = optional(list(string))<br>    empty_dir_size      = optional(string)<br>  }))</pre> | `{}` | no |
| <a name="input_vpc_connector_create"></a> [vpc\_connector\_create](#input\_vpc\_connector\_create) | Populate this to create a Serverless VPC Access connector. | <pre>object({<br>    ip_cidr_range = optional(string)<br>    machine_type  = optional(string)<br>    name          = optional(string)<br>    network       = optional(string)<br>    instances = optional(object({<br>      max = optional(number)<br>      min = optional(number)<br>      }), {}<br>    )<br>    throughput = optional(object({<br>      max = optional(number)<br>      min = optional(number)<br>      }), {}<br>    )<br>    subnet = optional(object({<br>      name       = optional(string)<br>      project_id = optional(string)<br>    }), {})<br>  })</pre> | `null` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cloud_run_job_details"></a> [cloud\_run\_job\_details](#output\_cloud\_run\_job\_details) | Display the details of the cloud run jobs. |
<!-- END_TF_DOCS -->

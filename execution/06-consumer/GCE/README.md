# Google Compute Engine

## Overview

This Terraform solution provides a streamlined way to deploy and manage Google Compute Engine (GCE) virtual machines (VMs) using YAML configuration files.

This solution utilizes a modular approach, with the gce.tf file defining a Terraform module that leverages the compute-vm module from the Google Cloud Foundation Fabric. The module encapsulates the logic for creating and configuring GCE VMs based on parameters provided in YAML files.

## Pre-Requisites

### Prior Step Completion :

- **Completed Prior Stages:** Successful deployment of GCE resources depends on the completion of the following stages:

    - **01-organization:** This stage handles the activation of required Google Cloud APIs for GCE.
    - **02-networking:** This stage sets up the necessary network infrastructure, such as VPCs and subnets, to support GCE connectivity.
    - **03-security/GCE:** This stage configures firewall rules to allow access to GCE instances on the appropriate ports and IP ranges.
    - **04-producer** : Optional however highly recommended for you to create a producer service such as CloudSQL, MRC or AlloyDB before you create a GCE instance.


### Enabled APIs:

Ensure the following Google Cloud APIs are enabled in your project:

- Compute Engine API
- Cloud Resource Manager API (for tag bindings)

### Permissions:

The user or service account executing Terraform must have the following roles (or equivalent permissions):

- Compute Admin (for managing VMs)
- Service Account User (if using service accounts)
- Tag Admin (if using tag bindings)

## Execution Steps

1. **Configuration** :

    - Create YAML configuration files (e.g., instance1.yaml, instance2.yaml) in the configs directory (or the path specified in the config_folder_path variable).
    - Edit the YAML files to specify the desired VM configurations. (See **Examples** below)

2. **Terraform Initialization**:

    - Open your terminal and navigate to the GCE directory containing the Terraform configuration.
    - Run the following command to initialize Terraform:

    ```
    terraform init
    ```
3. **Review the Execution Plan:**

    - Use the following command to generate an execution plan. This will show you the changes Terraform will make to your Google Cloud infrastructure:

    ```
    terraform plan
    ```

Carefully review the plan to ensure it aligns with your intended configuration.

4. **Apply the Configuration:**

    Once you're satisfied with the plan, execute the terraform apply command to provision your GCE instances:

    ```
    terraform apply
    ```

Terraform will read the YAML files from the `06-consumer/GCE/configs` folder by default and create the corresponding GCE instances in your Google Cloud project.

5. **Monitor and Manage:**

    * After the instances are created, you can monitor their status, performance, and logs through the Google Cloud Console or using the Google Cloud CLI.

    * Use Terraform to manage updates and changes to your AlloyDB instances as needed.

### Examples

- instance1.yaml : This is a sample YAML file that defines the configuration for a GCE instance. The file specifies the name, project ID, region, zone, image, and network of the instance.

  ```
  name: instance1
  project_id: <project-id>
  region: us-central1
  zone: us-central1-a
  image: ubuntu-os-cloud/ubuntu-2204-lts
  network: projects/<project-id>/global/networks/<network-name>
  subnetwork: projects/<project-id>/regions/us-central1/subnetworks/<subnetwork-name>
  ```

- instance2.yaml : This is another sample YAML file that defines the configuration for a GCE instance. The file specifies the name, project ID, region, zone, image, network, and instance type of the instance.

  ```
  name: instance1
  project_id: <project-id>
  region: us-central1
  zone: us-central1-a
  image: ubuntu-os-cloud/ubuntu-2204-lts
  network: projects/<project-id>/global/networks/<network-name>
  subnetwork: projects/<project-id>/regions/us-central1/subnetworks/<subnetwork-name>
  instance_type: e2-medium
  ```

- instance3.yaml : This is another sample YAML file that defines the configuration for a GCE instance. The file specifies the name, project ID, region, zone, image, network, and instance type of the instance with a startup script.

  ```
  name: instance3
  project_id: <project-id>
  region: us-central1
  zone: us-central1-a
  image: ubuntu-os-cloud/ubuntu-2204-lts
  network: projects/<project-id>/global/networks/<network-name>
  subnetwork: projects/<project-id>/regions/us-central1/subnetworks/<subnetwork-name>
  metadata:
    "startup_script": |
        #!/bin/bash
        echo Hello World
  ```

## Important Notes

- The solution assumes that the required network and subnetwork already exist in your project as you've gone through the previous steps.
- Carefully review and customize the YAML configuration files to match your specific requirements (e.g., project ID, region, zone, image, instance type, etc.).
- Be sure to provide the correct service account credentials (if applicable) to allow Terraform to interact with your Google Cloud project.
- Refer to the variables.tf file for a complete list of available variables and their descriptions.
- The Terraform module used in this solution (cloud-foundation-fabric/modules/compute-vm) might have additional configuration options and capabilities. Refer to its [documentation](https://github.com/GoogleCloudPlatform/cloud-foundation-fabric/tree/master/modules/compute-vm) for further customization.
- Remember to replace placeholders ( \<project-id>, \<network-name>, \<subnetwork-name>) with your actual values in the YAML files and Terraform configuration.



## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_vm"></a> [vm](#module\_vm) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/compute-vm | n/a |


## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_attached_disk_defaults"></a> [attached\_disk\_defaults](#input\_attached\_disk\_defaults) | Defaults for attached disks options. | <pre>object({<br>    auto_delete  = optional(bool, false)<br>    mode         = string<br>    replica_zone = string<br>    type         = string<br>  })</pre> | <pre>{<br>  "auto_delete": true,<br>  "mode": "READ_WRITE",<br>  "replica_zone": null,<br>  "type": "pd-balanced"<br>}</pre> | no |
| <a name="input_attached_disks"></a> [attached\_disks](#input\_attached\_disks) | Additional disks, if options is null defaults will be used in its place. Source type is one of 'image' (zonal disks in vms and template), 'snapshot' (vm), 'existing', and null. | <pre>list(object({<br>    name        = string<br>    device_name = optional(string)<br>    # TODO: size can be null when source_type is attach<br>    size              = string<br>    snapshot_schedule = optional(string)<br>    source            = optional(string)<br>    source_type       = optional(string)<br>    options = optional(<br>      object({<br>        auto_delete  = optional(bool, false)<br>        mode         = optional(string, "READ_WRITE")<br>        replica_zone = optional(string)<br>        type         = optional(string, "pd-balanced")<br>      }),<br>      {<br>        auto_delete  = true<br>        mode         = "READ_WRITE"<br>        replica_zone = null<br>        type         = "pd-balanced"<br>      }<br>    )<br>  }))</pre> | `[]` | no |
| <a name="input_boot_disk"></a> [boot\_disk](#input\_boot\_disk) | Boot disk properties. | <pre>object({<br>    auto_delete       = optional(bool, true)<br>    snapshot_schedule = optional(string)<br>    source            = optional(string)<br>    initialize_params = optional(object({<br>      image = optional(string, "projects/debian-cloud/global/images/family/debian-11")<br>      size  = optional(number, 10)<br>      type  = optional(string, "pd-balanced")<br>    }))<br>    use_independent_disk = optional(bool, false)<br>  })</pre> | <pre>{<br>  "initialize_params": {}<br>}</pre> | no |
| <a name="input_can_ip_forward"></a> [can\_ip\_forward](#input\_can\_ip\_forward) | Enable IP forwarding. | `bool` | `false` | no |
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding GCE configuration values. | `string` | `"./configs"` | no |
| <a name="input_create_template"></a> [create\_template](#input\_create\_template) | Create instance template instead of instances. | `bool` | `false` | no |
| <a name="input_description"></a> [description](#input\_description) | Description of a Compute Instance. | `string` | `"Managed by the compute-vm Terraform module."` | no |
| <a name="input_enable_display"></a> [enable\_display](#input\_enable\_display) | Enable virtual display on the instances. | `bool` | `false` | no |
| <a name="input_encryption"></a> [encryption](#input\_encryption) | Encryption options. Only one of kms\_key\_self\_link and disk\_encryption\_key\_raw may be set. If needed, you can specify to encrypt or not the boot disk. | <pre>object({<br>    encrypt_boot            = optional(bool, false)<br>    disk_encryption_key_raw = optional(string)<br>    kms_key_self_link       = optional(string)<br>  })</pre> | `null` | no |
| <a name="input_group"></a> [group](#input\_group) | Define this variable to create an instance group for instances. Disabled for template use. | <pre>object({<br>    named_ports = map(number)<br>  })</pre> | `null` | no |
| <a name="input_hostname"></a> [hostname](#input\_hostname) | Instance FQDN name. | `string` | `null` | no |
| <a name="input_iam"></a> [iam](#input\_iam) | IAM bindings in {ROLE => [MEMBERS]} format. | `map(list(string))` | `{}` | no |
| <a name="input_instance_schedule"></a> [instance\_schedule](#input\_instance\_schedule) | Assign or create and assign an instance schedule policy. Either resource policy id or create\_config must be specified if not null. Set active to null to dtach a policy from vm before destroying. | <pre>object({<br>    resource_policy_id = optional(string)<br>    create_config = optional(object({<br>      active          = optional(bool, true)<br>      description     = optional(string)<br>      expiration_time = optional(string)<br>      start_time      = optional(string)<br>      timezone        = optional(string, "UTC")<br>      vm_start        = optional(string)<br>      vm_stop         = optional(string)<br>    }))<br>  })</pre> | `null` | no |
| <a name="input_instance_type"></a> [instance\_type](#input\_instance\_type) | Instance type. | `string` | `"f1-micro"` | no |
| <a name="input_labels"></a> [labels](#input\_labels) | Instance labels. | `map(string)` | `{}` | no |
| <a name="input_metadata"></a> [metadata](#input\_metadata) | Instance metadata. | `map(string)` | `{}` | no |
| <a name="input_min_cpu_platform"></a> [min\_cpu\_platform](#input\_min\_cpu\_platform) | Minimum CPU platform. | `string` | `null` | no |
| <a name="input_network_attached_interfaces"></a> [network\_attached\_interfaces](#input\_network\_attached\_interfaces) | Network interfaces using network attachments. | `list(string)` | `[]` | no |
| <a name="input_options"></a> [options](#input\_options) | Instance options. | <pre>object({<br>    allow_stopping_for_update = optional(bool, true)<br>    deletion_protection       = optional(bool, false)<br>    node_affinities = optional(map(object({<br>      values = list(string)<br>      in     = optional(bool, true)<br>    })), {})<br>    spot               = optional(bool, false)<br>    termination_action = optional(string)<br>  })</pre> | <pre>{<br>  "allow_stopping_for_update": true,<br>  "deletion_protection": false,<br>  "spot": false,<br>  "termination_action": null<br>}</pre> | no |
| <a name="input_scratch_disks"></a> [scratch\_disks](#input\_scratch\_disks) | Scratch disks configuration. | <pre>object({<br>    count     = number<br>    interface = string<br>  })</pre> | <pre>{<br>  "count": 0,<br>  "interface": "NVME"<br>}</pre> | no |
| <a name="input_service_account"></a> [service\_account](#input\_service\_account) | Service account email and scopes. If email is null, the default Compute service account will be used unless auto\_create is true, in which case a service account will be created. Set the variable to null to avoid attaching a service account. | <pre>object({<br>    auto_create = optional(bool, false)<br>    email       = optional(string)<br>    scopes      = optional(list(string))<br>  })</pre> | `{}` | no |
| <a name="input_shielded_config"></a> [shielded\_config](#input\_shielded\_config) | Shielded VM configuration of the instances. | <pre>object({<br>    enable_secure_boot          = bool<br>    enable_vtpm                 = bool<br>    enable_integrity_monitoring = bool<br>  })</pre> | `null` | no |
| <a name="input_snapshot_schedules"></a> [snapshot\_schedules](#input\_snapshot\_schedules) | Snapshot schedule resource policies that can be attached to disks. | <pre>map(object({<br>    schedule = object({<br>      daily = optional(object({<br>        days_in_cycle = number<br>        start_time    = string<br>      }))<br>      hourly = optional(object({<br>        hours_in_cycle = number<br>        start_time     = string<br>      }))<br>      weekly = optional(list(object({<br>        day        = string<br>        start_time = string<br>      })))<br>    })<br>    description = optional(string)<br>    retention_policy = optional(object({<br>      max_retention_days         = number<br>      on_source_disk_delete_keep = optional(bool)<br>    }))<br>    snapshot_properties = optional(object({<br>      chain_name        = optional(string)<br>      guest_flush       = optional(bool)<br>      labels            = optional(map(string))<br>      storage_locations = optional(list(string))<br>    }))<br>  }))</pre> | `{}` | no |
| <a name="input_tag_bindings"></a> [tag\_bindings](#input\_tag\_bindings) | Resource manager tag bindings for this instance, in tag key => tag value format. | `map(string)` | `null` | no |
| <a name="input_tag_bindings_firewall"></a> [tag\_bindings\_firewall](#input\_tag\_bindings\_firewall) | Firewall (network scoped) tag bindings for this instance, in tag key => tag value format. | `map(string)` | `null` | no |
| <a name="input_tags"></a> [tags](#input\_tags) | Instance network tags for firewall rule targets. | `list(string)` | `[]` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_external_ip"></a> [external\_ip](#output\_external\_ip) | Instance main interface external IP addresses. |
| <a name="output_id"></a> [id](#output\_id) | Fully qualified instance id. |
| <a name="output_instances_self_links"></a> [instances\_self\_links](#output\_instances\_self\_links) | List of self-links for compute instances |
| <a name="output_internal_ip"></a> [internal\_ip](#output\_internal\_ip) | Instance main interface internal IP address. |
| <a name="output_internal_ips"></a> [internal\_ips](#output\_internal\_ips) | Instance interfaces internal IP addresses. |
| <a name="output_vm_instances"></a> [vm\_instances](#output\_vm\_instances) | Map of VM instance information |

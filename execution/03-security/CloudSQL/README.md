## Introduction

Secure your Cloud SQL databases with Google Cloud Firewall. This guide provides step-by-step instructions on configuring firewall rules to control access and protect your valuable data.

### Execution Steps

1. **Create your configuration .tfvars files:**

    * Create `cloudsql.tfvars` file defining the values for ingress rules and egress rules. Ensure these files are stored in the `configuration/security/cloudsql` folder.

    * For reference on how to structure your `cloudsql.tfvars` file , see the [example](#example) section below or refer to sample `terraform.tfvars.example` file . Each field and its structure is described in the [input section](#inputs) below.


2. **Initialize Terraform:**

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
    terraform apply -var-file="../../../configuration/security/cloudsql.tfvars"
    ```

### Example

```
project_id = "<your-project-id>"
network    = "projects/<your-project-id>/global/networks/<your-vpc-name>"
egress_rules = {
  allow-egress-cloudsql = {
    deny = false
    rules = [{
      protocol = "tcp"
      ports    = ["3306"]
    }]
  }
}
```

<!-- BEGIN_TF_DOCS -->

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cloudsql_firewall"></a> [cloudsql\_firewall](#module\_cloudsql\_firewall) | github.com/GoogleCloudPlatform/cloud-foundation-fabric//modules/net-vpc-firewall | v30.0.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_network"></a> [network](#input\_network) | Name of the VPC network or fully qualified network id. | `string` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | The ID of the google cloud project where this compute instance will be created. | `string` | n/a | yes |
| <a name="input_default_rules_config"></a> [default\_rules\_config](#input\_default\_rules\_config) | Optionally created convenience rules. Set the 'disabled' attribute to true, or individual rule attributes to empty lists to disable. | <pre>object({<br>    admin_ranges = optional(list(string))<br>    disabled     = optional(bool, true)<br>    http_ranges = optional(list(string), [<br>      "35.191.0.0/16", "130.211.0.0/22", "209.85.152.0/22", "209.85.204.0/22"]<br>    )<br>    http_tags = optional(list(string), ["http-server"])<br>    https_ranges = optional(list(string), [<br>      "35.191.0.0/16", "130.211.0.0/22", "209.85.152.0/22", "209.85.204.0/22"]<br>    )<br>    https_tags = optional(list(string), ["https-server"])<br>    ssh_ranges = optional(list(string), ["35.235.240.0/20"])<br>    ssh_tags   = optional(list(string), ["ssh"])<br>  })</pre> | <pre>{<br>  "disabled": true<br>}</pre> | no |
| <a name="input_egress_rules"></a> [egress\_rules](#input\_egress\_rules) | List of egress rule definitions, default to deny action. Null destination ranges will be replaced with 0/0. | <pre>map(object({<br>    deny               = optional(bool, true)<br>    description        = optional(string)<br>    destination_ranges = optional(list(string))<br>    disabled           = optional(bool, false)<br>    enable_logging = optional(object({<br>      include_metadata = optional(bool)<br>    }))<br>    priority             = optional(number, 1000)<br>    source_ranges        = optional(list(string))<br>    targets              = optional(list(string))<br>    use_service_accounts = optional(bool, false)<br>    rules = optional(list(object({<br>      protocol = string<br>      ports    = optional(list(string))<br>    })), [{ protocol = "all" }])<br>  }))</pre> | `{}` | no |
| <a name="input_ingress_rules"></a> [ingress\_rules](#input\_ingress\_rules) | List of ingress rule definitions, default to allow action. Null source ranges will be replaced with 0/0. | <pre>map(object({<br>    deny               = optional(bool, false)<br>    description        = optional(string)<br>    destination_ranges = optional(list(string), []) # empty list is needed as default to allow deletion after initial creation with a value. See https://github.com/hashicorp/terraform-provider-google/issues/14270<br>    disabled           = optional(bool, false)<br>    enable_logging = optional(object({<br>      include_metadata = optional(bool)<br>    }))<br>    priority             = optional(number, 1000)<br>    source_ranges        = optional(list(string))<br>    sources              = optional(list(string))<br>    targets              = optional(list(string))<br>    use_service_accounts = optional(bool, false)<br>    rules = optional(list(object({<br>      protocol = string<br>      ports    = optional(list(string))<br>    })), [{ protocol = "all" }])<br>  }))</pre> | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cloudsql_firewall_rules"></a> [cloudsql\_firewall\_rules](#output\_cloudsql\_firewall\_rules) | Map of firewall rules created. |
<!-- END_TF_DOCS -->

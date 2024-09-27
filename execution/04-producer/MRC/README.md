# Memorystore for Redis Cluster (MRC)

## Introduction

This Terraform module simplifies the creation and management of Memorystore for Redis Cluster (MRC) instances on Google Cloud Platform (GCP). MRC provides a fully managed, scalable, and highly available Redis service, ideal for caching, session management, and real-time analytics.

## Pre-Requisites

Before creating your first MRC instance, ensure you have completed the following prerequisites:

1. **Completed Prior Stages:** Successful deployment of MRC resources depends on the completion of the following stages:
    * **01-organization:** This stage handles the activation of required Google Cloud APIs for MRC.
    * **02-networking:** This stage sets up the necessary network infrastructure, such as VPCs and subnets, to support MRC connectivity. For MRC, ensure that you create a service connection policy. For the same, in networking tfvars configure this section carefully :

    ```
    create_scp_policy      = true
    subnets_for_scp_policy = ["subnet-name"] # List subnets here from the same region as the SCP
    ```
    * **03-security/MRC:** This stage configures firewall rules to allow access to MRC instances on the appropriate ports and IP ranges.

2. **API Enablement**: Ensure the following Google Cloud APIs have been enabled:

    * Compute Engine API (compute.googleapis.com)
    * Redis API (redis.googleapis.com)
    * Cloud Resource Manager API (cloudresourcemanager.googleapis.com)

3. **IAM Permissions**:  Grant yourself (or the appropriate users/service accounts) the following IAM roles at the project level (or higher):

    * Redis Admin (roles/redis.admin)

## Let's Get Started! 🚀

With the prerequisites in place and your MRC configuration files ready, you can now leverage Terraform to automate the creation of your MRC instances. Here's the workflow:

### Execution Steps

1. Create your configuration files:

Create YAML files defining the properties of each MRC instance you want to create. Ensure these files are stored in the config folder under the configuration/producer/MRC/config folder.

Each YAML file should map to a single MRC instance, providing details such as instance name, project ID, region, shard count, replica count, and network ID. Each field and its structure are described in the input section below.

For reference on how to structure your MRC configuration YAML files, see the example section below.

2. Initialize Terraform:

Open your terminal and navigate to the directory containing the Terraform configuration.

Run the following command to initialize Terraform:

```
terraform init
```

3. Review the Execution Plan:

Use the terraform plan command to generate an execution plan. This will show you the changes Terraform will make to your Google Cloud infrastructure:

```
terraform plan -var-file=../../../configuration/producer/MRC/mrc.tfvars
```

4. Apply the Configuration:

Once you're satisfied with the plan, execute the terraform apply command to provision your MRC instances:

```
terraform apply -var-file=../../../configuration/producer/MRC/mrc.tfvars
```

5. Monitor and Manage:

After the instances are created, you can monitor their status, performance, and logs through the Google Cloud Console or using the Google Cloud CLI. Use Terraform to manage updates and changes to your MRC instances as needed.

## Example YAML Configuration

```
redis_cluster_name: my-redis-cluster
project_id: test-project
region: us-central1
shard_count: 3
replica_count: 1
network_id: projects/test-project/global/networks/network-name
```

## Important Notes

- Refer to the official Google Cloud Memorystore for Redis documentation for the most up-to-date information and best practices: https://cloud.google.com/memorystore/docs/redis
- Order of Execution: Make sure to complete the necessary networking stages before attempting to create MRC instances. Terraform will leverage the resources and configurations established in these prior stages.
- Troubleshooting: If you encounter errors during the MRC creation process, verify that all prerequisites are satisfied and that the dependencies between stages are correctly configured.

<!-- BEGIN_TF_DOCS -->
## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | 5.32.0 |


## Resources

| Name | Type |
|------|------|
| [google_redis_cluster.cluster-ha](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/redis_cluster) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding MRC configuration values. | `string` | `"../../../configuration/producer/MRC/config"` | no |
| <a name="input_deletion_protection_enabled"></a> [deletion\_protection\_enabled](#input\_deletion\_protection\_enabled) | Indicates if the cluster is deletion protected or not. If the value if set to true, any delete cluster operation will fail. Default value is true. | `bool` | `true` | no |
| <a name="input_region"></a> [region](#input\_region) | The region in which to create the Redis cluster. | `string` | `"us-central1"` | no |
| <a name="input_replica_count"></a> [replica\_count](#input\_replica\_count) | Number of replicas per shard in the Redis cluster. | `number` | `1` | no |
| <a name="input_shard_count"></a> [shard\_count](#input\_shard\_count) | Number of shards (replicas) in the Redis cluster. | `number` | `3` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_redis_cluster_details"></a> [redis\_cluster\_details](#output\_redis\_cluster\_details) | Detailed information about each Redis cluster |
<!-- END_TF_DOCS -->

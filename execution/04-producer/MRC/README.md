## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | 5.32.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_redis_cluster.cluster-ha](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/redis_cluster) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_config_folder_path"></a> [config\_folder\_path](#input\_config\_folder\_path) | Location of YAML files holding MRC configuration values. | `string` | `"./config"` | no |
| <a name="input_region"></a> [region](#input\_region) | The region in which to create the Redis cluster. | `string` | `"us-central1"` | no |
| <a name="input_replica_count"></a> [replica\_count](#input\_replica\_count) | Number of replicas per shard in the Redis cluster. | `number` | `1` | no |
| <a name="input_shard_count"></a> [shard\_count](#input\_shard\_count) | Number of shards (replicas) in the Redis cluster. | `number` | `3` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_redis_cluster_details"></a> [redis\_cluster\_details](#output\_redis\_cluster\_details) | Detailed information about each Redis cluster |
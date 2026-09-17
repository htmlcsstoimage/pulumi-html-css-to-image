# StorageDestination

Configure HCTI to write images to an existing bucket or Space. This resource creates neither buckets nor IAM roles, and deleting it does not delete stored files. Supply exactly one provider object inside `connectionInfo`: `awsS3`, `cloudflareR2`, `backblazeB2`, `digitaloceanSpaces`, `wasabi`, `googleCloudStorage`, or `otherS3Compatible`. Required fields are specific to that object.

**Resource token:** `html-css-to-image:index:StorageDestination`

See [provider setup](../../README.md#development) and the [complete YAML example](../../examples/storage-destination/Pulumi.yaml). The full lifecycle needs `storage_destinations:read`, `storage_destinations:create_update`, and `storage_destinations:delete`, plus permission to use referenced resources.

AWS uses an existing role ARN; the other providers use access keys. On updates, omitting `secretAccessKey` retains the existing secret when provider and access-key identity are unchanged. Changing either requires a supplied secret. Credentials are secret outputs; import cannot recover them. See the [AWS example](../../examples/aws-storage/Pulumi.yaml) and [external-ID lookup](../functions/get_aws_storage_external_id.md) for role trust setup.

## Inputs

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| `connectionInfo` | [StorageDestinationConnectionInfo](#nested-objects) | Yes | Set exactly one provider object: aws_s3, cloudflare_r2, backblaze_b2, digitalocean_spaces, wasabi, google_cloud_storage, or other_s3_compatible. Switching objects updates in place. |
| `disabled` | Boolean | No | Disable this destination for new renders. Defaults to false. |
| `hctiStorageDisabled` | Boolean | No | Store output only at this destination. Defaults to false. True requires authenticated /store rendering and supplies no public HCTI CDN URL. |
| `name` | String | Yes | Display name, 3–255 characters after trimming. |

## Additional computed attributes

All inputs above are also readable resource outputs. Omitted nullable rendering inputs stay unset; they do not report effective defaults. The following are additional computed outputs.

| Attribute | Type | Description |
| --- | --- | --- |
| `id` | String | Resource identity used for import. |
| `createdAt` | String | Creation timestamp. |
| `lastTestError` | String | Latest API connection test error, or null. Marked sensitive because it is external service diagnostic text. |
| `lastTestSucceeded` | Boolean | Latest API connection test result, or null. |
| `lastTestedAt` | String | Timestamp of the latest API connection test, or null. Refresh does not run tests. |
| `updatedAt` | String | Last update timestamp. |

## Nested objects

### StorageDestinationConnectionInfo

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `awsS3` | [StorageDestinationConnectionInfoAwsS3](#nested-objects) | No | Connection settings for aws_s3. Set exactly one provider object inside connection_info.  |
| `backblazeB2` | [StorageDestinationConnectionInfoBackblazeB2](#nested-objects) | No | Connection settings for backblaze_b2. Set exactly one provider object inside connection_info.  |
| `cloudflareR2` | [StorageDestinationConnectionInfoCloudflareR2](#nested-objects) | No | Connection settings for cloudflare_r2. Set exactly one provider object inside connection_info.  |
| `digitaloceanSpaces` | [StorageDestinationConnectionInfoDigitaloceanSpaces](#nested-objects) | No | Connection settings for digitalocean_spaces. Set exactly one provider object inside connection_info.  |
| `googleCloudStorage` | [StorageDestinationConnectionInfoGoogleCloudStorage](#nested-objects) | No | Connection settings for google_cloud_storage. Set exactly one provider object inside connection_info.  |
| `otherS3Compatible` | [StorageDestinationConnectionInfoOtherS3Compatible](#nested-objects) | No | Connection settings for other_s3_compatible. Set exactly one provider object inside connection_info.  |
| `wasabi` | [StorageDestinationConnectionInfoWasabi](#nested-objects) | No | Connection settings for wasabi. Set exactly one provider object inside connection_info.  |

### StorageDestinationConnectionInfoAwsS3

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `bucket` | String | Yes | Existing bucket or Space name, up to 255 characters. The provider never creates or deletes the bucket.  |
| `keyPrefix` | String | No | Object key prefix, up to 1024 characters. The API trims surrounding whitespace and slashes; omission uses the bucket root.  |
| `region` | String | Yes | Bucket region supported by aws_s3.  |
| `roleArn` | String | Yes | AWS only. IAM role ARN, up to 2048 characters. Configure the organization external ID in its trust policy before saving.  |

### StorageDestinationConnectionInfoBackblazeB2

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `accessKeyId` | String | Yes | Required except for AWS. Up to 512 characters. Google Cloud Storage requires an HMAC access ID beginning with GOOG.  |
| `bucket` | String | Yes | Existing bucket or Space name, up to 255 characters. The provider never creates or deletes the bucket.  |
| `keyPrefix` | String | No | Object key prefix, up to 1024 characters. The API trims surrounding whitespace and slashes; omission uses the bucket root.  |
| `region` | String | Yes | Bucket region supported by backblaze_b2.  |
| `secretAccessKey` | String | No | Sensitive access key secret. Required on create or a change of provider/access key ID. Omit on unrelated updates to retain the existing secret. Empty is invalid. Up to 990 UTF-16 units and 996 UTF-8 bytes after trimming.  |

### StorageDestinationConnectionInfoCloudflareR2

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `accessKeyId` | String | Yes | Required except for AWS. Up to 512 characters. Google Cloud Storage requires an HMAC access ID beginning with GOOG.  |
| `bucket` | String | Yes | Existing bucket or Space name, up to 255 characters. The provider never creates or deletes the bucket.  |
| `cloudflareAccountId` | String | Yes | R2 only. Required 32-character hexadecimal account ID.  |
| `cloudflareJurisdiction` | String | No | R2 only. eu or fedramp; omit for default jurisdiction.  |
| `keyPrefix` | String | No | Object key prefix, up to 1024 characters. The API trims surrounding whitespace and slashes; omission uses the bucket root.  |
| `secretAccessKey` | String | No | Sensitive access key secret. Required on create or a change of provider/access key ID. Omit on unrelated updates to retain the existing secret. Empty is invalid. Up to 990 UTF-16 units and 996 UTF-8 bytes after trimming.  |

### StorageDestinationConnectionInfoDigitaloceanSpaces

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `accessKeyId` | String | Yes | Required except for AWS. Up to 512 characters. Google Cloud Storage requires an HMAC access ID beginning with GOOG.  |
| `bucket` | String | Yes | Existing bucket or Space name, up to 255 characters. The provider never creates or deletes the bucket.  |
| `keyPrefix` | String | No | Object key prefix, up to 1024 characters. The API trims surrounding whitespace and slashes; omission uses the bucket root.  |
| `region` | String | Yes | Bucket region supported by digitalocean_spaces.  |
| `secretAccessKey` | String | No | Sensitive access key secret. Required on create or a change of provider/access key ID. Omit on unrelated updates to retain the existing secret. Empty is invalid. Up to 990 UTF-16 units and 996 UTF-8 bytes after trimming.  |

### StorageDestinationConnectionInfoGoogleCloudStorage

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `accessKeyId` | String | Yes | Required except for AWS. Up to 512 characters. Google Cloud Storage requires an HMAC access ID beginning with GOOG.  |
| `bucket` | String | Yes | Existing bucket or Space name, up to 255 characters. The provider never creates or deletes the bucket.  |
| `keyPrefix` | String | No | Object key prefix, up to 1024 characters. The API trims surrounding whitespace and slashes; omission uses the bucket root.  |
| `secretAccessKey` | String | No | Sensitive access key secret. Required on create or a change of provider/access key ID. Omit on unrelated updates to retain the existing secret. Empty is invalid. Up to 990 UTF-16 units and 996 UTF-8 bytes after trimming.  |

### StorageDestinationConnectionInfoOtherS3Compatible

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `accessKeyId` | String | Yes | Required except for AWS. Up to 512 characters. Google Cloud Storage requires an HMAC access ID beginning with GOOG.  |
| `bucket` | String | Yes | Existing bucket or Space name, up to 255 characters. The provider never creates or deletes the bucket.  |
| `endpoint` | String | Yes | Required for other_s3_compatible. Public HTTPS origin, up to 512 characters, without credentials, path, query, or fragment.  |
| `forcePathStyle` | Boolean | No | Other S3-compatible only. Omitted or null uses true. No value is injected into the request.  |
| `keyPrefix` | String | No | Object key prefix, up to 1024 characters. The API trims surrounding whitespace and slashes; omission uses the bucket root.  |
| `region` | String | No | Optional signing region, up to 64 characters. Omitted or null uses the API default us-east-1.  |
| `secretAccessKey` | String | No | Sensitive access key secret. Required on create or a change of provider/access key ID. Omit on unrelated updates to retain the existing secret. Empty is invalid. Up to 990 UTF-16 units and 996 UTF-8 bytes after trimming.  |

### StorageDestinationConnectionInfoWasabi

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `accessKeyId` | String | Yes | Required except for AWS. Up to 512 characters. Google Cloud Storage requires an HMAC access ID beginning with GOOG.  |
| `bucket` | String | Yes | Existing bucket or Space name, up to 255 characters. The provider never creates or deletes the bucket.  |
| `keyPrefix` | String | No | Object key prefix, up to 1024 characters. The API trims surrounding whitespace and slashes; omission uses the bucket root.  |
| `region` | String | Yes | Bucket region supported by wasabi.  |
| `secretAccessKey` | String | No | Sensitive access key secret. Required on create or a change of provider/access key ID. Omit on unrelated updates to retain the existing secret. Empty is invalid. Up to 990 UTF-16 units and 996 UTF-8 bytes after trimming.  |

## Import

```sh
pulumi import html-css-to-image:index:StorageDestination example MANAGEMENT_ID
```

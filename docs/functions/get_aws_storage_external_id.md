# getAwsStorageExternalId

Read the organization’s external ID and HCTI writer role ARN for an AWS IAM role trust policy. This lookup has no inputs and returns `externalId` and `writerRoleArn`. It does not create a storage destination or an IAM role. The API requires `storage_destinations:create_update` for this lookup.

**Function token:** `html-css-to-image:index:getAwsStorageExternalId`

```yaml
variables:
  trust:
    fn::invoke:
      function: html-css-to-image:index:getAwsStorageExternalId
      arguments: {}
outputs:
  externalId: ${trust.externalId}
  writerRoleArn: ${trust.writerRoleArn}
```

Use `externalId` for the trust policy’s `sts:ExternalId` condition and `writerRoleArn` for `Principal.AWS`, allowing `sts:AssumeRole`. Create/update your role before applying the HCTI storage destination. The destination's `roleArn` is your role, not the returned HCTI writer role. See the [AWS example](../../examples/aws-storage/Pulumi.yaml).

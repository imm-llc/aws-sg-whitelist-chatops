# aws-sg-whitelist-chatops

ChatOps AWS security group whitelisting monorepo.

## Structure

The code in `handler/` handles the incoming slash command request.

The code in `internal/` handles checking whitelist entries, updating security groups, and adding entries to DynamoDB.

The code in `cleanup/` handles incoming DynamoDB Streams data.

## Docs

[EC2 Go SDK docs](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/)

[Security Group Ingress Docs](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2.AuthorizeSecurityGroupIngress)

[DynamoDB Streams Go example](https://docs.aws.amazon.com/lambda/latest/dg/with-ddb-create-package.html#with-ddb-example-deployment-pkg-go)

## Invocation 

Example: `/db staging rds 1.1.1.1`
## TODO

Pull port in from map or environment to make this tool more extensible.

## Environment Variables

### Handler

`STAGING_ALLOWED_SLACK_USER_IDS`: Comma-delimited list of Slack User IDs of people who are allowed to invoke this command to add entries for staging

`PRODUCTION_ALLOWED_SLACK_USER_IDS`: Comma-delimited list of Slack User IDs of people who are allowed to invoke this command to add entries for production

`STAGING_SECURITY_GROUP_ID`: Security Group ID to add entries to for staging

`PRODUCTION_SECURITY_GROUP_ID`: Security Group ID to add entries to for production

`DDB_TABLE_NAME`: DynamoDB table to add expiring entries to

`SLACK_VERIFICATION_TOKEN`: Verification token of Slack application

### Cleanup

None required

## Notes

### Slack App

Slack App ID:

### Examples

`ec2.DescribeSecurityGroups` output:

```json
{
  "SecurityGroups": [{
      "Description": "",
      "GroupId": "",
      "GroupName": "",
      "IpPermissions": [{
          "IpProtocol": "-1",
          "IpRanges": [{
              "CidrIp": "",
              "Description": ""
            },{
              "CidrIp": "",
              "Description": ""
            }]
        },{
          "FromPort": 443,
          "IpProtocol": "tcp",
          "IpRanges": [{
              "CidrIp": "0.0.0.0/0",
              "Description": "HTTPS"
            }],
          "ToPort": 443
        }],
      "IpPermissionsEgress": [{
          "IpProtocol": "-1",
          "IpRanges": [{
              "CidrIp": "0.0.0.0/0",
              "Description": "All"
            }]
        }],
      "OwnerId": "",
      "Tags": [{
          "Key": "app_name",
          "Value": ""
        },{
          "Key": "Name",
          "Value": ""
        }],
      "VpcId": ""
    }]
}
```

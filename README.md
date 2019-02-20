# Cloud-Security-Audit-Serverless


## Cloud Security Audit 
https://github.com/Appliscale/cloud-security-audit


## Installing Serverless
https://serverless.com/framework/docs/providers/aws/guide/installation/


## Deployment 

```bash
sls deploy --aws-profile <PROFILE>
```
After deployment is needed to increase lambda timeout - 6 second isn't sufficient, but about 1 minute will be. Also lambda must have access to some services - EC2, ELB, S3, IAM and KMS. Its role should have policies such as: 
AmazonEC2ReadOnlyAccess, AmazonS3ReadOnlyAccess, AWSKeyManagementServicePowerUser. These policies allows to read and list information about above services.
Instead of using ```--aws-profile``` flag, you can follow this tutorial https://serverless.com/framework/docs/providers/aws/guide/credentials/.

You can also use earlier prepared lambda _cloud-security-audit-serverless-dev-cloud-security-audit_.

## Running

```bash
sls invoke -f cloud-security-audit --aws-profile <PROFILE>
```
Output will be available in _CloudWatch/Logs_.




 






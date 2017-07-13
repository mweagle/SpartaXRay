# SpartaXRay
Sparta-based application that demonstrates how to enable [Lambda XRay Tracing](https://aws.amazon.com/blogs/aws/aws-x-ray-update-general-availability-including-lambda-integration/) as well as provision a [CloudWatch Dashboard](http://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch_Dashboards.html) together with your Sparta-based AWS Lambda service.


# Usage:

```
go run main.go provision --s3Bucket $MY_S3_BUCKET
```
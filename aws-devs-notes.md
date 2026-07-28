# AWS for Developers — Notes

## 1. Core Concepts

- **Region**: A geographic area (e.g., `us-east-1`) containing multiple data centers.
- **Availability Zone (AZ)**: Isolated data center(s) within a region. Deploy across AZs for high availability.
- **Edge Location**: Used by CloudFront (CDN) for caching content close to users.
- **IAM**: Identity and Access Management — controls *who* can do *what*.

---

## 2. IAM (Identity and Access Management)

- **User**: Represents a person/service with long-term credentials.
- **Group**: Collection of users sharing permissions.
- **Role**: Temporary permissions assumed by users, services, or applications (no long-term credentials).
- **Policy**: JSON document defining permissions.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::my-bucket/*"
    }
  ]
}
```

**Best practices**
- Never use the root account for daily tasks.
- Follow least privilege principle.
- Use roles for EC2/Lambda instead of hardcoding credentials.
- Enable MFA.

---

## 3. Compute

### EC2 (Elastic Compute Cloud)
- Virtual servers ("instances").
- Instance types: `t3.micro` (burstable), `m5.large` (general purpose), `c5.xlarge` (compute optimized), etc.
- Pricing models: On-Demand, Reserved, Spot, Savings Plans.
- **Security Groups**: Virtual firewall at the instance level (stateful).
- **AMI**: Amazon Machine Image — template used to launch instances.
- **Auto Scaling Group (ASG)**: Automatically adds/removes instances based on demand.
- **Elastic Load Balancer (ELB)**: Distributes traffic (ALB for HTTP/HTTPS, NLB for TCP/high performance).

### Lambda (Serverless)
- Run code without provisioning servers; pay per invocation/duration.
- Triggered by events: API Gateway, S3, DynamoDB Streams, SQS, EventBridge, etc.
- Limits: 15 min max execution time, memory 128 MB–10,240 MB.
- Cold starts are a common performance consideration.

```python
def lambda_handler(event, context):
    return {
        'statusCode': 200,
        'body': 'Hello from Lambda!'
    }
```

### Elastic Beanstalk
- PaaS: upload code, AWS handles provisioning, load balancing, scaling.

### ECS / EKS / Fargate
- **ECS**: Elastic Container Service — run Docker containers.
- **EKS**: Managed Kubernetes.
- **Fargate**: Serverless compute for containers (no EC2 management).

---

## 4. Storage

### S3 (Simple Storage Service)
- Object storage, organized into **buckets**.
- Storage classes: Standard, Intelligent-Tiering, Standard-IA, One Zone-IA, Glacier, Glacier Deep Archive.
- Versioning, lifecycle policies, and encryption (SSE-S3, SSE-KMS, SSE-C) available.
- Static website hosting supported.

```bash
aws s3 cp file.txt s3://my-bucket/
aws s3 sync ./local-folder s3://my-bucket/folder/
```

### EBS (Elastic Block Store)
- Persistent block storage attached to a single EC2 instance (like a virtual hard drive).

### EFS (Elastic File System)
- Managed NFS file system, shareable across multiple EC2 instances.

---

## 5. Databases

### RDS (Relational Database Service)
- Managed SQL databases: MySQL, PostgreSQL, MariaDB, Oracle, SQL Server, Aurora.
- Handles backups, patching, replication (Multi-AZ), read replicas.

### DynamoDB
- Managed NoSQL key-value/document database.
- Fully serverless, scales automatically.
- Primary key: Partition key (+ optional Sort key).
- On-demand or provisioned capacity modes.

```python
import boto3
table = boto3.resource('dynamodb').Table('Users')
table.put_item(Item={'id': '123', 'name': 'Alice'})
```

### ElastiCache
- Managed in-memory caching: Redis or Memcached.

---

## 6. Networking

### VPC (Virtual Private Cloud)
- Your own isolated network within AWS.
- **Subnet**: Public (has route to Internet Gateway) or Private.
- **Internet Gateway**: Allows internet access.
- **NAT Gateway**: Allows private subnet outbound internet access (not inbound).
- **Route Table**: Controls traffic routing between subnets.
- **NACL**: Network ACL — stateless firewall at subnet level.
- **Security Group**: Stateful firewall at instance level.

### Route 53
- DNS service; supports routing policies (simple, weighted, latency-based, failover, geolocation).

### CloudFront
- CDN — caches content at edge locations for lower latency.

### API Gateway
- Create, publish, and manage REST/HTTP/WebSocket APIs.
- Common integration: API Gateway → Lambda → DynamoDB.

---

## 7. Messaging & Integration

- **SQS (Simple Queue Service)**: Message queue for decoupling services. Standard (at-least-once) or FIFO (exactly-once, ordered).
- **SNS (Simple Notification Service)**: Pub/sub messaging — push notifications to subscribers (email, SMS, Lambda, SQS).
- **EventBridge**: Event bus for routing events between AWS services and apps based on rules.
- **Step Functions**: Orchestrate multi-step workflows (state machines) across Lambda, ECS, etc.

---

## 8. Monitoring & Logging

- **CloudWatch**: Metrics, logs, alarms, dashboards.
  - CloudWatch Logs: application/system logs.
  - CloudWatch Alarms: trigger actions (e.g., auto scaling, SNS notification) based on metric thresholds.
- **CloudTrail**: Logs all API calls/account activity for auditing.
- **X-Ray**: Distributed tracing for debugging microservices/serverless apps.

---

## 9. Security

- **KMS (Key Management Service)**: Create/manage encryption keys.
- **Secrets Manager**: Store and rotate secrets (DB credentials, API keys).
- **Systems Manager Parameter Store**: Store config data/secrets (cheaper alternative to Secrets Manager).
- **WAF**: Web Application Firewall — protects against common web exploits.
- **Shield**: DDoS protection.

---

## 10. CI/CD & Developer Tools

- **CodeCommit**: Managed Git repositories.
- **CodeBuild**: Compile/test source code.
- **CodeDeploy**: Automate deployments to EC2, Lambda, ECS.
- **CodePipeline**: Orchestrate full CI/CD pipeline.
- **CloudFormation**: Infrastructure as Code (IaC) using JSON/YAML templates.
- **CDK (Cloud Development Kit)**: Define infrastructure using real programming languages (Python, TypeScript, etc.), compiles to CloudFormation.
- **SAM (Serverless Application Model)**: Simplified CloudFormation for serverless apps.

```yaml
# Simple CloudFormation snippet
Resources:
  MyBucket:
    Type: AWS::S3::Bucket
    Properties:
      BucketName: my-example-bucket
```

---

## 11. AWS CLI & SDK Basics

```bash
aws configure                     # set up credentials
aws s3 ls                         # list buckets
aws ec2 describe-instances        # list EC2 instances
aws lambda invoke --function-name myFunc output.json
```

```python
# boto3 (Python SDK) example
import boto3
s3 = boto3.client('s3')
s3.upload_file('file.txt', 'my-bucket', 'file.txt')
```

---

## 12. Common Architecture Patterns

- **Serverless web app**: Route 53 → CloudFront → S3 (static frontend) → API Gateway → Lambda → DynamoDB.
- **Decoupled processing**: Producer → SQS → Lambda/EC2 consumers.
- **Event-driven**: S3 upload → Lambda trigger → processing → SNS notification.
- **Three-tier web app**: ALB → EC2 (Auto Scaling) → RDS (Multi-AZ).

---

## 13. Well-Architected Framework Pillars

1. **Operational Excellence**
2. **Security**
3. **Reliability**
4. **Performance Efficiency**
5. **Cost Optimization**
6. **Sustainability**

---

## 14. Quick Exam/Interview Tips

- SGs are stateful; NACLs are stateless.
- S3 is regional but globally accessible by name (bucket names are globally unique).
- Lambda scales automatically; concurrency limits can throttle it.
- Use IAM roles instead of access keys wherever possible.
- RTO (Recovery Time Objective) vs RPO (Recovery Point Objective) — key DR concepts.

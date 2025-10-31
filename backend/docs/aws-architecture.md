# Chess Coach AWS Architecture

**Cloud Provider:** Amazon Web Services (AWS)
**Team Size:** 2-3 developers
**Philosophy:** Cloud-native on AWS, but locally-first development
**Status:** Production Architecture Design

---

## Executive Summary

This document provides an AWS-native architecture for Chess Coach that balances:
1. **AWS-native services** for production reliability
2. **Local development** that mirrors AWS (Docker Compose)
3. **Small team efficiency** (2-3 developers can manage)
4. **Cost optimization** (~$150/mo MVP → $800/mo at scale)

### AWS Services Used

| Service | Purpose | Why This Service |
|---------|---------|------------------|
| **ECS Fargate** | Container hosting | Serverless, no EC2 management, auto-scaling |
| **Application Load Balancer** | WebSocket routing | Native WebSocket support, health checks |
| **ElastiCache (Redis)** | Real-time coordination | Managed Redis, automatic failover |
| **RDS PostgreSQL** | Data persistence | Managed DB, automated backups |
| **MSK (Managed Kafka)** | Message broker | AWS-managed Kafka alternative to NATS |
| **S3** | Game archives (PGN files) | Cheap storage, 99.999999999% durability |
| **CloudWatch** | Observability | Native AWS integration, logs + metrics |
| **Secrets Manager** | API keys, credentials | Automatic rotation, encryption |
| **CloudFront** | Frontend CDN | Global edge, WebSocket support |
| **Route 53** | DNS | Easy SSL with ACM |

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Service-by-Service Breakdown](#service-by-service-breakdown)
3. [Local Development Setup](#local-development-setup)
4. [Message Broker Decision: MSK vs NATS on ECS](#message-broker-decision)
5. [Networking & Security](#networking-security)
6. [CI/CD Pipeline](#cicd-pipeline)
7. [Cost Analysis](#cost-analysis)
8. [Monitoring & Alerting](#monitoring-alerting)
9. [Scaling Strategy](#scaling-strategy)
10. [Disaster Recovery](#disaster-recovery)
11. [Step-by-Step Implementation](#step-by-step-implementation)

---

## 1. Architecture Overview

### High-Level Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    AWS INFRASTRUCTURE                        │
└─────────────────────────────────────────────────────────────┘

Internet
    ↓
┌─────────────────────────────────────────┐
│  CloudFront (CDN)                       │  ← Frontend (React)
│  - Global edge locations                │
│  - WebSocket upgrade support            │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│  Application Load Balancer              │
│  - WebSocket support                    │
│  - SSL termination (ACM)                │
│  - Health checks                        │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│  ECS Fargate Cluster                    │
│  ┌────────────────────────────────┐    │
│  │  Backend Service (Auto-scaling)│    │
│  │  - WebSocket server            │    │
│  │  - Message router              │    │
│  │  - AI integration              │    │
│  │  - Engine coordination         │    │
│  │  (2-10 tasks based on load)    │    │
│  └────────────────────────────────┘    │
└─────────────────────────────────────────┘
    ↓                    ↓
┌──────────────┐    ┌──────────────────┐
│ ElastiCache  │    │  MSK (Kafka)     │
│ (Redis)      │    │  or               │
│ - Game rooms │    │  NATS on ECS     │
│ - Sessions   │    │  - Message queue │
│ - Pub/sub    │    │  - Event stream  │
└──────────────┘    └──────────────────┘
    ↓
┌─────────────────────────────────────────┐
│  RDS PostgreSQL (Multi-AZ)              │
│  - Games, users, chat history           │
│  - Automated backups                    │
└─────────────────────────────────────────┘
    ↓
┌─────────────────────────────────────────┐
│  S3                                     │
│  - PGN archives                         │
│  - User uploads                         │
└─────────────────────────────────────────┘
```

### AWS Regions Strategy

**Primary Region:** `us-east-1` (N. Virginia)
- Lowest latency to East Coast US
- Most AWS services available
- Cheapest pricing

**Future Multi-Region:**
- Add `eu-west-1` (Ireland) for European users
- Use Route 53 latency-based routing

---

## 2. Service-by-Service Breakdown

### 2.1 ECS Fargate (Backend Hosting)

**What:** Serverless container runtime (no EC2 management)

**Configuration:**

```yaml
# ecs-task-definition.json
{
  "family": "chess-coach-backend",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",      # 0.5 vCPU
  "memory": "1024",  # 1 GB
  "taskRoleArn": "arn:aws:iam::ACCOUNT:role/chess-coach-task-role",
  "executionRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskExecutionRole",
  "containerDefinitions": [
    {
      "name": "backend",
      "image": "ACCOUNT.dkr.ecr.us-east-1.amazonaws.com/chess-coach:latest",
      "portMappings": [
        {
          "containerPort": 3000,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {"name": "NODE_ENV", "value": "production"},
        {"name": "REDIS_URL", "value": "redis://elasticache-endpoint:6379"}
      ],
      "secrets": [
        {
          "name": "ANTHROPIC_API_KEY",
          "valueFrom": "arn:aws:secretsmanager:us-east-1:ACCOUNT:secret:chess-coach/anthropic-api-key"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/chess-coach",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "backend"
        }
      },
      "healthCheck": {
        "command": ["CMD-SHELL", "curl -f http://localhost:3000/health || exit 1"],
        "interval": 30,
        "timeout": 5,
        "retries": 3
      }
    }
  ]
}
```

**Auto-Scaling:**

```json
// auto-scaling-policy.json
{
  "TargetTrackingScalingPolicyConfiguration": {
    "TargetValue": 70.0,
    "PredefinedMetricSpecification": {
      "PredefinedMetricType": "ECSServiceAverageCPUUtilization"
    },
    "ScaleOutCooldown": 60,
    "ScaleInCooldown": 300
  }
}
```

**Why Fargate:**
- ✅ No server management
- ✅ Pay only for container runtime
- ✅ Auto-scaling built-in
- ✅ Perfect for small teams

**Cost:** $0.04/vCPU/hour + $0.004/GB/hour
- 1 task (0.5 vCPU, 1GB) = ~$15/mo running 24/7
- 3 tasks for HA = ~$45/mo

---

### 2.2 Application Load Balancer (ALB)

**What:** Layer 7 load balancer with WebSocket support

**Configuration:**

```yaml
# Target Group (WebSocket-enabled)
TargetGroup:
  Protocol: HTTP
  Port: 3000
  TargetType: ip
  VpcId: vpc-xxxxx
  HealthCheckEnabled: true
  HealthCheckPath: /health
  HealthCheckProtocol: HTTP
  HealthCheckIntervalSeconds: 30
  HealthCheckTimeoutSeconds: 5
  HealthyThresholdCount: 2
  UnhealthyThresholdCount: 3
  Matcher:
    HttpCode: 200
  # Enable connection draining
  DeregistrationDelay: 30
  TargetGroupAttributes:
    - Key: stickiness.enabled
      Value: true
    - Key: stickiness.type
      Value: lb_cookie
    - Key: stickiness.lb_cookie.duration_seconds
      Value: 86400  # 24 hours (keep WebSocket connections)

# Listener (HTTPS)
Listener:
  Protocol: HTTPS
  Port: 443
  Certificates:
    - CertificateArn: arn:aws:acm:us-east-1:ACCOUNT:certificate/xxxxx
  DefaultActions:
    - Type: forward
      TargetGroupArn: !Ref TargetGroup
```

**Why ALB:**
- ✅ Native WebSocket support (connection upgrade)
- ✅ SSL termination with ACM (free certificates)
- ✅ Health checks for Fargate tasks
- ✅ Sticky sessions (same client → same backend)

**Cost:** ~$16/mo (base) + $0.008/LCU-hour
- MVP: ~$20-25/mo
- Scale: ~$40-60/mo

---

### 2.3 ElastiCache (Redis)

**What:** Managed Redis for real-time coordination

**Use Cases:**
- WebSocket connection mapping (user → task)
- Game room state (multiplayer)
- Pub/sub for broadcasting moves
- Session storage
- Rate limiting

**Configuration:**

```yaml
CacheCluster:
  Engine: redis
  CacheNodeType: cache.t4g.micro  # 0.5GB RAM (MVP)
  NumCacheNodes: 1  # Single node for MVP
  # For HA: Use replication group
  AutoMinorVersionUpgrade: true
  SnapshotRetentionLimit: 5
  SnapshotWindow: "03:00-05:00"  # UTC
  PreferredMaintenanceWindow: "mon:05:00-mon:07:00"

# For Production (Multi-AZ):
ReplicationGroup:
  Engine: redis
  CacheNodeType: cache.t4g.small  # 1.37GB RAM
  NumCacheClusters: 2  # Primary + Replica
  AutomaticFailoverEnabled: true
  MultiAZEnabled: true
```

**Code Example:**

```typescript
import { createClient } from 'redis';

// Connection
const redis = createClient({
  url: process.env.REDIS_URL,
  socket: {
    reconnectStrategy: (retries) => Math.min(retries * 50, 500)
  }
});

// Store WebSocket connection mapping
await redis.hSet('ws:connections', userId, taskId);

// Pub/sub for game room
await redis.publish(`game:${gameId}`, JSON.stringify(move));

// Subscribe to game events
const subscriber = redis.duplicate();
await subscriber.subscribe(`game:${gameId}`, (message) => {
  broadcastToRoom(gameId, message);
});
```

**Why ElastiCache:**
- ✅ Managed service (automatic failover, patching)
- ✅ Multi-AZ for HA
- ✅ Automatic backups
- ✅ VPC isolation (secure)

**Cost:**
- MVP (t4g.micro): $12/mo
- Production (t4g.small, Multi-AZ): $50/mo

---

### 2.4 RDS PostgreSQL

**What:** Managed relational database

**Schema:**

```sql
-- Games table
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    pgn TEXT NOT NULL,
    white_player VARCHAR(255) NOT NULL,
    black_player VARCHAR(255) NOT NULL,
    result VARCHAR(10),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Chat history
CREATE TABLE chat_history (
    id SERIAL PRIMARY KEY,
    game_id INTEGER REFERENCES games(id),
    user_id VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,  -- 'user' | 'assistant'
    content TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW()
);

-- User sessions
CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    data JSONB,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_games_created_at ON games(created_at DESC);
CREATE INDEX idx_chat_game_id ON chat_history(game_id);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
```

**Configuration:**

```yaml
DBInstance:
  DBInstanceClass: db.t4g.micro  # 1 vCPU, 1GB RAM (MVP)
  Engine: postgres
  EngineVersion: "16.1"
  AllocatedStorage: 20  # GB (GP3 SSD)
  StorageType: gp3
  StorageEncrypted: true
  MultiAZ: false  # MVP (true for production)
  BackupRetentionPeriod: 7  # days
  PreferredBackupWindow: "03:00-04:00"  # UTC
  PreferredMaintenanceWindow: "mon:04:00-mon:05:00"
  DeletionProtection: true
  EnablePerformanceInsights: true
  PerformanceInsightsRetentionPeriod: 7

# For Production:
DBInstance:
  DBInstanceClass: db.t4g.small  # 2 vCPU, 2GB RAM
  MultiAZ: true
  AllocatedStorage: 100  # GB
```

**Connection Pooling (Application Side):**

```typescript
import { Pool } from 'pg';

const pool = new Pool({
  host: process.env.DB_HOST,
  port: 5432,
  database: process.env.DB_NAME,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  max: 20,  // Max connections per instance
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
  ssl: {
    rejectUnauthorized: false  // RDS SSL
  }
});
```

**Why RDS:**
- ✅ Automated backups and point-in-time recovery
- ✅ Automatic patching and upgrades
- ✅ Multi-AZ for HA (production)
- ✅ Performance Insights for query optimization

**Cost:**
- MVP (t4g.micro, Single-AZ): $15/mo
- Production (t4g.small, Multi-AZ): $65/mo

---

### 2.5 Message Broker: MSK vs NATS on ECS

#### Option A: Amazon MSK (Managed Kafka) ❌ **TOO EXPENSIVE FOR SMALL TEAMS**

**Pros:**
- Fully managed Kafka
- Integrates with AWS services
- Multi-AZ by default

**Cons:**
- **Minimum cost: $200/mo** (2x t3.small brokers)
- Overkill for <100K users
- Complex configuration

**Verdict:** Skip MSK for MVP

---

#### Option B: NATS on ECS Fargate ✅ **RECOMMENDED**

**Why:**
- ✅ Runs in ECS Fargate ($15-30/mo)
- ✅ Simpler than Kafka
- ✅ Perfect for real-time messaging
- ✅ Easy to manage for 2-3 developers

**ECS Task Definition:**

```json
{
  "family": "nats-jetstream",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "512",
  "memory": "1024",
  "containerDefinitions": [
    {
      "name": "nats",
      "image": "nats:2.10-alpine",
      "portMappings": [
        {"containerPort": 4222, "protocol": "tcp"},
        {"containerPort": 8222, "protocol": "tcp"}
      ],
      "command": [
        "--jetstream",
        "--store_dir=/data",
        "--max_memory_store=512MB",
        "--max_file_store=5GB"
      ],
      "mountPoints": [
        {
          "sourceVolume": "nats-data",
          "containerPath": "/data"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/nats",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "nats"
        }
      }
    }
  ],
  "volumes": [
    {
      "name": "nats-data",
      "efsVolumeConfiguration": {
        "fileSystemId": "fs-xxxxx",
        "transitEncryption": "ENABLED"
      }
    }
  ]
}
```

**Note:** NATS needs persistent storage (EFS) for JetStream

**EFS Configuration:**

```yaml
EFSFileSystem:
  PerformanceMode: generalPurpose
  ThroughputMode: bursting
  Encrypted: true
  LifecyclePolicies:
    - TransitionToIA: AFTER_30_DAYS
```

**Cost Breakdown:**
- ECS Fargate (NATS): $15/mo
- EFS storage (10GB): $3/mo
- Total: **$18/mo**

vs MSK: $200+/mo

**Verdict:** NATS on ECS is 10x cheaper and simpler

---

#### Option C: Amazon EventBridge + SQS ⚠️ **ALTERNATIVE**

For teams wanting pure AWS services:

```
Frontend → ALB → ECS Backend → EventBridge → SQS → ECS Workers
```

**Pros:**
- Pure AWS-native
- Serverless, pay-per-use
- No broker to manage

**Cons:**
- Higher latency (100-300ms)
- Not suitable for real-time chess moves
- More complex message routing

**Verdict:** Not ideal for real-time gaming

---

### 2.6 S3 (Game Archives)

**What:** Store PGN files, user-uploaded games

**Buckets:**

```
chess-coach-pgn-archives/
  ├── imported/
  │   └── {userId}/{gameId}.pgn
  └── user-uploads/
      └── {userId}/{timestamp}.pgn
```

**S3 Lifecycle Policy:**

```json
{
  "Rules": [
    {
      "Id": "MoveToGlacierAfter90Days",
      "Status": "Enabled",
      "Transitions": [
        {
          "Days": 90,
          "StorageClass": "GLACIER_INSTANT_RETRIEVAL"
        }
      ]
    }
  ]
}
```

**Cost:** ~$0.023/GB/month (S3 Standard)
- 100K games (10MB each) = 1TB = $23/mo
- After 90 days → Glacier = $4/mo

---

### 2.7 CloudWatch (Observability)

**What:** AWS-native monitoring, logs, metrics

**Log Groups:**
- `/ecs/chess-coach` (backend logs)
- `/ecs/nats` (NATS logs)
- `/aws/rds/instance/chess-coach/postgresql` (DB logs)

**Custom Metrics:**

```typescript
import { CloudWatch } from '@aws-sdk/client-cloudwatch';

const cw = new CloudWatch({ region: 'us-east-1' });

// Publish custom metric
await cw.putMetricData({
  Namespace: 'ChessCoach',
  MetricData: [
    {
      MetricName: 'ChessMovesProcessed',
      Value: 1,
      Unit: 'Count',
      Timestamp: new Date(),
      Dimensions: [
        { Name: 'Environment', Value: 'production' }
      ]
    }
  ]
});
```

**CloudWatch Alarms:**

```yaml
Alarms:
  - Name: HighCPUUtilization
    MetricName: CPUUtilization
    Namespace: AWS/ECS
    Statistic: Average
    Period: 300
    EvaluationPeriods: 2
    Threshold: 80
    ComparisonOperator: GreaterThanThreshold
    AlarmActions:
      - arn:aws:sns:us-east-1:ACCOUNT:alerts

  - Name: HighDatabaseConnections
    MetricName: DatabaseConnections
    Namespace: AWS/RDS
    Statistic: Average
    Period: 300
    EvaluationPeriods: 1
    Threshold: 80
    ComparisonOperator: GreaterThanThreshold
```

**Cost:**
- Free tier: 10 custom metrics, 5GB logs ingestion
- Beyond: ~$10-30/mo

---

### 2.8 CloudFront (Frontend CDN)

**What:** Global CDN for React frontend, supports WebSocket

**Configuration:**

```yaml
Distribution:
  Origins:
    # Static assets (S3)
    - Id: S3Origin
      DomainName: chess-coach-frontend.s3.amazonaws.com
      S3OriginConfig:
        OriginAccessIdentity: origin-access-identity/cloudfront/XXXXX

    # WebSocket API (ALB)
    - Id: ALBOrigin
      DomainName: chess-coach-alb-xxxxx.us-east-1.elb.amazonaws.com
      CustomOriginConfig:
        HTTPPort: 80
        HTTPSPort: 443
        OriginProtocolPolicy: https-only
        OriginKeepaliveTimeout: 60
        OriginReadTimeout: 60

  DefaultCacheBehavior:
    TargetOriginId: S3Origin
    ViewerProtocolPolicy: redirect-to-https
    CachePolicyId: 658327ea-f89d-4fab-a63d-7e88639e58f6  # CachingOptimized

  CacheBehaviors:
    # WebSocket endpoint (no caching)
    - PathPattern: /ws/*
      TargetOriginId: ALBOrigin
      ViewerProtocolPolicy: https-only
      CachePolicyId: 4135ea2d-6df8-44a3-9df3-4b5a84be39ad  # CachingDisabled
      AllowedMethods: [GET, HEAD, OPTIONS, PUT, POST, PATCH, DELETE]

  PriceClass: PriceClass_100  # US, Canada, Europe
  ViewerCertificate:
    AcmCertificateArn: arn:aws:acm:us-east-1:ACCOUNT:certificate/xxxxx
    SslSupportMethod: sni-only
```

**Why CloudFront:**
- ✅ Global edge locations (low latency)
- ✅ DDoS protection (AWS Shield Standard)
- ✅ Free SSL certificates (ACM)
- ✅ WebSocket support via custom origins

**Cost:**
- Free tier: 1TB data transfer/month
- Beyond: $0.085/GB (first 10TB)
- MVP: $0-10/mo

---

### 2.9 Secrets Manager

**What:** Secure storage for API keys, DB passwords

**Secrets:**
```
chess-coach/anthropic-api-key
chess-coach/db-password
chess-coach/jwt-secret
chess-coach/openai-api-key
```

**Automatic Rotation:**

```typescript
// Lambda function for password rotation
export const handler = async (event: any) => {
  const { SecretId, Token, Step } = event;

  if (Step === 'createSecret') {
    // Generate new password
    const newPassword = generateSecurePassword();
    await rds.modifyDBInstance({
      MasterUserPassword: newPassword
    });
  }

  if (Step === 'testSecret') {
    // Test new credentials
    await testDatabaseConnection(newPassword);
  }

  if (Step === 'finishSecret') {
    // Mark rotation complete
  }
};
```

**Cost:** $0.40/secret/month + $0.05/10K API calls
- 5 secrets = $2/mo

---

## 3. Local Development Setup

### The Key Principle: Local = AWS

**Goal:** Run everything locally using Docker Compose that mimics AWS

### docker-compose.yml

```yaml
version: '3.8'

services:
  # Simulates RDS
  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: chess_coach
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: dev
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U dev"]
      interval: 5s
      timeout: 3s
      retries: 3

  # Simulates ElastiCache
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 3

  # Simulates NATS on ECS
  nats:
    image: nats:2.10-alpine
    ports:
      - "4222:4222"
      - "8222:8222"
    command: --jetstream --http_port 8222 --store_dir=/data
    volumes:
      - nats-data:/data
    healthcheck:
      test: ["CMD", "wget", "-q", "-O-", "http://localhost:8222/healthz"]
      interval: 5s
      timeout: 3s
      retries: 3

  # Simulates S3 (LocalStack)
  localstack:
    image: localstack/localstack:latest
    ports:
      - "4566:4566"
    environment:
      SERVICES: s3,secretsmanager
      DEBUG: 1
    volumes:
      - localstack-data:/var/lib/localstack

  # Backend (your code)
  backend:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "3000:3000"
    environment:
      NODE_ENV: development
      REDIS_URL: redis://redis:6379
      DATABASE_URL: postgresql://dev:dev@postgres:5432/chess_coach
      NATS_URL: nats://nats:4222
      AWS_ENDPOINT: http://localstack:4566  # LocalStack
      AWS_REGION: us-east-1
      AWS_ACCESS_KEY_ID: test
      AWS_SECRET_ACCESS_KEY: test
    volumes:
      - ./src:/app/src  # Hot reload
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
      nats:
        condition: service_healthy
    command: bun run dev

volumes:
  postgres-data:
  nats-data:
  localstack-data:
```

### Start Everything

```bash
# One command starts entire stack
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop everything
docker-compose down
```

### Environment Variables (.env)

```bash
# Local development
NODE_ENV=development
PORT=3000

# Database (matches docker-compose)
DATABASE_URL=postgresql://dev:dev@localhost:5432/chess_coach

# Redis
REDIS_URL=redis://localhost:6379

# NATS
NATS_URL=nats://localhost:4222

# AWS (LocalStack for local dev)
AWS_REGION=us-east-1
AWS_ENDPOINT=http://localhost:4566
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test

# API Keys (get from Secrets Manager in production)
ANTHROPIC_API_KEY=sk-ant-xxx
```

### Testing AWS Services Locally

```typescript
import { S3Client, PutObjectCommand } from '@aws-sdk/client-s3';

// Works with both LocalStack and real AWS
const s3 = new S3Client({
  region: process.env.AWS_REGION,
  endpoint: process.env.AWS_ENDPOINT,  // Only set for LocalStack
  forcePathStyle: true,  // Required for LocalStack
});

await s3.send(new PutObjectCommand({
  Bucket: 'chess-coach-pgn-archives',
  Key: `imported/${userId}/${gameId}.pgn`,
  Body: pgnContent,
}));
```

---

## 4. Message Broker Decision

### Final Recommendation: NATS on ECS

**Architecture:**

```
ECS Cluster
├── chess-coach-backend (2-10 tasks)
│   └── Connects to NATS via internal DNS
└── nats-jetstream (1 task, always on)
    └── Persistent storage via EFS
```

**Why Not MSK (Kafka):**

| Factor | MSK | NATS on ECS |
|--------|-----|-------------|
| **Cost** | $200+/mo | $18/mo |
| **Setup** | Complex (2+ brokers, Zookeeper) | Single container |
| **Latency** | 20-50ms | 5-10ms |
| **Management** | AWS-managed but still complex | Simple config |
| **Overkill?** | Yes for <100K users | Perfect scale |

**When to migrate to MSK:**
- You reach 100K+ concurrent users
- You need multi-region replication
- You want AWS's SLA and support
- Budget >$500/mo for messaging

---

## 5. Networking & Security

### VPC Design

```
VPC: 10.0.0.0/16
├── Public Subnets (ALB)
│   ├── 10.0.1.0/24 (us-east-1a)
│   └── 10.0.2.0/24 (us-east-1b)
├── Private Subnets (ECS, NATS)
│   ├── 10.0.11.0/24 (us-east-1a)
│   └── 10.0.12.0/24 (us-east-1b)
└── Database Subnets (RDS, ElastiCache)
    ├── 10.0.21.0/24 (us-east-1a)
    └── 10.0.22.0/24 (us-east-1b)
```

### Security Groups

```yaml
# ALB Security Group
ALBSecurityGroup:
  Ingress:
    - IpProtocol: tcp
      FromPort: 443
      ToPort: 443
      CidrIp: 0.0.0.0/0  # HTTPS from anywhere
    - IpProtocol: tcp
      FromPort: 80
      ToPort: 80
      CidrIp: 0.0.0.0/0  # HTTP redirect
  Egress:
    - IpProtocol: -1
      CidrIp: 10.0.0.0/16  # To VPC only

# ECS Security Group
ECSSecurityGroup:
  Ingress:
    - IpProtocol: tcp
      FromPort: 3000
      ToPort: 3000
      SourceSecurityGroupId: !Ref ALBSecurityGroup
  Egress:
    - IpProtocol: -1
      CidrIp: 0.0.0.0/0  # Internet for API calls

# RDS Security Group
RDSSecurityGroup:
  Ingress:
    - IpProtocol: tcp
      FromPort: 5432
      ToPort: 5432
      SourceSecurityGroupId: !Ref ECSSecurityGroup
  Egress: []  # No outbound

# ElastiCache Security Group
CacheSecurityGroup:
  Ingress:
    - IpProtocol: tcp
      FromPort: 6379
      ToPort: 6379
      SourceSecurityGroupId: !Ref ECSSecurityGroup
```

### IAM Roles

```yaml
# ECS Task Role (what container can do)
TaskRole:
  AssumeRolePolicyDocument:
    Statement:
      - Effect: Allow
        Principal:
          Service: ecs-tasks.amazonaws.com
        Action: sts:AssumeRole
  Policies:
    - PolicyName: S3Access
      PolicyDocument:
        Statement:
          - Effect: Allow
            Action:
              - s3:GetObject
              - s3:PutObject
            Resource: arn:aws:s3:::chess-coach-pgn-archives/*
    - PolicyName: SecretsAccess
      PolicyDocument:
        Statement:
          - Effect: Allow
            Action:
              - secretsmanager:GetSecretValue
            Resource: arn:aws:secretsmanager:*:*:secret:chess-coach/*
    - PolicyName: CloudWatchMetrics
      PolicyDocument:
        Statement:
          - Effect: Allow
            Action:
              - cloudwatch:PutMetricData
            Resource: "*"

# ECS Execution Role (for ECS to pull image, get secrets)
ExecutionRole:
  ManagedPolicyArns:
    - arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy
```

---

## 6. CI/CD Pipeline

### GitHub Actions + AWS

**Workflow:**

```yaml
# .github/workflows/deploy.yml
name: Deploy to AWS

on:
  push:
    branches: [main]

env:
  AWS_REGION: us-east-1
  ECR_REPOSITORY: chess-coach
  ECS_SERVICE: chess-coach-backend
  ECS_CLUSTER: chess-coach-cluster

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v2
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ${{ env.AWS_REGION }}

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v1

      - name: Build, tag, and push image to ECR
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG .
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG
          docker tag $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG \
            $ECR_REGISTRY/$ECR_REPOSITORY:latest
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:latest

      - name: Update ECS service
        run: |
          aws ecs update-service \
            --cluster $ECS_CLUSTER \
            --service $ECS_SERVICE \
            --force-new-deployment

      - name: Wait for deployment
        run: |
          aws ecs wait services-stable \
            --cluster $ECS_CLUSTER \
            --services $ECS_SERVICE
```

### Deploy Time

- Build + Push to ECR: ~3 minutes
- ECS deployment (rolling): ~2 minutes
- **Total: ~5 minutes** from push to live

---

## 7. Cost Analysis

### MVP (0-10K users)

| Service | Configuration | Cost/Month |
|---------|--------------|------------|
| **ECS Fargate** | 2 tasks (0.5 vCPU, 1GB) | $30 |
| **ALB** | Base + data | $20 |
| **ElastiCache** | t4g.micro | $12 |
| **RDS** | t4g.micro, Single-AZ | $15 |
| **NATS (ECS)** | 1 task | $15 |
| **EFS** | 10GB | $3 |
| **S3** | 100GB | $2 |
| **CloudFront** | 500GB transfer | $5 |
| **CloudWatch** | Logs + metrics | $10 |
| **Secrets Manager** | 5 secrets | $2 |
| **Route 53** | 1 hosted zone | $1 |
| **AI API** | Claude | $100 |
| **Total** | | **~$215/mo** |

### Production (10K-100K users)

| Service | Configuration | Cost/Month |
|---------|--------------|------------|
| **ECS Fargate** | 5 tasks (0.5 vCPU, 1GB) | $75 |
| **ALB** | Higher traffic | $40 |
| **ElastiCache** | t4g.small, Multi-AZ | $50 |
| **RDS** | t4g.small, Multi-AZ | $65 |
| **NATS (ECS)** | 1 task | $15 |
| **EFS** | 50GB | $15 |
| **S3** | 1TB | $23 |
| **CloudFront** | 5TB transfer | $40 |
| **CloudWatch** | Logs + metrics | $30 |
| **Secrets Manager** | 5 secrets | $2 |
| **Route 53** | 1 hosted zone | $1 |
| **AI API** | Claude (volume) | $500 |
| **Total** | | **~$856/mo** |

### Cost Optimization Tips

1. **Reserved Instances:** Save 30-40% on RDS/ElastiCache
2. **Compute Savings Plans:** 15-20% off ECS Fargate
3. **S3 Lifecycle:** Move old PGNs to Glacier ($4/TB/mo)
4. **CloudFront:** Use S3 Transfer Acceleration
5. **Right-size:** Start small, scale based on metrics

---

## 8. Monitoring & Alerting

### CloudWatch Dashboards

```json
{
  "widgets": [
    {
      "type": "metric",
      "properties": {
        "title": "ECS CPU & Memory",
        "metrics": [
          ["AWS/ECS", "CPUUtilization", {"stat": "Average"}],
          [".", "MemoryUtilization", {"stat": "Average"}]
        ],
        "period": 300,
        "region": "us-east-1"
      }
    },
    {
      "type": "metric",
      "properties": {
        "title": "ALB Request Count",
        "metrics": [
          ["AWS/ApplicationELB", "RequestCount", {"stat": "Sum"}]
        ],
        "period": 60
      }
    },
    {
      "type": "metric",
      "properties": {
        "title": "RDS Connections",
        "metrics": [
          ["AWS/RDS", "DatabaseConnections", {"stat": "Average"}]
        ]
      }
    }
  ]
}
```

### SNS Alerts

```yaml
SNSTopic:
  TopicName: chess-coach-alerts
  Subscriptions:
    - Protocol: email
      Endpoint: team@example.com
    - Protocol: sms
      Endpoint: +1234567890

Alarms:
  - HighErrorRate:
      Metric: HTTPCode_Target_5XX_Count
      Threshold: 10
      Period: 300
      Actions: [!Ref SNSTopic]

  - DatabaseHighCPU:
      Metric: CPUUtilization
      Namespace: AWS/RDS
      Threshold: 80
      Period: 300

  - LowDiskSpace:
      Metric: FreeStorageSpace
      Namespace: AWS/RDS
      Threshold: 5GB
      ComparisonOperator: LessThanThreshold
```

---

## 9. Scaling Strategy

### Auto-Scaling Triggers

```yaml
# ECS Service Auto Scaling
ServiceScaling:
  MinCapacity: 2
  MaxCapacity: 10
  TargetCPUUtilization: 70
  TargetMemoryUtilization: 80
  ScaleOutCooldown: 60s
  ScaleInCooldown: 300s

# RDS Read Replicas (future)
ReadReplicas:
  Count: 2  # For read-heavy workloads
  InstanceClass: db.t4g.small
```

### Scaling Timeline

| Users | ECS Tasks | RDS | ElastiCache | NATS |
|-------|-----------|-----|-------------|------|
| 0-10K | 2 | t4g.micro | t4g.micro | 1 task |
| 10K-50K | 5 | t4g.small Multi-AZ | t4g.small Multi-AZ | 1 task |
| 50K-100K | 10 | t4g.medium Multi-AZ | t4g.medium Multi-AZ | 2 tasks |
| 100K+ | 20+ | r6g.large Multi-AZ + Replicas | r6g.large Multi-AZ | MSK (3 brokers) |

---

## 10. Disaster Recovery

### Backup Strategy

**RDS:**
- Automated backups: 7 days retention
- Manual snapshots: Before major changes
- Point-in-time recovery: Up to 5 minutes

**ElastiCache:**
- Daily snapshots: 3 days retention
- Manual snapshots: Before updates

**S3:**
- Versioning enabled
- Cross-region replication (optional)
- Lifecycle policies

**Recovery Time Objective (RTO):** 1 hour
**Recovery Point Objective (RPO):** 5 minutes

### Disaster Recovery Runbook

```bash
# 1. Restore RDS from snapshot
aws rds restore-db-instance-from-db-snapshot \
  --db-instance-identifier chess-coach-restored \
  --db-snapshot-identifier chess-coach-snap-2025-10-31

# 2. Update ECS task definition with new DB endpoint
aws ecs register-task-definition \
  --cli-input-json file://task-def-updated.json

# 3. Update ECS service
aws ecs update-service \
  --cluster chess-coach-cluster \
  --service chess-coach-backend \
  --task-definition chess-coach-backend:NEW_VERSION

# 4. Verify health
aws ecs describe-services \
  --cluster chess-coach-cluster \
  --services chess-coach-backend
```

---

## 11. Step-by-Step Implementation

### Phase 1: Foundation (Week 1)

**Day 1-2: AWS Account Setup**
- [ ] Create AWS account
- [ ] Set up billing alerts
- [ ] Configure IAM users (no root access)
- [ ] Enable CloudTrail
- [ ] Set up VPC with subnets

**Day 3-4: Database & Cache**
- [ ] Create RDS PostgreSQL (t4g.micro)
- [ ] Create ElastiCache Redis (t4g.micro)
- [ ] Run migrations
- [ ] Test connectivity from local

**Day 5-7: Container Setup**
- [ ] Create ECR repository
- [ ] Build Docker image
- [ ] Push to ECR
- [ ] Create ECS cluster
- [ ] Deploy NATS task (with EFS)

### Phase 2: Application (Week 2)

**Day 1-3: Backend Deployment**
- [ ] Create ALB with target group
- [ ] Configure SSL certificate (ACM)
- [ ] Deploy ECS service
- [ ] Test WebSocket connection
- [ ] Configure auto-scaling

**Day 4-5: Frontend Deployment**
- [ ] Build React app
- [ ] Upload to S3
- [ ] Configure CloudFront
- [ ] Update DNS (Route 53)

**Day 6-7: Integration**
- [ ] Connect backend to NATS
- [ ] Connect to RDS
- [ ] Test end-to-end flow
- [ ] Load testing (100 concurrent users)

### Phase 3: Production Hardening (Week 3)

**Day 1-2: Security**
- [ ] Move secrets to Secrets Manager
- [ ] Enable VPC Flow Logs
- [ ] Configure WAF (optional)
- [ ] Set up security groups

**Day 3-4: Monitoring**
- [ ] Create CloudWatch dashboards
- [ ] Set up alarms
- [ ] Configure SNS notifications
- [ ] Enable X-Ray tracing (optional)

**Day 5-7: CI/CD**
- [ ] Set up GitHub Actions
- [ ] Test deploy pipeline
- [ ] Document rollback procedure
- [ ] Create runbooks

### Phase 4: Launch (Week 4)

**Day 1-3: Beta Testing**
- [ ] Invite beta users
- [ ] Monitor metrics
- [ ] Fix critical bugs
- [ ] Load test (1K concurrent)

**Day 4-5: Go Live**
- [ ] Public launch
- [ ] Monitor closely
- [ ] Be ready for hotfixes

**Day 6-7: Optimization**
- [ ] Analyze metrics
- [ ] Optimize queries
- [ ] Right-size resources
- [ ] Plan next features

---

## Summary: AWS Architecture Checklist

### ✅ Key Services

| Service | Purpose | Config |
|---------|---------|--------|
| ECS Fargate | Backend hosting | 2-10 tasks, auto-scale |
| ALB | Load balancing + WebSocket | HTTPS, sticky sessions |
| ElastiCache | Real-time state | Redis, Multi-AZ (prod) |
| RDS | Data persistence | PostgreSQL, Multi-AZ (prod) |
| NATS on ECS | Messaging | JetStream + EFS |
| S3 | Game archives | Lifecycle policies |
| CloudFront | Frontend CDN | Global edge |
| CloudWatch | Monitoring | Logs + metrics + alarms |

### 💰 Cost Summary

- **MVP:** ~$215/mo (0-10K users)
- **Growth:** ~$856/mo (10K-100K users)
- **Scale:** ~$2K/mo (100K+ users)

### 🚀 Small Team Benefits

- ✅ No Kubernetes complexity
- ✅ AWS-managed services (less ops)
- ✅ Auto-scaling built-in
- ✅ Local development mirrors AWS
- ✅ 2-3 developers can manage

### 📝 Next Actions

1. Review this architecture with team
2. Set up AWS account and VPC
3. Start with Phase 1 (Foundation)
4. Deploy MVP in 3-4 weeks

---

**Document Version:** 1.0
**Last Updated:** October 31, 2025
**Cloud Provider:** AWS
**Team Size:** 2-3 developers

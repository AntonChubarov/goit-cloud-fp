# Project description

**GoIT Cloud Computing Course Final Project**

The application is a full-stack URL Shortener built with Golang and React.
The backend is a Golang web server that embeds and serves the compiled
React frontend using Go's runtime-embedded filesystem. It also manages
URL data storage in a PostgreSQL database and automatically applies
schema migrations at startup.

# Project Setup and Deployment Instructions

## Prerequisites
Ensure the following tools are installed and configured:
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html)
- [Docker](https://docs.docker.com/get-docker/)
- [cfn-lint](https://github.com/aws-cloudformation/cfn-lint)
- AWS CLI profile named `default` configured with appropriate credentials

## Local Development
To run the service locally using Docker Compose:

```bash
make run-locally
```

Visit [http://localhost:8080](http://localhost:8080) in your browser.

To stop the local service:

```bash
make stop-local-process
```

## Container Image Management

### Login to ECR
```bash
make login
```

### Create ECR Repository (if not already present)
```bash
make create-ecr
```

### Build Docker Image
```bash
make build
```

### Tag Docker Image
```bash
make tag
```

### Push Image to ECR
```bash
make push
```

### Combined Image Preparation Step
This step performs login, creates the repo (if needed), builds, tags, and pushes the image:
```bash
make prepare-image
```

## CloudFormation Template Linting
```bash
make lint
```

## Deployment to AWS
Deploy the CloudFormation stack:

```bash
make deploy
```

Monitor the deployment in AWS Console:
> CloudFormation → Stacks → `goit-cc-fp-stack` → Resources tab

## Stack Deletion
To tear down the stack and clean up resources:

```bash
make destroy
```
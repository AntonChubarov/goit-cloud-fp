APP_NAME=goit-cc-fp
AWS_REGION=us-east-1
AWS_PROFILE=default
STACK_NAME=$(APP_NAME)-stack
TEMPLATE_FILE=cloudformation.yml

ACCOUNT_ID=$(shell aws sts get-caller-identity --query Account --output text --profile $(AWS_PROFILE))
ECR_REPO=$(APP_NAME)
ECR_URI=$(ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/$(ECR_REPO)
IMAGE_TAG=latest

.PHONY: all login build push create-ecr deploy destroy

all: push deploy

login:
	@echo "Logging into ECR..."
	aws ecr get-login-password --region $(AWS_REGION) --profile $(AWS_PROFILE) | \
		docker login --username AWS --password-stdin $(ECR_URI)

create-ecr:
	@echo "Creating ECR repository if it doesn't exist..."
	@aws ecr describe-repositories --repository-names $(ECR_REPO) --region $(AWS_REGION) --profile $(AWS_PROFILE) >/dev/null 2>&1 || \
		aws ecr create-repository --repository-name $(ECR_REPO) --region $(AWS_REGION) --profile $(AWS_PROFILE) >/dev/null 2>&1

build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(IMAGE_TAG) .

tag:
	@echo "Tagging image for ECR..."
	docker tag $(APP_NAME):$(IMAGE_TAG) $(ECR_URI):$(IMAGE_TAG)

push: create-ecr login build tag
	@echo "Pushing image to ECR..."
	docker push $(ECR_URI):$(IMAGE_TAG)
	@echo "Image pushed: $(ECR_URI):$(IMAGE_TAG)"

lint:
	@echo "Linting CloudFormation template..."
	cfn-lint $(TEMPLATE_FILE)

deploy:
	@echo "Deploying CloudFormation stack: $(STACK_NAME)"
	aws cloudformation deploy \
		--stack-name $(STACK_NAME) \
		--template-file $(TEMPLATE_FILE) \
		--capabilities CAPABILITY_NAMED_IAM \
		--region $(AWS_REGION) \
		--profile $(AWS_PROFILE)
	@echo "Stack deployed: $(STACK_NAME)"

destroy:
	@echo "Deleting CloudFormation stack: $(STACK_NAME)"
	aws cloudformation delete-stack \
		--stack-name $(STACK_NAME) \
		--region $(AWS_REGION) \
		--profile $(AWS_PROFILE)
	@echo "Waiting for stack deletion..."
	aws cloudformation wait stack-delete-complete \
		--stack-name $(STACK_NAME) \
		--region $(AWS_REGION) \
		--profile $(AWS_PROFILE)
	@echo "Stack deleted: $(STACK_NAME)"

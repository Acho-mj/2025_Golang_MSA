.PHONY: help aws-login-admin aws-login-dev ecr-login docker-build-order docker-build-user docker-build docker-push-order docker-push-user docker-push helm-deploy-order helm-deploy-user helm-deploy kubeconfig

AWS_ACCOUNT_ID ?= 052747538895
AWS_REGION ?= ap-northeast-2
PROFILE_ADMIN ?= saas-dev-admin
PROFILE_DEV ?= saas-dev

IMAGE_TAG ?= $(shell git rev-parse --short HEAD)

ORDER_SERVICE_NAME ?= order-service
USER_SERVICE_NAME ?= user-service

ORDER_SERVICE_DIR ?= backend/services/order
USER_SERVICE_DIR ?= backend/services/user

ORDER_CHART_PATH ?= deploy/helm/order
USER_CHART_PATH ?= deploy/helm/user

KUBE_NAMESPACE ?= default
EKS_CLUSTER_NAME ?= saas-dev-cluster

ECR_REGISTRY := $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
ORDER_IMAGE := $(ECR_REGISTRY)/$(ORDER_SERVICE_NAME):$(IMAGE_TAG)
USER_IMAGE := $(ECR_REGISTRY)/$(USER_SERVICE_NAME):$(IMAGE_TAG)

help:
	@echo "사용 가능한 타겟:"
	@echo "  aws-login-admin     - $(PROFILE_ADMIN) 프로파일로 AWS SSO 로그인"
	@echo "  aws-login-dev       - $(PROFILE_DEV) 프로파일로 AWS SSO 로그인"
	@echo "  ecr-login           - ECR 로그인 (admin 프로파일)"
	@echo "  docker-build        - order/user 서비스 Docker 이미지 빌드"
	@echo "  docker-push         - order/user 서비스 Docker 이미지 ECR 푸시"
	@echo "  helm-deploy         - order/user Helm 차트 배포/업데이트"
	@echo "  kubeconfig          - EKS kubeconfig 업데이트"

aws-login-admin:
	aws sso login --profile $(PROFILE_ADMIN)

aws-login-dev:
	aws sso login --profile $(PROFILE_DEV)

ecr-login: aws-login-admin
	aws ecr get-login-password --profile $(PROFILE_ADMIN) --region $(AWS_REGION) | docker login --username AWS --password-stdin $(ECR_REGISTRY)

docker-build-order:
	docker build \
		-f $(ORDER_SERVICE_DIR)/Dockerfile \
		-t $(ORDER_SERVICE_NAME):$(IMAGE_TAG) \
		-t $(ORDER_IMAGE) \
		.

docker-build-user:
	docker build \
		-f $(USER_SERVICE_DIR)/Dockerfile \
		-t $(USER_SERVICE_NAME):$(IMAGE_TAG) \
		-t $(USER_IMAGE) \
		.

docker-build: docker-build-order docker-build-user

docker-push-order: docker-build-order ecr-login
	docker push $(ORDER_IMAGE)

docker-push-user: docker-build-user ecr-login
	docker push $(USER_IMAGE)

docker-push: docker-push-order docker-push-user

helm-deploy-order:
	helm upgrade --install $(ORDER_SERVICE_NAME) $(ORDER_CHART_PATH) \
		--namespace $(KUBE_NAMESPACE) \
		--set image.repository=$(ECR_REGISTRY)/$(ORDER_SERVICE_NAME) \
		--set image.tag=$(IMAGE_TAG)

helm-deploy-user:
	helm upgrade --install $(USER_SERVICE_NAME) $(USER_CHART_PATH) \
		--namespace $(KUBE_NAMESPACE) \
		--set image.repository=$(ECR_REGISTRY)/$(USER_SERVICE_NAME) \
		--set image.tag=$(IMAGE_TAG)

helm-deploy: helm-deploy-order helm-deploy-user

kubeconfig: aws-login-dev
	aws eks update-kubeconfig \
		--name $(EKS_CLUSTER_NAME) \
		--region $(AWS_REGION) \
		--profile $(PROFILE_DEV)


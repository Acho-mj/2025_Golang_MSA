# 2025 Golang MSA

Go 기반 마이크로서비스 아키텍처 실습 프로젝트입니다. 주문(`order-service`)과 사용자(`user-service`) 두 개의 서비스를 중심으로 하며, 다음과 같은 기술 스택을 사용합니다.

- **Connect RPC (gRPC compatible)**: 서비스 간 통신을 위한 RPC 프레임워크
- **Amazon DynamoDB**: 주문/사용자 데이터를 저장하는 NoSQL 데이터베이스
- **Amazon ECR + EKS**: Docker 이미지 관리와 Kubernetes 배포 환경
- **Helm**: Kubernetes 리소스를 선언적으로 배포
- **IRSA (IAM Roles for Service Accounts)**: 파드별로 AWS 권한을 분리

## 구조 요약

| 계층 | 구성 요소 | 설명 |
| --- | --- | --- |
| 소스/빌드 | Makefile | `docker-push`, `helm-deploy`, `kubeconfig` 등 배포 자동화 명령 제공 |
| 컨테이너 레지스트리 | Amazon ECR | `order-service`, `user-service` Docker 이미지 저장소 |
| 배포 플랫폼 | Amazon EKS | Helm으로 배포된 Pod, Service가 실행되는 쿠버네티스 클러스터 |
| 서비스 디스커버리 | Kubernetes Service | `order-service-order-service`, `user-service-user-service` ClusterIP 제공 |
| 서비스 간 통신 | Connect RPC | `order-service` → `user-service` RPC 호출 (USER_SERVICE_URL 환경 변수 기반) |
| 데이터 저장소 | DynamoDB | `order`/`user` 테이블, IRSA (`eks-dynamodb-role-irsa`)로 접근 제어 |

## 주요 기능

- 주문 서비스는 사용자 서비스를 RPC로 호출하여 사용자 정보를 검증한 뒤 주문을 생성합니다.
- `make docker-push` 및 `make helm-deploy`를 통해 이미지 빌드/푸시와 배포를 자동화할 수 있습니다.
- Helm 차트(`deploy/helm/order`, `deploy/helm/user`)에서 환경 변수, 리소스 한도, 프로브 등을 쉽게 조정할 수 있습니다.

## 로컬/배포 워크플로우

1. **ECR에 이미지 푸시**
   ```bash
   make aws-login-admin
   make ecr-login
   make docker-push
   ```
2. **EKS 배포**
   ```bash
   make aws-login-dev
   make kubeconfig
   make helm-deploy
   ```
3. **테스트 (DNS & RPC)**
   ```bash
   kubectl run curl-test --image=curlimages/curl --rm -it -n default -- sh
   curl -s http://user-service-user-service.default.svc.cluster.local:8080/healthz
   curl -s -X POST -H "Content-Type: application/json" \
     -d '{"user_id":"user1"}' \
     http://user-service-user-service.default.svc.cluster.local:8080/user.UserService/GetUser
   ```

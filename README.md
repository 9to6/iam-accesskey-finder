# IAM Accesskey finder

IAM access key의 생성시간관리를 위한 API

## 사용기술

- golang 1.22 + gin
- aws sdk v2
- makefile

## 지원 가능 아키텍쳐

- arm64
- amd64

## 리포지토리 설명

```
├── Dockerfile
├── Makefile
├── README.md
├── _k8s
│   ├── deployment.yaml
│   ├── service.yaml
│   └── serviceaccount.yaml
├── aws.go
├── go.mod
├── go.sum
├── main.go

```

- golang 파일
  - 아래 두개의 파일로 구성
  - `main.go`
  - `aws.go`
    - AWS 리소스 제어를 위한 코드 
- Kubernetes manifests
  - `_k8s` 폴더 참조
- 기타
  - `Dockerfile`
  - `Makefile`

## Endpoint

`GET` method 만 지원

- `/health-check`
- `/expired-keys`

## 배포방법

도커 배포는 개인 `dockerhub` 계정을 이용한다.  

https://hub.docker.com/repository/docker/9to5/iam-accesskey-finder/general

### 도커 빌드
```sh
make docker-build
```

### 도커 배포
```sh
make docker-push
make docker-manifest
make docker-manifest-push
```

### EKS 배포

#### 배포시 주의사항
1. `ServiceAccount`의 `annotations`에 **IRSA**(IAM Roles for Service Account) 적용된 Role ARN을 추가한다.

```sh
kubectl apply -f _k8s
```
2. Access Key 만료시간 기준 변경하기  

- API 실행시 환경변수를 통해 값을 가져온다.  
- 지원 포멧은 [이 곳](https://pkg.go.dev/time#ParseDuration)을 통해 확인한다.  
- [해당 변수 확인](https://github.com/9to6/iam-accesskey-finder/blob/8b8e23b4cacacb7cffa084633a38380dc3418732/_k8s/deployment.yaml#L20)

> kubernetes manifest 적용후 EKS에 `80` 포트로 엔드포인트가 생성되고 서비스가 외부 노출된다.

## 로컬 테스트

1. 환경변수 등록  

`AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `ACCESS_KEY_EXPIRE_TIME` 를 환경변수에 등록

```sh
export ACCESS_KEY_EXPIRE_TIME=1440h
export AWS_REGION=ap-northeast-2
export AWS_ACCESS_KEY_ID=AKIblahblah
export AWS_SECRET_ACCESS_KEY=secret
```

2. 도커 실행  
```sh
docker run -it --rm -p8080:8080 -e "ACCESS_KEY_EXPIRE_TIME=$ACCESS_KEY_EXPIRE_TIME" -e "AWS_REGION=$AWS_REGION" -e "AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID" -e "AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY" 9to5/iam-accesskey-finder
```

3. API call  

```sh
curl localhost:8080/expired-keys | jq '.keys|length'  
```

> 16개 검색됨

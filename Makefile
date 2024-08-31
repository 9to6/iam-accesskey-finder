APP := iam-accesskey-finder
REPO := 9to5/$(APP)

IMAGE_TAG := latest
URL := $(REPO):$(IMAGE_TAG)

docker-build: docker-build-amd64 docker-build-arm64

docker-push:
	docker push $(URL)-arm64
	docker push $(URL)-amd64

docker-build-amd64:
	docker build --platform linux/amd64 -t $(URL)-amd64 .

docker-build-arm64:
	docker build --platform linux/arm64 -t $(URL)-arm64 .

docker-manifest-push: docker-push
	docker manifest push $(REPO):$(IMAGE_TAG)

docker-manifest:
	docker manifest create $(REPO):$(IMAGE_TAG) --amend $(URL)-arm64 --amend $(URL)-amd64

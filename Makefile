export AWS_PROFILE := $(AWS_PROFILE)
export AWS_DEFAULT_REGION := $(shell aws configure get region --profile $(AWS_PROFILE))

YAML := config.yaml
LAMBDA := "lambda.go"
INTERVAL ?= "rate(15 minutes)"
TF_ARGS := -var=interval=$(INTERVAL)

build:
	GOOS=linux go build $(LAMBDA)
	GOOS=linux go build -o $$GOPATH/bin/checker main.go
	zip terraform/checker.zip lambda $(YAML)
	rm -f lambda

tf_apply: tf_plan
	cd terraform && terraform_13 apply $(TF_ARGS) -auto-approve

tf_destroy:
	cd terraform && terraform_13 destroy $(TF_ARGS)

tf_plan:
	cd terraform && terraform_13 plan $(TF_ARGS)

tf_clean:
	rm -f terraform/checker.zip

destroy: tf_destroy tf_clean encrypt
apply: build tf_apply tf_clean encrypt
plan: build tf_plan

encrypt:
	git secret hide -m

decrypt:
	git secret reveal -f

push: encrypt
	git add .
	git commit
	git push

prepare:
	git secret remove -c
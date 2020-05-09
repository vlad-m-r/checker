export AWS_PROFILE := $(AWS_PROFILE)
export AWS_DEFAULT_REGION := $(shell aws configure get region --profile $(AWS_PROFILE))

YAML := config.yaml
LAMBDA := "lambda.go"
INTERVAL ?= "rate(15 minutes)"
TF_ARGS := -var=interval=$(INTERVAL)

build: decrypt
	GOOS=linux go build $(LAMBDA)
	zip terraform/checker.zip lambda $(YAML)
	rm -f lambda

tf_apply: tf_plan
	cd terraform && terraform apply $(TF_ARGS) -auto-approve

tf_destroy:
	cd terraform && terraform destroy $(TF_ARGS)

tf_plan:
	cd terraform && terraform plan $(TF_ARGS)

tf_clean:
	rm -f terraform/checker.zip

destroy: tf_destroy tf_clean encrypt
apply: build tf_apply tf_clean encrypt
plan: build tf_plan

encrypt:
	git secret hide -m

decrypt: encrypt
	git secret reveal -f

push: encrypt
	git add .
	git commit
	git push

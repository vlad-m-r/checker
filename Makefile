export AWS_PROFILE := $(AWS_PROFILE)
export AWS_DEFAULT_REGION := $(shell aws configure get region --profile $(AWS_PROFILE))

YAML := config.yaml
LAMBDA := "lambda.go"
INTERVAL ?= "rate(15 minutes)"
TF_ARGS := -var=interval=$(INTERVAL)

build:
	GOOS=linux go build $(LAMBDA)
	zip checker.zip lambda $(YAML)
	rm -f lambda
	mv checker.zip terraform

tf_apply: build
	cd terraform && terraform apply $(TF_ARGS) -auto-approve

tf_destroy:
	cd terraform && terraform destroy $(TF_ARGS)

tf_plan: build
	cd terraform && terraform plan $(TF_ARGS)

tf_clean:
	cd terraform && rm -f checker.zip

clean_tf_state:
	rm -f ./terraform/terraform.tfstate-*
	rm -f ./terraform/terraform.tfstate.backup-*

destroy: tf_destroy clean_tf_state tf_clean encrypt
apply: tf_apply clean_tf_state tf_clean encrypt
plan: tf_plan

CONFIG := $(YAML)-$(shell md5sum $(YAML))
TF_STATE  := ./terraform/terraform.tfstate-$(shell md5sum ./terraform/terraform.tfstate)
TF_STATE_BACKUP := ./terraform/terraform.tfstate.backup-$(shell md5sum ./terraform/terraform.tfstate.backup)

.PHONY: encrypt
encrypt: $(CONFIG) $(TF_STATE) $(TF_STATE_BACKUP)

$(TF_STATE):
	git secret add ./terraform/terraform.tfstate
	git secret hide
	git add $@
	touch $@

$(TF_STATE_BACKUP):
	git secret add ./terraform/terraform.tfstate.backup
	git secret hide
	git add $@
	touch $@

$(CONFIG):
	rm -f $(YAML)-*
	git secret add $(YAML)
	git secret hide
	git add $@
	touch $@

decrypt:
	git secret reveal -f

push: encrypt
	git add .
	git commit
	git push

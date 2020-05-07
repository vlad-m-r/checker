export AWS_PROFILE := $(AWS_PROFILE)
export AWS_DEFAULT_REGION := $(shell aws configure get region --profile $(AWS_PROFILE))

CONFIG := config.yaml
LAMBDA := "lambda.go"
INTERVAL ?= "rate(15 minutes)"

build: decrypt
	GOOS=linux go build $(LAMBDA)
	zip checker.zip lambda $(CONFIG)
	rm -f lambda
	mv checker.zip terraform

apply: build
	cd terraform && terraform apply -var=interval=$(INTERVAL) -auto-approve && rm -f checker.zip

destroy:
	cd terraform && terraform destroy

plan: build
	cd terraform && terraform plan -var=interval=$(INTERVAL) && rm -f checker.zip

CONFIG := $(CONFIG)-$(shell md5sum $(CONFIG))
TF_STATE  := ./terraform/terraform.tfstate-$(shell md5sum ./terraform/terraform.tfstate)
TF_STATE_BACKUP := ./terraform/terraform.tfstate.backup-$(shell md5sum ./terraform/terraform.tfstate.backup)

.PHONY: encrypt
encrypt: $(CHECKSUM) $(TF_STATE) $(TF_STATE_BACKUP)

$(TF_STATE):
	git secret add ./terraform/terraform.tfstate
	git secret hide
	rm -f ./terraform/terraform.tfstate-*
	touch $@

$(TF_STATE_BACKUP):
	git secret add ./terraform/terraform.tfstate.backup
	git secret hide
	rm -f ./terraform/terraform.tfstate.backup-*
	touch $@

$(CHECKSUM):
	git secret add $(CONFIG)
	git secret hide
	rm -f ./$(CONFIG)-*
	touch $@

decrypt:
	git secret reveal -f

push: encrypt
	git add .
	git commit
	git push

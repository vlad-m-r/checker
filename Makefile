export AWS_PROFILE := $(AWS_PROFILE)
export AWS_DEFAULT_REGION := $(shell aws configure get region --profile $(AWS_PROFILE))

INTERVAL ?= "*/10 * * * *"

build: decrypt
	GOOS=linux go build main.go
	zip checker.zip main config.yaml
	rm -f main
	mv checker.zip terraform

apply: build
	cd terraform && terraform apply -var=interval=$(INTERVAL) -auto-approve && rm -f checker.zip

destroy:
	cd terraform && terraform destroy

plan: build
	cd terraform && terraform plan -var=interval=$(INTERVAL) && rm -f checker.zip

encrypt:
	git secret add config.yaml
	git secret hide

decrypt:
	git secret reveal -f


push: encrypt
	git add .
	git commit -F

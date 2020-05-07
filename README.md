# Summary

checker is a simple app for http endpoint monitoring controlled by a yaml configuration file

Does http requests against preconfigured endpoints and notifies via email in case of failures
 
Has an integration with AWS to send emails via SES

Intended to be run as a lambda in AWS with Cloudwatch event rules as a cron 

## Requirements

For emails:
* Requires a verified email address or domain in AWS SES

For AWS

## Setup

```shell script
git clone git@github.com:vlad-m-r/checker.git
```

```shell script
go mod vendor
```

## Prepare config

The app requires yaml configuration yaml to function





## Usage

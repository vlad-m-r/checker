package controllers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/vlad-m-r/checker/api/models"
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

const (
	CharSet = "UTF-8"
)

func NewAwsClient(yamlContent []byte) *AwsClient {
	awsClient := AwsClient{}
	awsClient.UnmarshalYaml(yamlContent)
	awsClient.createSession()
	return &awsClient
}

type AwsClient struct {
	Yaml        models.AwsClientConfig `yaml:"AwsClient"`
	session     client.ConfigProvider
	EmailClient *EmailClient
}

func (a *AwsClient) UnmarshalYaml(yamlContent []byte) {
	err := yaml.Unmarshal(yamlContent, &a)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
}

func (a *AwsClient) createSession() {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(a.Yaml.Region),
		Credentials: credentials.NewSharedCredentials(a.Yaml.CredsFile, a.Yaml.Profile),
	})

	if err != nil {
		log.Fatal("Failed to create AWS credentials session:", err)
		return
	}

	a.session = s
}

func (a *AwsClient) sendSESMail(results []*CheckResult, emailClient *EmailClient) {
	log.Println("Sending AWS email")

	// Create an SES session.
	svc := ses.New(a.session)

	// Assemble the email.
	input := a.constructEmail(results, emailClient)

	log.Println(input)

	// Attempt to send the email.
	_, err := svc.SendEmail(input)

	// Display error messages if they occur.
	a.verifyEmailResponse(err)

}

func (a *AwsClient) constructEmail(results []*CheckResult, emailClient *EmailClient) *ses.SendEmailInput {
	var messages []string

	for _, checkResult := range results {
		for _, err := range checkResult.err {
			messages = append(messages, checkResult.Name+":"+err.Error())
		}
	}

	emailInput := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{
				aws.String(emailClient.Yaml.CC),
			},
			ToAddresses: []*string{
				aws.String(emailClient.Yaml.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(strings.Join(messages, "\n")),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(emailClient.Yaml.Subject),
			},
		},
		Source: aws.String(emailClient.Yaml.From),
	}

	return emailInput
}

func (a *AwsClient) verifyEmailResponse(err error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return
	}
}
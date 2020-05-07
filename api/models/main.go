package models

import "time"

type LambdaEvent struct {
	Name string `json:"name"`
}

type EmailClientConfig struct {
	From    string `yaml:"From"`
	To      string `yaml:"To"`
	CC      string `yaml:"CC"`
	Subject string `yaml:"Subject"`
}

type ConcurrencyConfig struct {
	Limit int `yaml:"Limit"`
}

type EmailClient struct {
	Yaml EmailClientConfig `yaml:"EmailClient"`
}

type AwsClientConfig struct {
	Profile   string `yaml:"profile"`
	Region    string `yaml:"region"`
	CredsFile string `yaml:"credsfile"`
}

type AwsClient struct {
	Yaml AwsClientConfig `yaml:"AwsClient"`
}

type HttpClientConfig struct {
	InsecureSkipVerify    bool          `yaml:"InsecureSkipVerify"`
	ResponseHeaderTimeout time.Duration `yaml:"ResponseHeaderTimeout"`
	ConcurrencyLimit      int           `yaml:"ConcurrencyLimit"`
}

type HttpClient struct {
	HttpClient HttpClientConfig `yaml:"HttpClient"`
}

type Check struct {
	Name     string    `yaml:"name"`
	URL      string    `yaml:"url"`
	Requests []Request `yaml:"requests"`
}

type Request struct {
	Method  string   `yaml:"method"`
	Payload string   `yaml:"payload"`
	Asserts []Assert `yaml:"asserts"`
}

type Assert struct {
	Type string `yaml:"type"`
	Key  string `yaml:"key"`
}

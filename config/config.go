package config

import (
	"strings"
)

// Just some ideas. It maybe better, instead of having a Config as a
// Name/Value, to really have the Config structure generated, and have
// the environment variable names as variables in Go.

// Pending to think how to generate this file (and the code to
// load it from env variables)

// Here a simple example to put env things into structs:
//https://github.com/02amanag/environment/blob/main/environment.go

func BuildCharmConfig(environ []string) (config CharmConfig) {
	// TODO done in a simple way. Maybe using struct tags like `env: NAME` or
	// any other idea?

	var vars map[string]string = make(map[string]string)
	for _, fullvar := range environ {
		name, value, found := strings.Cut(fullvar, "=")
		if found {
			vars[name] = value
		}
	}

	if val, ok := vars["POSTGRESQL_DB_CONNECT_STRING"]; ok {
		config.Integrations.PostgresqlUrl = &val
		delete(vars, "POSTGRESQL_DB_CONNECT_STRING")
	}

	// Put all variables in Configs for now.
	for name, value := range vars {
		config.Configs = append(config.Configs, Config{Name: name, Value: value})
	}

	return
}

type CharmConfig struct {
	Configs      []Config
	Integrations Integrations
}

type Config struct {
	Name  string
	Value string
}

type Integrations struct {
	MongoUrl      *string
	PostgresqlUrl *string
	RedisUrl      *string
	MysqlUrl      *string
	S3            *S3Config
}

type S3Config struct {
	AccessKey       string
	SecretKey       string
	Bucket          string
	Region          *string
	StorageClass    *string
	Endpoint        *string
	Path            *string
	ApiVersion      *string
	UriStyle        *string
	AddressingStyle *string
	Attributes      *string // this is really a dict
	TlsCaChain      *string // this is really a slice of strings
}

type SamlConfig struct {
	EntityID                string
	MetadataURL             string
	SigningCertificate      string
	SingleSignOnRedirectURL string
}

// This are the current env values for integrations:
// REDIS_DB_CONNECT_STRING
// MYSQL_DB_CONNECT_STRING
// POSTGRESQL_DB_CONNECT_STRING
// MONGODB_DB_CONNECT_STRING
// S3_ACCESS_KEY
// S3_SECRET_KEY
// S3_REGION
// S3_STORAGE_CLASS
// S3_BUCKET
// S3_ENDPOINT
// S3_PATH
// S3_API_VERSION
// S3_URI_STYLE
// S3_ADDRESSING_STYLE
// S3_ATTRIBUTES
// S3_TLS_CA_CHAIN
// SAML_ENTITY_ID
// SAML_METADATA_URL
// SAML_SINGLE_SIGN_ON_REDIRECT_URL
// SAML_SIGNING_CERTIFICATE

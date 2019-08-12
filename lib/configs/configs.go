package configs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var sess = session.Must(session.NewSession())

// Get returns the value of a config key
func Get(configKey, defaultConfigValue string, required bool) (configValue string) {

	ssmPath := filepath.Join("/", EnvOr("SSM_PATH", ""), configKey)

	if os.Getenv("SSM_PATH") != "" {
		// Check for path in SSM
		svc := ssm.New(sess, &aws.Config{
			Region: aws.String(EnvOr("AWS_DEFAULT_REGION", "us-east-1")),
		})
		o, err := svc.GetParameter(&ssm.GetParameterInput{
			Name:           &ssmPath,
			WithDecryption: aws.Bool(true),
		})

		if err != nil {
			fmt.Printf("Config %s is not in SSM\n", ssmPath)
		}

		if o.Parameter != nil {
			configValue = *o.Parameter.Value
		}
	}

	// Override with environment variable if it is set
	configValue = EnvOr(configKey, configValue)

	if configValue == "" {
		configValue = defaultConfigValue
	}

	if configValue == "" && required {
		log.Fatalf("Configuration value \"%s\" required but not set", configKey)
	}

	return configValue
}

// EnvOr returns the OS environment variable's value or a default value
func EnvOr(s, e string) string {
	envVar := os.Getenv(s)
	if envVar != "" {
		return envVar
	}
	return e
}

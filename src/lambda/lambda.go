package lambda

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"os"
	"strings"
)

type RepoPayload struct {
	RepoID   int
	URL      string
	RepoName string
	RepoPass string
	Branch   string
}

func Trigger(payloadValues interface{}, funcName string, awsRegion string) error {
	payload, err := json.Marshal(payloadValues)
	if err != nil {
		return err
	}
	mySession := session.Must(session.NewSession())
	svc := lambda.New(mySession, aws.NewConfig().WithRegion(awsRegion))
	input := &lambda.InvokeInput{
		InvocationType: aws.String("Event"),
		FunctionName:   aws.String(funcName),
		Payload:        payload,
		LogType:        aws.String("Tail"),
	}
	_, err = svc.Invoke(input)
	return err
}

func GetLoadFullRepoLambdaFunc() string {
	env := os.Getenv("GW_DEPLOY_ENV")
	if env == "" {
		env = "dev"
	}
	return "gitwize-lambda-" + strings.ToLower(env) + "-load_full_one_repo"
}

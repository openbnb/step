package to

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws/arn"
)

// LambdaArn takes a name OR arn and returns Arn defaulting to AWS Environment variables
func LambdaArn(region *string, account_id *string, name_or_arn *string) *string {
	return createArn("arn:aws:lambda:%v:%v:function:%v", region, account_id, name_or_arn)
}

// StepArn takes a name OR arn and returns Arn defaulting to AWS Environment variables
func StepArn(region *string, account_id *string, name_or_arn *string) *string {
	return createArn("arn:aws:states:%v:%v:stateMachine:%v", region, account_id, name_or_arn)
}

func RoleArn(account_id *string, name_or_arn *string) *string {
	return createArn("arn:aws:iam::%v%v:role/%v", account_id, Strp(""), name_or_arn)
}

// InterpolateArnVariables replaces any resource parameter templates with the appropriate values
func InterpolateArnVariables(state_machine *string, region *string, account_id *string, name_or_arn *string) *string {
	variableTemplate := map[string]*string{
		"{{aws_account}}": account_id,
		"{{aws_region}}":  region,
		"{{lambda_name}}": name_or_arn,
	}
	for k, v := range variableTemplate {
		*state_machine = strings.Replace(*state_machine, k, *v, -1)
	}
	return state_machine
}

func ArnPath(arn string) string {
	_, _, res := ArnRegionAccountResource(arn)

	path := strings.Split(res, "/")

	switch len(path) {
	case 0:
		return "/"
	case 1:
		return "/"
	case 2:
		return "/"
	default:
		return fmt.Sprintf("/%v/", strings.Join(path[1:len(path)-1], "/"))
	}
}

func AwsRegionAccountFromContext(ctx context.Context) (*string, *string) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return nil, nil
	}

	if lc == nil {
		return nil, nil
	}

	arn := lc.InvokedFunctionArn
	region, account, _ := ArnRegionAccountResource(arn)
	return &region, &account
}

func ArnRegionAccountResource(arnstr string) (string, string, string) {
	a, err := arn.Parse(arnstr)
	if err != nil {
		return "", "", ""
	}
	return a.Region, a.AccountID, a.Resource
}

func createArn(arn_str string, region *string, account_id *string, name_or_arn *string) *string {
	if len(*name_or_arn) < 5 || (*name_or_arn)[:4] == "arn:" {
		return name_or_arn
	}

	if region == nil || account_id == nil || name_or_arn == nil {
		return name_or_arn
	}

	arn := fmt.Sprintf(arn_str, *region, *account_id, *name_or_arn)
	return &arn
}

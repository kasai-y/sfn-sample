package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/fatih/color"
	"github.com/kasai-y/sfn-sample/config"
	"github.com/urfave/cli"
	"os"
	"time"
)

var configFile string

func main() {
	app := cli.NewApp()
	app.Name = "sfn-sample"
	app.Usage = "aws StepDFunctions sample."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config,c",
			Usage:       "config file.",
			Destination: &configFile,
		},
	}
	app.Action = action
	_ = app.Run(os.Args)
}

func action(_ *cli.Context) {

	cfg, err := config.Load(configFile)
	if err != nil {
		color.Red("%+v", err)
	}

	sess := session.Must(session.NewSession(aws.NewConfig().
		WithRegion(cfg.AWS.Region).
		WithCredentials(credentials.NewStaticCredentials(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		))))

	c := sfn.New(sess)

	// ステートマシンを実行
	input := `{"Comment": "Insert your JSON here"}`
	name := fmt.Sprintf("sfn-sample-%s", time.Now().Format("20060102-150405"))
	_, err = c.StartExecution(&sfn.StartExecutionInput{
		Input:           aws.String(input),
		Name:            aws.String(name),
		StateMachineArn: aws.String(cfg.StepFunctions.Arn),
	})
	if err != nil {
		color.Red("%+v", err)
	}

	// ステートマシンの一覧を取得
	listExecutions, err := c.ListExecutions(&sfn.ListExecutionsInput{
		MaxResults:      nil,
		NextToken:       nil,
		StateMachineArn: aws.String(cfg.StepFunctions.Arn),
		StatusFilter:    nil,
	})
	if err != nil {
		color.Red("%+v", err)
	}
	for _, exc := range listExecutions.Executions {
		color.Cyan("name:%s, stats:%s", *exc.Name, *exc.Status)
	}
}

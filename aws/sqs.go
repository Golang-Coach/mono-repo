package aws

import (
	"github.com/golang-coach/mono-repo/constants"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Queue struct {
	sqs *sqs.SQS
}

func NewQueue(sqs *sqs.SQS) Queue {
	return Queue{sqs: sqs}
}

func (q Queue) Send(message string) error {
	_, err := q.sqs.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(1),
		MessageBody:  aws.String(message),
		QueueUrl:     aws.String(os.Getenv(constants.QueueName)),
	})

	return err
}

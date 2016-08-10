package mock_test

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/revboss/go-mock"
	"testing"
)

func TestSQS(t *testing.T) {
	queue := &mock.SQS{}
	_, e := queue.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String("testing-queue"),
		MessageBody: aws.String(`{}`),
	})

	if e != nil {
		t.Error(e)
		t.FailNow()
	}

	if len(queue.Messages) != 1 {
		t.Errorf("Expected 1 SQS message")
		t.FailNow()
	}

	messages, e := queue.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: aws.String("testing-queue"),
	})
	if e != nil {
		t.Error(e)
		t.FailNow()
	}

	if len(messages.Messages) != 1 {
		t.Errorf("Expected 1 message")
		t.FailNow()
	}

	_, e = queue.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String("testing-queue"),
		ReceiptHandle: messages.Messages[0].ReceiptHandle,
	})
	if e != nil {
		t.Error(e)
		t.FailNow()
	}

	messages, e = queue.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: aws.String("testing-queue"),
	})
	if e != nil {
		t.Error(e)
		t.FailNow()
	}

	if len(messages.Messages) != 0 {
		t.Errorf("Expected 0 messages")
	}
}

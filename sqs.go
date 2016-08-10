package mock

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"strconv"
	"testing"
)

type SQS struct {
	Messages []*sqs.Message
}

func (s *SQS) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	id, e := strconv.Atoi(*input.ReceiptHandle)
	if e != nil {
		return nil, e
	}

	if id > len(s.Messages) || id < 0 {
		return nil, fmt.Errorf("Unknown ReceiptHandle: %s", *input.ReceiptHandle)
	}

	messages := []*sqs.Message{}
	if id > 0 {
		messages = append(messages, s.Messages[0:id-1]...)
		messages = append(messages, s.Messages[id+1:]...)
	}

	s.Messages = messages

	return &sqs.DeleteMessageOutput{}, nil
}

func (s *SQS) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	output := &sqs.ReceiveMessageOutput{}

	if len(s.Messages) == 0 {
		return output, nil
	}

	output.Messages = []*sqs.Message{s.Messages[0]}

	return output, nil
}

func (s *SQS) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	s.Messages = append(s.Messages, &sqs.Message{
		Body:          input.MessageBody,
		ReceiptHandle: aws.String(strconv.Itoa(len(s.Messages))),
	})
	return &sqs.SendMessageOutput{}, nil
}

func TestSQS(t *testing.T) {
	queue := &SQS{}
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

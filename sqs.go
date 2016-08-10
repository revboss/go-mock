package mock

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"strconv"
	"sync"
)

type SQS struct {
	Messages []*sqs.Message

	mut sync.Mutex
}

func (s *SQS) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	s.mut.Lock()
	defer s.mut.Unlock()
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
	s.mut.Lock()
	defer s.mut.Unlock()
	output := &sqs.ReceiveMessageOutput{}

	if len(s.Messages) == 0 {
		return output, nil
	}

	output.Messages = []*sqs.Message{s.Messages[0]}

	return output, nil
}

func (s *SQS) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	s.mut.Lock()
	defer s.mut.Unlock()
	s.Messages = append(s.Messages, &sqs.Message{
		Body:          input.MessageBody,
		ReceiptHandle: aws.String(strconv.Itoa(len(s.Messages))),
	})
	return &sqs.SendMessageOutput{}, nil
}

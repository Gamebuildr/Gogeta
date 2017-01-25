package poller

import "github.com/aws/aws-sdk-go/service/sqs"

type ProcessFunc func(msg *sqs.Message) error

func (f ProcessFunc) ProcessMessage(msg *sqs.Message) error {
	return f(msg)
}

type Process interface {
	ProcessMessage(msg *sqs.Message) error
}

func Start(process Process) {
	// region, _ := os.Getenv()
	// sqsURL, _ := os.Getenv()

	// amazonQueue := AmazonQueue{
	// 	Client: sqs.New(session.New(), &aws.Config{Region: aws.String(region)}),
	// 	Region: region,
	// 	URL:    sqsURL,
	// }
}

// func InboundMessages(session *sqs.SQS, messages []*sqs.Message, process Process) {
// 	for i := range messages {
// 		go func(message *sqs.Message) {
// 			if err := ProcessInbound(session, message, process); err != nil {
// 				logger.Warning(err.Error())
// 			}
// 		}(messages[i])
// 	}
// }

// func ProcessInbound(session *sqs.SQS, m *sqs.Message, process Process) error {
// 	var err error
// 	err = process.ProcessMessage(m)
// 	if err != nil {
// 		RemoveMessageFromPoller(session, m)
// 		return err
// 	}
// 	RemoveMessageFromPoller(session, m)
// 	return nil
// }

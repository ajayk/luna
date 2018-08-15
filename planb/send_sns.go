package planb

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func PublishMessage(region, subject, message, planb_topic_arn string) {
	if planb_topic_arn == "" {
		log.Println("AWS SNS Publish SKIP!!")
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	svc := sns.New(sess)
	params := &sns.PublishInput{
		Message: aws.String(message), // Required
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"Key": &sns.MessageAttributeValue{ // Required
				DataType:    aws.String("String"), // Required
				StringValue: aws.String("String"),
			},
		},
		MessageStructure: aws.String("messageStructure"),
		Subject:          aws.String(subject),
		TopicArn:         aws.String(planb_topic_arn),
	}

	resp, err := svc.Publish(params)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	log.Println(resp)
}

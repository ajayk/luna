package awsobject

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func ReleaseElasticip(region string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create an EC2 service client.
	svc := ec2.New(sess)
	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("domain"),
				Values: aws.StringSlice([]string{"vpc"}),
			},
		},
	})
	if err != nil {
		log.Fatal("Unable to elastic IP address ", err)
	}

	if len(result.Addresses) == 0 {
		log.Printf("No elastic IPs for %s region\n", *svc.Config.Region)
	} else {
		for _, addr := range result.Addresses {
			if addr.AssociationId == nil {
				// log.Println("********** should release", addr.String())
				// log.Println("InstancdId=", addr.InstanceId)
				releaseElasticip(svc, *addr.AllocationId)
			}
		}
	}
}

func releaseElasticip(svc *ec2.EC2, allocationID string) {
	_, err := svc.ReleaseAddress(&ec2.ReleaseAddressInput{
		AllocationId: aws.String(allocationID),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "InvalidAllocationID.NotFound" {
			log.Fatalln("Allocation ID %s does not exist: ", allocationID)
		}
		log.Fatalln("Unable to release IP address for allocation: ", allocationID, err)
	}

	log.Println("Successfully released allocation ID ", allocationID)
}

func AllocateIP(region string, instanceID string) string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create an EC2 service client.
	svc := ec2.New(sess)

	//Call AllocateAddress, passing in "vpc" as the Domain value.
	allocRes, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
		Domain: aws.String("vpc"),
	})
	if err != nil {
		log.Fatal("Unable to allocate IP address", allocRes, err)
	}

	assocRes, err := svc.AssociateAddress(&ec2.AssociateAddressInput{
		AllocationId: allocRes.AllocationId,
		InstanceId:   aws.String(instanceID),
	})
	if err != nil {
		log.Fatal("Unable to associate IP address with ", instanceID, err)
	}

	log.Printf("Successfully allocated %s with instance %s.\n\tallocation id: %s, association id: %s\n",
		*allocRes.PublicIp, instanceID, *allocRes.AllocationId, *assocRes.AssociationId)

	return *allocRes.PublicIp

}

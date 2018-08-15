package awsobject

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func DescribeSecurityGroups(svc *ec2.EC2) *ec2.DescribeSecurityGroupsOutput {
	output, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		log.Println("error at describing security group", err)
	}
	return output
}

func UpdateSecurityGroup(aws_region, revokeip, addip string) error {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(aws_region)})

	var groupname string = ""

	// get all attached securityGroups
	output := DescribeSecurityGroups(svc)

	// find security group which has (should be revoked)ip
	for _, sg := range output.SecurityGroups {
		for _, ipper := range sg.IpPermissions {
			for _, iprs := range ipper.IpRanges {
				if *iprs.CidrIp == revokeip+"/32" {
					groupname = *sg.GroupName
					break
				}
			}
			if groupname != "" {
				break
			}
		}
		if groupname != "" {
			break
		}
	}

	if groupname == "" {
		//return errors.New("Security group that contain revoke ip not found error")
		log.Println("Security group that contain revoke ip not found error")
		log.Println("so add ip to groupname =" + *output.SecurityGroups[0].GroupName)
		groupname = *output.SecurityGroups[0].GroupName
	} else {
		log.Println("found security group: ", groupname)
	}

	//revoke ip from security group's ingress
	_, err := svc.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
		GroupName:  aws.String(groupname),
		CidrIp:     aws.String(revokeip + "/32"),
		FromPort:   aws.Int64(26656),
		ToPort:     aws.Int64(26656),
		IpProtocol: aws.String("tcp"),
	})
	if err != nil {
		log.Println("revoke error:", err)
	} else {
		log.Println("revoke done")
	}

	err = nil

	// add ip to securitygroup's ingress
	_, err = svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupName: aws.String(groupname),
		IpPermissions: []*ec2.IpPermission{
			{
				FromPort:   aws.Int64(26656),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String(addip + "/32"),
						Description: aws.String("Sentry"),
					},
				},
				ToPort: aws.Int64(26656),
			},
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println("add to securtiygroup error:", err.Error(), "ip:", addip)
		}
	} else {
		log.Println("add to ingress done")
	}

	return err
}

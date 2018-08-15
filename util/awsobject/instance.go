package awsobject

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/lunamint/luna/util/shell"
	//"github.com/aws/aws-sdk-go/service/sns"
)

func GetInstanceID(region string) string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	ec2m := ec2metadata.New(sess)
	doc, err := ec2m.GetInstanceIdentityDocument()
	if err != nil {
		log.Fatal("doc error :", doc, err)
	}

	return doc.InstanceID
}

func GetRegion() string {
	cmd := "curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone"
	cmdres, err := shell.SimpleShellCall("", cmd, false)
	if err != nil {
		log.Println("EC2_REGION error:", err)
	}

	//log.Println("cmd2res=", cmdres)
	trimed := strings.TrimSpace(cmdres)
	return trimed[:len(trimed)-1]

}

func GetAccountID(region string) string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	ec2m := ec2metadata.New(sess)
	doc, err := ec2m.GetInstanceIdentityDocument()
	if err != nil {
		log.Fatal("doc error :", doc, err)
	}

	return doc.AccountID
}

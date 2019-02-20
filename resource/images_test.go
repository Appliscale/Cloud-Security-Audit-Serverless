package resource

import (
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/configuration"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession/clientfactory/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
)

func TestLoadImagesFromAWS(t *testing.T) {
	config := configuration.GetTestConfig(t)
	defer config.ClientFactory.(*mocks.ClientFactoryMock).Destroy()

	ec2Client, _ := config.ClientFactory.GetEc2Client(csasession.SessionConfig{})
	ec2Client.(*mocks.MockEC2Client).
		EXPECT().
		DescribeImages(&ec2.DescribeImagesInput{
			Owners: []*string{aws.String("self")},
		}).
		Times(1).
		Return(&ec2.DescribeImagesOutput{}, nil)

	LoadResource(&Images{}, &config, "region")
}

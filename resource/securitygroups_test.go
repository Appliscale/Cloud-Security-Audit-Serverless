package resource

import (
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/configuration"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession/clientfactory/mocks"
	"github.com/aws/aws-sdk-go/service/ec2"
	"testing"
)

func TestLoadSecurityGroupsFromAWS(t *testing.T) {
	config := configuration.GetTestConfig(t)
	defer config.ClientFactory.(*mocks.ClientFactoryMock).Destroy()

	ec2Client, _ := config.ClientFactory.GetEc2Client(csasession.SessionConfig{})
	ec2Client.(*mocks.MockEC2Client).
		EXPECT().
		DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{}).
		Times(1).
		Return(&ec2.DescribeSecurityGroupsOutput{}, nil)

	LoadResource(&SecurityGroups{}, &config, "region")
}

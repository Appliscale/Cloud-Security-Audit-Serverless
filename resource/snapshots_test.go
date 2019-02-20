package resource

import (
	"testing"

	"github.com/Appliscale/Cloud-Security-Audit-Serverless/configuration"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession/clientfactory/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TestLoadSnapshotsFromAWS(t *testing.T) {
	config := configuration.GetTestConfig(t)
	defer config.ClientFactory.(*mocks.ClientFactoryMock).Destroy()

	ec2Client, _ := config.ClientFactory.GetEc2Client(csasession.SessionConfig{})
	ec2Client.(*mocks.MockEC2Client).
		EXPECT().
		DescribeSnapshots(&ec2.DescribeSnapshotsInput{
			OwnerIds: []*string{aws.String("self")},
		}).
		Times(1).
		Return(&ec2.DescribeSnapshotsOutput{}, nil)

	LoadResource(&Snapshots{}, &config, "region")
}

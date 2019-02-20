package resource

import (
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/configuration"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type SecurityGroups map[string][]*ec2.IpPermission

func (s *SecurityGroups) LoadFromAWS(config *configuration.Config, region string) error {

	ec2API, err := config.ClientFactory.GetEc2Client(csasession.SessionConfig{Profile: config.Profile, Region: region})
	if err != nil {
		return err
	}

	q := &ec2.DescribeSecurityGroupsInput{}

	for {
		result, err := ec2API.DescribeSecurityGroups(q)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case "OptInRequired":
					break
				default:
					return err
				}
			} else {
				return err
			}
		}
		for _, sg := range result.SecurityGroups {
			(*s)[*sg.GroupId] = append((*s)[*sg.GroupId], sg.IpPermissions...)
		}
		if result.NextToken == nil {
			break
		}
		q.NextToken = result.NextToken
	}
	return nil
}

func (s *SecurityGroups) GetIpPermissionsByID(groupID string) []*ec2.IpPermission {

	if sg, ok := (*s)[groupID]; ok {
		return sg
	}
	return nil
}

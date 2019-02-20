package clientfactory

import (
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ClientFactory interface {
	GetKmsClient(config csasession.SessionConfig) (KmsClient, error)
	GetEc2Client(config csasession.SessionConfig) (EC2Client, error)
	GetS3Client(config csasession.SessionConfig) (S3Client, error)
}

// GetKmsClient creates a new KMS client from cached session.
func (factory *ClientFactoryAWS) GetKmsClient(config csasession.SessionConfig) (KmsClient, error) {
	sess, err := factory.sessionFactory.GetSession(config)
	if err != nil {
		return nil, err
	}

	client := kms.New(sess)
	return AWSKmsClient{api: client}, nil
}

// GetEc2Client creates a new EC2 client from cached session.
func (factory *ClientFactoryAWS) GetEc2Client(config csasession.SessionConfig) (EC2Client, error) {
	sess, err := factory.sessionFactory.GetSession(config)
	if err != nil {
		return nil, err
	}

	client := ec2.New(sess)
	return AWSEC2Client{api: client}, nil
}

// GetS3Client creates a new S3 client from cached session.
func (factory *ClientFactoryAWS) GetS3Client(config csasession.SessionConfig) (S3Client, error) {
	sess, err := factory.sessionFactory.GetSession(config)
	if err != nil {
		return nil, err
	}

	client := s3.New(sess)
	return AWSS3Client{api: client}, nil
}

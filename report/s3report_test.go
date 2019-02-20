package report

import (
	"fmt"
	"testing"

	"github.com/Appliscale/Cloud-Security-Audit-Serverless/resource"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

type Permissions struct {
	in  string
	out string
}

func TestS3Report_WhenSSEAlgorithmIsAES256CheckEncryptionTypeReturnsAES256(t *testing.T) {

	s3BucketReport := &S3BucketReport{}
	AES256Rule := s3.ServerSideEncryptionByDefault{
		SSEAlgorithm: aws.String("AES256"),
	}
	s3BucketReport.CheckEncryptionType(AES256Rule, resource.NewKMSKeys())
	assert.Equal(t, AES256, s3BucketReport.EncryptionType)
}

func TestS3Report_WhenSSEAlgorithmIsCustomAWSKMSCheckEncryptionTypeReturnsCKMS(t *testing.T) {
	kmsKeyArn := "arn:aws:kms:us-east-1:126286021559:key/2fdaec7f-6f04-4b2c-b6ea-a1a6d8437c3e"
	kmsKeys := resource.NewKMSKeys()
	kmsKeys.Values[kmsKeyArn] = &resource.KMSKey{
		Custom: true,
	}

	s3BucketReport := &S3BucketReport{}
	customKMSKeyRule := s3.ServerSideEncryptionByDefault{
		KMSMasterKeyID: &kmsKeyArn,
		SSEAlgorithm:   aws.String("aws:kms"),
	}

	s3BucketReport.CheckEncryptionType(customKMSKeyRule, kmsKeys)
	assert.Equalf(t, CKMS, s3BucketReport.EncryptionType, fmt.Sprintf("Expected %s, got %s", CKMS.String(), s3BucketReport.EncryptionType))
}

func TestS3Report_WhenSSEAlgorithmIsDefaultAWSKMSCheckEncryptionTypeReturnsDKMS(t *testing.T) {
	kmsKeyArn := "arn:aws:kms:us-east-1:126286021559:key/2fdaec7f-6f04-4b2c-b6ea-a1a6d8437c3e"
	kmsKeys := resource.NewKMSKeys()
	kmsKeys.Values[kmsKeyArn] = &resource.KMSKey{
		Custom: false,
	}
	s3BucketReport := &S3BucketReport{}
	customKMSKeyRule := s3.ServerSideEncryptionByDefault{
		KMSMasterKeyID: &kmsKeyArn,
		SSEAlgorithm:   aws.String("aws:kms"),
	}
	s3BucketReport.CheckEncryptionType(customKMSKeyRule, kmsKeys)
	assert.Equal(t, DKMS, s3BucketReport.EncryptionType)
}

func TestGetTypeOfAccessACL(t *testing.T) {
	var permissions = []Permissions{
		{"WRITE_ACP", "W"},
		{"READ", "R"},
		{"DELETE", "D"},
		{"FULL_CONTROL", "RWD"},
	}

	for _, permission := range permissions {
		got := getTypeOfAccessACL(permission.in)
		assert.Equalf(t, got, permission.out, fmt.Sprintf("Expected %s, got %s", permission.out, got))
	}

}

func TestGetTypeOfAccessPolicy(t *testing.T) {
	var actions = resource.Actions{"s3:DeleteObject", "s3:GetObjectVersion", "s3:PutObjectAcl"}
	got := getTypeOfAccessPolicy(actions)
	expected := "[DRW]"
	assert.Equalf(t, got, expected, fmt.Sprintf("Expected %s, got %s", expected, got))

}

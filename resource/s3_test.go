package resource

import (
	"fmt"
	"testing"

	"github.com/Appliscale/Cloud-Security-Audit-Serverless/configuration"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession/clientfactory/mocks"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestS3_ActionUnmarshalJSONCreateSliceOfStringsFromJsonArray(t *testing.T) {
	b := []byte(`["s3:GetObject","s3:GetBucketLocation"]`)
	actions := Actions{}
	err := actions.UnmarshalJSON(b)
	assert.Nilf(t, err, "UnmarshalJSON should not return error for array of actions.")
	assert.Equalf(t, 2, len(actions), "Actions should contain two elements.")
}

func TestS3_ActionUnmarshalJSONCreateSliceOfStringsFromJsonString(t *testing.T) {
	b := []byte(`"s3:GetObject"`)
	actions := Actions{}
	err := actions.UnmarshalJSON(b)
	assert.Nilf(t, err, "UnmarshalJSON should not return error for string in actions object.")
	assert.Equalf(t, 1, len(actions), "Actions should contain two elements.")
}

func TestS3_ActionUnmarshalJSONReturnsErrorFromJsonMap(t *testing.T) {
	b := []byte(`{"something":{"s3":"GetObject"}}`)
	actions := Actions{}
	err := actions.UnmarshalJSON(b)
	assert.NotNilf(t, err, "UnmarshalJSON should return error for Json Map")
}

func TestS3_PrincipalUnmarshalJSONCreatesMapOfSlicesIfJSONPropertiesAreMapOfArrays(t *testing.T) {
	b := []byte(`{"AWS": ["something","something2"]}`)
	principal := Principal{}
	err := principal.UnmarshalJSON(b)
	assert.Nilf(t, err, "This should not return error")
	assert.Equal(t, 2, len(principal.Map["AWS"]))
}

func TestS3_PrincipalUnmarshalJSONCreateMapOfSliceIfJSonPropertyIsMap(t *testing.T) {
	b := []byte(`{"Service":"blabla"}`)
	principal := Principal{}
	err := principal.UnmarshalJSON(b)
	assert.Nilf(t, err, "This should not return error")
	assert.Equal(t, 1, len(principal.Map["Service"]))
}

func TestS3_PrincipalUnmarshalJSONAssignWildcardIfJsonPropertyIsString(t *testing.T) {
	b := []byte(`"*"`)
	principal := Principal{}
	err := principal.UnmarshalJSON(b)
	assert.Nilf(t, err, "This should not return error")
	fmt.Printf("\n%v\n", principal)
	assert.Equal(t, "*", principal.Wildcard)
}

func TestS3Buckets_LoadNames(t *testing.T) {
	config := configuration.GetTestConfig(t)
	defer config.ClientFactory.(*mocks.ClientFactoryMock).Destroy()

	ec2Client, _ := config.ClientFactory.GetS3Client(csasession.SessionConfig{})
	ec2Client.(*mocks.MockS3Client).
		EXPECT().
		ListBuckets(&s3.ListBucketsInput{}).
		Times(1).
		Return(&s3.ListBucketsOutput{}, nil)

	s3Bucket := &S3Buckets{}
	s3Bucket.LoadNames(&config, "region")

}

func TestResourcesWithSlice(t *testing.T) {
	b := []byte(`["arn:aws:s3:::examplebucket/*"]`)
	resources := Resources{}
	err := resources.UnmarshalJSON(b)
	assert.Nilf(t, err, "UnmarshalJSON should not return error for array of actions.")
}

func TestResourcesWithTwoElements(t *testing.T) {
	b := []byte(`["arn:aws:iam::111122223333:root","arn:aws:iam::444455556666:root"]`)
	resources := Resources{}
	err := resources.UnmarshalJSON(b)
	assert.Nilf(t, err, "UnmarshalJSON should not return error for array of resources.")
	assert.Equalf(t, 2, len(resources), "Resources should contain two elements.")
}

func TestResourcesWithOneElement(t *testing.T) {
	b := []byte(`"s3:GetObject"`)
	resources := Resources{}
	err := resources.UnmarshalJSON(b)
	assert.Nilf(t, err, "UnmarshalJSON should not return error for string in actions object.")
	assert.Equalf(t, 1, len(resources), "Resources should contain one element.")
}

func TestResourcesWithMap(t *testing.T) {
	b := []byte(`{"something":{"s3":"GetObject"}}`)
	resources := Resources{}
	err := resources.UnmarshalJSON(b)
	assert.NotNilf(t, err, "UnmarshalJSON should return error for Json Map.")
}

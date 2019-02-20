package configuration

import (
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession/clientfactory/mocks"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession/sessionfactory"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/logger"
	"testing"
)

func GetTestConfig(t *testing.T) (config Config) {
	myLogger := logger.CreateQuietLogger()
	config.Logger = &myLogger
	clientFactory := mocks.NewClientFactoryMock(t)
	config.ClientFactory = &clientFactory
	config.SessionFactory = sessionfactory.New()

	return config
}

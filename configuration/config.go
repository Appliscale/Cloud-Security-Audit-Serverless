package configuration

import (
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession/clientfactory"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/csasession/sessionfactory"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/logger"
)

type Config struct {
	Regions        *[]string
	Services       *[]string
	Profile        string
	SessionFactory *sessionfactory.SessionFactory
	ClientFactory  clientfactory.ClientFactory
	Logger         *logger.Logger
	Mfa            bool
	MfaDuration    int64
}

func GetConfig() (config Config) {
	myLogger := logger.CreateDefaultLogger()
	config.Logger = &myLogger
	config.SessionFactory = sessionfactory.New()
	config.ClientFactory = clientfactory.New(config.SessionFactory)

	return config
}

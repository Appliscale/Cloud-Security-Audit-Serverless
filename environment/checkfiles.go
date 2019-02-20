package environment

import (
	"bufio"
	"github.com/Appliscale/Cloud-Security-Audit-Serverless/configuration"
	"github.com/Appliscale/perun/helpers"
	"os"
	"strconv"
	"strings"
)

var regions = []string{
	"us-east-1",
	"us-east-2",
	"us-west-1",
	"us-west-2",
	"ca-central-1",
	"ca-central-1",
	"eu-west-1",
	"eu-west-2",
	"ap-northeast-1",
	"ap-northeast-2",
	"ap-southeast-1",
	"ap-southeast-2",
	"ap-south-1",
	"sa-east-1",
}

func CheckAWSConfigFiles(config *configuration.Config) bool {
	homeDir, pathError := GetUserHomeDir()
	if pathError != nil {
		config.Logger.Error(pathError.Error())
		return false
	}

	configAWSExists, configError := isAWSConfigPresent(homeDir)
	if configError != nil {
		config.Logger.Error(configError.Error())
	}

	credentialsExists, credentialsError := isCredentialsPresent(homeDir)
	if credentialsError != nil {
		config.Logger.Error(credentialsError.Error())
	}

	profile := config.Profile

	if configAWSExists {
		profilesInConfig := getProfilesFromFile(config, homeDir+"/.aws/config")
		if !helpers.SliceContains(profilesInConfig, profile) {
			var ans string
			config.Logger.GetInput("You don't have the "+profile+" profile in your config file. Would you like to create one? *Y* / *N*", &ans)
			if strings.ToUpper(ans) == "Y" {
				region := getUserRegion(config)
				CreateAWSConfigFile(config, profile, region, "")
			} else {
				config.Logger.Info("You can use another profile by setting the \"-p\" argument or specify a different default profile by setting the AWS_PROFILE variable")
				return false
			}
		}
		if credentialsExists {
			addProfileToCredentials(profile, homeDir, config)
		} else {
			CreateAWSCredentialsFile(config, profile)
		}
	} else {
		if credentialsExists {
			var ans string
			config.Logger.GetInput("File .aws/config does not exist, but .aws/credentials has been found. Do you want to create config file using one of the profiles in the .aws/credentias? *Y* / *N*", &ans)
			if strings.ToUpper(ans) == "Y" {
				createConfigProfileFromCredentials(homeDir, config, profile)
				return true

			} else {
				profile = setProfileInfoAndCreateConfigFile(config)
				CreateAWSCredentialsFile(config, profile)
			}
		} else {
			config.Logger.Info("File .aws/config does not exist.")
			profile = setProfileInfoAndCreateConfigFile(config)
			CreateAWSCredentialsFile(config, profile)
		}
	}
	return true
}

func isAWSConfigPresent(homePath string) (bool, error) {
	_, credentialsError := os.Open(homePath + "/.aws/config")
	if credentialsError != nil {
		return false, nil
	}
	return true, nil
}

func isCredentialsPresent(homePath string) (bool, error) {
	_, credentialsError := os.Open(homePath + "/.aws/credentials")
	if credentialsError != nil {
		return false, nil
	}
	return true, nil
}

func getUserRegion(config *configuration.Config) string {
	showAvailableRegions(config)
	var numberRegion int
	config.Logger.GetInput("Region", &numberRegion)

	for numberRegion < 0 || numberRegion >= 14 {
		config.Logger.Always("Try again, invalid region")
		config.Logger.GetInput("Region", &numberRegion)
	}
	region := regions[numberRegion]
	config.Logger.Always("Your region is: " + region)
	return region
}

func showAvailableRegions(config *configuration.Config) {
	config.Logger.Always("Available Regions:")
	for i := 0; i < len(regions); i++ {
		pom := strconv.Itoa(i)
		config.Logger.Always("Number " + pom + " region " + regions[i])
	}
}

func getUserOutput(config *configuration.Config) string {
	var output string
	config.Logger.GetInput("Input the output format [json, text, table]", &output)
	for !helpers.SliceContains([]string{"json", "text", "table"}, output) {
		config.Logger.Always("Try again, invalid output")
		config.Logger.GetInput("Input the output format [json, text, table]", &output)
	}
	config.Logger.Always("Your output is: " + output)
	return output
}

func getUserProfile(config *configuration.Config) string {
	var profile string
	config.Logger.GetInput("Input name of profile", &profile)
	for profile == "" {
		config.Logger.Always("Try again, invalid profile")
		config.Logger.GetInput("Input name of profile", &profile)
	}
	config.Logger.Always("Your region is: " + profile)
	return profile
}

func getProfilesFromFile(config *configuration.Config, path string) []string {
	credentials, credentialsError := os.Open(path)
	if credentialsError != nil {
		config.Logger.Error(credentialsError.Error())
		return []string{}
	}
	defer credentials.Close()
	profiles := make([]string, 0)
	scanner := bufio.NewScanner(credentials)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "[") {
			profile := strings.TrimPrefix(scanner.Text(), "[")
			profile = strings.TrimSuffix(profile, "]")
			if strings.Contains(profile, "profile ") {
				profile = strings.TrimPrefix(profile, "profile ")
			}
			if strings.Contains(profile, "-long-term") {
				profile = strings.TrimSuffix(profile, "-long-term")
			}
			profiles = append(profiles, profile)
		}
	}
	return profiles
}

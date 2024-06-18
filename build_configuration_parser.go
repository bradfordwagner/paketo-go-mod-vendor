package gomodvendor

import (
	"os"
)

type BuildConfiguration struct {
	GoProxy   string
	GoNoProxy string
	GoSumDB   string
	GoNoSumDB string
	GoPrivate string
}

type BuildConfigurationParser struct {
}

func NewBuildConfigurationParser() BuildConfigurationParser {
	return BuildConfigurationParser{}
}

func (p BuildConfigurationParser) Parse() (buildConfiguration BuildConfiguration, err error) {
	// read configs for: GOPROXY, GONOPROXY, GOSUMDB, GONOSUMDB, and GOPRIVATE
	if val, ok := os.LookupEnv("GOPROXY"); ok {
		buildConfiguration.GoProxy = val
	}
	if val, ok := os.LookupEnv("GONOPROXY"); ok {
		buildConfiguration.GoNoProxy = val
	}
	if val, ok := os.LookupEnv("GOSUMDB"); ok {
		buildConfiguration.GoSumDB = val
	}
	if val, ok := os.LookupEnv("GONOSUMDB"); ok {
		buildConfiguration.GoNoSumDB = val
	}
	if val, ok := os.LookupEnv("GOPRIVATE"); ok {
		buildConfiguration.GoPrivate = val
	}
	return buildConfiguration, nil
}

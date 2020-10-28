package config

import "github.com/Solar-2020/GoUtils/common"

var (
	Config config
)

type config struct {
	common.SharedConfig
	InterviewDataBaseConnectionString  string `envconfig:"INTERVIEW_DB_CONNECTION_STRING" default:"-"`
}

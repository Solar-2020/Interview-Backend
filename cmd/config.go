package main

type config struct {
	Port string `envconfig:"PORT" default:"8099"`
	InterviewDataBaseConnectionString  string `envconfig:"INTERVIEW_DB_CONNECTION_STRING" default:"-"`
}

package main

import (
	aw "github.com/deanishe/awgo"
)

const (
	projectNameConfigKey = "projectName"
	tokenConfigKey       = "token"
)

// Service workflow service
type Service struct {
	wf *aw.Workflow
}

// NewService return a create service
func NewService() *Service {
	return &Service{
		wf: aw.New(),
	}
}

func (s *Service) getToken() string {
	return s.wf.Config.Get(tokenConfigKey)
}

func (s *Service) getProjectName() string {
	return s.wf.Config.Get(projectNameConfigKey)
}

package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/google/subcommands"
)

type setupCommand struct {
	*Service

	token       string
	projectName string
}

var _ subcommands.Command = (*setupCommand)(nil)

func (s *setupCommand) Name() string     { return "setup" }
func (s *setupCommand) Synopsis() string { return "setup scrapbox" }
func (s *setupCommand) Usage() string {
	return `setup:
`
}

func (s *setupCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.token, "t", "", "scrapbox connect.sid")
	f.StringVar(&s.projectName, "p", "", "scrapbox project name")
}

func (s *setupCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	args := f.Args()
	if len(args) == 2 {
		s.projectName = args[0]
		s.token = args[1]
	}

	s.wf.Config.Set(projectNameConfigKey, s.projectName, false)
	s.wf.Config.Set(tokenConfigKey, s.token, false)

	if err := s.wf.Config.Do(); err != nil {
		s.wf.FatalError(err)
		return subcommands.ExitFailure
	}

	s.wf.NewItem(fmt.Sprintf("projectName = %s token = %s",
		s.wf.Config.GetString(projectNameConfigKey),
		s.wf.Config.GetString(tokenConfigKey),
	))
	s.wf.SendFeedback()

	return subcommands.ExitSuccess
}

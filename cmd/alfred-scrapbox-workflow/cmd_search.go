package main

import (
	"context"
	"flag"
	"path"
	"time"

	"github.com/google/subcommands"
	"github.com/kyokomi/go-scrapbox"
)

const (
	cacheKeyPages   = "pages"
	scrapboxBaseURL = "https://scrapbox.io/"
)

// searchCommand 検索コマンド
type searchCommand struct {
	*Service

	query string
}

var _ subcommands.Command = (*searchCommand)(nil)

func (s *searchCommand) Name() string     { return "search" }
func (s *searchCommand) Synopsis() string { return "search scrapbox pages" }
func (s *searchCommand) Usage() string {
	return `search:
`
}

func (s *searchCommand) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.query, "query", "", "search query")
}

func (s *searchCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if len(s.query) <= 1 {
		s.wf.Fatal("No pages were found.")
		return subcommands.ExitUsageError
	}

	res, err := s.pageListWithCache(ctx)
	if err != nil {
		s.wf.FatalError(err)
		return subcommands.ExitFailure
	}

	for _, page := range res.Pages {
		//if strings.Index(strings.ToLower(page.Title), strings.ToLower(query)) == -1 {
		//	continue
		//}

		s.wf.NewItem(page.Title).
			Arg(path.Join(scrapboxBaseURL, s.getProjectName(), page.Title)).
			Valid(true)

		// TODO: icon取得
		//found, iconURL, err := client.Page.IconURL(ctx, s.projectName, page.Title)
		//if err != nil {
		//	log.Fatalf("%+v", err)
		//}
	}

	s.wf.Filter(s.query)
	s.wf.WarnEmpty("No pages were found.", "Try different query.")
	s.wf.SendFeedback()

	return subcommands.ExitSuccess
}

func (s *searchCommand) pageListWithCache(ctx context.Context) (scrapbox.PageListResponse, error) {
	// Cacheは1分
	if !s.wf.Cache.Expired(cacheKeyPages, 1*time.Minute) {
		var res scrapbox.PageListResponse
		if err := s.wf.Cache.LoadJSON(cacheKeyPages, &res); err != nil {
			return res, err
		}
		return res, nil
	}

	client := scrapbox.NewClient(s.getToken())

	res, err := client.Page.ListAll(ctx, s.getProjectName())
	if err != nil {
		s.wf.FatalError(err)
	}

	// cacheのファイル名
	if err := s.wf.Cache.StoreJSON(cacheKeyPages, &res); err != nil {
		s.wf.FatalError(err)
	}

	return res, nil
}

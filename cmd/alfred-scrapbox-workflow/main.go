package main

import (
	"context"
	"flag"
	"log"
	"path"
	"time"

	aw "github.com/deanishe/awgo"
	"github.com/kyokomi/go-scrapbox"
)

const scrapboxBaseURL = "https://scrapbox.io/"

func main() {
	log.SetFlags(log.Llongfile)
	token := flag.String("t", "<token>", "scrapbox connect.sid")
	projectName := flag.String("p", "kyokomi", "scrapbox project name")
	query := flag.String("q", "query", "query args")
	flag.Parse()

	s := NewService()

	// TODO: token設定のcommandにする
	s.wf.Config.Set("projectName", *projectName, false)
	s.wf.Config.Set("token", *token, false)
	if err := s.wf.Config.Do(); err != nil {
		s.wf.FatalError(err)
		return
	}

	// TODO: 検索のcommandにする
	s.SearchPage(*query)
}

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
	return s.wf.Config.Get("token")
}

func (s *Service) getProjectName() string {
	return s.wf.Config.Get("projectName")
}

// SearchPage ページ検索
func (s *Service) SearchPage(query string) {
	if len(query) <= 1 {
		s.wf.Fatal("No pages were found.")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := s.pageListWithCache(ctx)
	if err != nil {
		s.wf.FatalError(err)
		return
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

	s.wf.Filter(query)
	s.wf.WarnEmpty("No pages were found.", "Try different query.")
	s.wf.SendFeedback()
}

const cacheKeyPages = "pages"

func (s *Service) pageListWithCache(ctx context.Context) (scrapbox.PageListResponse, error) {
	// Cacheは1分
	if !s.wf.Cache.Expired(cacheKeyPages, 1*time.Minute) {
		var res scrapbox.PageListResponse
		if err := s.wf.Cache.LoadJSON(cacheKeyPages, &res); err != nil {
			return res, err
		}
		return res, nil
	}

	client := scrapbox.NewClient(s.getToken())

	offset := uint(0)
	limit := uint(1000) // TODO: 次ページなくなるまで取得?

	res, err := client.Page.List(ctx, s.getProjectName(), offset, limit)
	if err != nil {
		s.wf.FatalError(err)
	}

	// cacheのファイル名
	if err := s.wf.Cache.StoreJSON(cacheKeyPages, &res); err != nil {
		s.wf.FatalError(err)
	}

	return res, nil
}

package base

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	"github.com/urandom/readeef/content"
	"github.com/urandom/readeef/content/data"
	"github.com/urandom/readeef/parser"
)

type Feed struct {
	ArticleSorting
	Error
	RepoRelated

	data           data.Feed
	parsedArticles []content.Article
}

type UserFeed struct {
	ArticleSearch
	UserRelated
}

func (f Feed) String() string {
	return f.data.Title + " " + strconv.FormatInt(int64(f.data.Id), 10)
}

func (f *Feed) Data(d ...data.Feed) data.Feed {
	if f.HasErr() {
		return data.Feed{}
	}

	if len(d) > 0 {
		f.data = d[0]
	}

	return f.data
}

func (f Feed) Validate() error {
	if f.data.Link == "" {
		return ValidationError{errors.New("Feed has no link")}
	}

	if u, err := url.Parse(f.data.Link); err != nil || !u.IsAbs() {
		return ValidationError{errors.New("Feed has no link")}
	}

	return nil
}

func (f *Feed) Refresh(pf parser.Feed) {
	if f.HasErr() {
		return
	}

	in := f.Data()

	in.Title = pf.Title
	in.Description = pf.Description
	in.SiteLink = pf.SiteLink
	in.HubLink = pf.HubLink

	f.parsedArticles = make([]content.Article, len(pf.Articles))

	for i := range pf.Articles {
		ai := data.Article{
			Title:       pf.Articles[i].Title,
			Description: pf.Articles[i].Description,
			Link:        pf.Articles[i].Link,
			Date:        pf.Articles[i].Date,
		}
		ai.FeedId = in.Id

		if pf.Articles[i].Guid != "" {
			ai.Guid.Valid = true
			ai.Guid.String = pf.Articles[i].Guid
		}

		a := &Article{data: ai}

		f.parsedArticles[i] = a
	}

	f.Data(in)
}

func (f *Feed) ParsedArticles() (a []content.Article) {
	if f.HasErr() {
		return
	}

	return f.parsedArticles
}

func (f Feed) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.data)
}

func (uf UserFeed) Validate() error {
	if uf.user.Data().Login == "" {
		return ValidationError{errors.New("UserFeed has no user")}
	}

	return nil
}

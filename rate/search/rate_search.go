package search

import (
	"context"
	"database/sql"
	"database/sql/driver"
	. "github.com/core-go/sql"
	"reflect"
)

type RateCommentSearchService interface {
	Search(ctx context.Context, rf *RateFilter) ([]Rate, int64, error)
}

type rateCommentSearchService struct {
	Database       *sql.DB
	BuildQuery     func(sm interface{}) (string, []interface{})
	BuildFromQuery func(ctx context.Context, db *sql.DB, fieldsIndex map[string]int, models interface{}, query string, params []interface{}, limit int64, offset int64, toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}, options ...func(context.Context, interface{}) (interface{}, error)) (int64, error)
	ModelType   reflect.Type
	Map         func(ctx context.Context, model interface{}) (interface{}, error)
	fieldsIndex map[string]int
	ToArray     func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	queryInfo func(ids []string) ([]Info, error)
	getOffset func(limit int64, page int64, opts ...int64) int64
}

func NewRateSearchService(Database *sql.DB,
	BuildQuery func(sm interface{}) (string, []interface{}),
	ToArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}, queryInfo func(ids []string) ([]Info, error),
	buildFromQuery func(ctx context.Context, db *sql.DB, fieldsIndex map[string]int, models interface{}, query string, params []interface{}, limit int64, offset int64, toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}, options ...func(context.Context, interface{}) (interface{}, error)) (int64, error),
	getOffset func(limit int64, page int64, opts ...int64) int64,
) (*rateCommentSearchService, error) {
	modelType := reflect.TypeOf(Rate{})
	fieldsIndex, _ := GetColumnIndexes(modelType)
	return &rateCommentSearchService{
		Database:       Database,
		BuildQuery:     BuildQuery,
		BuildFromQuery: buildFromQuery,
		ModelType:      modelType,
		Map:            nil,
		fieldsIndex:    fieldsIndex,
		queryInfo:      queryInfo,
		getOffset:      getOffset,
		ToArray:        ToArray,
	}, nil
}

func (f *rateCommentSearchService) Search(ctx context.Context, rf *RateFilter) ([]Rate, int64, error) {
	sql, params := f.BuildQuery(rf)
	rates := make([]Rate, 0)
	if rf.Page == 0 {
		rf.Page = 1
	}
	offset := f.getOffset(rf.Limit, rf.Page)

	total1, er2 := f.BuildFromQuery(ctx, f.Database, f.fieldsIndex, &rates, sql, params, rf.Limit, offset, f.ToArray, f.Map)
	if er2 != nil {
		return rates, total1, er2
	}
	if f.queryInfo == nil {
		return rates, total1, nil
	}
	ids := make([]string, 0)
	for _, r := range rates {
		ids = append(ids, r.Author)
	}
	if f.queryInfo == nil {
		return rates, total1, nil
	}
	infos, err := f.queryInfo(ids)
	if err != nil {
		return rates, total1, err
	}
	for k, _ := range rates {
		c := rates[k]
		i := BinarySearch(infos, c.Author)
		if i >= 0 && !c.Anonymous {
			if rates[k].Author == infos[i].Id {
				rates[k].AuthorURL = &infos[i].Url
				rates[k].AuthorName = &infos[i].Name
			}
		}
	}
	if err != nil {
		return rates, total1, err
	}
	return rates, total1, nil
}

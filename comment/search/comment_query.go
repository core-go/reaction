package search

import (
	"fmt"
	q "github.com/core-go/sql"
	"strings"
)

func BuildCommentQuery(filter interface{}) (query string, params []interface{}) {
	query = `select * from ratecomment`
	s := filter.(*CommentFilter)
	var where []string

	i := 1
	if len(s.CommentId) > 0 {
		where = append(where, fmt.Sprintf(`commentid = %s`, q.BuildDollarParam(i)))
		params = append(params, s.CommentId)
		i++
	}
	if len(s.Id) > 0 {
		where = append(where, fmt.Sprintf(`id = %s`, q.BuildDollarParam(i)))
		params = append(params, s.Id)
		i++
	}
	if len(s.Author) > 0 {
		where = append(where, fmt.Sprintf(`author = %s`, q.BuildDollarParam(i)))
		params = append(params, s.Author)
		i++
	}
	if len(s.UserId) > 0 {
		where = append(where, fmt.Sprintf(`userid = %s`, q.BuildDollarParam(i)))
		params = append(params, s.UserId)
		i++
	}
	if len(s.Comment) > 0 {
		where = append(where, fmt.Sprintf(`comment ilike %s`, q.BuildDollarParam(i)))
		params = append(params, "%" + s.Comment + "%")
		i++
	}
	if s.Time != nil {
		if s.Time.Min != nil {
			where = append(where, fmt.Sprintf(`time >= %s`, q.BuildDollarParam(i)))
			params = append(params, s.Time.Min)
			i++
		}
		if s.Time.Max != nil {
			where = append(where, fmt.Sprintf(`time <= %s`, q.BuildDollarParam(i)))
			params = append(params, s.Time.Max)
			i++
		}
	}

	if len(where) > 0 {
		query = query + ` where ` + strings.Join(where, " and ")
	}
	return
}

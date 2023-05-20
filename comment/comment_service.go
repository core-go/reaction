package comment

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type CommentService interface {
	Load(ctx context.Context, id string, author string) ([]Response, error)
	Create(ctx context.Context, id string, commentId string, userId string, author string, comment Request) (int64, error)
	Update(ctx context.Context, id string, commentId string, userId string, author string, comment Request) (int64, error)
	Delete(ctx context.Context, id string, commentId string, author string) (int64, error)
}

func NewCommentService(db *sql.DB, commentTable string, commentIdCol string, idCol string, authorCol string, userIdCol string, commentCol string, anonymousCol string, timeCol string, updatedAtCol string, rateTable string, rateIdCol string, rateAuthorCol string, commentCountCol string, userTable string, userIdUserCol string, imageUrlUserCol string, UsernameUserCol string, queryInfo func(ids []string) ([]Info, error), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) CommentService {
	return &commentService{
		DB:              db,
		CommentTable:    commentTable,
		CommentIdCol:    commentIdCol,
		IdCol:           idCol,
		AuthorCol:       authorCol,
		UserIdCol:       userIdCol,
		CommentCol:      commentCol,
		AnonymousCol:    anonymousCol,
		TimeCol:         timeCol,
		UpdatedAtCol:    updatedAtCol,
		RateTable:       rateTable,
		RateIdCol:       rateIdCol,
		RateAuthorCol:   rateAuthorCol,
		CommentCountCol: commentCountCol,
		userTable:       userTable,
		imageUrlUserCol: imageUrlUserCol,
		userIdUserCol:   userIdUserCol,
		ToArray:         toArray,
		QueryInfo:       queryInfo,
		UsernameUserCol: UsernameUserCol,
	}
}

type commentService struct {
	DB              *sql.DB
	CommentTable    string
	CommentIdCol    string
	IdCol           string
	AuthorCol       string
	UserIdCol       string
	CommentCol      string
	TimeCol         string
	UpdatedAtCol    string
	RateTable       string
	RateIdCol       string
	RateAuthorCol   string
	AnonymousCol    string
	CommentCountCol string
	userTable       string
	userIdUserCol   string
	imageUrlUserCol string
	UsernameUserCol string
	QueryInfo       func(ids []string) ([]Info, error)
	ToArray         func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (s *commentService) Load(ctx context.Context, id string, author string) ([]Response, error) {
	var comments []Comment
	var rs []Response
	query := fmt.Sprintf(
		"select s.%s, s.%s, s.%s, s.%s, s.%s, s.%s, s.%s, s.%s, s.histories from %s s where s.%s = $1 and s.%s = $2",
		s.CommentIdCol, s.IdCol, s.AuthorCol, s.UserIdCol, s.CommentCol, s.AnonymousCol, s.TimeCol, s.UpdatedAtCol, s.CommentTable, s.IdCol, s.AuthorCol)
	rows, err := s.DB.QueryContext(ctx, query, id, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.CommentId, &comment.Id, &comment.Author, &comment.UserId, &comment.Comment, &comment.Anonymous, &comment.Time, &comment.UpdatedAt, s.ToArray(&comment.Histories))
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if len(comments) == 0 {
		return rs, nil
	}
	ids := make([]string, 0)
	for _, r := range comments {
		ids = append(ids, r.UserId)
	}
	if s.QueryInfo == nil {
		for k, _ := range comments {
			c := comments[k]
			r := toResponse(c)
			rs = append(rs, r)
		}
		return rs, nil
	}
	infos, err := s.QueryInfo(ids)
	if err != nil {
		return nil, err
	}
	for k, _ := range comments {
		c := comments[k]
		r := toResponse(c)
		i := BinarySearch(infos, c.UserId)
		if i >= 0 && !c.Anonymous {
			if infos[i].Id == c.UserId {
				r.AuthorURL = &infos[i].Url
				r.AuthorName = &infos[i].Name
			}
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func (s *commentService) Create(ctx context.Context, id string, commentId string, userId string, author string, rq Request) (int64, error) {
	var t = time.Now()
	comment := Comment{Id: id, CommentId: commentId, Author: author, UserId: userId, Comment: rq.Comment, Anonymous: rq.Anonymous, Time: &t}
	query1 := fmt.Sprintf(
		"insert into %s(%s, %s, %s, %s, %s, %s, %s) values ($1, $2, $3, $4, $5, $6, $7)",
		s.CommentTable, s.CommentIdCol, s.IdCol, s.AuthorCol, s.UserIdCol, s.CommentCol, s.AnonymousCol, s.TimeCol)
	stmt1, err := s.DB.Prepare(query1)
	if err != nil {
		return -1, err
	}
	res1, err := stmt1.ExecContext(ctx, comment.CommentId, comment.Id, comment.Author, comment.UserId, comment.Comment, comment.Anonymous, comment.Time)
	if err != nil {
		return -1, err
	}

	query2 := fmt.Sprintf(
		"update %s set %s = %s.%s + 1 where %s = $1 and %s = $2",
		s.RateTable, s.CommentCountCol, s.RateTable, s.CommentCountCol, s.RateIdCol, s.RateAuthorCol)
	stmt2, err := s.DB.Prepare(query2)
	if err != nil {
		return -1, err
	}
	stmt2.ExecContext(ctx, comment.Id, comment.Author)

	return res1.RowsAffected()
}

func (s *commentService) Update(ctx context.Context, id string, commentId string, userId string, author string, req Request) (int64, error) {
	t := time.Now()
	var comment = Comment{Id: id, CommentId: commentId, UserId: userId, Author: author,
		Comment: req.Comment, Anonymous: req.Anonymous, UpdatedAt: &t}
	var oldComment Comment
	query1 := fmt.Sprintf("select %s, %s, %s, histories from %s where %s = $1 limit 1", s.TimeCol, s.UpdatedAtCol, s.CommentCol, s.CommentTable, s.CommentIdCol)
	rows, _ := s.DB.QueryContext(ctx, query1, comment.CommentId)
	for rows.Next() {
		err := rows.Scan(&oldComment.Time, &oldComment.UpdatedAt, &oldComment.Comment, s.ToArray(&oldComment.Histories))
		if err != nil {
			return 0, err
		}
	}
	rows.Close()

	if oldComment.Histories != nil {
		comment.Histories = append(oldComment.Histories, Histories{Time: oldComment.UpdatedAt, Comment: oldComment.Comment})
	} else {
		comment.Histories = append(oldComment.Histories, Histories{Time: oldComment.Time, Comment: oldComment.Comment})
	}

	query := fmt.Sprintf(
		"update %s set %s = $1, %s = $2, histories = $3 where %s = $4;",
		s.CommentTable, s.CommentCol, s.UpdatedAtCol, s.CommentIdCol)
	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return -1, err
	}
	res, err := stmt.ExecContext(ctx, comment.Comment, time.Now(), s.ToArray(comment.Histories), comment.CommentId)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

func (s *commentService) Delete(ctx context.Context, id string, commentId string, author string) (int64, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()
	query1 := fmt.Sprintf("delete from %s where %s = $1", s.CommentTable, s.CommentIdCol)
	stmt1, er0 := tx.Prepare(query1)
	if er0 != nil {
		return -1, nil
	}
	res1, er1 := stmt1.ExecContext(ctx, commentId)
	if er1 != nil {
		return -1, er1
	}

	query2 := fmt.Sprintf(
		"update %s set %s = %s.%s - 1 where %s = $1 and %s = $2",
		s.RateTable, s.CommentCountCol, s.RateTable, s.CommentCountCol, s.RateIdCol, s.RateAuthorCol)
	stmt2, err := tx.Prepare(query2)
	if err != nil {
		return -1, err
	}
	_, err1 := stmt2.ExecContext(ctx, id, author)
	if err1 != nil {
		return -1, err
	}
	r, _ := res1.RowsAffected()
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return r, nil
}

package comment

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type CommentService interface {
	Load(ctx context.Context, id string, author string) ([]*Comment, error)
	Create(ctx context.Context, comment *Comment) (int64, error)
	Update(ctx context.Context, comment *Comment) (int64, error)
	Delete(ctx context.Context, commentId string, id string, author string) (int64, error)
}

func NewCommentService(
	db *sql.DB,
	commentTable string,
	commentIdCol string,
	idCol string,
	authorCol string,
	userIdCol string,
	commentCol string,
	timeCol string,
	updatedAtCol string,
	rateTable string,
	rateIdCol string,
	rateAuthorCol string,
	commentCountCol string,
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	},
) CommentService {
	return &commentService{
		DB:              db,
		CommentTable:    commentTable,
		CommentIdCol:    commentIdCol,
		IdCol:           idCol,
		AuthorCol:       authorCol,
		UserIdCol:       userIdCol,
		CommentCol:      commentCol,
		TimeCol:         timeCol,
		UpdatedAtCol:    updatedAtCol,
		RateTable:       rateTable,
		RateIdCol:       rateIdCol,
		RateAuthorCol:   rateAuthorCol,
		CommentCountCol: commentCountCol,
		ToArray:         toArray,
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
	CommentCountCol string
	ToArray         func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (s *commentService) Load(ctx context.Context, id string, author string) ([]*Comment, error) {
	var comments []*Comment
	query := fmt.Sprintf(
		"select %s, %s, %s, %s, %s, %s, %s, histories from %s where %s = $1 and %s = $2",
		s.CommentIdCol, s.IdCol, s.AuthorCol, s.UserIdCol, s.CommentCol, s.TimeCol, s.UpdatedAtCol, s.CommentTable, s.IdCol, s.AuthorCol)
	fmt.Println(query)
	rows, err := s.DB.QueryContext(ctx, query, id, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.CommentId, &comment.Id, &comment.Author, &comment.UserId, &comment.Comment, &comment.Time, &comment.UpdatedAt, s.ToArray(&comment.Histories))
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	return comments, nil
}

func (s *commentService) Create(ctx context.Context, comment *Comment) (int64, error) {
	query1 := fmt.Sprintf(
		"insert into %s(%s, %s, %s, %s, %s, %s) values ($1, $2, $3, $4, $5, $6)",
		s.CommentTable, s.CommentIdCol, s.IdCol, s.AuthorCol, s.UserIdCol, s.CommentCol, s.TimeCol)
	fmt.Println(query1)
	stmt1, err := s.DB.Prepare(query1)
	if err != nil {
		return -1, err
	}
	res1, err := stmt1.ExecContext(ctx, comment.CommentId, comment.Id, comment.Author, comment.UserId, comment.Comment, comment.Time)
	if err != nil {
		return -1, err
	}

	query2 := fmt.Sprintf(
		"update %s set %s = %s.%s + 1 where %s = $1 and %s = $2",
		s.RateTable, s.CommentCountCol, s.RateTable, s.CommentCountCol, s.RateIdCol, s.RateAuthorCol)
	fmt.Println(query2)
	stmt2, err := s.DB.Prepare(query2)
	if err != nil {
		return -1, err
	}
	stmt2.ExecContext(ctx, comment.Id, comment.Author)

	return res1.RowsAffected()
}

func (s *commentService) Update(ctx context.Context, comment *Comment) (int64, error) {
	var oldComment Comment
	query1 := fmt.Sprintf("select %s, %s, %s, histories from %s where %s = $1 limit 1", s.TimeCol, s.UpdatedAtCol, s.CommentCol, s.CommentTable, s.CommentIdCol)
	fmt.Println(query1)
	rows, _ := s.DB.QueryContext(ctx, query1, comment.CommentId)
	for rows.Next() {
		rows.Scan(&oldComment.Time, &oldComment.UpdatedAt, &oldComment.Comment, s.ToArray(&oldComment.Histories))
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
	fmt.Println(query)
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

func (s *commentService) Delete(ctx context.Context, commentId string, id string, author string) (int64, error) {
	query1 := fmt.Sprintf("delete from %s where %s = $1", s.CommentTable, s.CommentIdCol)
	fmt.Println(query1)
	stmt1, er0 := s.DB.Prepare(query1)
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
	fmt.Println(query2)
	stmt2, err := s.DB.Prepare(query2)
	if err != nil {
		return -1, err
	}
	_, err1 := stmt2.ExecContext(ctx, id, author)
	if err1 != nil {
		return -1, err
	}
	return res1.RowsAffected()
}

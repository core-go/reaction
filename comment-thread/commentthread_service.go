package commentthread

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type CommentThreadService interface {
	Comment(ctx context.Context, comment CommentThread) (int64, error)
	Update(ctx context.Context, comment CommentThread) (int64, error)
	Remove(ctx context.Context, commentId string) (int64, error)
	Load(ctx context.Context, commentId string) (*CommentThread, error)
}

func NewCommentThreadService(
	db *sql.DB,
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	},
	threadTable string,
	commentIdThreadCol string,
	idThreadCol string,
	authorThreadCol string,
	historiesThreadCol string,
	commentThreadCol string,
	timeThreadCol string,
	userIdThreadCol string,
	updatedAtCol string,
	threadReplyTable string,
	commentIdThreadReplyCol string,
	commentThreadIdReplyCol string,
	threadInfoTable string,
	commentIdthreadInfo string,
	threadReplyInfoTable string,
	commentIdThreadReplyInfoCol string,
	reactionTable string,
	commentIdReactionCol string,
	reactionReplyTable string,
	commentIdReactionRelyCol string,
) CommentThreadService {
	return &commentThreadService{
		db:                          db,
		toArray:                     toArray,
		threadTable:                 threadTable,
		commentIdThreadCol:          commentIdThreadCol,
		idThreadCol:                 idThreadCol,
		authorThreadCol:             authorThreadCol,
		historiesThreadCol:          historiesThreadCol,
		commentThreadCol:            commentThreadCol,
		timeThreadCol:               timeThreadCol,
		updatedAtCol:                updatedAtCol,
		userIdThreadCol:             userIdThreadCol,
		threadReplyTable:            threadReplyTable,
		commentIdThreadReplyCol:     commentIdThreadReplyCol,
		commentThreadIdReplyCol:     commentThreadIdReplyCol,
		threadInfoTable:             threadInfoTable,
		commentIdthreadInfo:         commentIdthreadInfo,
		threadReplyInfoTable:        threadReplyInfoTable,
		commentIdThreadReplyInfoCol: commentIdThreadReplyInfoCol,
		reactionTable:               reactionTable,
		commentIdReactionCol:        commentIdReactionCol,
		reactionReplyTable:          reactionReplyTable,
		commentIdReactionRelyCol:    commentIdReactionRelyCol,
	}
}

type commentThreadService struct {
	db      *sql.DB
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	threadTable                 string
	commentIdThreadCol          string
	idThreadCol                 string
	authorThreadCol             string
	historiesThreadCol          string
	commentThreadCol            string
	timeThreadCol               string
	updatedAtCol                string
	userIdThreadCol             string
	threadReplyTable            string
	commentIdThreadReplyCol     string
	commentThreadIdReplyCol     string
	threadInfoTable             string
	commentIdthreadInfo         string
	threadReplyInfoTable        string
	commentIdThreadReplyInfoCol string
	reactionTable               string
	commentIdReactionCol        string
	reactionReplyTable          string
	commentIdReactionRelyCol    string
}

func (s *commentThreadService) Remove(ctx context.Context, commentId string) (int64, error) {
	var rowResult int64
	rowResult = 0
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()
	qr1 := fmt.Sprintf("Delete from %s where %s = $1", s.threadTable, s.commentIdThreadCol)
	res, err := tx.ExecContext(ctx, qr1, commentId)
	if err != nil {
		return -1, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowResult += rowsAffected

	qr2 := fmt.Sprintf("Delete from %s where %s = $1", s.threadReplyTable, s.commentThreadIdReplyCol)
	res, err = tx.ExecContext(ctx, qr2, commentId)
	if err != nil {
		return -1, err
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowResult += rowsAffected

	qr3 := fmt.Sprintf("Delete from %s where %s = $1", s.threadInfoTable, s.commentIdthreadInfo)
	res, err = tx.ExecContext(ctx, qr3, commentId)
	if err != nil {
		return -1, err
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowResult += rowsAffected
	idQr := fmt.Sprintf("select %s from %s where %s = $1", s.commentIdThreadReplyCol, s.threadReplyTable, s.commentThreadIdReplyCol)
	rows, err := s.db.QueryContext(ctx, idQr, commentId)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	var ids = []string{}
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}
	qr4 := fmt.Sprintf("delete from %s a where a.%s = ANY($1)",
		s.threadReplyInfoTable, s.commentIdThreadReplyInfoCol)

	res, err = tx.ExecContext(ctx, qr4, s.toArray(ids))
	if err != nil {
		return -1, err
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowResult += rowsAffected
	qr5 := fmt.Sprintf("delete from %s a where %s = $1",
		s.reactionTable, s.commentIdReactionCol)

	res, err = tx.ExecContext(ctx, qr5, commentId)
	if err != nil {
		return -1, err
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowResult += rowsAffected
	qr6 := fmt.Sprintf("delete from %s a where a.%s = ANY($1)",
		s.reactionReplyTable, s.commentIdReactionRelyCol)

	res2, err := tx.ExecContext(ctx, qr6, s.toArray(ids))
	if err != nil {
		return -1, err
	}
	rowsAffected, err = res2.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowResult += rowsAffected

	if err = tx.Commit(); err != nil {
		return -1, err
	}
	return rowResult, nil
}

func (s *commentThreadService) Load(ctx context.Context, commentId string) (*CommentThread, error) {
	qr1 := fmt.Sprintf("Select * from %s where %s = $1", s.threadTable, s.commentIdThreadCol)
	res, err := s.db.QueryContext(ctx, qr1, commentId)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	for res.Next() {
		var item CommentThread
		err = res.Scan(
			&item.CommentId,
			&item.Id,
			&item.Author,
			&item.Comment,
			&item.Time,
			&item.UpdatedAt,
			s.toArray(&item.Histories),
			&item.UserId)
		if err != nil {
			return nil, err
		}
		return &item, nil
	}
	return nil, nil
}

func (s *commentThreadService) Update(ctx context.Context, comment CommentThread) (int64, error) {
	exist, err := s.Load(ctx, comment.CommentId)
	if err != nil {
		return -1, err
	}
	if exist != nil {
		updatedTime := time.Now()
		exist.Histories = append(exist.Histories, History{Comment: comment.Comment, Time: updatedTime})
		qr1 := fmt.Sprintf("update %s set %s = $1, %s = $2, %s = $3 where %s = $4",
			s.threadTable, s.commentThreadCol, s.updatedAtCol, s.historiesThreadCol, s.commentIdThreadCol)
		res, err := s.db.ExecContext(ctx, qr1, comment.Comment, updatedTime, s.toArray(exist.Histories), comment.CommentId)
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	}
	return -1, nil
}

func (s *commentThreadService) Comment(ctx context.Context, comment CommentThread) (int64, error) {
	qr1 := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s,%s,%s) values($1, $2, $3, $4, $5, $6, $7)",
		s.threadTable, s.commentIdThreadCol, s.idThreadCol, s.authorThreadCol, s.userIdThreadCol, s.commentThreadCol, s.timeThreadCol, s.historiesThreadCol)
	res, err := s.db.ExecContext(ctx, qr1, comment.CommentId, comment.Id, comment.Author, comment.UserId, comment.Comment, comment.Time, s.toArray([]History{}))
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

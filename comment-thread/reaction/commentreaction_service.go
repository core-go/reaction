package reaction

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CommentReactionService interface {
	Save(ctx context.Context, commentId string, author string, userId string, reaction int) (int64, error)
	Remove(ctx context.Context, commentId string, author string, userId string) (int64, error)
}

type commentReactionService struct {
	db              *sql.DB
	reactionTable   string
	commentIdCol    string
	authorCol       string
	userIdCol       string
	reactionCol     string
	timeCol         string
	parentTable     string
	parentIdCol     string
	parentUsefulCol string
}

func NewCommentReactionService(db *sql.DB, reactionTable string, commentIdCol string,
	authorCol string, userIdCol string, timeCol string, reactionCol string, parentTable string, parentIdCol string, parentUsefulCol string) CommentReactionService {

	return &commentReactionService{
		db:              db,
		reactionTable:   reactionTable,
		commentIdCol:    commentIdCol,
		authorCol:       authorCol,
		userIdCol:       userIdCol,
		reactionCol:     reactionCol,
		timeCol:         timeCol,
		parentTable:     parentTable,
		parentIdCol:     parentIdCol,
		parentUsefulCol: parentUsefulCol,
	}
}

func (s *commentReactionService) Remove(ctx context.Context, commentId string, author string, userId string) (int64, error) {
	result := int64(0)
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return -1, err
	}
	qr := fmt.Sprintf(`delete from %s where  %s= $1 and %s = $2 and %s= $3`,
		s.reactionTable, s.commentIdCol, s.authorCol, s.userIdCol)
	rows, err := tx.ExecContext(ctx, qr, commentId, author, userId)
	if err != nil {
		return -1, err
	}
	numRows, err := rows.RowsAffected()
	if err != nil {
		return -1, err
	}
	result += numRows
	qr = fmt.Sprintf(`update %s set %s = %s - 1 where %s = $1`,
		s.parentTable, s.parentUsefulCol, s.parentUsefulCol, s.parentIdCol)
	rows, err = tx.ExecContext(ctx, qr, commentId)
	if err != nil {
		return -1, err
	}
	numRows, err = rows.RowsAffected()
	if err != nil {
		return -1, err
	}
	result += numRows
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return numRows, nil
}

func (s *commentReactionService) Save(ctx context.Context, commentId string, author string, userId string, reaction int) (int64, error) {
	result := int64(0)
	tx, err := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return -1, err
	}
	qr := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s) values($1,$2,$3,$4,$5)",
		s.reactionTable, s.commentIdCol, s.authorCol, s.userIdCol, s.timeCol, s.reactionCol)
	res, err := tx.ExecContext(ctx, qr, commentId, author, userId, time.Now(), reaction)
	if err != nil {
		return -1, err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	result += numRows
	qr2 := fmt.Sprintf("insert into %s(%s,%s) values($1,1) on conflict(%s) do update set %s = %s.%s + 1 where %s.%s = $1",
		s.parentTable, s.parentIdCol, s.parentUsefulCol, s.parentIdCol, s.parentUsefulCol, s.parentTable, s.parentUsefulCol, s.parentTable, s.parentIdCol)

	rows, err := tx.ExecContext(ctx, qr2, commentId)
	if err != nil {
		return -1, err
	}
	numRows, err = rows.RowsAffected()
	if err != nil {
		return -1, err
	}
	result += numRows
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return result, nil
}

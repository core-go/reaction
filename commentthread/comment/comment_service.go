package comment

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type CommentService interface {
	GetComments(ctx context.Context, commentThreadId string, userId *string) ([]Comment, error)
	Create(ctx context.Context, comment Comment) (int64, error)
	Update(ctx context.Context, comment Comment) (int64, error)
	Remove(ctx context.Context, commentId string, commentThreadId string) (int64, error)
}

func NewCommentService(db *sql.DB, replyTable string, commentIdCol string, authorCol string, idCol string, updatedAtCol string, commentCol string, userIdCol string, timeCol string, parentCol string, historiesCol string, commentThreadIdCol string, reactionCol string, commentReactionTable string, commentIdReactionCol string, userTable string, userIdUserCol string, usernameUserCol string, avatarUserCol string, commentInfoTable string, userfulCountInfoCol string, commentIdInfoCol string, commentThreadInfoTable string, commentIdCommentThreadInfoCol string, replyCountCommentThreadInfoCol string, usefulCountCommentThreadInfoCol string, queryInfo func(ids []string) ([]Info, error), toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) CommentService {
	return &commentService{
		db:                              db,
		ReplyTable:                      replyTable,
		commentIdCol:                    commentIdCol,
		authorCol:                       authorCol,
		idCol:                           idCol,
		parentCol:                       parentCol,
		updatedAtCol:                    updatedAtCol,
		historiesCol:                    historiesCol,
		commentCol:                      commentCol,
		reactionCol:                     reactionCol,
		userIdCol:                       userIdCol,
		timeCol:                         timeCol,
		commentThreadIdCol:              commentThreadIdCol,
		commentReactionTable:            commentReactionTable,
		commentIdReactionCol:            commentIdReactionCol,
		userTable:                       userTable,
		userIdUserCol:                   userIdUserCol,
		usernameUserCol:                 usernameUserCol,
		avatarUserCol:                   avatarUserCol,
		commentInfoTable:                commentInfoTable,
		userfulCountInfoCol:             userfulCountInfoCol,
		commentIdInfoCol:                commentIdInfoCol,
		commentThreadInfoTable:          commentThreadInfoTable,
		commentIdCommentThreadInfoCol:   commentIdCommentThreadInfoCol,
		replyCountCommentThreadInfoCol:  replyCountCommentThreadInfoCol,
		usefulCountCommentThreadInfoCol: usefulCountCommentThreadInfoCol,
		queryInfo:                       queryInfo,
		toArray:                         toArray,
	}
}

type commentService struct {
	db                              *sql.DB
	ReplyTable                      string
	commentIdCol                    string
	authorCol                       string
	idCol                           string
	parentCol                       string
	updatedAtCol                    string
	historiesCol                    string
	commentCol                      string
	reactionCol                     string
	userIdCol                       string
	timeCol                         string
	commentThreadIdCol              string
	commentReactionTable            string
	commentIdReactionCol            string
	userTable                       string
	userIdUserCol                   string
	usernameUserCol                 string
	avatarUserCol                   string
	commentInfoTable                string
	userfulCountInfoCol             string
	commentIdInfoCol                string
	commentThreadInfoTable          string
	commentIdCommentThreadInfoCol   string
	replyCountCommentThreadInfoCol  string
	usefulCountCommentThreadInfoCol string

	queryInfo func(ids []string) ([]Info, error)
	toArray   func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (s *commentService) Update(ctx context.Context, comment Comment) (int64, error) {
	qr := fmt.Sprintf("select %s, %s,%s from %s where %s = $1", s.commentIdCol, s.commentCol, s.historiesCol, s.ReplyTable, s.commentIdCol)
	rows := s.db.QueryRow(qr, comment.CommentId)
	var exist = Comment{}
	err := rows.Scan(&exist.CommentId, &exist.Comment, s.toArray(&exist.Histories))
	if err != nil {
		return -1, err
	}
	updatedTime := time.Now()
	exist.Histories = append(exist.Histories, History{
		Comment: exist.Comment,
		Time:    updatedTime,
	})
	qr1 := fmt.Sprintf("update %s set %s = $1, %s = $2, %s = $3", s.ReplyTable, s.commentCol, s.historiesCol, s.updatedAtCol)
	rows2, er2 := s.db.ExecContext(ctx, qr1, comment.Comment, s.toArray(exist.Histories), updatedTime)
	if er2 != nil {
		return -1, er2
	}
	return rows2.RowsAffected()
}

// Delete implements CommentThreadReplyService
func (s *commentService) Remove(ctx context.Context, commentId string, commentThreadId string) (int64, error) {
	rowsAffected := int64(0)
	qr1 := fmt.Sprintf("delete from %s where %s = $1", s.ReplyTable, s.commentIdCol)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	rows, err := tx.ExecContext(ctx, qr1, commentId)
	if err != nil {
		return -1, err
	}
	numberRows, err := rows.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowsAffected += numberRows
	qr2 := fmt.Sprintf("delete from %s where %s = $1", s.commentInfoTable, s.commentIdInfoCol)
	rows2, err := tx.ExecContext(ctx, qr2, commentId)
	if err != nil {
		return -1, err
	}
	numberRows, err = rows2.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowsAffected += numberRows
	qr3 := fmt.Sprintf("delete from %s where %s = $1", s.commentReactionTable, s.commentIdReactionCol)
	rows3, err := tx.ExecContext(ctx, qr3, commentId)
	if err != nil {
		return -1, err
	}
	numberRows, err = rows3.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowsAffected += numberRows
	qr4 := fmt.Sprintf("update %s set %s = %s - 1 where %s = $1", s.commentThreadInfoTable, s.replyCountCommentThreadInfoCol, s.replyCountCommentThreadInfoCol, s.commentIdCommentThreadInfoCol)
	rows4, err := tx.ExecContext(ctx, qr4, commentThreadId)
	if err != nil {
		return -1, err
	}
	numberRows, err = rows4.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowsAffected += numberRows
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return rowsAffected, nil
}

// Create implements CommentThreadReplyService
func (s *commentService) Create(ctx context.Context, comment Comment) (int64, error) {
	rowsAffected := int64(0)
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()

	qr := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s,%s,%s,%s,%s) values($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		s.ReplyTable, s.commentIdCol, s.idCol, s.authorCol, s.userIdCol, s.commentCol, s.parentCol, s.timeCol, s.historiesCol, s.commentThreadIdCol)
	rows1, err := tx.ExecContext(ctx, qr, comment.CommentId, comment.Id, comment.Author, comment.UserId, comment.Comment, comment.Parent, time.Now(), s.toArray([]interface{}{}), comment.CommentThreadId)
	if err != nil {
		return -1, err
	}
	numberRows, err := rows1.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowsAffected += numberRows
	qr2 := fmt.Sprintf("insert into %s(%s,%s,%s) values($1, 1, 0) on conflict(%s) do update set %s = %s.%s + 1 where %s.%s = $1",
		s.commentThreadInfoTable,
		s.commentIdCommentThreadInfoCol,
		s.replyCountCommentThreadInfoCol,
		s.usefulCountCommentThreadInfoCol,
		s.commentIdCommentThreadInfoCol,
		s.replyCountCommentThreadInfoCol,
		s.commentThreadInfoTable,
		s.replyCountCommentThreadInfoCol,
		s.commentThreadInfoTable,
		s.commentIdCommentThreadInfoCol)
	rows2, err := tx.ExecContext(ctx, qr2, comment.CommentThreadId)
	if err != nil {
		return -1, err
	}
	numberRows, err = rows2.RowsAffected()
	if err != nil {
		return -1, err
	}
	rowsAffected += numberRows
	err = tx.Commit()
	if err != nil {
		return -1, err
	}
	return rowsAffected, nil
}

func (s *commentService) GetComments(ctx context.Context, commentThreadId string, userId *string) ([]Comment, error) {
	qr := ""
	qr2 := ""
	arr := []interface{}{}
	if userId != nil && len(*userId) > 0 {
		arr = append(arr, userId)
		qr = fmt.Sprintf(`, case when d.%s = 1 then true else false end as disable`, s.reactionCol)
		qr2 = fmt.Sprintf(`left join %s d on a.%s = d.%s and a.%s = $1`,
			s.commentReactionTable, s.commentIdCol, s.commentIdReactionCol, s.userIdCol)
	}
	param := "1"
	if userId != nil && len(*userId) > 0 {
		param = "2"
	}
	query := fmt.Sprintf(`select a.*,c.%s%s from %s a
                                  left join %s c on a.%s = c.%s %s 
                                  where a.%s = $%s`, s.userfulCountInfoCol, qr, s.ReplyTable,
		s.commentInfoTable, s.commentIdCol, s.commentIdInfoCol, qr2,
		s.commentThreadIdCol, param)
	arr = append(arr, commentThreadId)
	fmt.Println(query)
	rows, err := s.db.QueryContext(ctx, query, arr...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.CommentId,
			&comment.CommentThreadId,
			&comment.Id,
			&comment.Author,
			&comment.UserId,
			&comment.Comment,
			&comment.Parent,
			&comment.Time,
			&comment.UpdatedAt,
			s.toArray(&comment.Histories),
			&comment.UsefulCount,
			&comment.Disable,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}
	if len(comments) == 0 {
		return comments, nil
	}
	ids := make([]string, 0)
	for _, r := range comments {
		ids = append(ids, r.Author)
	}
	infos, err := s.queryInfo(ids)
	if err != nil {
		return nil, err
	}
	for k, _ := range comments {
		c := comments[k]
		i := BinarySearch(infos, c.Author)
		if i >= 0 {
			comments[k].AuthorURL = &infos[i].Url
			comments[k].AuthorName = &infos[i].Name
		}
	}

	return comments, nil
}

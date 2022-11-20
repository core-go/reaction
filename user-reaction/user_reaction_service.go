package user_reaction

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type UserReactionService interface {
	React(ctx context.Context, id string, author string, reaction string) (int64, error)
	Unreact(ctx context.Context, id string, author string, reaction string) (int64, error)
	CheckReaction(ctx context.Context, id string, author string) (int64, error)
}

func NewUserReactionService(
	DB *sql.DB,
	userReactionTable string,
	id string,
	author string,
	reaction string,

	userinfoTable string,
	infoId string,
	reactionCount string,
	prefix string,
	suffix string,
) UserReactionService {
	return &userReactionService{
		DB:                DB,
		userReactionTable: userReactionTable,
		id:                id,
		author:            author,
		reaction:          reaction,
		prefix:            prefix,
		suffix:            suffix,
		userinfoTable:     userinfoTable,
		infoId:            infoId,
		reactionCount:     reactionCount,
	}
}

type userReactionService struct {
	DB                *sql.DB
	userReactionTable string
	id                string
	author            string
	reaction          string
	prefix            string
	suffix            string
	userinfoTable     string
	infoId            string
	reactionCount     string
}

func (s *userReactionService) CheckReaction(ctx context.Context, id string, author string) (int64, error) {
	stmt := fmt.Sprintf(`select reaction from %s where %s = $1 and %s = $2`, s.userReactionTable, s.id, s.author)
	rows, err := s.DB.QueryContext(ctx, stmt, id, author)
	if err != nil {
		return -1, err
	}
	defer rows.Close()
	for rows.Next() {
		var reaction int64
		rows.Scan(&reaction)
		return reaction, nil
	}
	return -1, nil
}

func (s *userReactionService) React(ctx context.Context, id string, author string, reaction string) (int64, error) {
	userReaction, err := s.CheckReaction(ctx, id, author)
	if err != nil {
		return -1, err
	}
	if userReaction <= 0 {
		obj := map[string]int{
			"l1": 0,
			"l2": 0,
			"l3": 0,
		}
		obj["l"+reaction] = 1
		query := fmt.Sprintf("insert into %s(%s, %s, %s) values ($1, $2, $3)", s.userReactionTable, s.id, s.author, s.reaction)
		query2 := fmt.Sprintf("insert into %s(%s, %s1%s, %s2%s,%s3%s, %s) values ($1, $2, $3, $4, 1) on conflict (%s) do update set %s1%s = %s.%s1%s + $2, %s2%s = %s.%s2%s + $3, %s3%s =%s.%s3%s + $4, %s = %s.%s + 1",
			s.userinfoTable, s.infoId, s.prefix, s.suffix, s.prefix, s.suffix, s.prefix, s.suffix, s.reactionCount,
			s.infoId, s.prefix, s.suffix, s.userinfoTable, s.prefix, s.suffix, s.prefix, s.suffix, s.userinfoTable, s.prefix, s.suffix, s.prefix, s.suffix, s.userinfoTable, s.prefix, s.suffix, s.reactionCount, s.userinfoTable, s.reactionCount,
		)
		var rowCount int64
		rowCount = 0
		tx, err := s.DB.BeginTx(ctx, nil)
		if err != nil {
			return -1, err
		}
		defer tx.Rollback()
		res, err := tx.ExecContext(ctx, query, id, author, reaction)
		if err != nil {
			return -1, err
		}

		rowaffected, err := res.RowsAffected()
		if err != nil {
			return -1, err
		}

		rowCount += rowaffected

		res2, err := tx.ExecContext(ctx, query2, id, obj["l1"], obj["l2"], obj["l3"])
		if err != nil {
			return -1, err
		}

		rowaffected2, err := res2.RowsAffected()
		if err != nil {
			return -1, err
		}

		rowCount += rowaffected2

		if err = tx.Commit(); err != nil {
			return -1, err
		}
		return rowCount, nil
	} else {
		if _reaction, _ := strconv.Atoi(reaction); userReaction != int64(_reaction) {
			query := fmt.Sprintf("Update %s set %s = $1 where %s = $2 and %s = $3",
				s.userReactionTable, s.reaction, s.id, s.author)
			query2 := fmt.Sprintf("Update %s set %s%d%s = %s%d%s - 1, %s%s%s = %s%s%s + 1 where %s = $1",
				s.userinfoTable, s.prefix, userReaction, s.suffix, s.prefix, userReaction, s.suffix,
				s.prefix, reaction, s.suffix, s.prefix, reaction, s.suffix, s.infoId)
			var rowCount int64
			rowCount = 0
			tx, err := s.DB.BeginTx(ctx, nil)
			if err != nil {
				return -1, err
			}
			defer tx.Rollback()
			res, err := tx.ExecContext(ctx, query, reaction, id, author)
			if err != nil {
				return -1, err
			}

			rowaffected, err := res.RowsAffected()
			if err != nil {
				return -1, err
			}

			rowCount += rowaffected

			res2, err := tx.ExecContext(ctx, query2, id)
			if err != nil {
				return -1, err
			}

			rowaffected2, err := res2.RowsAffected()
			if err != nil {
				return -1, err
			}

			rowCount += rowaffected2

			if err = tx.Commit(); err != nil {
				return -1, err
			}
			return rowCount, nil
		}
		return 0, nil
	}
}

func (s *userReactionService) Unreact(ctx context.Context, id string, author string, reaction string) (int64, error) {
	query := fmt.Sprintf("DELETE from %s WHERE %s = $1 and %s = $2 and %s = $3", s.userReactionTable, s.id, s.author, s.reaction)
	query2 := fmt.Sprintf("UPDATE %s SET %s%s%s = %s%s%s - 1, %s = %s - 1 WHERE %s = $1 ", s.userinfoTable, s.prefix, reaction, s.suffix, s.prefix, reaction, s.suffix, s.reactionCount, s.reactionCount, s.infoId)

	var rowCount int64
	rowCount = 0
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, query, id, author, reaction)
	if err != nil {
		return -1, err
	}

	rowaffected, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	if rowaffected == 0 {
		return rowCount, nil
	}
	rowCount += rowaffected

	res2, err := tx.ExecContext(ctx, query2, id)
	if err != nil {
		return -1, err
	}

	rowaffected2, err := res2.RowsAffected()
	if err != nil {
		return -1, err
	}

	rowCount += rowaffected2

	if err = tx.Commit(); err != nil {
		return -1, err
	}
	return rowCount, nil
}

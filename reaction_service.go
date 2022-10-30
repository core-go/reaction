package reaction

import (
	"context"
	"database/sql"
	"fmt"
)

type ReactionService interface {
	Insert(ctx context.Context, reaction *Reaction) (int64, error)
	Delete(ctx context.Context, reaction *Reaction) (int64, error)
}

func NewReactionService(
	db *sql.DB,
	table string,
	id string,
	author string,
	userId string,
	time string,
	reaction string,
	rateTable string,
	rateId string,
	rateAuthor string,
	usefulCount string,
) ReactionService {
	return &reactionService{
		DB:          db,
		Table:       table,
		Id:          id,
		Author:      author,
		UserId:      userId,
		Time:        time,
		Reaction:    reaction,
		RateTable:   rateTable,
		RateId:      rateId,
		RateAuthor:  rateAuthor,
		UsefulCount: usefulCount,
	}
}

type reactionService struct {
	DB          *sql.DB
	Table       string
	Id          string
	Author      string
	UserId      string
	Time        string
	Reaction    string
	RateTable   string
	RateId      string
	RateAuthor  string
	UsefulCount string
}

func (s *reactionService) Insert(ctx context.Context, reaction *Reaction) (int64, error) {
	query1 := fmt.Sprintf("insert into %s(%s , %s, %s, %s,  %s) values ($1, $2, $3, $4, $5)",
		s.Table, s.Id, s.Author, s.UserId, s.Time, s.Reaction)
	fmt.Println(query1)
	stmt1, err := s.DB.Prepare(query1)
	if err != nil {
		return -1, err
	}
	res1, err := stmt1.ExecContext(ctx, reaction.Id, reaction.Author, reaction.UserId, reaction.Time, reaction.Type)
	if err != nil {
		return -1, err
	}

	query2 := fmt.Sprintf("update %s set %s = %s.%s + 1 where %s = $1 and %s = $2;",
		s.RateTable, s.UsefulCount, s.RateTable, s.UsefulCount, s.RateId, s.RateAuthor)
	fmt.Println(query2)
	stmt2, err := s.DB.Prepare(query2)
	if err != nil {
		return -1, err
	}
	stmt2.ExecContext(ctx, reaction.Id, reaction.Author)

	return res1.RowsAffected()
}

func (s *reactionService) Delete(ctx context.Context, reaction *Reaction) (int64, error) {
	query1 := fmt.Sprintf("delete from %s where %s = $1 and %s = $2 and %s = $3;",
		s.Table, s.Id, s.Author, s.UserId)
	fmt.Println(query1)
	stmt1, err := s.DB.Prepare(query1)
	if err != nil {
		return -1, err
	}
	res1, err := stmt1.ExecContext(ctx, reaction.Id, reaction.Author, reaction.UserId)
	if err != nil {
		return -1, err
	}

	query2 := fmt.Sprintf("update %s set %s = %s.%s - 1 where %s = $1 and %s = $2;",
		s.RateTable, s.UsefulCount, s.RateTable, s.UsefulCount, s.RateId, s.RateAuthor)
	fmt.Println(query2)
	stmt2, err := s.DB.Prepare(query2)
	if err != nil {
		return -1, err
	}
	stmt2.ExecContext(ctx, reaction.Id, reaction.Author)

	return res1.RowsAffected()
}

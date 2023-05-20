package rate

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type RateService interface {
	Load(ctx context.Context, id string, author string) (*Rate, error)
	Rate(ctx context.Context, id string, author string, req Request) (int64, error)
}

func NewRateService(
	db *sql.DB,
	rateTable string,
	idCol string,
	authorCol string,
	anonymousCol string,
	rateCol string,
	reviewCol string,
	timeCol string,
	usefulCountCol string,
	replyCountCol string,
	infoTable string,
	infoIdCol string,
	infoRateCol string,
	rateCountCol string,
	rateScoreCol string,
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	},
) RateService {
	return &rateService{
		DB:             db,
		RateTable:      rateTable,
		IdCol:          idCol,
		AuthorCol:      authorCol,
		AnonymousCol:   anonymousCol,
		RateCol:        rateCol,
		ReviewCol:      reviewCol,
		TimeCol:        timeCol,
		UsefulCountCol: usefulCountCol,
		ReplyCountCol:  replyCountCol,
		InfoTable:      infoTable,
		InfoIdCol:      infoIdCol,
		InfoRateCol:    infoRateCol,
		RateCountCol:   rateCountCol,
		RateScoreCol:   rateScoreCol,
		ToArray:        toArray,
	}
}

type rateService struct {
	DB           *sql.DB
	RateTable    string
	IdCol        string
	AuthorCol    string
	AnonymousCol string
	RateCol      string
	ReviewCol    string
	string
	TimeCol        string
	UsefulCountCol string
	ReplyCountCol  string
	InfoTable      string
	InfoIdCol      string
	InfoRateCol    string
	RateCountCol   string
	RateScoreCol   string
	ToArray        func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (s *rateService) Load(ctx context.Context, id string, author string) (*Rate, error) {
	query := fmt.Sprintf("select %s, %s, %s, %s, %s, %s, %s, histories from %s where %s = $1 and %s = $2 limit 1",
		s.IdCol, s.AuthorCol, s.RateCol, s.ReviewCol, s.TimeCol, s.UsefulCountCol, s.ReplyCountCol, s.RateTable, s.IdCol, s.AuthorCol)
	rows, err := s.DB.QueryContext(ctx, query, id, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rate Rate
		err = rows.Scan(&rate.Id, &rate.Author, &rate.Rate, &rate.Review, &rate.Time, &rate.UsefulCount, &rate.ReplyCount, s.ToArray(&rate.Histories))
		if err != nil {
			return nil, err
		}
		return &rate, nil
	}
	return nil, nil
}

func (s *rateService) Rate(ctx context.Context, id string, author string, req Request) (int64, error) {
	var rate = Rate{Id: id, Author: author, Review: req.Review, Rate: req.Rate, Anonymous: req.Anonymous}
	oldRate, _ := s.Load(ctx, rate.Id, rate.Author)
	query1 := fmt.Sprintf("insert into %s(%s, %s, %s%d, %s, %s) values ($1, %d, 1, 1, %d) on conflict (%s) do update set ",
		s.InfoTable, s.InfoIdCol, s.InfoRateCol, s.InfoRateCol, rate.Rate, s.RateCountCol, s.RateScoreCol, rate.Rate, rate.Rate, s.InfoIdCol)
	if oldRate != nil {
		if oldRate.Rate != rate.Rate {
			query1 += fmt.Sprintf(
				"%s%d = %s.%s%d - 1, %s%d = %s.%s%d + 1, %s = %s.%s + %d - %d, %s = (%s.%s + %d - %d) / %s.%s",
				s.InfoRateCol, oldRate.Rate, s.InfoTable, s.InfoRateCol, oldRate.Rate,
				s.InfoRateCol, rate.Rate, s.InfoTable, s.InfoRateCol, rate.Rate,
				s.RateScoreCol, s.InfoTable, s.RateScoreCol, rate.Rate, oldRate.Rate,
				s.InfoRateCol, s.InfoTable, s.RateScoreCol, rate.Rate, oldRate.Rate, s.InfoTable, s.RateCountCol)
		} else if oldRate.Rate == rate.Rate && oldRate.Review != rate.Review {
			query1 = ""
		} else {
			return 0, nil
		}
		rate.Histories = append(oldRate.Histories, Histories{Time: oldRate.Time, Rate: oldRate.Rate, Review: oldRate.Review})
	} else {
		query1 += fmt.Sprintf(
			"%s = %s.%s + 1, %s%d = %s.%s%d + 1, %s = %s.%s + %d, %s = (%s.%s + %d) / (%s.%s + 1)",
			s.RateCountCol, s.InfoTable, s.RateCountCol,
			s.InfoRateCol, rate.Rate, s.InfoTable, s.InfoRateCol, rate.Rate,
			s.RateScoreCol, s.InfoTable, s.RateScoreCol, rate.Rate,
			s.InfoRateCol, s.InfoTable, s.RateScoreCol, rate.Rate, s.InfoTable, s.RateCountCol)
	}

	stmt1, err := s.DB.Prepare(query1)
	if err != nil {
		return -1, err
	}
	stmt1.ExecContext(ctx, rate.Id)

	query2 := fmt.Sprintf(
		"insert into %s(%s, %s, %s, %s, %s, %s, histories) values ($1, $2, $3, $4, $5, $6, $7) on conflict (%s, %s) do update set %s = $3,  %s = $4, %s = $5, %s = $6, histories = $7",
		s.RateTable, s.IdCol, s.AuthorCol, s.AnonymousCol, s.RateCol, s.ReviewCol, s.TimeCol, s.IdCol, s.AuthorCol, s.AnonymousCol, s.RateCol, s.ReviewCol, s.TimeCol)
	stmt, err := s.DB.Prepare(query2)
	if err != nil {
		return -1, err
	}
	res2, err := stmt.ExecContext(ctx, rate.Id, rate.Author, rate.Anonymous, rate.Rate, rate.Review, rate.Time, s.ToArray(rate.Histories))
	if err != nil {
		return -1, err
	}

	return res2.RowsAffected()
}

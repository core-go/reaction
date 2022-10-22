package response

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type ResponseService interface {
	Load(ctx context.Context, id string, author string) (*Response, error)
	Response(ctx context.Context, response *Response) (int64, error)
}

func NewResponseService(
	db *sql.DB,
	responseTable string,
	idCol string,
	authorCol string,
	descriptionCol string,
	timeCol string,
	usefulCountCol string,
	commentCountCol string,
	infoTable string,
	infoIdCol string,
	responseCountCol string,
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	},
) ResponseService {
	return &responseService{
		DB:               db,
		ResponseTable:    responseTable,
		IdCol:            idCol,
		AuthorCol:        authorCol,
		DescriptionCol:   descriptionCol,
		TimeCol:          timeCol,
		UsefulCountCol:   usefulCountCol,
		CommentCountCol:  commentCountCol,
		InfoTable:        infoTable,
		InfoIdCol:        infoIdCol,
		ResponseCountCol: responseCountCol,
		ToArray:          toArray,
	}
}

type responseService struct {
	DB               *sql.DB
	ResponseTable    string
	IdCol            string
	AuthorCol        string
	DescriptionCol   string
	TimeCol          string
	UsefulCountCol   string
	CommentCountCol  string
	InfoTable        string
	InfoIdCol        string
	ResponseCountCol string
	ToArray          func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (s *responseService) Load(ctx context.Context, id string, author string) (*Response, error) {
	query := fmt.Sprintf("select %s, %s, %s, %s, %s, %s, histories from %s where %s = $1 and %s = $2 limit 1",
		s.IdCol, s.AuthorCol, s.DescriptionCol, s.TimeCol, s.UsefulCountCol, s.CommentCountCol, s.ResponseTable, s.IdCol, s.AuthorCol)
	fmt.Println(query)
	rows, err := s.DB.QueryContext(ctx, query, id, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var response Response
		err = rows.Scan(&response.Id, &response.Author, &response.Description, &response.Time, &response.UsefulCount, &response.CommentCount, s.ToArray(&response.Histories))
		if err != nil {
			return nil, err
		}
		return &response, nil
	}
	return nil, nil
}

func (s *responseService) Response(ctx context.Context, response *Response) (int64, error) {
	oldResponse, _ := s.Load(ctx, response.Id, response.Author)
	if oldResponse != nil {
		if oldResponse.Description == response.Description {
			return 0, nil
		}
		response.Histories = append(oldResponse.Histories, Histories{Time: oldResponse.Time, Description: oldResponse.Description})
	} else {
		query1 := fmt.Sprintf(
			"insert into %s(%s, %s) values ($1, 1) on conflict (%s) do update set %s = %s.%s + 1",
			s.InfoTable, s.InfoIdCol, s.ResponseCountCol, s.InfoIdCol, s.ResponseCountCol, s.InfoTable, s.ResponseCountCol)
		fmt.Println(query1)
		stmt1, err := s.DB.Prepare(query1)
		if err != nil {
			return -1, err
		}
		stmt1.ExecContext(ctx, response.Id)
	}

	query2 := fmt.Sprintf(
		"insert into %s(%s, %s, %s, %s, histories) values ($1, $2, $3, $4, $5) on conflict (%s, %s) do update set %s = $3, %s = $4, histories = $5",
		s.ResponseTable, s.IdCol, s.AuthorCol, s.DescriptionCol, s.TimeCol, s.IdCol, s.AuthorCol, s.DescriptionCol, s.TimeCol)
	fmt.Println(query2)
	stmt, err := s.DB.Prepare(query2)
	if err != nil {
		return -1, err
	}
	res2, err := stmt.ExecContext(ctx, response.Id, response.Author, response.Description, response.Time, s.ToArray(response.Histories))
	if err != nil {
		return -1, err
	}

	return res2.RowsAffected()
}

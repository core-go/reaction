package rates

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type RatesService interface {
	Rate(ctx context.Context, id string, author string, rate *Request) (int64, error)
}

func NewRatesService(
	db *sql.DB,
	max int,
	tableName string,
	idCol string,
	rateCol string,
	ratesCol string,
	reviewCol string,
	authorCol string,
	anonymousCol string,
	timeCol string,
	usefulCol string,
	replyCol string,

	fullInfoTableName string,
	fullInfoIdCol string,
	fullInfoScoreCol string,
	fullInfoCountCol string,
	fullInfoRateCol string,

	infoTablesName []string,
	infoIdCol string,
	infoRateCol string,
	infoCountCol string,
	infoScoreCol string,
	ToArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	},
) RatesService {
	return &ratesService{
		DB:           db,
		Max:          max,
		TableName:    tableName,
		IdCol:        idCol,
		RateCol:      rateCol,
		RatesCol:     ratesCol,
		ReviewCol:    reviewCol,
		UsefulCol:    usefulCol,
		ReplyCol:     replyCol,
		AuthorCol:    authorCol,
		AnonymousCol: anonymousCol,
		TimeCol:      timeCol,

		FullInfoTableName: fullInfoTableName,
		FullInfoIdCol:     fullInfoIdCol,
		FullInfoRateCol:   fullInfoRateCol,
		FullCountCol:      fullInfoCountCol,
		FullScoreCol:      fullInfoScoreCol,

		InfoTablesName: infoTablesName,
		InfoIdCol:      infoIdCol,
		InfoRateCol:    infoRateCol,
		InfoCountCol:   infoCountCol,
		InfoScoreCol:   infoScoreCol,
		ToArray:        ToArray,
	}
}

type ratesService struct {
	DB  *sql.DB
	Max int

	TableName    string
	IdCol        string
	RateCol      string
	RatesCol     string
	ReviewCol    string
	UsefulCol    string
	ReplyCol     string
	AuthorCol    string
	AnonymousCol string
	TimeCol      string

	FullInfoTableName string
	FullInfoIdCol     string
	FullInfoRateCol   string
	FullCountCol      string
	FullScoreCol      string

	InfoTablesName []string
	InfoIdCol      string
	InfoRateCol    string
	InfoCountCol   string
	InfoScoreCol   string
	ToArray        func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (s *ratesService) Rate(ctx context.Context, id string, author string, req *Request) (int64, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	defer tx.Rollback()

	t := time.Now()
	rate := Rates{Id: id, Author: author, Rate: req.Rate, Rates: req.Rates, Review: req.Review, Anonymous: req.Anonymous, Time: &t}
	if req.Rates != nil && len(req.Rates) > 0 {
		rate.Rate = avg(req.Rates)
	}
	// load rates
	oldRate, _ := s.load(ctx, rate.Id, rate.Author)
	existRate := oldRate != nil
	//  loop all rate and then upsert info table
	_, err = s.upsertInfoTables(ctx, tx, oldRate, rate, s.Max, s.InfoTablesName)
	if err != nil {
		return -1, err
	}
	// upsert full info table
	_, err = s.upsertFullInfoTable(ctx, tx, oldRate, rate, s.Max, s.FullInfoTableName, s.InfoTablesName, existRate)
	if err != nil {
		return -1, err
	}
	// upsert table rate
	queryRate := fmt.Sprintf(
		"insert into %s(%s, %s, %s, %s, %s, %s, %s, histories) values ($1, $2, $3, $4, $5, $6, $7, $8) on conflict (%s, %s) do update set %s = $3,  %s = $4, %s = $5, %s = $6, %s = $7, histories = $8",
		s.TableName, s.IdCol, s.AuthorCol, s.AnonymousCol, s.RateCol, s.RatesCol, s.ReviewCol, s.TimeCol, s.IdCol, s.AuthorCol, s.AnonymousCol, s.RateCol, s.RatesCol, s.ReviewCol, s.TimeCol)
	fmt.Println(queryRate)
	stmt, err := tx.Prepare(queryRate)
	if err != nil {
		return -1, err
	}
	res2, err := stmt.ExecContext(ctx, rate.Id, rate.Author, rate.Anonymous, rate.Rate, s.ToArray(rate.Rates), rate.Review, rate.Time, s.ToArray(rate.Histories))
	if err != nil {
		return -1, err
	}

	r, err1 := res2.RowsAffected()
	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return r, err1
}
func (s *ratesService) load(ctx context.Context, id string, author string) (*Rates, error) {
	query := fmt.Sprintf("select %s,%s,%s,%s,%s,%s,%s,%s,histories from %s where %s = $1 and %s = $2",
		s.IdCol, s.AuthorCol, s.RateCol, s.RatesCol, s.TimeCol, s.ReviewCol, s.UsefulCol, s.ReplyCol,
		s.TableName, s.IdCol, s.AuthorCol)
	rows, err := s.DB.QueryContext(ctx, query, id, author)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rate Rates
		err := rows.Scan(&rate.Id, &rate.Author, &rate.Rate, pq.Array(&rate.Rates), &rate.Time, &rate.Review, &rate.UsefulCount, &rate.ReplyCount, pq.Array(&rate.Histories))
		if err != nil {
			return nil, err
		}
		return &rate, nil
	}
	return nil, nil
}
func (s *ratesService) upsertInfoTables(ctx context.Context, tx *sql.Tx, oldRate *Rates, rate Rates, max int, infoTablesNames []string) (int64, error) {
	queries := make([]string, 0)
	params := make([][]interface{}, 0)
	for i := 1; i <= max; i++ {
		index := i - 1
		rateValue := int(rate.Rates[index])
		infoTable := infoTablesNames[index]

		query1 := fmt.Sprintf("insert into %s(%s, %s, %s%d, %s, %s) values ($1, %d, 1, 1, %d) on conflict (%s) do update set ",
			infoTable, s.InfoIdCol, s.InfoRateCol, s.InfoRateCol, rateValue, s.InfoCountCol, s.InfoScoreCol, rateValue, rateValue, s.InfoIdCol)
		if oldRate != nil {
			if oldRate.Rates[index] != rate.Rates[index] {
				oRate := int(oldRate.Rates[index])
				query1 += fmt.Sprintf(
					"%s%d = %s.%s%d - 1, %s%d = %s.%s%d + 1, %s = %s.%s + %d - %d, %s = (%s.%s + %d - %d) / %s.%s",
					s.InfoRateCol, oRate, infoTable, s.InfoRateCol, oRate,
					s.InfoRateCol, rateValue, infoTable, s.InfoRateCol, rateValue,
					s.InfoScoreCol, infoTable, s.InfoScoreCol, rateValue, oRate,
					s.InfoRateCol, infoTable, s.InfoScoreCol, rateValue, oRate, infoTable, s.InfoCountCol)
			} else if oldRate.Rates[index] != rate.Rates[index] && oldRate.Review != rate.Review {
				query1 = ""
			} else {
				break
			}
			rate.Histories = append(oldRate.Histories, Histories{Time: oldRate.Time, Rate: oldRate.Rate, Review: oldRate.Review})
		} else {
			query1 += fmt.Sprintf(
				"%s = %s.%s + 1, %s%d = %s.%s%d + 1, %s = %s.%s + %d, %s = (%s.%s + %d) / (%s.%s + 1)",
				s.InfoCountCol, infoTable, s.InfoCountCol,
				s.InfoRateCol, rateValue, infoTable, s.InfoRateCol, rateValue,
				s.InfoScoreCol, infoTable, s.InfoScoreCol, rateValue,
				s.InfoRateCol, infoTable, s.InfoScoreCol, rateValue, infoTable, s.InfoCountCol)
		}
		fmt.Println(query1)
		queries = append(queries, query1)
		params = append(params, []interface{}{rate.Id})
	}
	fmt.Println(queries)
	return s.ExecBatch(ctx, queries, params)
}
func (s *ratesService) upsertFullInfoTable(ctx context.Context, tx *sql.Tx, oldRate *Rates, rate Rates, max int, fullInfoTableName string, infoTablesName []string, existRate bool) (int64, error) {
	updatedRateRangeAverage := sum(rate.Rates)
	countOfUserRatedMore := 1 // default is new rate
	if existRate {
		updatedRateRangeAverage -= sum(oldRate.Rates)
		countOfUserRatedMore = 0
	}
	updatedRateRangeAverage = updatedRateRangeAverage / float32(max)
	queryi, paramsi := s.buildQueryInsertFullInfo(&rate, fullInfoTableName, infoTablesName)
	nextIndex := len(paramsi) + 1
	queryu, paramsu := s.buildQueryUpdateFullInfo(&rate, updatedRateRangeAverage, countOfUserRatedMore, fullInfoTableName, infoTablesName, nextIndex)
	queryMerged := fmt.Sprintf("%s  on conflict (%s) do %s", queryi, s.FullInfoIdCol, queryu)
	stmt, err := tx.Prepare(queryMerged)
	if err != nil {
		return -1, err
	}
	pm := append(paramsi, paramsu...)
	fmt.Println(queryMerged, pm)
	rs, err := stmt.ExecContext(ctx, pm...)
	if err != nil {
		return -1, err
	}

	return rs.RowsAffected()
}
func (s *ratesService) buildQueryInsertFullInfo(rate *Rates, fullInfoTableName string, infoTablesName []string) (string, []interface{}) {
	rateCols := []string{}
	params := []interface{}{rate.Id, rate.Rate, rate.Rate, rate.Id}
	rateQuerys := []string{}
	for i := 1; i <= len(infoTablesName); i++ {
		rateCols = append(rateCols, fmt.Sprintf("%s%d", s.RateCol, i))
		rateQuerys = append(rateQuerys, fmt.Sprintf("(select avg(%s) from %s where %s = $4 group by %s)",
			s.InfoRateCol, infoTablesName[i-1], s.InfoIdCol, s.InfoIdCol))
	}
	query := fmt.Sprintf("insert into %s(%s, %s, %s, %s, %s)values ($1, $2, 1, $3, %s)",
		fullInfoTableName, s.FullInfoIdCol, s.FullInfoRateCol, s.FullCountCol, s.FullScoreCol,
		strings.Join(rateCols, ", "), strings.Join(rateQuerys, ","))
	return query, params
}
func (s *ratesService) buildQueryUpdateFullInfo(rate *Rates, score float32, count int, fullInfoTableName string, infoTablesName []string, index int) (string, []interface{}) {
	if len(infoTablesName) > 0 {
		ss := []string{}
		params := []interface{}{}
		for i := 1; i <= len(infoTablesName); i++ {
			ss = append(ss, fmt.Sprintf("%s%d = (select avg(%s) from %s where %s = $%d group by %s)", s.InfoRateCol, i, s.InfoRateCol, infoTablesName[i-1],
				s.InfoIdCol, index+1, s.InfoIdCol))
		}

		sql := fmt.Sprintf("update set %[1]s = (%[2]s.%[3]s + $%[4]d)/(%[2]s.%[5]s + %[11]d), %[6]s = %[2]s.%[7]s + $%[4]d, %[8]s = %[2]s.%[9]s + %[11]d, %[10]s",
			s.FullInfoRateCol, fullInfoTableName, s.FullScoreCol, index, s.FullCountCol, s.FullScoreCol,
			s.FullScoreCol, s.FullCountCol, s.FullCountCol, strings.Join(ss, ", "), count)
		params = append(params, score, rate.Id)
		return sql, params
	}
	return "", nil
}
func (s *ratesService) ExecBatch(ctx context.Context, stmts []string, params [][]interface{}) (int64, error) {
	var rowResult int64
	rowResult = 0
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	defer tx.Rollback()

	for index, _ := range stmts {
		fmt.Println(stmts[index])
		fmt.Println(params[index]...)
		stmt, err := tx.PrepareContext(ctx, stmts[index])
		if err != nil {
			return -1, err
		}

		res, err := stmt.ExecContext(ctx, params[index]...)
		if err != nil {
			return -1, err
		}
		rowAffected, err := res.RowsAffected()
		rowResult += rowAffected
	}
	if err = tx.Commit(); err != nil {
		return -1, err
	}
	return rowResult, nil
}
func sum(arr []float32) float32 {
	total := float32(0)
	for _, val := range arr {
		total += val
	}
	return total
}
func avg(numbers []float32) float32 {
	var res float32
	res = 0
	for i := 0; i < len(numbers); i++ {
		res += numbers[i]
	}
	return res / float32(len(numbers))
}

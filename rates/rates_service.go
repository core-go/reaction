package rates

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

type RatesService interface {
	Rate(ctx context.Context, rate *Rates) (int64, error)
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
) RatesService {
	return &ratesService{
		DB:        db,
		Max:       max,
		TableName: tableName,
		IdCol:     idCol,
		RateCol:   rateCol,
		RatesCol:  ratesCol,
		ReviewCol: reviewCol,
		UsefulCol: usefulCol,
		ReplyCol:  replyCol,
		AuthorCol: authorCol,
		TimeCol:   timeCol,

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
	}
}

type ratesService struct {
	DB  *sql.DB
	Max int

	TableName string
	IdCol     string
	RateCol   string
	RatesCol  string
	ReviewCol string
	UsefulCol string
	ReplyCol  string
	AuthorCol string
	TimeCol   string

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
}

func (s *ratesService) Rate(ctx context.Context, rate *Rates) (int64, error) {
	// check exist
	var info int
	query := fmt.Sprintf("select count(*) from %s where %s = $1", s.FullInfoTableName, s.FullInfoIdCol)
	err := s.DB.QueryRow(query, rate.Id).Scan(&info)
	if err != nil {
		return -1, err
	}
	if rate.Rates != nil && len(rate.Rates) > 0 {
		rate.Rate = avg(rate.Rates)
	}

	*rate.Time = time.Now()

	// if info doesn't exist, insert rate table
	if info == 0 {
		res, err := s.insert(ctx, rate, s.FullInfoTableName, true)
		if err != nil {
			return -1, err
		}
		return res, nil
	}
	// if info exist, get ratecriteria from ratetable
	exists, err := s.load(ctx, rate.Id, rate.Author)
	if err != nil {
		return -1, err
	}
	if exists == nil {
		res, err := s.insert(ctx, rate, s.FullInfoTableName, false)
		if err != nil {
			return -1, err
		}
		return res, nil
	}
	if len(exists.Histories) > 0 {
		history := exists.Histories
		history = append(history, Histories{Time: exists.Time, Rate: exists.Rate, Review: exists.Review})
		rate.Histories = history
	} else {
		rate.Histories = append(rate.Histories, Histories{Time: exists.Time, Rate: exists.Rate, Review: exists.Review})
	}
	// update
	res, err := s.update(ctx, *rate, exists.Rate)
	if err != nil {
		return -1, err
	}
	return res, nil
	// return 0, nil
}

func avg(numbers []float32) float32 {
	var res float32
	res = 0
	for i := 0; i < len(numbers); i++ {
		res += numbers[i]
	}
	return res / float32(len(numbers))
}
func (s *ratesService) insert(ctx context.Context, rate *Rates, fullInfoTableName string, isNewInfo bool) (int64, error) {

	stmts := []string{}
	params := [][]interface{}{}

	mainStmt := fmt.Sprintf("insert into %s(%s, %s, %s, %s, %s, %s, histories) values ($1, $2, $3, $4, $5, $6, $7)", s.TableName, s.IdCol, s.AuthorCol, s.RateCol, s.RatesCol, s.ReviewCol, s.TimeCol)
	mainParam := []interface{}{
		rate.Id, rate.Author, rate.Rate, pq.Array(rate.Rates), rate.Review, rate.Time, pq.Array(rate.Histories),
	}

	if isNewInfo == true {
		// insert info tables
		stmts, params = s.insertInfo(rate, s.InfoTablesName, stmts, params)

		// insert ratefullinfo table
		stmt3, ratefulParams := s.insertFullInfo(rate, s.FullInfoTableName, s.InfoTablesName)

		params = append(params, ratefulParams)
		stmts = append(stmts, stmt3)
		stmts = append(stmts, mainStmt)
		params = append(params, mainParam)

		// RETURN
		res, err := s.ExecBatch(ctx, stmts, params)
		if err != nil {
			return res, err
		}
		return res, nil
		/////
	} else {
		fullstmt, fullparam := s.updateFullInfo(rate, s.FullInfoTableName, s.InfoTablesName)
		stmts = append(stmts, fullstmt)
		params = append(params, fullparam)
		for i := 0; i < len(rate.Rates); i++ {
			stmt, param := s.updateNewInfo(rate.Id, rate.Rates[i], s.InfoTablesName[i])
			stmts = append(stmts, stmt)
			params = append(params, param)
		}
		stmts = append(stmts, mainStmt)
		params = append(params, mainParam)
		// RETURN
		res, err := s.ExecBatch(ctx, stmts, params)
		if err != nil {
			return res, err
		}
		return res, nil
		/////
	}
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
func (s *ratesService) insertFullInfo(rate *Rates, fullInfoTableName string, infoTablesName []string) (string, []interface{}) {
	rateCols := []string{}
	ratefulParams := []interface{}{rate.Id, rate.Rate, rate.Rate}
	rateQuerys := []string{}
	for i := 1; i <= len(infoTablesName); i++ {
		rateCols = append(rateCols, fmt.Sprintf("%s%d", s.RateCol, i))
		rateQuerys = append(rateQuerys, fmt.Sprintf("(select avg(%s) from %s where %s = $4 group by %s)", s.InfoRateCol, infoTablesName[i-1], s.InfoIdCol, s.InfoIdCol))

	}
	stmt3 := fmt.Sprintf("insert into %s(%s, %s, %s, %s, %s)values ($1, $2, 1, $3, %s)", fullInfoTableName, s.FullInfoIdCol, s.FullInfoRateCol, s.FullCountCol, s.FullScoreCol, strings.Join(rateCols, ", "), strings.Join(rateQuerys, ","))
	ratefulParams = append(ratefulParams, rate.Id)
	fmt.Println(stmt3)
	return stmt3, ratefulParams
}

func (s *ratesService) insertInfo(rate *Rates, infoTableNames []string, stmts []string, params [][]interface{}) ([]string, [][]interface{}) {
	var infoParam []interface{}
	var rateCols []string
	var rateColParams []interface{}
	for i := 0; i < len(rate.Rates); i++ {
		infoParam = []interface{}{}
		rateCols = []string{}
		rateColParams = []interface{}{}
		for j := 1; j <= s.Max; j++ {
			rateCols = append(rateCols, fmt.Sprintf("%s%d", s.RateCol, j))
			if float32(j) == rate.Rates[i] {
				rateColParams = append(rateColParams, 1)
			} else {
				rateColParams = append(rateColParams, 0)
			}
		}

		rpr := []string{}
		for i := 1; i <= s.Max; i++ {
			rpr = append(rpr, fmt.Sprintf(`$%d`, i+3))
		}
		stmt2 := fmt.Sprintf("insert into %s(%s,%s,%s,%s,%s) values($1, $2, 1, $3, %s)",
			infoTableNames[i], s.InfoIdCol, s.InfoRateCol, s.InfoCountCol, s.InfoScoreCol, strings.Join(rateCols, ","), strings.Join(rpr, ", "))
		fmt.Println(stmt2)
		infoParam = append(infoParam, rate.Id, rate.Rates[i], rate.Rates[i])
		infoParam = append(infoParam, rateColParams...)

		stmts = append(stmts, stmt2)
		params = append(params, infoParam)
	}
	return stmts, params
}

func (s *ratesService) updateFullInfo(rate *Rates, fullInfoTableName string, infoTablesName []string) (string, []interface{}) {
	if len(infoTablesName) > 0 {
		ss := []string{}
		params := []interface{}{}
		// ps := []interface{}{}
		for i := 1; i <= len(infoTablesName); i++ {
			ss = append(ss, fmt.Sprintf("%s%d = (select avg(%s) from %s where %s = $2 group by %s)", s.InfoRateCol, i, s.InfoRateCol, infoTablesName[i-1], s.InfoIdCol, s.InfoIdCol))
		}
		// update companyfullinfo set rate = (3 + score)/(count+1), score = score + 3, count = count + 1, companyrateinfo1 =
		sql := fmt.Sprintf("update %s set %s = (%s + $1)/(%s + 1), %s = %s + $1, %s = %s + 1, %s",
			fullInfoTableName, s.FullInfoRateCol, s.FullScoreCol, s.FullCountCol, s.FullScoreCol, s.FullScoreCol, s.FullCountCol, s.FullCountCol, strings.Join(ss, ", "))
		params = append(params, rate.Rate, rate.Id)
		// params = append(params, ps)
		return sql, params
	} else {
		params := []interface{}{}
		sql := fmt.Sprintf("update %s set %s = (%s + $1)/(%s + 1), %s = %s + $1, %s = %s + 1",
			fullInfoTableName, s.FullInfoRateCol, s.FullScoreCol, s.FullCountCol, s.FullScoreCol, s.FullScoreCol, s.FullCountCol, s.FullCountCol)
		params = append(params, rate.Rate)
		return sql, params
	}
}

func (s *ratesService) updateNewInfo(id string, rate float32, tableInfoName string) (string, []interface{}) {
	sql := fmt.Sprintf("update %s set %s = (%s + $1)/(%s + 1), %s = %s + 1, %s = %s + $1",
		tableInfoName, s.InfoRateCol, s.InfoScoreCol, s.InfoCountCol, s.InfoCountCol, s.InfoCountCol, s.InfoScoreCol, s.InfoScoreCol)
	//, %s%d = %s%d + 1 where %s = $2",
	//, s.InfoRateCol, int32(rate.Rate), s.InfoRateCol, int32(rate.Rate), s.InfoIdCol
	for i := 1; i <= (s.Max); i++ {
		if rate == float32(i) {
			sql += fmt.Sprintf(", %s%d = %s%d + 1 ", s.InfoRateCol, int32(rate), s.InfoRateCol, int32(rate))
		}
	}
	sql += fmt.Sprintf(" where %s = $2", s.InfoIdCol)
	fmt.Println(sql)
	params := []interface{}{rate, id}
	return sql, params
}

func (s *ratesService) updateOldInfo(id string, newRate float32, oldRate float32, fullInfotableName string, infoTableNames []string) (string, []interface{}) {
	delta := newRate - oldRate
	ss := []string{}
	params := []interface{}{}
	var stmt string
	if len(infoTableNames) > 0 {
		for i := 1; i <= len(infoTableNames); i++ {
			sql := fmt.Sprintf("%s%d = (select avg(%s) from %s where %s = $2 group by %s)",
				s.FullInfoRateCol, i, s.InfoRateCol, infoTableNames[i-1], s.InfoIdCol, s.InfoIdCol)
			ss = append(ss, sql)
		}
		stmt = fmt.Sprintf("update %s set %s = ($1 + %s)/%s, %s = %s + $1, %s",
			s.FullInfoTableName, s.FullInfoRateCol, s.FullScoreCol, s.FullCountCol, s.FullScoreCol, s.FullScoreCol, strings.Join(ss, ", "))
		params = append(params, delta)
	} else {
		stmt = fmt.Sprintf("update %s set %s =($1 + %s)/%s , %s = %s + $1 ",
			s.FullInfoTableName, s.FullInfoRateCol, s.FullInfoRateCol, s.FullCountCol, s.FullScoreCol, s.FullScoreCol)
		params = append(params, delta)
	}
	params = append(params, id)
	return stmt, params
}

func (s *ratesService) update(ctx context.Context, newRate Rates, oldRate float32) (int64, error) {
	rates := newRate.Rates
	r := newRate.Rate
	stmts := []string{}
	params := [][]interface{}{}

	stmt1 := fmt.Sprintf("update %s set %s = $1, %s = $2, %s = $3, %s = $4, %s = $5, %s = $6, %s = $7, %s = $8, histories = $9",
		s.TableName, s.IdCol, s.AuthorCol, s.RateCol, s.RatesCol, s.TimeCol, s.ReviewCol, s.ReplyCol, s.UsefulCol)

	param1 := []interface{}{newRate.Id, newRate.Author, newRate.Rate, pq.Array(newRate.Rates), newRate.Time, newRate.Review, newRate.ReplyCount, newRate.UsefulCount, pq.Array(newRate.Histories)}

	if len(rates) > 0 && r > 0 {
		for i := 0; i < len(newRate.Rates); i++ {
			sql2, param2 := s.updateNewInfo(newRate.Id, rates[i], s.InfoTablesName[i])
			stmts = append(stmts, sql2)
			params = append(params, param2)
		}
		sql3, param3 := s.updateOldInfo(newRate.Id, newRate.Rate, oldRate, s.FullInfoTableName, s.InfoTablesName)
		stmts = append(stmts, sql3)
		params = append(params, param3)

		stmts = append(stmts, stmt1)
		params = append(params, param1)
		// RETURN EXACBATCH (stmts, true)
		res, err := s.ExecBatch(ctx, stmts, params)
		if err != nil {
			return -1, err
		}
		return res, nil
		///
	} else {
		stmts = append(stmts, stmt1)
		params = append(params, param1)
		// RETURN EXACBATCH (stmts, true)
		res, err := s.ExecBatch(ctx, stmts, params)
		if err != nil {
			return -1, err
		}
		return res, nil
		////////
	}
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
	// return 1, nil
	return rowResult, nil
}

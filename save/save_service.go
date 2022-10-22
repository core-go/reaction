package save

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type SaveService interface {
	Load(ctx context.Context, id string) ([]*Save, error)
	Save(ctx context.Context, id string, item string) (int64, error)
	Remove(ctx context.Context, id string, item string) (int64, error)
}

func NewSaveService(
	db *sql.DB,
	table string,
	idCol string,
	itemCol string,
	max int,
) SaveService {
	return &saveService{
		DB:      db,
		Table:   table,
		IdCol:   idCol,
		ItemCol: itemCol,
		Max:     max,
	}
}

type saveService struct {
	DB      *sql.DB
	Table   string
	IdCol   string
	ItemCol string
	Max     int
}

func (s *saveService) Load(ctx context.Context, id string) ([]*Save, error) {
	var saveList []*Save
	query := fmt.Sprintf("select %s as id, %s as item from %s where %s = $1", s.IdCol, s.ItemCol, s.Table, s.IdCol)
	fmt.Println(query)
	rows, err := s.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var save Save
		err = rows.Scan(&save.Id, pq.Array(&save.Item))
		if err != nil {
			return nil, err
		}
		saveList = append(saveList, &save)
	}
	return saveList, nil
}

func (s *saveService) Save(ctx context.Context, id string, item string) (int64, error) {
	var items []string
	query0 := fmt.Sprintf("select %s as items from %s where %s = $1", s.ItemCol, s.Table, s.IdCol)
	rows, err0 := s.DB.QueryContext(ctx, query0, id)
	if err0 != nil {
		return -1, err0
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(pq.Array(&items))
		if err != nil {
			return -1, err
		}

	}
	if len(items) == 0 {
		query := fmt.Sprintf("insert into %s(%s, %s) values ($1, $2)", s.Table, s.IdCol, s.ItemCol)
		stmt, er0 := s.DB.Prepare(query)
		if er0 != nil {
			return -1, nil
		}
		items = append(items, item)
		res, err := stmt.ExecContext(ctx, id, pq.Array(items))
		if err != nil {
			return -1, err
		}

		return res.RowsAffected()
	} else {
		query := fmt.Sprintf("update %s set %s = $1 where %s = $2", s.Table, s.ItemCol, s.IdCol)
		stmt, er0 := s.DB.Prepare(query)
		if er0 != nil {
			return -1, nil
		}
		for _, v := range items {
			if v == item {
				return -1, nil
			}
		}
		items = append(items, item)
		if len(items) > s.Max {
			items = items[1:]
		}
		res, err := stmt.ExecContext(ctx, pq.Array(items), id)
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	}
}

func (s *saveService) Remove(ctx context.Context, id string, item string) (int64, error) {
	query0 := fmt.Sprintf("select %s as id, %s as items from %s where %s = $1", s.IdCol, s.ItemCol, s.Table, s.IdCol)
	rows, err0 := s.DB.QueryContext(ctx, query0, id)
	if err0 != nil {
		return -1, err0
	}
	var items []Save
	defer rows.Close()
	for rows.Next() {
		var item Save
		err := rows.Scan(&item.Id, pq.Array(&item.Item))
		if err != nil {
			return -1, err
		}
		items = append(items, item)
	}
	if items == nil {
		return 0, nil
	}
	newItems := []string{}
	for i := 0; i < len(items[0].Item); i++ {
		if items[0].Item[i] != item {
			newItems = append(newItems, items[0].Item[i])
		}
	}
	query := fmt.Sprintf("update %s set %s = $1 where %s = $2", s.Table, s.ItemCol, s.IdCol)
	fmt.Println(query)
	stmt, er0 := s.DB.Prepare(query)
	if er0 != nil {
		return -1, er0
	}
	res, err := stmt.ExecContext(ctx, pq.Array(&newItems), id)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

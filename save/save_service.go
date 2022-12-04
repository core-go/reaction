package save

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
)

type SaveService interface {
	Load(ctx context.Context, id string, listResult interface{}) error
	Save(ctx context.Context, id string, item string) (int64, error)
	Remove(ctx context.Context, id string, item string) (int64, error)
}

func NewSaveService(
	db *sql.DB,
	modelType reflect.Type,
	table string,
	idCol string,
	itemCol string,
	max int,
	targetTable string,
	idTargetCol string,
	toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
},

) SaveService {
	return &saveService{
		DB:          db,
		table:       table,
		idCol:       idCol,
		itemCol:     itemCol,
		modelType:   modelType,
		max:         max,
		targetTable: targetTable,
		idTargetCol: idTargetCol,
		toArray:     toArray,
	}
}

type saveService struct {
	DB          *sql.DB
	table       string
	idCol       string
	itemCol     string
	max         int
	targetTable string
	idTargetCol string
	modelType   reflect.Type
	toArray     func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
}

func (s *saveService) Load(ctx context.Context, id string, listResult interface{}) error {
	var saveList []Items
	query := fmt.Sprintf("select %s as id, %s as items from %s where %s = $1", s.idCol, s.itemCol, s.table, s.idCol)
	rows, err := s.DB.QueryContext(ctx, query, id)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var save Items
		err = rows.Scan(&save.Id, s.toArray(&save.Items))
		if err != nil {
			return err
		}
		saveList = append(saveList, save)
	}
	if len(saveList) == 0 {
		return nil
	}
	if len(saveList[0].Items) == 0 {
		return nil
	}

	query2 := fmt.Sprintf("SELECT * FROM %s WHERE %s = ANY($1)", s.targetTable, s.idTargetCol)
	stmt, err0 := s.DB.PrepareContext(ctx, query2)
	if err0 != nil {
		return err0
	}
	rows, err1 := stmt.QueryContext(ctx, s.toArray(saveList[0].Items))
	if err1 != nil {
		return err1
	}

	slicePtrValue := reflect.ValueOf(listResult)
	sliceValue := reflect.Indirect(slicePtrValue)
	list, err := scan(rows, s.modelType, nil, s.toArray)
	if err != nil {
		return err
	}
	sliceValue.Set(reflect.ValueOf(list))
	return nil

}

func (s *saveService) Save(ctx context.Context, id string, item string) (int64, error) {
	var items []string
	query0 := fmt.Sprintf("select %s as items from %s where %s = $1", s.itemCol, s.table, s.idCol)
	rows, err0 := s.DB.QueryContext(ctx, query0, id)
	if err0 != nil {
		return -1, err0
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(s.toArray(&items))
		if err != nil {
			return -1, err
		}
	}

	if items == nil {
		query := fmt.Sprintf("insert into %s(%s, %s) values ($1, $2)", s.table, s.idCol, s.itemCol)
		stmt, er0 := s.DB.Prepare(query)
		if er0 != nil {
			return -1, er0
		}
		items = append(items, item)
		res, err := stmt.ExecContext(ctx, id, s.toArray(items))
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	} else {
		query := fmt.Sprintf("update %s set %s = $1 where %s = $2", s.table, s.itemCol, s.idCol)
		stmt, er0 := s.DB.Prepare(query)
		if er0 != nil {
			return -1, er0
		}
		for _, v := range items {
			if v == item {
				return -1, nil
			}
		}
		items = append(items, item)
		if len(items) > s.max {
			items = items[1:]
		}
		res, err := stmt.ExecContext(ctx, s.toArray(items), id)
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	}
}

func (s *saveService) Remove(ctx context.Context, id string, item string) (int64, error) {
	query0 := fmt.Sprintf("select %s as id, %s as items from %s where %s = $1", s.idCol, s.itemCol, s.table, s.idCol)
	rows, err0 := s.DB.QueryContext(ctx, query0, id)
	if err0 != nil {
		return -1, err0
	}
	var items []Items
	defer rows.Close()
	for rows.Next() {
		var item Items
		err := rows.Scan(&item.Id, s.toArray(&item.Items))
		if err != nil {
			return -1, err
		}
		items = append(items, item)
	}
	if items == nil {
		return 0, nil
	}
	newItems := []string{}
	for i := 0; i < len(items[0].Items); i++ {
		if items[0].Items[i] != item {
			newItems = append(newItems, items[0].Items[i])
		}
	}
	query := fmt.Sprintf("update %s set %s = $1 where %s = $2", s.table, s.itemCol, s.idCol)
	fmt.Println(query)
	stmt, er0 := s.DB.Prepare(query)
	if er0 != nil {
		return -1, er0
	}
	res, err := stmt.ExecContext(ctx, s.toArray(&newItems), id)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

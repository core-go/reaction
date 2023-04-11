package comment

import (
	"context"
	"database/sql"
	"fmt"
	q "github.com/core-go/sql"
	"github.com/lib/pq"
	"reflect"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	_ "unicode/utf8"
)

var collator = collate.New(language.Und)

type queryInfo struct {
	db          *sql.DB
	table       string
	url         string
	id          string
	name        string
	displayName string
}

var colsInfo, _ = q.GetColumnIndexes(reflect.TypeOf(Info{}))

func NewQueryInfo(db *sql.DB, table string, url string, id string, name string, displayName string) queryInfo {
	return queryInfo{db: db, table: table, url: url, id: id, name: name, displayName: displayName}
}

func (i queryInfo) Load(ids []string) ([]Info, error) {
	rs := make([]Info, 0)
	if len(ids) == 0 {
		return rs, nil
	}
	ids = distinct(ids)
	querysql := fmt.Sprintf(`select %s as id, %s as url, %s as name, %s as displayname from %s where %s = any(%s) and %s is not null order by %s`,
		i.id, i.url, i.name, i.displayName, i.table, i.id, q.BuildDollarParam(1), i.url, i.id)
	r := make([]Info, 0)
	err := q.Query(context.Background(), i.db, colsInfo, &r, querysql, pq.Array(ids))
	if err != nil {
		return rs, err
	}
	return r, nil
}

func ObjectAssign(dst interface{}, src interface{}) error {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		return fmt.Errorf("dst is not a pointer")
	}

	srcValue := reflect.ValueOf(src)
	if !srcValue.IsValid() {
		return fmt.Errorf("src is not valid")
	}

	dstElem := dstValue.Elem()
	if !dstElem.CanSet() {
		return fmt.Errorf("dst is not settable")
	}

	dstType := dstElem.Type()
	srcType := srcValue.Type()

	if !srcType.AssignableTo(dstType) {
		return fmt.Errorf("src type %s is not assignable to dst type %s", srcType, dstType)
	}

	dstElem.Set(srcValue)
	return nil
}

func BinarySearch(ar []Info, el string) int {
	m := 0
	n := len(ar) - 1

	for m <= n {
		k := (n + m) >> 1
		cmp := compare(el, ar[k].Id)
		if cmp > 0 {
			m = k + 1
		} else if cmp < 0 {
			n = k - 1
		} else {
			return k
		}
	}
	return -m - 1
}

func distinct(arr []string) []string {
	// Sort the input array
	sort.Strings(arr)
	// Create a new array to store distinct elements
	distinctArr := make([]string, 0, len(arr))
	// Iterate through the sorted array and append only distinct elements to the new array
	for i := 0; i < len(arr); i++ {
		if i == 0 || arr[i] != arr[i-1] {
			distinctArr = append(distinctArr, arr[i])
		}
	}
	return distinctArr
}

func compare(a, b string) int {
	return collator.CompareString(a, b)
}

package comment

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	_ "unicode/utf8"
)

var collator = collate.New(language.Und)

type queryInfo struct {
	db      *sql.DB
	toArray func(interface{}) interface {
		driver.Valuer
		sql.Scanner
	}
	table       string
	url         string
	id          string
	name        string
	displayName string
}

func NewQueryInfo(db *sql.DB, table string, url string, id string, name string, displayName string, toArray func(interface{}) interface {
	driver.Valuer
	sql.Scanner
}) queryInfo {
	return queryInfo{db: db, table: table, url: url, id: id, name: name, displayName: displayName, toArray: toArray}
}

func (i queryInfo) Load(ids []string) ([]Info, error) {
	rs := make([]Info, 0)
	if len(ids) == 0 {
		return rs, nil
	}
	ids = distinct(ids)
	querysql := fmt.Sprintf(`select %s as id, %s as url, COALESCE(%s,%s) as name from %s where %s = any(%s) and %s is not null order by %s`,
		i.id, i.url, i.displayName, i.name, i.table, i.id, "$1", i.url, i.id)
	r := make([]Info, 0)
	rows, err := i.db.Query(querysql, i.toArray(ids))
	for rows.Next() {
		var info Info
		err := rows.Scan(&info.Id, &info.Url, &info.Name)
		if err != nil {
			return nil, err
		}
		r = append(r, info)
	}
	if err != nil {
		return rs, err
	}
	return r, nil
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

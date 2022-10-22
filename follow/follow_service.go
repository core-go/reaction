package follow

import (
	"context"
	"database/sql"
	"fmt"
)

type FollowService interface {
	Follow(ctx context.Context, id string, target string) (int64, error)
	UnFollow(ctx context.Context, id string, target string) (int64, error)
	CheckFollow(ctx context.Context, id string, target string) (int, error)
}

func NewFollowService(
	db *sql.DB,
	followerTable string,
	followerIdCol string,
	followerCol string,
	followingTable string,
	followingIdCol string,
	followingCol string,
	userInfoTable string,
	userInfoIdCol string,
	followerCountCol string,
	followingCountCol string,
) FollowService {
	return &followService{
		DB:                db,
		FollowerTable:     followerTable,
		FollowerIdCol:     followerIdCol,
		FollowerCol:       followerCol,
		FollowingTable:    followingTable,
		FollowingIdCol:    followingIdCol,
		FollowingCol:      followingCol,
		UserInfoTable:     userInfoTable,
		UserInfoIdCol:     userInfoIdCol,
		FollowerCountCol:  followerCountCol,
		FollowingCountCol: followingCountCol,
	}
}

type followService struct {
	DB                *sql.DB
	FollowerTable     string
	FollowerIdCol     string
	FollowerCol       string
	FollowingTable    string
	FollowingIdCol    string
	FollowingCol      string
	UserInfoTable     string
	UserInfoIdCol     string
	FollowerCountCol  string
	FollowingCountCol string
}

func (s *followService) CheckFollow(ctx context.Context, id string, target string) (int, error) {
	query := fmt.Sprintf("select %s from %s where %s = $1 and %s = $2", s.FollowingIdCol, s.FollowingTable, s.FollowingIdCol, s.FollowingCol)
	rows, err := s.DB.QueryContext(ctx, query, id, target)
	if err != nil {
		return -1, err
	}
	var result []string

	defer rows.Close()
	for rows.Next() {
		var val string
		err := rows.Scan(&val)
		if err != nil {
			return -1, err
		}
		result = append(result, val)
	}

	return len(result), nil
}

func (s *followService) Follow(ctx context.Context, id string, target string) (int64, error) {
	query1 := fmt.Sprintf("insert into %s (%s, %s) values ($1, $2);", s.FollowerTable, s.FollowerIdCol, s.FollowerCol)
	stmt1, er0 := s.DB.Prepare(query1)
	if er0 != nil {
		return -1, nil
	}
	stmt1.ExecContext(ctx, id, target)

	query2 := fmt.Sprintf("insert into %s (%s, %s) values ($1, $2);", s.FollowingTable, s.FollowingIdCol, s.FollowingCol)
	fmt.Println(query2)
	stmt2, er0 := s.DB.Prepare(query2)
	if er0 != nil {
		return -1, nil
	}
	stmt2.ExecContext(ctx, target, id)

	query3 := fmt.Sprintf(
		"insert into %s(%s, %s, %s) values ($1, 0, 1) on conflict (%s) do update set %s = %s.%s + 1;",
		s.UserInfoTable, s.UserInfoIdCol, s.FollowerCountCol, s.FollowingCountCol, s.UserInfoIdCol, s.FollowingCountCol, s.FollowingCountCol, s.UserInfoTable)
	fmt.Println(query3)
	stmt3, er0 := s.DB.Prepare(query3)
	if er0 != nil {
		return -1, nil
	}
	stmt3.ExecContext(ctx, id)

	query4 := fmt.Sprintf(
		"insert into %s(%s, %s, %s) values ($1, 1, 0) on conflict (%s) do update set %s = %s.%s + 1;",
		s.UserInfoTable, s.UserInfoIdCol, s.FollowerCountCol, s.FollowingCountCol, s.UserInfoIdCol, s.FollowerCountCol, s.FollowerCountCol, s.UserInfoTable)
	fmt.Println(query4)
	stmt4, er0 := s.DB.Prepare(query4)
	if er0 != nil {
		return -1, nil
	}
	res4, er4 := stmt4.ExecContext(ctx, target)
	if er4 != nil {
		return -1, er4
	}

	return res4.RowsAffected()
}

func (s *followService) UnFollow(ctx context.Context, id string, target string) (int64, error) {
	query1 := fmt.Sprintf("delete from %s where %s = $1 and %s = $2;", s.FollowerTable, s.FollowerIdCol, s.FollowerCol)
	fmt.Println(query1)
	stmt1, er0 := s.DB.Prepare(query1)
	if er0 != nil {
		return -1, er0
	}
	stmt1.ExecContext(ctx, id, target)

	query2 := fmt.Sprintf("delete from %s where %s = $1 and %s = $2;", s.FollowingTable, s.FollowingIdCol, s.FollowingCol)
	fmt.Println(query2)
	stmt2, er0 := s.DB.Prepare(query2)
	if er0 != nil {
		return -1, er0
	}
	stmt2.ExecContext(ctx, target, id)

	query3 := fmt.Sprintf("update %s set %s = %s - 1 where %s = $1;", s.UserInfoTable, s.FollowingCountCol, s.FollowingCountCol, s.UserInfoIdCol)
	fmt.Println(query3)
	stmt3, er0 := s.DB.Prepare(query3)
	if er0 != nil {
		return -1, er0
	}
	stmt3.ExecContext(ctx, id)

	query4 := fmt.Sprintf("update %s set %s = %s - 1 where %s = $1;", s.UserInfoTable, s.FollowerCountCol, s.FollowerCountCol, s.UserInfoIdCol)
	fmt.Println(query4)
	stmt4, er0 := s.DB.Prepare(query4)
	if er0 != nil {
		return -1, er0
	}
	res, err := stmt4.ExecContext(ctx, target)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

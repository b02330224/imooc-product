package repositories

import (
	"database/sql"
	"errors"
	"imooc-product/common"
	"imooc-product/datamodels"
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (*datamodels.User, error)
	Insert(user *datamodels.User) (userId int64, err error)
}

func NewUserRepository(table string , db *sql.DB) IUserRepository {
	return &UserManagerRepository{table:table, mysqlConn:db}
}

type UserManagerRepository struct {
	table string
	mysqlConn *sql.DB
}

func (u *UserManagerRepository) Conn() (err error)  {
	if u.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}

	return
}

func (u *UserManagerRepository) Select (userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空")
	}

	if err = u.Conn();err != nil {
		return &datamodels.User{}, err
	}

	sql := "select * from " + u.table + " where userName = ?"
	rows, err := u.mysqlConn.Query(sql, userName)
	if err != nil {
		return &datamodels.User{}, err
	}

	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}


	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)

	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}

	sql := "insert into " + u.table + "(nickname, userName, passWord) values (?, ?, ?)"
	stmt, err := u.mysqlConn.Prepare(sql)
	if err != nil {
		return userId, err
	}

	result, err := stmt.Exec(user.Nickname, user.UserName, user.HashPassword)
	if err != nil {
		return userId, err
	}

	userId, err = result.LastInsertId()

	return
}

func (u *UserManagerRepository) SelectById(userId int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}

	sql := "select * from " + u.table + " where id=" + strconv.FormatInt(userId, 10)
	row, err := u.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.User{}, err
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在")
	}

	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}
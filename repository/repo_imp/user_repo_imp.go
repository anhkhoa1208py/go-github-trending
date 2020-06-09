package repo_imp

import (
	"backend-github-trending/banana"
	"backend-github-trending/db"
	"backend-github-trending/log"
	"backend-github-trending/model"
	"backend-github-trending/model/req"
	"backend-github-trending/repository"
	"context"
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type UserRepoIml struct {
	sql *db.Sql
}

func NewUserRep(sql *db.Sql) repository.UserRepo {
	return &UserRepoIml{
		sql: sql,
	}
}

func (u *UserRepoIml) SaveUser(context context.Context, user model.User) (model.User, error){
	stament := `
		INSERT INTO users(user_id, email, password, role, full_name, created_at, updated_at)
		VALUES(:user_id, :email, :password, :role, :full_name, :created_at, :updated_at)
	`
	user.CreateAt = time.Now()
	user.UpdateAt = time.Now()

	_, err := u.sql.Db.NamedExecContext(context, stament, user)
	if err != nil {
		log.Error(err.Error())
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation"{
				return user, banana.UserConflict
			}
		}
		return user, banana.SignUpFail
	}

	return user, nil
}

func (u *UserRepoIml) CheckLogin(context context.Context, loginReq req.RequestSignIn) (model.User, error)  {
	var user = model.User{}
	err := u.sql.Db.GetContext(context, &user, "SELECT * FROM users WHERE email=$1", loginReq.Email)

	if err != nil{
		if err == sql.ErrNoRows{
			return user, banana.UserNotFound
		}
		log.Error(err.Error())
		return user, err
	}

	return user, nil
}
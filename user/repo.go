package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const LIMIT int = 25

type repo struct {
	db  *sql.DB
	log *logrus.Entry
}

func newRepo(db *sql.DB, table string) *repo {
	return &repo{
		db:  db,
		log: logger.GetLogger().WithField("component", "user_type"),
	}
}

func (r *repo) InsertSingle(ctx context.Context, user User) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}
	id := uuid.New()
	now := time.Now().In(loc)
	hmac_sha256 := hmac.New(sha256.New, []byte(os.Getenv("JWT_SECRET_KEY")))
	hmac_sha256.Write([]byte(user.Password))

	password, err := bcrypt.GenerateFromPassword(hmac_sha256.Sum(nil), 12)
	if err != nil {
		r.log.Error(err)

		return err
	}

	sql_query :=
		`INSERT INTO
			user ()
		VALUES (
			unhex(replace(?, '-', '')),
			UPPER(?),
			?,
			?,
			?,
			?,
			?
		)`

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not begin a transacation: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, sql_query, id, id.String(), user.Email, string(password), user.Fullname, now)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) FindAll(ctx context.Context, max_limit int) ([]User, error) {
	var users []User
	var User User
	var sql_query string
	var rows *sql.Rows
	var err error

	sql_query = "SELECT * FROM user ORDER BY fullname ASC LIMIT ?"
	rows, err = r.db.QueryContext(ctx, sql_query, max_limit+1)
	if err != nil {
		r.log.Error(err)

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&User.Id, &User.Fullname, &User.Email, &User.CreatedAt)
		if err != nil {
			r.log.Error(err)

			return nil, err
		}

		users = append(users, User)
	}

	r.log.Info("users data successfully retrieved from database")
	return users, nil
}

// func (model *repo) FindAfterId(ctx context.Context, id string, max_limit int) ([]User, error) {
// 	var users []User
// 	var User User
// 	var sql_query string
// 	var rows *sql.Rows
// 	var err error

// 	sql_query = "SELECT fullname FROM user WHERE id_text = ?"
// 	if err = r.db.QueryRowContext(ctx, sql_query, id).Scan(&User.Fullname); err != nil {
// 		r.log.Error(err)

// 		return nil, err
// 	}

// 	sql_query = fmt.Sprintf("SELECT %s FROM user WHERE fullname > ? ORDER BY fullname ASC LIMIT ?", r.fields)
// 	rows, err = r.db.QueryContext(ctx, sql_query, User.Fullname, max_limit+1)
// 	if err != nil {
// 		r.log.Error(err)

// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		err = rows.Scan(&User.Id, &User.Fullname, &User.Gender, &User.Email, &User.Username, &User.DateOfBirth, &User.CreatedAt)
// 		if err != nil {
// 			r.log.Error(err)

// 			return nil, err
// 		}

// 		users = append(users, User)
// 	}

// 	return users, nil
// }

// func (model *repo) FindBeforeId(ctx context.Context, id string, max_limit int) ([]User, error) {
// 	var users []User
// 	var User User
// 	var sql_query string
// 	var rows *sql.Rows
// 	var err error

// 	sql_query = "SELECT fullname FROM user WHERE id_text = ?"
// 	if err := r.db.QueryRowContext(ctx, sql_query, id).Scan(&User.Fullname); err != nil {
// 		r.log.Error(err)

// 		return nil, err
// 	}

// 	sql_query = fmt.Sprintf(`
// 	SELECT a.* FROM (
// 		SELECT %s FROM user WHERE fullname < ? ORDER BY fullname DESC LIMIT ?
// 	) a
// 	ORDER BY a.fullname ASC;`, r.fields)
// 	rows, err = r.db.QueryContext(ctx, sql_query, User.Fullname, max_limit+1)
// 	if err != nil {
// 		r.log.Error(err)

// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		err = rows.Scan(&User.Id, &User.Fullname, &User.Gender, &User.Email, &User.Username, &User.DateOfBirth, &User.CreatedAt)
// 		if err != nil {
// 			r.log.Error(err)

// 			return nil, err
// 		}

// 		users = append(users, User)
// 	}

// 	return users, nil
// }

// func (model *repo) FindById(ctx context.Context, id uuid.UUID) (User, error) {
// 	var user User

// 	sql_query := fmt.Sprintf("SELECT %s FROM user WHERE id_text=?", r.fields)
// 	rows, err := r.db.QueryContext(ctx, sql_query, id)
// 	if err != nil {
// 		r.log.Error(err)

// 		return user, err
// 	}
// 	defer rows.Close()

// 	if rows.Next() {
// 		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Fullname, &user.Gender, &user.DateOfBirth, &user.CreatedAt)
// 		if err != nil {
// 			r.log.Error(err)

// 			return user, err
// 		}
// 	}

// 	return user, nil
// }

package customer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mmiftahrzki/customer/auth"
	"github.com/mmiftahrzki/customer/customer/address"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

const limit int = 25

type repo struct {
	db  *sql.DB
	log *logrus.Entry
}

func newRepo(db *sql.DB) repo {
	return repo{
		db:  db,
		log: logger.GetLogger().WithField("component", "customerRepo"),
	}
}

func (r *repo) SelectAll(ctx context.Context) ([]modelSQL, error) {
	var sqlModel modelSQL
	var sqlModels []modelSQL
	const sqlQuery string = `SELECT a.id,
			a.email,
			a.first_name,
			a.last_name,
			a.address_id,
			a.active,
			a.created_at,
			b.id,
			b.address,
			b.district,
			b.city_id,
			b.postal_code
		FROM customer a
			JOIN address b ON b.id = a.address_id
		WHERE a.active = true
		ORDER BY a.id ASC
		LIMIT ?`

	select {
	case <-ctx.Done():
		r.log.Info("deadline exceeded from repo layer")

		return sqlModels, ctx.Err()
	default:
		rows, sqlErr := r.db.QueryContext(ctx, sqlQuery, limit+1)
		if sqlErr != nil {
			return sqlModels, sqlErr
		}
		defer rows.Close()

		for rows.Next() {
			rowScanErr := rows.Scan(
				&sqlModel.id,
				&sqlModel.email,
				&sqlModel.firstName,
				&sqlModel.lastName,
				&sqlModel.addressId,
				&sqlModel.active,
				&sqlModel.createdAt,
				&sqlModel.address.Id,
				&sqlModel.address.Address,
				&sqlModel.address.District,
				&sqlModel.address.CityId,
				&sqlModel.address.PostalCode,
			)
			if rowScanErr != nil {
				return sqlModels, rowScanErr
			}

			sqlModels = append(sqlModels, sqlModel)
		}

		r.log.Info("customers data successfully retrieved from database")

		return sqlModels, nil
	}
}

func (r *repo) SelectAllPrev(ctx context.Context, customer modelRead) (modelSQLs []modelSQL, err error) {
	var modelSQL modelSQL
	const sqlQuery string = `SELECT a.id,
			a.email,
			a.first_name,
			a.last_name,
			a.address_id,
			a.active,
			a.created_at,
			b.id,
			b.address,
			b.district,
			b.city_id,
			b.postal_code
		FROM customer a
			JOIN address b ON b.id = a.address_id
		WHERE a.active = TRUE
			AND a.id < ?
		ORDER BY a.id DESC
      LIMIT ?`

	rows, err := r.db.QueryContext(ctx, sqlQuery, customer.Id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&modelSQL.id,
			&modelSQL.email,
			&modelSQL.firstName,
			&modelSQL.lastName,
			&modelSQL.addressId,
			&modelSQL.active,
			&modelSQL.createdAt,
			&modelSQL.address.Id,
			&modelSQL.address.Address,
			&modelSQL.address.District,
			&modelSQL.address.CityId,
			&modelSQL.address.PostalCode,
		)

		if err != nil {
			return nil, err
		}

		modelSQLs = append(modelSQLs, modelSQL)
	}

	return modelSQLs, nil
}

func (r *repo) SelectAllNext(ctx context.Context, customer modelRead) (modelSQLs []modelSQL, err error) {
	var modelSQL modelSQL
	const sqlQuery string = `SELECT a.id,
			a.email,
			a.first_name,
			a.last_name,
			a.address_id,
			a.active,
			a.created_at,
			b.id,
			b.address,
			b.district,
			b.city_id,
			b.postal_code
		FROM customer a
			JOIN address b ON b.id = a.address_id
		WHERE a.active = TRUE
			AND a.id > ?
		ORDER BY a.id ASC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, sqlQuery, customer.Id, limit+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&modelSQL.id,
			&modelSQL.email,
			&modelSQL.firstName,
			&modelSQL.lastName,
			&modelSQL.addressId,
			&modelSQL.active,
			&modelSQL.createdAt,
			&modelSQL.address.Id,
			&modelSQL.address.Address,
			&modelSQL.address.District,
			&modelSQL.address.CityId,
			&modelSQL.address.PostalCode,
		)

		if err != nil {
			return nil, err
		}

		modelSQLs = append(modelSQLs, modelSQL)
	}

	return modelSQLs, nil
}

func (r *repo) SelectSingleById(ctx context.Context, id int) (modelSQL, error) {
	var modelSQL modelSQL
	const sqlQuery string = `
		SELECT a.id,
			a.email,
			a.first_name,
			a.last_name,
			a.address_id,
			a.active,
			a.created_at,
			b.id,
			b.address,
			b.district,
			b.city_id,
			b.postal_code
		FROM customer a
			JOIN address b ON b.id = a.address_id
		WHERE a.active = true
			AND a.id=?`
	rows, err := r.db.QueryContext(ctx, sqlQuery, id)
	if err != nil {
		return modelSQL, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&modelSQL.id,
			&modelSQL.email,
			&modelSQL.firstName,
			&modelSQL.lastName,
			&modelSQL.addressId,
			&modelSQL.active,
			&modelSQL.createdAt,
			&modelSQL.address.Id,
			&modelSQL.address.Address,
			&modelSQL.address.District,
			&modelSQL.address.CityId,
			&modelSQL.address.PostalCode,
		)
		if err != nil {
			return modelSQL, err
		}
	}

	return modelSQL, nil
}

func (r *repo) UpdateSingleById(ctx context.Context, id int, payload modelUpdate) error {
	const sqlQuery string = "UPDATE customer SET first_name=?, last_name=?, email=? WHERE id=?"
	_, dbErr := r.db.ExecContext(ctx, sqlQuery, payload.FirstName, payload.LastName, payload.Email, id)
	if dbErr != nil {
		return dbErr
	}

	return nil
}

func (r *repo) DeleteSingleById(ctx context.Context, id int) error {
	JWTContext := ctx.Value(auth.JWTContextKey)
	claim, ok := JWTContext.(*auth.ModelClaim)
	if !ok {
		return errors.New("asd")
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start a transaction: %w", err)
	}
	defer tx.Rollback()

	sqlQuery := "DELETE FROM customer a WHERE a.id = ? AND a.created_by = ?"
	_, err = tx.ExecContext(ctx, sqlQuery, id, claim.Email)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) InsertSingle(ctx context.Context, payload modelCreate) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	now := time.Now().In(loc)
	sqlQuery :=
		`INSERT INTO customer (
				first_name,
				last_name,
				email,
				created_at,
				address_id
			)
		VALUES (?, ?, ?, ?, ?);`

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not begin a transacation: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, sqlQuery, payload.FirstName, payload.LastName, payload.Email, now, 1)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) UpdateSingleAddressByCustomerId(ctx context.Context, id uint16, payload address.ModelUpdate) error {
	fields := []string{}
	structFields := []any{}

	if payload.Address != nil {
		fields = append(fields, "address=?")
		structFields = append(structFields, payload.Address)
	}

	if payload.Address2 != nil {
		fields = append(fields, "address2=?")
		structFields = append(structFields, payload.Address2)
	}

	if payload.District != nil {
		fields = append(fields, "district=?")
		structFields = append(structFields, payload.District)
	}

	if payload.PostalCode != nil {
		fields = append(fields, "postal_code=?")
		structFields = append(structFields, payload.PostalCode)
	}

	fields = append(fields, "last_update=?")
	structFields = append(structFields, time.Now())

	fieldsStr := strings.Join(fields, ", ")
	structFields = append(structFields, id)

	sqlQuery := fmt.Sprintf("UPDATE address SET %s WHERE address_id=?", fieldsStr)
	_, err := r.db.ExecContext(ctx, sqlQuery, structFields...)
	if err != nil {
		return err
	}

	return nil
}

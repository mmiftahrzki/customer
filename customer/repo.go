package customer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

const limit int = 25

type Repo interface {
	SelectAll(ctx context.Context) ([]modelSQL, error)
	SelectAllPrev(ctx context.Context, customer ModelRead) (modelSQLs []modelSQL, err error)
	SelectAllNext(ctx context.Context, customer ModelRead) (modelSQLs []modelSQL, err error)
	SelectSingleById(ctx context.Context, id int) (modelSQL, error)
	InsertSingle(ctx context.Context, payload modelCreate) error
	UpdateSingleById(ctx context.Context, id int, payload modelUpdate) error
	DeleteSingleById(ctx context.Context, id int) error
}

type repo struct {
	db  *sql.DB
	log *logrus.Entry
}

func newRepo(db *sql.DB) Repo {
	return &repo{
		db:  db,
		log: logger.GetLogger().WithField("component", "customerRepo"),
	}
}

func handleRows(rows *sql.Rows) ([]modelSQL, error) {
	var modelSql modelSQL
	var modelSqls []modelSQL
	var err error

	for rows.Next() {
		err = rows.Scan(
			&modelSql.id,
			&modelSql.email,
			&modelSql.firstName,
			&modelSql.lastName,
			&modelSql.addressId,
			&modelSql.active,
			&modelSql.createdAt,
			&modelSql.address.Id,
			&modelSql.address.Address,
			&modelSql.address.District,
			&modelSql.address.CityId,
			&modelSql.address.City,
			&modelSql.address.PostalCode,
			&modelSql.address.CountryId,
			&modelSql.address.Country,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}

		modelSqls = append(modelSqls, modelSql)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return modelSqls, nil
}

func (r *repo) SelectAll(ctx context.Context) ([]modelSQL, error) {
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
			c.city,
			b.postal_code,
			c.country_id,
			d.country
		FROM customer a
			LEFT JOIN address b ON b.id = a.address_id
			JOIN city c ON c.id = b.city_id
			JOIN country d ON d.id = c.country_id
		WHERE a.active = true
		ORDER BY a.id ASC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, sqlQuery, limit+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers, err := handleRows(rows)
	if err != nil {
		return nil, err
	}

	r.log.Info("customers data successfully retrieved from database")

	return customers, nil
}

func (r *repo) SelectAllPrev(ctx context.Context, customer ModelRead) (modelSQLs []modelSQL, err error) {
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
			c.city,
			b.postal_code,
			c.country_id,
			d.country
		FROM customer a
			LEFT JOIN address b ON b.id = a.address_id
			JOIN city c ON c.id = b.city_id
			JOIN country d ON d.id = c.country_id
		WHERE a.active = TRUE
			AND a.id < ?
		ORDER BY a.id DESC
      LIMIT ?`

	rows, err := r.db.QueryContext(ctx, sqlQuery, customer.Id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers, err := handleRows(rows)
	if err != nil {
		return nil, err
	}

	r.log.Info("customers data successfully retrieved from database")

	return customers, nil
}

func (r *repo) SelectAllNext(ctx context.Context, customer ModelRead) (modelSQLs []modelSQL, err error) {
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
			c.city,
			b.postal_code,
			c.country_id,
			d.country
		FROM customer a
			LEFT JOIN address b ON b.id = a.address_id
			JOIN city c ON c.id = b.city_id
			JOIN country d ON d.id = c.country_id
		WHERE a.active = TRUE
			AND a.id > ?
		ORDER BY a.id ASC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, sqlQuery, customer.Id, limit+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers, err := handleRows(rows)
	if err != nil {
		return nil, err
	}

	r.log.Info("customers data successfully retrieved from database")

	return customers, nil
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
			c.city,
			b.postal_code,
			c.country_id,
			d.country
		FROM customer a
			LEFT JOIN address b ON b.id = a.address_id
			JOIN city c ON c.id = b.city_id
			JOIN country d ON d.id = c.country_id
		WHERE a.active = true
			AND a.id=?`
	row := r.db.QueryRowContext(ctx, sqlQuery, id)
	err := row.Scan(
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
		&modelSQL.address.City,
		&modelSQL.address.PostalCode,
		&modelSQL.address.CountryId,
		&modelSQL.address.Country,
	)
	if err != nil {
		return modelSQL, err
	}

	return modelSQL, nil
}

func (r *repo) UpdateSingleById(ctx context.Context, id int, payload modelUpdate) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	now := time.Now().In(loc).Format(time.RFC3339)

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("couldnot begin a transaction: %w", err)
	}
	defer tx.Rollback()

	const sqlQuerySelectAddressId string = "SELECT address_id FROM customer WHERE id=?"
	rows, err := tx.QueryContext(ctx, sqlQuerySelectAddressId, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	addressId := 0
	for rows.Next() {
		err = rows.Scan(&addressId)

		if err != nil {
			return err
		}
	}

	const sqlQueryUpdateCustomer string = "UPDATE customer SET first_name=?, last_name=?, email=?, last_update=? WHERE id=?"
	_, err = tx.ExecContext(ctx, sqlQueryUpdateCustomer, payload.FirstName, payload.LastName, payload.Email, now, id)
	if err != nil {
		return err
	}

	const sqlQueryUpdateAddress string = "UPDATE address SET address=?, district=?, postal_code=?, city_id=?, last_update=? WHERE id=?"
	_, err = tx.ExecContext(ctx, sqlQueryUpdateAddress, payload.Address.Address, payload.Address.District, payload.Address.PostalCode, payload.Address.CityId, now, addressId)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteSingleById(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start a transaction: %w", err)
	}
	defer tx.Rollback()

	sqlQuery := "DELETE FROM customer a WHERE a.id = ?"
	_, err = tx.ExecContext(ctx, sqlQuery, id)
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

	now := time.Now().In(loc).Format(time.RFC3339)
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

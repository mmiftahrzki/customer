package customer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mmiftahrzki/customer/customer/address"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

const LIMIT int = 25

type repo struct {
	db  *sql.DB
	log *logrus.Entry
}

func newRepo(db *sql.DB) *repo {
	return &repo{
		db:  db,
		log: logger.GetLogger().WithField("component", "customer_repo"),
	}
}

func (r *repo) SelectAll(ctx context.Context) (sql_models []customerSQLModel, err error) {
	var sql_model customerSQLModel

	sql_query := `
		SELECT a.customer_id,
			a.email,
			a.first_name,
			a.last_name,
			a.address_id,
			a.active,
			a.create_date,
			b.address_id,
			b.address,
			b.address2,
			b.district,
			b.city_id,
			b.postal_code
		FROM customer a
			JOIN address b ON b.address_id = a.address_id
		WHERE a.active = true
		ORDER BY a.customer_id ASC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, sql_query, LIMIT+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&sql_model.customer_id,
			&sql_model.email,
			&sql_model.first_name,
			&sql_model.last_name,
			&sql_model.address_id,
			&sql_model.active,
			&sql_model.create_date,
			&sql_model.address.AddressId,
			&sql_model.address.Address,
			&sql_model.address.Address2,
			&sql_model.address.District,
			&sql_model.address.CityId,
			&sql_model.address.PostalCode,
		)

		if err != nil {
			return nil, err
		}

		sql_models = append(sql_models, sql_model)
	}

	r.log.Info("customers data successfully retrieved from database")

	return sql_models, nil
}

// func (r *repo) SelectAllPrev(ctx context.Context, customer customerReadModel) (sql_models []customerSQLModel, err error) {
// 	var sql_model customerSQLModel

// 	sql_query := `
// 		SELECT a.customer_id,
// 			a.email,
// 			a.first_name,
// 			a.last_name,
// 			a.address_id,
// 			a.active,
// 			a.create_date,
// 			b.address_id,
// 			b.address,
// 			b.address2,
// 			b.district,
// 			b.city_id,
// 			b.postal_code
// 		FROM customer a
// 			JOIN address b ON b.address_id = a.address_id
// 		WHERE a.active = true
// 			AND a.customer_id < ?
// 		ORDER BY a.customer_id ASC
// 		LIMIT ?`

// 	rows, err := r.db.QueryContext(ctx, sql_query, customer.Id, LIMIT+1)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		err = rows.Scan(
// 			&sql_model.customer_id,
// 			&sql_model.email,
// 			&sql_model.first_name,
// 			&sql_model.last_name,
// 			&sql_model.address_id,
// 			&sql_model.active,
// 			&sql_model.create_date,
// 			&sql_model.address.AddressId,
// 			&sql_model.address.Address,
// 			&sql_model.address.Address2,
// 			&sql_model.address.District,
// 			&sql_model.address.CityId,
// 			&sql_model.address.PostalCode,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		sql_models = append(sql_models, sql_model)
// 	}

// 	return sql_models, nil
// }

func (r *repo) SelectAllNext(ctx context.Context, customer customerReadModel) (sql_models []customerSQLModel, err error) {
	var sql_model customerSQLModel

	sql_query :=
		`SELECT a.customer_id,
			a.email,
			a.first_name,
			a.last_name,
			a.address_id,
			a.active,
			a.create_date,
			b.address_id,
			b.address,
			b.address2,
			b.district,
			b.city_id,
			b.postal_code
		FROM customer a
			JOIN address b ON b.address_id = a.address_id
		WHERE a.active = true
			AND a.customer_id > ?
		ORDER BY a.customer_id ASC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, sql_query, customer.Id, LIMIT+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&sql_model.customer_id,
			&sql_model.email,
			&sql_model.first_name,
			&sql_model.last_name,
			&sql_model.address_id,
			&sql_model.active,
			&sql_model.create_date,
			&sql_model.address.AddressId,
			&sql_model.address.Address,
			&sql_model.address.Address2,
			&sql_model.address.District,
			&sql_model.address.CityId,
			&sql_model.address.PostalCode,
		)

		if err != nil {
			return nil, err
		}

		sql_models = append(sql_models, sql_model)
	}

	return sql_models, nil
}

func (r *repo) SelectSingleById(ctx context.Context, id int) (customer_sql_model customerSQLModel, err error) {
	sql_query := `
		SELECT a.customer_id,
			a.email,
			a.first_name,
			a.last_name,
			a.address_id,
			a.active,
			a.create_date,
			b.address_id,
			b.address,
			b.address2,
			b.district,
			b.city_id,
			b.postal_code
		FROM customer a
			JOIN address b ON b.address_id = a.address_id
		WHERE a.active = true
			AND a.customer_id=?`
	rows, err := r.db.QueryContext(ctx, sql_query, id)
	if err != nil {
		return customer_sql_model, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&customer_sql_model.customer_id,
			&customer_sql_model.email,
			&customer_sql_model.first_name,
			&customer_sql_model.last_name,
			&customer_sql_model.address_id,
			&customer_sql_model.active,
			&customer_sql_model.create_date,
			&customer_sql_model.address.AddressId,
			&customer_sql_model.address.Address,
			&customer_sql_model.address.Address2,
			&customer_sql_model.address.District,
			&customer_sql_model.address.CityId,
			&customer_sql_model.address.PostalCode,
		)
		if err != nil {
			return customer_sql_model, err
		}
	}

	return customer_sql_model, nil
}

func (r *repo) UpdateSingleById(ctx context.Context, id int, payload CustomerUpdateModel) error {
	sql_query := "UPDATE customer SET first_name=?, last_name=?, email=? WHERE customer_id=?"
	_, err := r.db.ExecContext(ctx, sql_query, payload.FirstName, payload.LastName, payload.Email, id)
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

	sql_query := "DELETE FROM customer a WHERE a.customer_id = ?"
	_, err = tx.ExecContext(ctx, sql_query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) InsertSingle(ctx context.Context, payload CustomerCreateModel) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}

	now := time.Now().In(loc)
	sql_query := `
		INSERT INTO customer (
				first_name,
				last_name,
				email,
				create_date,
				store_id,
				address_id
			)
		VALUES (?, ?, ?, ?, ?, ?)`

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not begin a transacation: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, sql_query, payload.FirstName, payload.LastName, payload.Email, now, 1, 1)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) UpdateSingleAddressByCustomerId(ctx context.Context, id uint16, payload address.AddressUpdateModel) error {
	fields := []string{}
	struct_fields := []any{}

	if payload.Address != nil {
		fields = append(fields, "address=?")
		struct_fields = append(struct_fields, payload.Address)
	}

	if payload.Address2 != nil {
		fields = append(fields, "address2=?")
		struct_fields = append(struct_fields, payload.Address2)
	}

	if payload.District != nil {
		fields = append(fields, "district=?")
		struct_fields = append(struct_fields, payload.District)
	}

	if payload.PostalCode != nil {
		fields = append(fields, "postal_code=?")
		struct_fields = append(struct_fields, payload.PostalCode)
	}

	fields = append(fields, "last_update=?")
	struct_fields = append(struct_fields, time.Now())

	fields_string := strings.Join(fields, ", ")
	struct_fields = append(struct_fields, id)

	sql_query := fmt.Sprintf("UPDATE address SET %s WHERE address_id=?", fields_string)
	_, err := r.db.ExecContext(ctx, sql_query, struct_fields...)
	if err != nil {
		return err
	}

	return nil
}

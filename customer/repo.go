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
		log: logger.GetLogger().WithField("component", "customer_repo"),
	}
}

func (r *repo) SelectAll(ctx context.Context) (sql_models []modelSQL, err error) {
	var sql_model modelSQL

	sql_query :=
		`SELECT a.id,
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

	rows, err := r.db.QueryContext(ctx, sql_query, limit+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&sql_model.id,
			&sql_model.email,
			&sql_model.firstName,
			&sql_model.lastName,
			&sql_model.addressId,
			&sql_model.active,
			&sql_model.createdAt,
			&sql_model.address.Id,
			&sql_model.address.Address,
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

func (r *repo) SelectAllPrev(ctx context.Context, customer modelRead) (sql_models []modelSQL, err error) {
	var sql_model modelSQL

	sql_query :=
		`SELECT a.id,
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

	rows, err := r.db.QueryContext(ctx, sql_query, customer.Id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&sql_model.id,
			&sql_model.email,
			&sql_model.firstName,
			&sql_model.lastName,
			&sql_model.addressId,
			&sql_model.active,
			&sql_model.createdAt,
			&sql_model.address.Id,
			&sql_model.address.Address,
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

func (r *repo) SelectAllNext(ctx context.Context, customer modelRead) (sql_models []modelSQL, err error) {
	var sql_model modelSQL

	sql_query :=
		`SELECT a.id,
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

	rows, err := r.db.QueryContext(ctx, sql_query, customer.Id, limit+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&sql_model.id,
			&sql_model.email,
			&sql_model.firstName,
			&sql_model.lastName,
			&sql_model.addressId,
			&sql_model.active,
			&sql_model.createdAt,
			&sql_model.address.Id,
			&sql_model.address.Address,
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

func (r *repo) SelectSingleById(ctx context.Context, id int) (sql_model modelSQL, err error) {
	sql_query := `
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
	rows, err := r.db.QueryContext(ctx, sql_query, id)
	if err != nil {
		return sql_model, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&sql_model.id,
			&sql_model.email,
			&sql_model.firstName,
			&sql_model.lastName,
			&sql_model.addressId,
			&sql_model.active,
			&sql_model.createdAt,
			&sql_model.address.Id,
			&sql_model.address.Address,
			&sql_model.address.District,
			&sql_model.address.CityId,
			&sql_model.address.PostalCode,
		)
		if err != nil {
			return sql_model, err
		}
	}

	return sql_model, nil
}

func (r *repo) UpdateSingleById(ctx context.Context, id int, payload updateModel) error {
	sql_query := "UPDATE customer SET first_name=?, last_name=?, email=? WHERE id=?"
	_, err := r.db.ExecContext(ctx, sql_query, payload.FirstName, payload.LastName, payload.Email, id)
	if err != nil {
		return err
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

	sql_query := "DELETE FROM customer a WHERE a.id = ? AND a.created_by = ?"
	_, err = tx.ExecContext(ctx, sql_query, id, claim.Email)
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
	sql_query :=
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

	_, err = tx.ExecContext(ctx, sql_query, payload.FirstName, payload.LastName, payload.Email, now, 1)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) UpdateSingleAddressByCustomerId(ctx context.Context, id uint16, payload address.ModelUpdate) error {
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

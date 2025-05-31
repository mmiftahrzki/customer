package customer

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/mmiftahrzki/customer/customer/address"
)

const LIMIT int = 25

type repo struct {
	DB *sql.DB
}

func newRepo(db *sql.DB) *repo {
	return &repo{
		DB: db,
	}
}

func (r *repo) SelectAll(ctx context.Context) ([]Customer, error) {
	var entity CustomerSql
	var model Customer
	var models []Customer
	var sql_address address.AddressSql

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

	rows, err := r.DB.QueryContext(ctx, sql_query, LIMIT+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&entity.customer_id,
			&entity.email,
			&entity.first_name,
			&entity.last_name,
			&entity.address_id,
			&entity.active,
			&entity.create_date,
			&sql_address.AddressId,
			&sql_address.Address,
			&sql_address.Address2,
			&sql_address.District,
			&sql_address.CityId,
			&sql_address.PostalCode,
		)
		if err != nil {
			return nil, err
		}

		if entity.customer_id.Valid {
			model.Id = uint8(entity.customer_id.Int16)
		}

		if entity.email.Valid {
			model.Email = entity.email.String
		}

		if entity.first_name.Valid {
			model.FullName = entity.first_name.String
		}

		if entity.last_name.Valid {
			if len(model.FullName) > 0 {
				model.FullName += " " + entity.last_name.String
			}
		}

		if entity.address_id.Valid {
			model.Address.Id = uint8(entity.address_id.Int16)
		}

		if sql_address.Address.Valid {
			model.Address.Address = sql_address.Address.String
		}

		if sql_address.Address2.Valid {
			model.Address.Address2 = sql_address.Address2.String
		}

		if sql_address.District.Valid {
			model.Address.District = sql_address.District.String
		}

		if sql_address.CityId.Valid {
			model.Address.CityId = uint8(sql_address.CityId.Int16)
		}

		if sql_address.PostalCode.Valid {
			model.Address.PostalCode = sql_address.PostalCode.String
		}

		if entity.create_date.Valid {
			model.CreatedAt = entity.create_date.Time
		}

		models = append(models, model)
	}

	return models, nil
}

func (r *repo) SelectAllPrev(ctx context.Context, customer Customer) ([]Customer, error) {
	var model Customer
	var entity CustomerSql
	var entity_address address.AddressSql
	var models []Customer

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
			AND a.customer_id < ?
		ORDER BY a.customer_id ASC
		LIMIT ?`

	rows, err := r.DB.QueryContext(ctx, sql_query, customer.Id, LIMIT+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&entity.customer_id,
			&entity.email,
			&entity.first_name,
			&entity.last_name,
			&entity.address_id,
			&entity.active,
			&entity.create_date,
			&entity_address.AddressId,
			&entity_address.Address,
			&entity_address.Address2,
			&entity_address.District,
			&entity_address.CityId,
			&entity_address.PostalCode,
		)

		if err != nil {
			return nil, err
		}

		if entity.customer_id.Valid {
			model.Id = uint8(entity.customer_id.Int16)
		}

		if entity.email.Valid {
			model.Email = entity.email.String
		}

		if entity.first_name.Valid {
			model.FullName = entity.first_name.String
		}

		if entity.last_name.Valid {
			model.FullName += " " + entity.last_name.String
		}

		if entity.address_id.Valid {
			model.Address.Id = uint8(entity.address_id.Int16)
		}

		if entity_address.Address.Valid {
			model.Address.Address = entity_address.Address.String
		}

		if entity_address.Address2.Valid {
			model.Address.Address2 = entity_address.Address2.String
		}

		if entity_address.District.Valid {
			model.Address.District = entity_address.District.String
		}

		if entity_address.CityId.Valid {
			model.Address.CityId = uint8(entity_address.CityId.Int16)
		}

		if entity_address.PostalCode.Valid {
			model.Address.PostalCode = entity_address.PostalCode.String
		}

		if entity.create_date.Valid {
			model.CreatedAt = entity.create_date.Time
		}

		models = append(models, model)
	}

	return models, nil
}

func (r *repo) SelectAllNext(ctx context.Context, customer Customer) ([]Customer, error) {
	var model Customer
	var entity CustomerSql
	var entity_address address.AddressSql
	var models []Customer

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
			AND a.customer_id > ?
		ORDER BY a.customer_id ASC
		LIMIT ?`

	rows, err := r.DB.QueryContext(ctx, sql_query, customer.Id, LIMIT+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&entity.customer_id,
			&entity.email,
			&entity.first_name,
			&entity.last_name,
			&entity.address_id,
			&entity.active,
			&entity.create_date,
			&entity_address.AddressId,
			&entity_address.Address,
			&entity_address.Address2,
			&entity_address.District,
			&entity_address.CityId,
			&entity_address.PostalCode,
		)

		if err != nil {
			return nil, err
		}

		if entity.customer_id.Valid {
			model.Id = uint8(entity.customer_id.Int16)
		}

		if entity.email.Valid {
			model.Email = entity.email.String
		}

		if entity.first_name.Valid {
			model.FullName = entity.first_name.String
		}

		if entity.last_name.Valid {
			model.FullName += " " + entity.last_name.String
		}

		if entity.address_id.Valid {
			model.Address.Id = uint8(entity.address_id.Int16)
		}

		if entity_address.Address.Valid {
			model.Address.Address = entity_address.Address.String
		}

		if entity_address.Address2.Valid {
			model.Address.Address2 = entity_address.Address2.String
		}

		if entity_address.District.Valid {
			model.Address.District = entity_address.District.String
		}

		if entity_address.CityId.Valid {
			model.Address.CityId = uint8(entity_address.CityId.Int16)
		}

		if entity_address.PostalCode.Valid {
			model.Address.PostalCode = entity_address.PostalCode.String
		}

		if entity.create_date.Valid {
			model.CreatedAt = entity.create_date.Time
		}

		models = append(models, model)
	}

	return models, nil
}

func (r *repo) SelectSingleById(ctx context.Context, id uint) (Customer, error) {
	var entity CustomerSql
	var model Customer
	var sql_address address.AddressSql

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
			AND a.customer_id = ?`
	rows, err := r.DB.QueryContext(ctx, sql_query, id)
	if err != nil {
		return model, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(
			&entity.customer_id,
			&entity.email,
			&entity.first_name,
			&entity.last_name,
			&entity.address_id,
			&entity.active,
			&entity.create_date,
			&sql_address.AddressId,
			&sql_address.Address,
			&sql_address.Address2,
			&sql_address.District,
			&sql_address.CityId,
			&sql_address.PostalCode,
		)
		if err != nil {
			return model, err
		}

		if entity.customer_id.Valid {
			model.Id = uint8(entity.customer_id.Int16)
		}

		if entity.email.Valid {
			model.Email = entity.email.String
		}

		// if entity.first_name.Valid {
		// 	model.FirstName = entity.first_name.String
		// }

		// if entity.last_name.Valid {
		// 	model.LastName = entity.last_name.String
		// }

		if entity.first_name.Valid {
			model.FullName = entity.first_name.String
		}

		if entity.last_name.Valid {
			model.FullName += " " + entity.last_name.String
		}

		if entity.address_id.Valid {
			model.Address.Id = uint8(entity.address_id.Int16)
		}

		if sql_address.Address.Valid {
			model.Address.Address = sql_address.Address.String
		}

		if sql_address.Address2.Valid {
			model.Address.Address2 = sql_address.Address2.String
		}

		if sql_address.District.Valid {
			model.Address.District = sql_address.District.String
		}

		if sql_address.CityId.Valid {
			model.Address.CityId = uint8(sql_address.CityId.Int16)
		}

		if sql_address.PostalCode.Valid {
			model.Address.PostalCode = sql_address.PostalCode.String
		}

		if entity.create_date.Valid {
			model.CreatedAt = entity.create_date.Time
		}
	}

	return model, nil
}

func (r *repo) UpdateSingleById(ctx context.Context, id uint, payload CustomerUpdateModel) error {
	fields := []string{}
	struct_fields := []any{}

	if payload.FirstName != "" {
		fields = append(fields, "first_name=?")
		struct_fields = append(struct_fields, payload.FirstName)
	}

	if payload.LastName != "" {
		fields = append(fields, "last_name=?")
		struct_fields = append(struct_fields, payload.LastName)
	}

	if payload.Email != "" {
		fields = append(fields, "email=?")
		struct_fields = append(struct_fields, payload.Email)
	}

	fields_string := strings.Join(fields, ", ")
	struct_fields = append(struct_fields, id)

	sql_query := fmt.Sprintf("UPDATE customer SET %s WHERE customer_id = ?", fields_string)
	_, err := r.DB.ExecContext(ctx, sql_query, struct_fields...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteSingleById(ctx context.Context, id uint) error {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{})
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
		INSERT INTO customer(
				first_name,
				last_name,
				email,
				create_date,
				store_id,
				address_id
			)
		VALUES (?, ?, ?, ?, ?, ?)`

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{})
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

func (r *repo) UpdateSingleAddressByCustomerId(ctx context.Context, id uint8, payload address.AddressUpdateModel) error {
	fields := []string{}
	struct_fields := []any{}

	if payload.Address != "" {
		fields = append(fields, "address=?")
		struct_fields = append(struct_fields, payload.Address)
	}

	if payload.Address2 != "" {
		fields = append(fields, "address2=?")
		struct_fields = append(struct_fields, payload.Address2)
	}

	if payload.District != "" {
		fields = append(fields, "district=?")
		struct_fields = append(struct_fields, payload.District)
	}

	if payload.CityId != 0 {
		fields = append(fields, "city_id=?")
		struct_fields = append(struct_fields, payload.CityId)
	}

	if payload.PostalCode != "" {
		fields = append(fields, "postal_code=?")
		struct_fields = append(struct_fields, payload.PostalCode)
	}

	fields = append(fields, "last_update=?")
	struct_fields = append(struct_fields, time.Now())

	fields_string := strings.Join(fields, ", ")
	struct_fields = append(struct_fields, id)

	sql_query := fmt.Sprintf("UPDATE address SET %s WHERE address_id = ?", fields_string)
	_, err := r.DB.ExecContext(ctx, sql_query, struct_fields...)
	if err != nil {
		return err
	}

	return nil
}

package customer

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/mmiftahrzki/customer/config"
	"github.com/mmiftahrzki/customer/database"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/stretchr/testify/assert"
)

var db *sql.DB

func init() {
	var err error = nil
	logger := logger.GetLogger()
	cfg_db := config.DatabaseConfig{
		Host:          "localhost",
		Port:          3306,
		User:          "root",
		Password:      "toor",
		Name:          "portfolio",
		MaxConnection: 10,
	}

	db, err = database.New(cfg_db)
	if err != nil {
		logger.Fatalf("Database Error: %v\n", err)
	}
}

func excpectedStr(expected, got any) string {
	return fmt.Sprintf("Expected: %v but got: %v instead.", expected, got)
}

func TestCustomerRepository(t *testing.T) {
	defer db.Close()
	repo := newRepo(db)
	var customer_id uint = 3
	ctx := context.Background()

	t.Run("read customers", func(t *testing.T) {
		customers, err := repo.SelectAll(ctx)
		assert.NoError(t, err, excpectedStr(customers, err))
	})

	t.Run("create new customer", func(t *testing.T) {
		new_customer := CustomerCreateModel{
			FirstName: "Muhammad Miftah",
			LastName:  "Rizki",
			Email:     "mmiftahr@gmail.com",
		}

		err := repo.InsertSingle(ctx, new_customer)
		assert.NoError(t, err, excpectedStr(nil, err))
	})

	t.Run("read customer by id", func(t *testing.T) {
		created_at, err := time.Parse(time.RFC3339, "2006-02-14T22:04:36Z")
		assert.NoError(t, err, excpectedStr(nil, err))

		name, offset := created_at.Zone()
		assert.Equal(t, 2006, created_at.Year(), excpectedStr(2006, created_at.Year()))
		assert.Equal(t, time.Month(2), created_at.Month(), excpectedStr(time.Month(2), created_at.Month()))
		assert.Equal(t, 14, created_at.Day(), excpectedStr(14, created_at.Day()))
		assert.Equal(t, 22, created_at.Hour(), excpectedStr(22, created_at.Hour()))
		assert.Equal(t, 04, created_at.Minute(), excpectedStr(04, created_at.Minute()))
		assert.Equal(t, 36, created_at.Second(), excpectedStr(36, created_at.Second()))
		assert.Equal(t, "UTC", name, excpectedStr("UTC", name))
		assert.Equal(t, 0*3600, offset, excpectedStr(0, offset))

		customer := CustomerReadModel{
			Id:       3,
			Email:    "LINDA.WILLIAMS@sakilacustomer.org",
			FullName: "LINDA WILLIAMS",
			// Address: address.Address{
			// 	Id:         7,
			// 	Address:    "692 Joliet Street",
			// 	Address2:   "",
			// 	District:   "Attika",
			// 	CityId:     38,
			// 	PostalCode: "83579",
			// },
			CreatedAt: created_at,
		}

		db_customer, err := repo.SelectSingleById(ctx, customer_id)
		assert.NoError(t, err, excpectedStr(nil, err))
		assert.Equal(t, customer.Id, db_customer.Id, excpectedStr(customer.Id, db_customer.Id))
		assert.Equal(t, customer.Email, db_customer.Email, excpectedStr(customer.Email, db_customer.Email))
		assert.Equal(
			t,
			customer.FullName,
			db_customer.FullName,
			excpectedStr(customer.FullName, db_customer.FullName),
		)
		assert.Equal(
			t,
			customer.CreatedAt,
			db_customer.CreatedAt,
			excpectedStr(customer.CreatedAt, db_customer.CreatedAt),
		)
	})

	t.Run("update customer by id", func(t *testing.T) {
		updated_customer := CustomerUpdateModel{
			Email:     "LINDA.WILLIAMSUpdatedviaTEST@sakilacustomer.org",
			FirstName: "LINDA",
			LastName:  "WILLIAMS Updated via TEST",
		}
		fullname := fmt.Sprintf("%s %s", updated_customer.FirstName, updated_customer.LastName)

		err := repo.UpdateSingleById(ctx, customer_id, updated_customer)
		assert.NoError(t, err, excpectedStr(nil, err))
		db_customer, err := repo.SelectSingleById(ctx, customer_id)

		assert.NoError(t, err, excpectedStr(nil, err))
		assert.Equal(
			t,
			updated_customer.Email,
			db_customer.Email,
			excpectedStr(updated_customer.Email, db_customer.Email),
		)
		assert.Equal(
			t,
			fullname,
			db_customer.FullName,
			excpectedStr(fullname, db_customer.FullName),
		)
	})

	t.Run("delete customer by id", func(t *testing.T) {
		err := repo.DeleteSingleById(ctx, customer_id)
		assert.NoError(t, err, excpectedStr(nil, err))
	})
}

package customer

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/mmiftahrzki/customer/customer/address"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

type service struct {
	repo repo
	log  *logrus.Entry
}

var ErrCustomerAlreadyExists = fmt.Errorf("customer already exists")
var ErrCustomerNotFound = fmt.Errorf("customer not found")
var ErrInvalidCustomerAddressMismatch = fmt.Errorf("customer address mismatch")

func NewService(r repo) *service {
	svc := &service{
		repo: r,
		log:  logger.GetLogger().WithField("component", "customer_service"),
	}

	return svc
}

func (svc *service) GetMultiple(ctx context.Context) (customers []customerReadModel, err error) {
	customer_sqls, err := svc.repo.SelectAll(ctx)
	if err != nil {
		return
	}

	for _, customer_sql := range customer_sqls {
		customer := NewCustomerReadModel(customer_sql)
		customers = append(customers, customer)
	}

	return
}

// func (svc *service) GetPrev(ctx context.Context, id int) (customers []customerReadModel, err error) {
// 	customer, err := svc.GetSingleById(ctx, id)
// 	if err != nil {
// 		return
// 	}

// 	if reflect.ValueOf(customer).IsZero() {
// 		return nil, errors.New("implement me")
// 	}

// 	customer_sqls, err := svc.repo.SelectAllPrev(ctx, customer)
// 	if err != nil {
// 		return
// 	}

// 	for _, customer_sql := range customer_sqls {
// 		customer := NewCustomerReadModel(customer_sql)
// 		customers = append(customers, customer)
// 	}

// 	return
// }

func (svc *service) GetMultipleNext(ctx context.Context, id int) (customers []customerReadModel, err error) {
	customer, err := svc.GetSingleById(ctx, id)
	if err != nil {
		return
	}

	if reflect.ValueOf(customer).IsZero() {
		return nil, errors.New("implement me")
	}

	customer_sqls, err := svc.repo.SelectAllNext(ctx, customer)
	if err != nil {
		return
	}

	for _, customer_sql := range customer_sqls {
		customer := NewCustomerReadModel(customer_sql)
		customers = append(customers, customer)
	}

	return
}

func (svc *service) GetSingleById(ctx context.Context, id int) (customerReadModel, error) {
	customer := customerReadModel{}
	empty_customer_sql := customerSQLModel{}

	customer_sql, err := svc.repo.SelectSingleById(ctx, id)
	if err != nil {
		return customer, err
	}

	if customer_sql == empty_customer_sql {
		return customer, ErrCustomerNotFound
	}

	customer = NewCustomerReadModel(customer_sql)

	return customer, nil
}

func (svc *service) CreateNewSingle(ctx context.Context, new_customer CustomerCreateModel) error {
	if err := svc.repo.InsertSingle(ctx, new_customer); err != nil {
		if mysql_error, ok := err.(*mysql.MySQLError); ok {
			if mysql_error.Number == 1062 {
				return ErrCustomerAlreadyExists
			}

			return mysql_error
		}

		return err
	}

	return nil
}

func (svc *service) ModifySingleById(ctx context.Context, id int, modified_customer CustomerUpdateModel) error {
	if _, err := svc.GetSingleById(ctx, id); err != nil {
		return err
	}

	err := svc.repo.UpdateSingleById(ctx, id, modified_customer)
	if err != nil {
		return err
	}

	return nil
}

func (svc *service) DeleteSingleById(ctx context.Context, id int) error {
	if err := svc.repo.DeleteSingleById(ctx, id); err != nil {
		return err
	}

	return nil
}

func (svc *service) ModifySingleAddressById(ctx context.Context, customer_id int, address_id uint16, modified_customer_address address.AddressUpdateModel) error {
	customer_sql, err := svc.repo.SelectSingleById(ctx, customer_id)
	if err != nil {
		return err
	}

	if uint16(customer_sql.address_id.Int16) != address_id {
		return ErrInvalidCustomerAddressMismatch
	}

	if err := svc.repo.UpdateSingleAddressByCustomerId(ctx, address_id, modified_customer_address); err != nil {
		return err
	}

	return nil
}

package customer

import (
	"context"
	"errors"
	"reflect"
	"sort"

	"github.com/go-sql-driver/mysql"
	"github.com/mmiftahrzki/customer/customer/address"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

type service struct {
	repo repo
	log  *logrus.Entry
}

var errCustomerAlreadyExists = errors.New("customer already exists")
var errCustomerNotFound = errors.New("customer not found")
var errInvalidCustomerAddressMismatch = errors.New("customer address mismatch")

func newService(r repo) service {
	svc := service{
		repo: r,
		log:  logger.GetLogger().WithField("component", "customer_service"),
	}

	return svc
}

func (svc *service) GetMultiple(ctx context.Context) (customers []modelRead, err error) {
	customer_sqls, err := svc.repo.SelectAll(ctx)
	if err != nil {
		return
	}

	for _, customer_sql := range customer_sqls {
		customer := newReadModel(customer_sql)
		customers = append(customers, customer)
	}

	return
}

func (svc *service) GetMultiplePrev(ctx context.Context, id int) (customers []modelRead, err error) {
	customer, err := svc.GetSingleById(ctx, id)
	if err != nil {
		return
	}

	if reflect.ValueOf(customer).IsZero() {
		return nil, errors.New("implement me")
	}

	customer_sqls, err := svc.repo.SelectAllPrev(ctx, customer)
	if err != nil {
		return
	}

	for _, customer_sql := range customer_sqls {
		customer := newReadModel(customer_sql)
		customers = append(customers, customer)
	}

	sort.SliceStable(customers, func(i, j int) bool {
		return customers[i].Id < customers[j].Id
	})

	return
}

func (svc *service) GetMultipleNext(ctx context.Context, id int) (customers []modelRead, err error) {
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
		customer := newReadModel(customer_sql)
		customers = append(customers, customer)
	}

	return
}

func (svc *service) GetSingleById(ctx context.Context, id int) (modelRead, error) {
	customer := modelRead{}
	empty_customer_sql := modelSQL{}

	customer_sql, err := svc.repo.SelectSingleById(ctx, id)
	if err != nil {
		return customer, err
	}

	if customer_sql == empty_customer_sql {
		return customer, errCustomerNotFound
	}

	customer = newReadModel(customer_sql)

	return customer, nil
}

func (svc *service) CreateNewSingle(ctx context.Context, new_customer modelCreate) error {
	if err := svc.repo.InsertSingle(ctx, new_customer); err != nil {
		if mysql_error, ok := err.(*mysql.MySQLError); ok {
			if mysql_error.Number == 1062 {
				return errCustomerAlreadyExists
			}

			return mysql_error
		}

		return err
	}

	return nil
}

func (svc *service) ModifySingleById(ctx context.Context, id int, modified_customer updateModel) error {
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

func (svc *service) ModifySingleAddressById(ctx context.Context, customer_id int, address_id uint16, modified_customer_address address.ModelUpdate) error {
	customer_sql, err := svc.repo.SelectSingleById(ctx, customer_id)
	if err != nil {
		return err
	}

	if uint16(customer_sql.addressId.Int16) != address_id {
		return errInvalidCustomerAddressMismatch
	}

	if err := svc.repo.UpdateSingleAddressByCustomerId(ctx, address_id, modified_customer_address); err != nil {
		return err
	}

	return nil
}

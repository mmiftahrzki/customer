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
		log:  logger.GetLogger().WithField("component", "customerService"),
	}

	return svc
}

func (svc *service) GetMultiple(ctx context.Context) ([]modelRead, error) {
	var customers []modelRead

	select {
	case <-ctx.Done():
		svc.log.Info("deadline exceeded from service layer")

		return customers, ctx.Err()
	default:
		customerSqls, repoErr := svc.repo.SelectAll(ctx)
		if repoErr != nil {
			return customers, repoErr
		}

		for _, customerSql := range customerSqls {
			customer := newReadModel(customerSql)
			customers = append(customers, customer)
		}
	}

	return customers, nil
}

func (svc *service) GetMultiplePrev(ctx context.Context, id int) (customers []modelRead, err error) {
	customer, err := svc.GetSingleById(ctx, id)
	if err != nil {
		return
	}

	if reflect.ValueOf(customer).IsZero() {
		return nil, errors.New("implement me")
	}

	customerSqls, err := svc.repo.SelectAllPrev(ctx, customer)
	if err != nil {
		return
	}

	for _, customerSql := range customerSqls {
		customer := newReadModel(customerSql)
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

	customerSqls, err := svc.repo.SelectAllNext(ctx, customer)
	if err != nil {
		return
	}

	for _, customerSql := range customerSqls {
		customer := newReadModel(customerSql)
		customers = append(customers, customer)
	}

	return
}

func (svc *service) GetSingleById(ctx context.Context, id int) (modelRead, error) {
	var customer modelRead
	var emptyCustomerSql modelSQL

	customerSql, repoErr := svc.repo.SelectSingleById(ctx, id)
	if repoErr != nil {
		return customer, repoErr
	}

	if customerSql == emptyCustomerSql {
		return customer, errCustomerNotFound
	}

	customer = newReadModel(customerSql)

	return customer, nil
}

func (svc *service) CreateNewSingle(ctx context.Context, newCustomer modelCreate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		repoErr := svc.repo.InsertSingle(ctx, newCustomer)
		if repoErr != nil {
			mysqlErr, ok := repoErr.(*mysql.MySQLError)
			if ok && mysqlErr.Number == 1062 {
				return errCustomerAlreadyExists
			}

			return repoErr
		}
	}

	return nil
}

func (svc *service) ModifySingleById(ctx context.Context, id int, modifiedCustomer modelUpdate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err := svc.GetSingleById(ctx, id)
		if err != nil {
			return err
		}

		return svc.repo.UpdateSingleById(ctx, id, modifiedCustomer)
	}
}

func (svc *service) DeleteSingleById(ctx context.Context, id int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return svc.repo.DeleteSingleById(ctx, id)
	}
}

func (svc *service) ModifySingleAddressById(ctx context.Context, customerId int, addressId uint16, modifiedCustomerAddress address.ModelUpdate) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		customerSql, repoErr := svc.repo.SelectSingleById(ctx, customerId)
		if repoErr != nil {
			return repoErr
		}

		if uint16(customerSql.addressId.Int16) != addressId {
			return errInvalidCustomerAddressMismatch
		}

		return svc.repo.UpdateSingleAddressByCustomerId(ctx, addressId, modifiedCustomerAddress)
	}
}

package customer

import (
	"context"
	"database/sql"
	"errors"
	"sort"

	"github.com/go-sql-driver/mysql"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

type Service interface {
	CustomerList(ctx context.Context) ([]ModelRead, error)
	GetMultiplePrev(ctx context.Context, id int) ([]ModelRead, error)
	GetMultipleNext(ctx context.Context, id int) ([]ModelRead, error)

	CustomerDetails(ctx context.Context, id int) (ModelRead, error)

	RegisterCustomer(ctx context.Context, newCustomer modelCreate) error
	ModifySingleById(ctx context.Context, id int, modifiedCustomer modelUpdate) error
	DeleteSingleById(ctx context.Context, id int) error
}

type service struct {
	repo Repo
	log  *logrus.Entry
}

func newService(r Repo) Service {
	svc := &service{
		repo: r,
		log:  logger.GetLogger().WithField("component", "customerService"),
	}

	return svc
}

func (svc *service) CustomerList(ctx context.Context) ([]ModelRead, error) {
	err := ctx.Err()
	if err != nil {
		return nil, err
	}

	customers := []ModelRead{}
	customerSqls, err := svc.repo.SelectAll(ctx)
	if err != nil {
		return customers, err
	}

	for _, customerSql := range customerSqls {
		customer := newReadModel(customerSql)
		customers = append(customers, customer)
	}

	return customers, nil
}

func (svc *service) GetMultiplePrev(ctx context.Context, id int) ([]ModelRead, error) {
	err := ctx.Err()
	if err != nil {
		return nil, err
	}

	customers := []ModelRead{}
	customer, err := svc.CustomerDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	customerSqls, err := svc.repo.SelectAllPrev(ctx, customer)
	if err != nil {
		return nil, err
	}

	for _, customerSql := range customerSqls {
		customer := newReadModel(customerSql)
		customers = append(customers, customer)
	}

	sort.SliceStable(customers, func(i, j int) bool {
		return customers[i].Id < customers[j].Id
	})

	return customers, nil
}

func (svc *service) GetMultipleNext(ctx context.Context, id int) ([]ModelRead, error) {
	err := ctx.Err()
	if err != nil {
		return nil, err
	}

	customers := []ModelRead{}
	customer, err := svc.CustomerDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	customerSqls, err := svc.repo.SelectAllNext(ctx, customer)
	if err != nil {
		return nil, err
	}

	for _, customerSql := range customerSqls {
		customer := newReadModel(customerSql)
		customers = append(customers, customer)
	}

	return customers, nil
}

func (svc *service) CustomerDetails(ctx context.Context, id int) (ModelRead, error) {
	var customer ModelRead

	err := ctx.Err()
	if err != nil {
		return customer, err
	}

	customerSql, err := svc.repo.SelectSingleById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customer, errCustomerNotFound
		}

		return customer, err
	}

	customer = newReadModel(customerSql)

	return customer, nil
}

func (svc *service) RegisterCustomer(ctx context.Context, newCustomer modelCreate) error {
	err := ctx.Err()
	if err != nil {
		return err
	}

	err = svc.repo.InsertSingle(ctx, newCustomer)
	if err != nil {
		var mysqlErr *mysql.MySQLError

		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return errCustomerAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (svc *service) ModifySingleById(ctx context.Context, id int, modifiedCustomer modelUpdate) error {
	err := ctx.Err()
	if err != nil {
		return err
	}

	customer, err := svc.CustomerDetails(ctx, id)
	if err != nil {
		return err
	}

	return svc.repo.UpdateSingleById(ctx, customer.Id, modifiedCustomer)
}

func (svc *service) DeleteSingleById(ctx context.Context, id int) error {
	err := ctx.Err()
	if err != nil {
		return err
	}

	return svc.repo.DeleteSingleById(ctx, id)
}

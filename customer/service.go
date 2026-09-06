package customer

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/mmiftahrzki/customer/logger"
	"github.com/sirupsen/logrus"
)

const LIMIT = 25
const MAX_LIMIT = 100

type Service interface {
	CustomerList(ctx context.Context, take int, page int, order OrderBy) ([]ModelRead, error)

	CustomerDetails(ctx context.Context, id int) (ModelRead, error)

	RegisterCustomer(ctx context.Context, newCustomer modelCreate) error
	ModifySingleById(ctx context.Context, id int, modifiedCustomer modelUpdate) error
	DeleteSingleById(ctx context.Context, id int) error
}

func ParseTake(s string) (int, error) {
	if s == "" {
		return 25, nil
	}

	take, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("invalid take")
	}

	if take < LIMIT {
		return 0, errors.New("invalid take: min '25'")
	}

	if take > MAX_LIMIT {
		return 0, errors.New("invalid take: max '100'")
	}

	return take, nil
}

func ParsePage(s string) (int, error) {
	if s == "" {
		return 1, nil
	}

	page, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("invalid page")
	}

	if page <= 0 {
		return 0, errors.New("invalid page: min '1'")
	}

	return page, nil
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

func (svc *service) CustomerList(ctx context.Context, take int, page int, orderBy OrderBy) ([]ModelRead, error) {
	err := ctx.Err()
	if err != nil {
		return nil, err
	}

	orderByStr := orderBy.String()
	offset := (page - 1) * take
	customers := []ModelRead{}
	customerSqls, err := svc.repo.SelectAll(ctx, orderByStr, offset, take)
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

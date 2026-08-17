package customer

import (
	"context"
	"errors"
	"testing"
)

type mockRepo struct {
	insertSingleFunc     func(context.Context, modelCreate) error
	selectAllFunc        func(context.Context) ([]modelSQL, error)
	selectAllPrevFunc    func(context.Context, ModelRead) ([]modelSQL, error)
	selectAllNextFunc    func(context.Context, ModelRead) ([]modelSQL, error)
	selectSingleByIdFunc func(context.Context, int) (modelSQL, error)
	updateSingleByIdFunc func(context.Context, int, modelUpdate) error
	deleteSingleByIdFunc func(context.Context, int) error
}

func (m *mockRepo) InsertSingle(ctx context.Context, payload modelCreate) error {
	return m.insertSingleFunc(ctx, payload)
}

func (m *mockRepo) SelectAll(ctx context.Context) ([]modelSQL, error) {
	return m.selectAllFunc(ctx)
}

func (m *mockRepo) SelectAllPrev(ctx context.Context, customer ModelRead) (modelSQLs []modelSQL, err error) {
	return m.selectAllPrevFunc(ctx, customer)
}

func (m *mockRepo) SelectAllNext(ctx context.Context, customer ModelRead) (modelSQLs []modelSQL, err error) {
	return m.selectAllNextFunc(ctx, customer)
}

func (m *mockRepo) SelectSingleById(ctx context.Context, id int) (modelSQL, error) {
	return m.selectSingleByIdFunc(ctx, id)
}

func (m *mockRepo) UpdateSingleById(ctx context.Context, id int, payload modelUpdate) error {
	return m.updateSingleByIdFunc(ctx, id, payload)
}

func (m *mockRepo) DeleteSingleById(ctx context.Context, id int) error {
	return m.deleteSingleByIdFunc(ctx, id)
}

func TestRegisterCustomer_Success(t *testing.T) {
	var repo = &mockRepo{
		insertSingleFunc: func(ctx context.Context, mc modelCreate) error { return nil },
	}
	var svc = newService(repo)

	payload := modelCreate{}

	err := svc.RegisterCustomer(t.Context(), payload)
	if err != nil {
		t.Fatalf("expected no error got: %v", err)
	}
}

func TestRegisterCustomer_DuplicateEmail(t *testing.T) {
	var repo = &mockRepo{
		insertSingleFunc: func(ctx context.Context, mc modelCreate) error { return nil },
	}
	var svc = newService(repo)

	payload := modelCreate{}

	err := svc.RegisterCustomer(t.Context(), payload)
	if !errors.Is(err, errCustomerAlreadyExists) {
		t.Fatalf("expected errCustomerAlreadyexists, got: %v", err)
	}
}

func TestRegisterCustomer_ContextCancelled(t *testing.T) {
	repo := &mockRepo{
		insertSingleFunc: func(ctx context.Context, mc modelCreate) error { return nil },
	}
	svc := newService(repo)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	payload := modelCreate{}

	err := svc.RegisterCustomer(ctx, payload)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected contextCanceled, got: %v", err)
	}
}

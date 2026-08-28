package customer

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

var repoErr = errors.New("database unavailable")

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
	repo := &mockRepo{
		insertSingleFunc: func(ctx context.Context, mc modelCreate) error { return nil },
	}
	svc := newService(repo)

	payload := modelCreate{Email: "john@example.com", FirstName: "John", LastName: "Doe"}

	err := svc.RegisterCustomer(t.Context(), payload)
	if err != nil {
		t.Fatalf("expected no error got: %v\n", err)
	}
}

func TestRegisterCustomer_DuplicateEmail(t *testing.T) {
	repo := &mockRepo{
		insertSingleFunc: func(ctx context.Context, mc modelCreate) error { return &mysql.MySQLError{Number: 1062} },
	}
	svc := newService(repo)

	payload := modelCreate{Email: "john@example.com", FirstName: "John", LastName: "Doe"}

	err := svc.RegisterCustomer(t.Context(), payload)
	if !errors.Is(err, errCustomerAlreadyExists) {
		t.Fatalf("expected errCustomerAlreadyExists, got: %v\n", err)
	}
}

func TestRegisterCustomer_RepositoryError(t *testing.T) {
	repo := &mockRepo{
		insertSingleFunc: func(ctx context.Context, mc modelCreate) error { return repoErr },
	}
	svc := newService(repo)

	payload := modelCreate{Email: "john@example.com", FirstName: "John", LastName: "Doe"}

	err := svc.RegisterCustomer(t.Context(), payload)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repoErr, got: %v\n", err)
	}
}

func TestRegisterCustomer_ContextCancelled(t *testing.T) {
	repoCalled := false
	repo := &mockRepo{
		insertSingleFunc: func(ctx context.Context, mc modelCreate) error {
			repoCalled = true

			return nil
		},
	}
	svc := newService(repo)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	payload := modelCreate{Email: "john@example.com", FirstName: "John", LastName: "Doe"}

	err := svc.RegisterCustomer(ctx, payload)

	if repoCalled {
		t.Fatal("repository should not be called when context is cancelled")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected contextCanceled, got: %v\n", err)
	}
}

func TestCustomerDetails_Success(t *testing.T) {
	expected := modelSQL{
		id:        sql.NullInt16{Int16: 1, Valid: true},
		email:     sql.NullString{String: "john@example.com", Valid: true},
		firstName: sql.NullString{String: "John", Valid: true},
		lastName:  sql.NullString{String: "Doe", Valid: true},
		active:    sql.NullBool{Bool: true, Valid: true},
		createdAt: sql.NullTime{Time: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
	}
	repo := &mockRepo{
		selectSingleByIdFunc: func(ctx context.Context, i int) (modelSQL, error) {
			return expected, nil
		},
	}
	svc := newService(repo)
	got, err := svc.CustomerDetails(t.Context(), 1)
	if err != nil {
		t.Fatalf("expected no error, got: %v\n", err)
	}

	if got.Id != int(expected.id.Int16) {
		t.Errorf("expected id: %d, got: %d", expected.id.Int16, got.Id)
	}

	if got.Email != expected.email.String {
		t.Errorf("expected email: %s, got: %s", expected.email.String, got.Email)
	}

	expectedFullName := expected.firstName.String
	if expected.lastName.String != "" {
		expectedFullName += " " + expected.lastName.String
	}

	if got.FullName != expectedFullName {
		t.Errorf("expected full name: %s, got: %s", expectedFullName, got.FullName)
	}
}

func TestCustomerDetails_NotFound(t *testing.T) {
	repo := &mockRepo{
		selectSingleByIdFunc: func(ctx context.Context, i int) (modelSQL, error) {
			return modelSQL{}, sql.ErrNoRows
		},
	}
	svc := newService(repo)
	_, err := svc.CustomerDetails(t.Context(), 1)

	if !errors.Is(err, errCustomerNotFound) {
		t.Fatalf("expected errCustomerNotFound, got: %v\n", err)
	}
}

func TestCustomerDetails_RepositoryError(t *testing.T) {
	repo := &mockRepo{
		selectSingleByIdFunc: func(ctx context.Context, i int) (modelSQL, error) {
			return modelSQL{}, repoErr
		},
	}
	svc := newService(repo)
	_, err := svc.CustomerDetails(t.Context(), 1)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repoErr, got: %v\n", err)
	}
}

func TestCustomerDetails_ContextCancelled(t *testing.T) {
	repoCalled := false
	repo := &mockRepo{
		selectSingleByIdFunc: func(ctx context.Context, i int) (modelSQL, error) {
			repoCalled = true
			model := modelSQL{}

			return model, nil
		},
	}
	svc := newService(repo)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, err := svc.CustomerDetails(ctx, 1)
	if repoCalled {
		t.Fatal("repository should not be called when context is cancelled")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected contextCanceled, got: %v\n", err)
	}
}

func TestModifySingleById_Success(t *testing.T) {
	var gotId int
	var gotPayload modelUpdate

	repo := &mockRepo{
		selectSingleByIdFunc: func(ctx context.Context, id int) (modelSQL, error) {
			return modelSQL{
				id:        sql.NullInt16{Int16: 1, Valid: true},
				email:     sql.NullString{String: "john@example.com", Valid: true},
				firstName: sql.NullString{String: "John", Valid: true},
				lastName:  sql.NullString{String: "Doe", Valid: true},
				active:    sql.NullBool{Bool: true, Valid: true},
				createdAt: sql.NullTime{Time: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			}, nil
		},
		updateSingleByIdFunc: func(ctx context.Context, id int, mu modelUpdate) error {
			gotId = id
			gotPayload = mu

			return nil
		},
	}
	svc := newService(repo)

	firstName := "Jane"
	lastName := "Doe"
	email := "jane@example.com"
	payload := modelUpdate{FirstName: &firstName, LastName: &lastName, Email: &email}

	err := svc.ModifySingleById(t.Context(), 1, payload)
	if err != nil {
		t.Fatalf("expected no error, got: %v\n", err)
	}

	if gotId != 1 {
		t.Errorf("expected update id: 1, got: %v\n", gotId)
	}

	if gotPayload.FirstName != &firstName {
		t.Errorf("expected update first name: Jane, got: %v\n", gotPayload.FirstName)
	}

	if gotPayload.LastName != &lastName {
		t.Errorf("expected update last name: Doe, got: %v\n", gotPayload.LastName)
	}

	if gotPayload.Email != &email {
		t.Errorf("expected update email: jane@example.com, got: %v\n", gotPayload.Email)
	}
}

func TestModifySingleById_CustomerNotFound(t *testing.T) {
	isUpdateCalled := false
	repo := &mockRepo{
		selectSingleByIdFunc: func(ctx context.Context, i int) (modelSQL, error) {
			return modelSQL{}, sql.ErrNoRows
		},
		updateSingleByIdFunc: func(ctx context.Context, i int, mu modelUpdate) error {
			isUpdateCalled = true
			return nil
		},
	}
	svc := newService(repo)
	err := svc.ModifySingleById(t.Context(), 1, modelUpdate{})
	if !errors.Is(err, errCustomerNotFound) {
		t.Fatalf("expected errCustomerNotFound, got: %v\n", err)
	}

	if isUpdateCalled {
		t.Fatalf("update repository should not be called when customer is not found")
	}
}

func TestModifySingleById_SelectRepositoryError(t *testing.T) {
	isUpdateCalled := false
	repo := &mockRepo{
		selectSingleByIdFunc: func(ctx context.Context, i int) (modelSQL, error) {
			return modelSQL{}, repoErr
		},
		updateSingleByIdFunc: func(ctx context.Context, i int, mu modelUpdate) error {
			isUpdateCalled = true

			return nil
		},
	}
	svc := newService(repo)
	err := svc.ModifySingleById(t.Context(), 1, modelUpdate{})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repoErr, got: %v\n", err)
	}

	if isUpdateCalled {
		t.Fatalf("update repository should not be called when select fails")
	}
}

func TestModifySingleById_UpdateRepositoryError(t *testing.T) {
	repo := &mockRepo{
		selectSingleByIdFunc: func(ctx context.Context, i int) (modelSQL, error) {
			return modelSQL{
				id:        sql.NullInt16{Int16: 1, Valid: true},
				email:     sql.NullString{String: "john@example.com", Valid: true},
				firstName: sql.NullString{String: "John", Valid: true},
				lastName:  sql.NullString{String: "Doe", Valid: true},
				active:    sql.NullBool{Bool: true, Valid: true},
				createdAt: sql.NullTime{Time: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			}, nil
		},
		updateSingleByIdFunc: func(ctx context.Context, i int, mu modelUpdate) error {
			return repoErr
		},
	}

	svc := newService(repo)
	err := svc.ModifySingleById(t.Context(), 1, modelUpdate{})

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repoErr, got: %v\n", err)
	}
}

func TestModifySingleById_ContextCancelled(t *testing.T) {
	isSelectCalled := false
	isUpdateCalled := false
	repo := &mockRepo{
		selectSingleByIdFunc: func(ctx context.Context, i int) (modelSQL, error) {
			isSelectCalled = true

			return modelSQL{}, nil
		},
		updateSingleByIdFunc: func(ctx context.Context, i int, mu modelUpdate) error {
			isUpdateCalled = true

			return nil
		},
	}
	svc := newService(repo)
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := svc.ModifySingleById(ctx, 1, modelUpdate{})
	if isSelectCalled {
		t.Fatalf("select repository should not be called when operation is cancelled\n")
	}

	if isUpdateCalled {
		t.Fatalf("update repository should not be called when operation is cancelled\n")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected contextCanceled, got: %v\n", err)
	}
}

func TestDeleteSingleById_Success(t *testing.T) {
	var gotId int
	repo := &mockRepo{
		deleteSingleByIdFunc: func(ctx context.Context, i int) error {
			gotId = i

			return nil
		},
	}
	svc := newService(repo)
	err := svc.DeleteSingleById(t.Context(), 1)
	if err != nil {
		t.Fatalf("expected no error, got: %v\n", err)
	}

	if gotId != 1 {
		t.Errorf("expected delete id: 1, got: %d", gotId)
	}
}

func TestDeleteSingleById_RepositoryError(t *testing.T) {
	repo := &mockRepo{
		deleteSingleByIdFunc: func(ctx context.Context, i int) error {
			return repoErr
		},
	}

	svc := newService(repo)
	err := svc.DeleteSingleById(t.Context(), 0)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repoErr, got: %v\n", err)
	}
}

func TestDeleteSingleById_ContextCancelled(t *testing.T) {
	repoCalled := false
	repo := &mockRepo{
		deleteSingleByIdFunc: func(ctx context.Context, i int) error {
			repoCalled = true

			return nil
		},
	}
	svc := newService(repo)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := svc.DeleteSingleById(ctx, 0)
	if repoCalled {
		t.Fatal("delete repository should not be called when operation is cancelled\n")
	}

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected contextCanceled, got: %v", err)
	}
}

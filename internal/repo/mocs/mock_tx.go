// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/common/common_utils_test.go
//
// Generated by this command:
//
//	mockgen -source=internal/services/common/common_utils_test.go -destination=internal/repo/mocs/mock_tx.go
//

// Package mock_common is a generated GoMock package.
package mock_tx

import (
	context "context"
	reflect "reflect"

	pgx "github.com/jackc/pgx/v5"
	pgconn "github.com/jackc/pgx/v5/pgconn"
	gomock "go.uber.org/mock/gomock"
)

// MockTestTx is a mock of TestTx interface.
type MockTestTx struct {
	ctrl     *gomock.Controller
	recorder *MockTestTxMockRecorder
}

// MockTestTxMockRecorder is the mock recorder for MockTestTx.
type MockTestTxMockRecorder struct {
	mock *MockTestTx
}

// NewMockTestTx creates a new mock instance.
func NewMockTestTx(ctrl *gomock.Controller) *MockTestTx {
	mock := &MockTestTx{ctrl: ctrl}
	mock.recorder = &MockTestTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTestTx) EXPECT() *MockTestTxMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockTestTx) Begin(ctx context.Context) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin", ctx)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Begin indicates an expected call of Begin.
func (mr *MockTestTxMockRecorder) Begin(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockTestTx)(nil).Begin), ctx)
}

// Commit mocks base method.
func (m *MockTestTx) Commit(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockTestTxMockRecorder) Commit(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockTestTx)(nil).Commit), ctx)
}

// Conn mocks base method.
func (m *MockTestTx) Conn() *pgx.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Conn")
	ret0, _ := ret[0].(*pgx.Conn)
	return ret0
}

// Conn indicates an expected call of Conn.
func (mr *MockTestTxMockRecorder) Conn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Conn", reflect.TypeOf((*MockTestTx)(nil).Conn))
}

// CopyFrom mocks base method.
func (m *MockTestTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CopyFrom", ctx, tableName, columnNames, rowSrc)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CopyFrom indicates an expected call of CopyFrom.
func (mr *MockTestTxMockRecorder) CopyFrom(ctx, tableName, columnNames, rowSrc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CopyFrom", reflect.TypeOf((*MockTestTx)(nil).CopyFrom), ctx, tableName, columnNames, rowSrc)
}

// Exec mocks base method.
func (m *MockTestTx) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, sql}
	for _, a := range arguments {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exec", varargs...)
	ret0, _ := ret[0].(pgconn.CommandTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exec indicates an expected call of Exec.
func (mr *MockTestTxMockRecorder) Exec(ctx, sql any, arguments ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, sql}, arguments...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockTestTx)(nil).Exec), varargs...)
}

// LargeObjects mocks base method.
func (m *MockTestTx) LargeObjects() pgx.LargeObjects {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LargeObjects")
	ret0, _ := ret[0].(pgx.LargeObjects)
	return ret0
}

// LargeObjects indicates an expected call of LargeObjects.
func (mr *MockTestTxMockRecorder) LargeObjects() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LargeObjects", reflect.TypeOf((*MockTestTx)(nil).LargeObjects))
}

// Prepare mocks base method.
func (m *MockTestTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", ctx, name, sql)
	ret0, _ := ret[0].(*pgconn.StatementDescription)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Prepare indicates an expected call of Prepare.
func (mr *MockTestTxMockRecorder) Prepare(ctx, name, sql any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockTestTx)(nil).Prepare), ctx, name, sql)
}

// Query mocks base method.
func (m *MockTestTx) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(pgx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockTestTxMockRecorder) Query(ctx, sql any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockTestTx)(nil).Query), varargs...)
}

// QueryRow mocks base method.
func (m *MockTestTx) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	m.ctrl.T.Helper()
	varargs := []any{ctx, sql}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(pgx.Row)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockTestTxMockRecorder) QueryRow(ctx, sql any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, sql}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockTestTx)(nil).QueryRow), varargs...)
}

// Rollback mocks base method.
func (m *MockTestTx) Rollback(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockTestTxMockRecorder) Rollback(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockTestTx)(nil).Rollback), ctx)
}

// SendBatch mocks base method.
func (m *MockTestTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendBatch", ctx, b)
	ret0, _ := ret[0].(pgx.BatchResults)
	return ret0
}

// SendBatch indicates an expected call of SendBatch.
func (mr *MockTestTxMockRecorder) SendBatch(ctx, b any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendBatch", reflect.TypeOf((*MockTestTx)(nil).SendBatch), ctx, b)
}

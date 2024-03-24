// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repo/auth/db.go
//
// Generated by this command:
//
//	mockgen -source=internal/repo/auth/db.go -destination=internal/repo/auth/mocks/mock_db.go
//

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	context "context"
	model "projects/fb-server/pkg/model"
	reflect "reflect"

	pgx "github.com/jackc/pgx/v5"
	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	gomock "go.uber.org/mock/gomock"
	zap "go.uber.org/zap"
)

// MockFbAuthRepo is a mock of FbAuthRepo interface.
type MockFbAuthRepo struct {
	ctrl     *gomock.Controller
	recorder *MockFbAuthRepoMockRecorder
}

// MockFbAuthRepoMockRecorder is the mock recorder for MockFbAuthRepo.
type MockFbAuthRepoMockRecorder struct {
	mock *MockFbAuthRepo
}

// NewMockFbAuthRepo creates a new mock instance.
func NewMockFbAuthRepo(ctrl *gomock.Controller) *MockFbAuthRepo {
	mock := &MockFbAuthRepo{ctrl: ctrl}
	mock.recorder = &MockFbAuthRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFbAuthRepo) EXPECT() *MockFbAuthRepoMockRecorder {
	return m.recorder
}

// BeginTx mocks base method.
func (m *MockFbAuthRepo) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeginTx", ctx, txOptions)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeginTx indicates an expected call of BeginTx.
func (mr *MockFbAuthRepoMockRecorder) BeginTx(ctx, txOptions any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginTx", reflect.TypeOf((*MockFbAuthRepo)(nil).BeginTx), ctx, txOptions)
}

// ConfirmCredentialsToken mocks base method.
func (m *MockFbAuthRepo) ConfirmCredentialsToken(ctx context.Context, tx pgx.Tx, req model.UserCredentialsRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConfirmCredentialsToken", ctx, tx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConfirmCredentialsToken indicates an expected call of ConfirmCredentialsToken.
func (mr *MockFbAuthRepoMockRecorder) ConfirmCredentialsToken(ctx, tx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConfirmCredentialsToken", reflect.TypeOf((*MockFbAuthRepo)(nil).ConfirmCredentialsToken), ctx, tx, req)
}

// ConnectDBPool mocks base method.
func (m *MockFbAuthRepo) ConnectDBPool(ctx context.Context) (*pgxpool.Pool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectDBPool", ctx)
	ret0, _ := ret[0].(*pgxpool.Pool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConnectDBPool indicates an expected call of ConnectDBPool.
func (mr *MockFbAuthRepoMockRecorder) ConnectDBPool(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectDBPool", reflect.TypeOf((*MockFbAuthRepo)(nil).ConnectDBPool), ctx)
}

// DebugLogSqlErr mocks base method.
func (m *MockFbAuthRepo) DebugLogSqlErr(q string, err error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DebugLogSqlErr", q, err)
	ret0, _ := ret[0].(error)
	return ret0
}

// DebugLogSqlErr indicates an expected call of DebugLogSqlErr.
func (mr *MockFbAuthRepoMockRecorder) DebugLogSqlErr(q, err any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DebugLogSqlErr", reflect.TypeOf((*MockFbAuthRepo)(nil).DebugLogSqlErr), q, err)
}

// DeleteRecords mocks base method.
func (m *MockFbAuthRepo) DeleteRecords(ctx context.Context, tableName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRecords", ctx, tableName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRecords indicates an expected call of DeleteRecords.
func (mr *MockFbAuthRepoMockRecorder) DeleteRecords(ctx, tableName any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRecords", reflect.TypeOf((*MockFbAuthRepo)(nil).DeleteRecords), ctx, tableName)
}

// FindUser mocks base method.
func (m *MockFbAuthRepo) FindUser(ctx context.Context, req *model.UserRequest) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", ctx, req)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUser indicates an expected call of FindUser.
func (mr *MockFbAuthRepoMockRecorder) FindUser(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockFbAuthRepo)(nil).FindUser), ctx, req)
}

// FindUserCredentials mocks base method.
func (m *MockFbAuthRepo) FindUserCredentials(ctx context.Context, req model.UserCredentialsRequest) (model.UserCredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserCredentials", ctx, req)
	ret0, _ := ret[0].(model.UserCredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserCredentials indicates an expected call of FindUserCredentials.
func (mr *MockFbAuthRepoMockRecorder) FindUserCredentials(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserCredentials", reflect.TypeOf((*MockFbAuthRepo)(nil).FindUserCredentials), ctx, req)
}

// GetLogger mocks base method.
func (m *MockFbAuthRepo) GetLogger() *zap.SugaredLogger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogger")
	ret0, _ := ret[0].(*zap.SugaredLogger)
	return ret0
}

// GetLogger indicates an expected call of GetLogger.
func (mr *MockFbAuthRepoMockRecorder) GetLogger() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogger", reflect.TypeOf((*MockFbAuthRepo)(nil).GetLogger))
}

// GetPool mocks base method.
func (m *MockFbAuthRepo) GetPool() *pgxpool.Pool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPool")
	ret0, _ := ret[0].(*pgxpool.Pool)
	return ret0
}

// GetPool indicates an expected call of GetPool.
func (mr *MockFbAuthRepoMockRecorder) GetPool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPool", reflect.TypeOf((*MockFbAuthRepo)(nil).GetPool))
}

// GetPoolConfig mocks base method.
func (m *MockFbAuthRepo) GetPoolConfig() (*pgxpool.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPoolConfig")
	ret0, _ := ret[0].(*pgxpool.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPoolConfig indicates an expected call of GetPoolConfig.
func (mr *MockFbAuthRepoMockRecorder) GetPoolConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPoolConfig", reflect.TypeOf((*MockFbAuthRepo)(nil).GetPoolConfig))
}

// GracefulShutdown mocks base method.
func (m *MockFbAuthRepo) GracefulShutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GracefulShutdown")
}

// GracefulShutdown indicates an expected call of GracefulShutdown.
func (mr *MockFbAuthRepoMockRecorder) GracefulShutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GracefulShutdown", reflect.TypeOf((*MockFbAuthRepo)(nil).GracefulShutdown))
}

// PerformUsersRequestQuery mocks base method.
func (m *MockFbAuthRepo) PerformUsersRequestQuery(req *model.UsersRequest) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PerformUsersRequestQuery", req)
	ret0, _ := ret[0].([]string)
	return ret0
}

// PerformUsersRequestQuery indicates an expected call of PerformUsersRequestQuery.
func (mr *MockFbAuthRepoMockRecorder) PerformUsersRequestQuery(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PerformUsersRequestQuery", reflect.TypeOf((*MockFbAuthRepo)(nil).PerformUsersRequestQuery), req)
}

// ResetPassword mocks base method.
func (m *MockFbAuthRepo) ResetPassword(ctx context.Context, req *model.UserCredentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockFbAuthRepoMockRecorder) ResetPassword(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockFbAuthRepo)(nil).ResetPassword), ctx, req)
}

// SanitizeString mocks base method.
func (m *MockFbAuthRepo) SanitizeString(s string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SanitizeString", s)
	ret0, _ := ret[0].(string)
	return ret0
}

// SanitizeString indicates an expected call of SanitizeString.
func (mr *MockFbAuthRepoMockRecorder) SanitizeString(s any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SanitizeString", reflect.TypeOf((*MockFbAuthRepo)(nil).SanitizeString), s)
}

// SearchUsers mocks base method.
func (m *MockFbAuthRepo) SearchUsers(ctx context.Context, req *model.UsersRequest) ([]*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUsers", ctx, req)
	ret0, _ := ret[0].([]*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchUsers indicates an expected call of SearchUsers.
func (mr *MockFbAuthRepoMockRecorder) SearchUsers(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUsers", reflect.TypeOf((*MockFbAuthRepo)(nil).SearchUsers), ctx, req)
}

// TxCreateUser mocks base method.
func (m *MockFbAuthRepo) TxCreateUser(ctx context.Context, tx pgx.Tx, u model.User) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxCreateUser", ctx, tx, u)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TxCreateUser indicates an expected call of TxCreateUser.
func (mr *MockFbAuthRepoMockRecorder) TxCreateUser(ctx, tx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxCreateUser", reflect.TypeOf((*MockFbAuthRepo)(nil).TxCreateUser), ctx, tx, u)
}

// TxNewAuthCredentials mocks base method.
func (m *MockFbAuthRepo) TxNewAuthCredentials(ctx context.Context, tx pgx.Tx, uc model.UserCredentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxNewAuthCredentials", ctx, tx, uc)
	ret0, _ := ret[0].(error)
	return ret0
}

// TxNewAuthCredentials indicates an expected call of TxNewAuthCredentials.
func (mr *MockFbAuthRepoMockRecorder) TxNewAuthCredentials(ctx, tx, uc any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxNewAuthCredentials", reflect.TypeOf((*MockFbAuthRepo)(nil).TxNewAuthCredentials), ctx, tx, uc)
}

// UpdatePassword mocks base method.
func (m *MockFbAuthRepo) UpdatePassword(ctx context.Context, tx pgx.Tx, req model.UserCredentials) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePassword", ctx, tx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePassword indicates an expected call of UpdatePassword.
func (mr *MockFbAuthRepoMockRecorder) UpdatePassword(ctx, tx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePassword", reflect.TypeOf((*MockFbAuthRepo)(nil).UpdatePassword), ctx, tx, req)
}

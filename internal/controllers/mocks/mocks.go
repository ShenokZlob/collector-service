// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"github.com/ShenokZlob/collector-service/domain"
	"github.com/ShenokZlob/collector-service/pkg/contracts"
	mock "github.com/stretchr/testify/mock"
)

// NewMockAuthUsecase creates a new instance of MockAuthUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuthUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuthUsecase {
	mock := &MockAuthUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockAuthUsecase is an autogenerated mock type for the AuthUsecase type
type MockAuthUsecase struct {
	mock.Mock
}

type MockAuthUsecase_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAuthUsecase) EXPECT() *MockAuthUsecase_Expecter {
	return &MockAuthUsecase_Expecter{mock: &_m.Mock}
}

// Login provides a mock function for the type MockAuthUsecase
func (_mock *MockAuthUsecase) Login(data *dto.LoginRequest) (string, string, *domain.ResponseErr) {
	ret := _mock.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Login")
	}

	var r0 string
	var r1 string
	var r2 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(*dto.LoginRequest) (string, string, *domain.ResponseErr)); ok {
		return returnFunc(data)
	}
	if returnFunc, ok := ret.Get(0).(func(*dto.LoginRequest) string); ok {
		r0 = returnFunc(data)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(*dto.LoginRequest) string); ok {
		r1 = returnFunc(data)
	} else {
		r1 = ret.Get(1).(string)
	}
	if returnFunc, ok := ret.Get(2).(func(*dto.LoginRequest) *domain.ResponseErr); ok {
		r2 = returnFunc(data)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*domain.ResponseErr)
		}
	}
	return r0, r1, r2
}

// MockAuthUsecase_Login_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Login'
type MockAuthUsecase_Login_Call struct {
	*mock.Call
}

// Login is a helper method to define mock.On call
//   - data
func (_e *MockAuthUsecase_Expecter) Login(data interface{}) *MockAuthUsecase_Login_Call {
	return &MockAuthUsecase_Login_Call{Call: _e.mock.On("Login", data)}
}

func (_c *MockAuthUsecase_Login_Call) Run(run func(data *dto.LoginRequest)) *MockAuthUsecase_Login_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.LoginRequest))
	})
	return _c
}

func (_c *MockAuthUsecase_Login_Call) Return(s string, s1 string, responseErr *domain.ResponseErr) *MockAuthUsecase_Login_Call {
	_c.Call.Return(s, s1, responseErr)
	return _c
}

func (_c *MockAuthUsecase_Login_Call) RunAndReturn(run func(data *dto.LoginRequest) (string, string, *domain.ResponseErr)) *MockAuthUsecase_Login_Call {
	_c.Call.Return(run)
	return _c
}

// Logout provides a mock function for the type MockAuthUsecase
func (_mock *MockAuthUsecase) Logout(token string) *domain.ResponseErr {
	ret := _mock.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for Logout")
	}

	var r0 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(string) *domain.ResponseErr); ok {
		r0 = returnFunc(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ResponseErr)
		}
	}
	return r0
}

// MockAuthUsecase_Logout_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Logout'
type MockAuthUsecase_Logout_Call struct {
	*mock.Call
}

// Logout is a helper method to define mock.On call
//   - token
func (_e *MockAuthUsecase_Expecter) Logout(token interface{}) *MockAuthUsecase_Logout_Call {
	return &MockAuthUsecase_Logout_Call{Call: _e.mock.On("Logout", token)}
}

func (_c *MockAuthUsecase_Logout_Call) Run(run func(token string)) *MockAuthUsecase_Logout_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockAuthUsecase_Logout_Call) Return(responseErr *domain.ResponseErr) *MockAuthUsecase_Logout_Call {
	_c.Call.Return(responseErr)
	return _c
}

func (_c *MockAuthUsecase_Logout_Call) RunAndReturn(run func(token string) *domain.ResponseErr) *MockAuthUsecase_Logout_Call {
	_c.Call.Return(run)
	return _c
}

// Refresh provides a mock function for the type MockAuthUsecase
func (_mock *MockAuthUsecase) Refresh(token string) (string, string, *domain.ResponseErr) {
	ret := _mock.Called(token)

	if len(ret) == 0 {
		panic("no return value specified for Refresh")
	}

	var r0 string
	var r1 string
	var r2 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(string) (string, string, *domain.ResponseErr)); ok {
		return returnFunc(token)
	}
	if returnFunc, ok := ret.Get(0).(func(string) string); ok {
		r0 = returnFunc(token)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(string) string); ok {
		r1 = returnFunc(token)
	} else {
		r1 = ret.Get(1).(string)
	}
	if returnFunc, ok := ret.Get(2).(func(string) *domain.ResponseErr); ok {
		r2 = returnFunc(token)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*domain.ResponseErr)
		}
	}
	return r0, r1, r2
}

// MockAuthUsecase_Refresh_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Refresh'
type MockAuthUsecase_Refresh_Call struct {
	*mock.Call
}

// Refresh is a helper method to define mock.On call
//   - token
func (_e *MockAuthUsecase_Expecter) Refresh(token interface{}) *MockAuthUsecase_Refresh_Call {
	return &MockAuthUsecase_Refresh_Call{Call: _e.mock.On("Refresh", token)}
}

func (_c *MockAuthUsecase_Refresh_Call) Run(run func(token string)) *MockAuthUsecase_Refresh_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockAuthUsecase_Refresh_Call) Return(s string, s1 string, responseErr *domain.ResponseErr) *MockAuthUsecase_Refresh_Call {
	_c.Call.Return(s, s1, responseErr)
	return _c
}

func (_c *MockAuthUsecase_Refresh_Call) RunAndReturn(run func(token string) (string, string, *domain.ResponseErr)) *MockAuthUsecase_Refresh_Call {
	_c.Call.Return(run)
	return _c
}

// Register provides a mock function for the type MockAuthUsecase
func (_mock *MockAuthUsecase) Register(data *dto.RegisterRequest) (string, string, *domain.ResponseErr) {
	ret := _mock.Called(data)

	if len(ret) == 0 {
		panic("no return value specified for Register")
	}

	var r0 string
	var r1 string
	var r2 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(*dto.RegisterRequest) (string, string, *domain.ResponseErr)); ok {
		return returnFunc(data)
	}
	if returnFunc, ok := ret.Get(0).(func(*dto.RegisterRequest) string); ok {
		r0 = returnFunc(data)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(*dto.RegisterRequest) string); ok {
		r1 = returnFunc(data)
	} else {
		r1 = ret.Get(1).(string)
	}
	if returnFunc, ok := ret.Get(2).(func(*dto.RegisterRequest) *domain.ResponseErr); ok {
		r2 = returnFunc(data)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*domain.ResponseErr)
		}
	}
	return r0, r1, r2
}

// MockAuthUsecase_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type MockAuthUsecase_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - data
func (_e *MockAuthUsecase_Expecter) Register(data interface{}) *MockAuthUsecase_Register_Call {
	return &MockAuthUsecase_Register_Call{Call: _e.mock.On("Register", data)}
}

func (_c *MockAuthUsecase_Register_Call) Run(run func(data *dto.RegisterRequest)) *MockAuthUsecase_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*dto.RegisterRequest))
	})
	return _c
}

func (_c *MockAuthUsecase_Register_Call) Return(s string, s1 string, responseErr *domain.ResponseErr) *MockAuthUsecase_Register_Call {
	_c.Call.Return(s, s1, responseErr)
	return _c
}

func (_c *MockAuthUsecase_Register_Call) RunAndReturn(run func(data *dto.RegisterRequest) (string, string, *domain.ResponseErr)) *MockAuthUsecase_Register_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterTelegram provides a mock function for the type MockAuthUsecase
func (_mock *MockAuthUsecase) RegisterTelegram(user *domain.User) (string, string, *domain.ResponseErr) {
	ret := _mock.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for RegisterTelegram")
	}

	var r0 string
	var r1 string
	var r2 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(*domain.User) (string, string, *domain.ResponseErr)); ok {
		return returnFunc(user)
	}
	if returnFunc, ok := ret.Get(0).(func(*domain.User) string); ok {
		r0 = returnFunc(user)
	} else {
		r0 = ret.Get(0).(string)
	}
	if returnFunc, ok := ret.Get(1).(func(*domain.User) string); ok {
		r1 = returnFunc(user)
	} else {
		r1 = ret.Get(1).(string)
	}
	if returnFunc, ok := ret.Get(2).(func(*domain.User) *domain.ResponseErr); ok {
		r2 = returnFunc(user)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*domain.ResponseErr)
		}
	}
	return r0, r1, r2
}

// MockAuthUsecase_RegisterTelegram_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterTelegram'
type MockAuthUsecase_RegisterTelegram_Call struct {
	*mock.Call
}

// RegisterTelegram is a helper method to define mock.On call
//   - user
func (_e *MockAuthUsecase_Expecter) RegisterTelegram(user interface{}) *MockAuthUsecase_RegisterTelegram_Call {
	return &MockAuthUsecase_RegisterTelegram_Call{Call: _e.mock.On("RegisterTelegram", user)}
}

func (_c *MockAuthUsecase_RegisterTelegram_Call) Run(run func(user *domain.User)) *MockAuthUsecase_RegisterTelegram_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.User))
	})
	return _c
}

func (_c *MockAuthUsecase_RegisterTelegram_Call) Return(s string, s1 string, responseErr *domain.ResponseErr) *MockAuthUsecase_RegisterTelegram_Call {
	_c.Call.Return(s, s1, responseErr)
	return _c
}

func (_c *MockAuthUsecase_RegisterTelegram_Call) RunAndReturn(run func(user *domain.User) (string, string, *domain.ResponseErr)) *MockAuthUsecase_RegisterTelegram_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCardsServicer creates a new instance of MockCardsServicer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCardsServicer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCardsServicer {
	mock := &MockCardsServicer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockCardsServicer is an autogenerated mock type for the CardsServicer type
type MockCardsServicer struct {
	mock.Mock
}

type MockCardsServicer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCardsServicer) EXPECT() *MockCardsServicer_Expecter {
	return &MockCardsServicer_Expecter{mock: &_m.Mock}
}

// AddCardToCollection provides a mock function for the type MockCardsServicer
func (_mock *MockCardsServicer) AddCardToCollection(collectionId string, card *domain.Card) *domain.ResponseErr {
	ret := _mock.Called(collectionId, card)

	if len(ret) == 0 {
		panic("no return value specified for AddCardToCollection")
	}

	var r0 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(string, *domain.Card) *domain.ResponseErr); ok {
		r0 = returnFunc(collectionId, card)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ResponseErr)
		}
	}
	return r0
}

// MockCardsServicer_AddCardToCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddCardToCollection'
type MockCardsServicer_AddCardToCollection_Call struct {
	*mock.Call
}

// AddCardToCollection is a helper method to define mock.On call
//   - collectionId
//   - card
func (_e *MockCardsServicer_Expecter) AddCardToCollection(collectionId interface{}, card interface{}) *MockCardsServicer_AddCardToCollection_Call {
	return &MockCardsServicer_AddCardToCollection_Call{Call: _e.mock.On("AddCardToCollection", collectionId, card)}
}

func (_c *MockCardsServicer_AddCardToCollection_Call) Run(run func(collectionId string, card *domain.Card)) *MockCardsServicer_AddCardToCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*domain.Card))
	})
	return _c
}

func (_c *MockCardsServicer_AddCardToCollection_Call) Return(responseErr *domain.ResponseErr) *MockCardsServicer_AddCardToCollection_Call {
	_c.Call.Return(responseErr)
	return _c
}

func (_c *MockCardsServicer_AddCardToCollection_Call) RunAndReturn(run func(collectionId string, card *domain.Card) *domain.ResponseErr) *MockCardsServicer_AddCardToCollection_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteCardFromCollection provides a mock function for the type MockCardsServicer
func (_mock *MockCardsServicer) DeleteCardFromCollection(collectionId string, card *domain.Card) *domain.ResponseErr {
	ret := _mock.Called(collectionId, card)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCardFromCollection")
	}

	var r0 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(string, *domain.Card) *domain.ResponseErr); ok {
		r0 = returnFunc(collectionId, card)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ResponseErr)
		}
	}
	return r0
}

// MockCardsServicer_DeleteCardFromCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteCardFromCollection'
type MockCardsServicer_DeleteCardFromCollection_Call struct {
	*mock.Call
}

// DeleteCardFromCollection is a helper method to define mock.On call
//   - collectionId
//   - card
func (_e *MockCardsServicer_Expecter) DeleteCardFromCollection(collectionId interface{}, card interface{}) *MockCardsServicer_DeleteCardFromCollection_Call {
	return &MockCardsServicer_DeleteCardFromCollection_Call{Call: _e.mock.On("DeleteCardFromCollection", collectionId, card)}
}

func (_c *MockCardsServicer_DeleteCardFromCollection_Call) Run(run func(collectionId string, card *domain.Card)) *MockCardsServicer_DeleteCardFromCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*domain.Card))
	})
	return _c
}

func (_c *MockCardsServicer_DeleteCardFromCollection_Call) Return(responseErr *domain.ResponseErr) *MockCardsServicer_DeleteCardFromCollection_Call {
	_c.Call.Return(responseErr)
	return _c
}

func (_c *MockCardsServicer_DeleteCardFromCollection_Call) RunAndReturn(run func(collectionId string, card *domain.Card) *domain.ResponseErr) *MockCardsServicer_DeleteCardFromCollection_Call {
	_c.Call.Return(run)
	return _c
}

// ListCardsInCollection provides a mock function for the type MockCardsServicer
func (_mock *MockCardsServicer) ListCardsInCollection(collectionId string) ([]domain.Card, *domain.ResponseErr) {
	ret := _mock.Called(collectionId)

	if len(ret) == 0 {
		panic("no return value specified for ListCardsInCollection")
	}

	var r0 []domain.Card
	var r1 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(string) ([]domain.Card, *domain.ResponseErr)); ok {
		return returnFunc(collectionId)
	}
	if returnFunc, ok := ret.Get(0).(func(string) []domain.Card); ok {
		r0 = returnFunc(collectionId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Card)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string) *domain.ResponseErr); ok {
		r1 = returnFunc(collectionId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.ResponseErr)
		}
	}
	return r0, r1
}

// MockCardsServicer_ListCardsInCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListCardsInCollection'
type MockCardsServicer_ListCardsInCollection_Call struct {
	*mock.Call
}

// ListCardsInCollection is a helper method to define mock.On call
//   - collectionId
func (_e *MockCardsServicer_Expecter) ListCardsInCollection(collectionId interface{}) *MockCardsServicer_ListCardsInCollection_Call {
	return &MockCardsServicer_ListCardsInCollection_Call{Call: _e.mock.On("ListCardsInCollection", collectionId)}
}

func (_c *MockCardsServicer_ListCardsInCollection_Call) Run(run func(collectionId string)) *MockCardsServicer_ListCardsInCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCardsServicer_ListCardsInCollection_Call) Return(cards []domain.Card, responseErr *domain.ResponseErr) *MockCardsServicer_ListCardsInCollection_Call {
	_c.Call.Return(cards, responseErr)
	return _c
}

func (_c *MockCardsServicer_ListCardsInCollection_Call) RunAndReturn(run func(collectionId string) ([]domain.Card, *domain.ResponseErr)) *MockCardsServicer_ListCardsInCollection_Call {
	_c.Call.Return(run)
	return _c
}

// SetCardCountInCollection provides a mock function for the type MockCardsServicer
func (_mock *MockCardsServicer) SetCardCountInCollection(collectionId string, card *domain.Card) *domain.ResponseErr {
	ret := _mock.Called(collectionId, card)

	if len(ret) == 0 {
		panic("no return value specified for SetCardCountInCollection")
	}

	var r0 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(string, *domain.Card) *domain.ResponseErr); ok {
		r0 = returnFunc(collectionId, card)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ResponseErr)
		}
	}
	return r0
}

// MockCardsServicer_SetCardCountInCollection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetCardCountInCollection'
type MockCardsServicer_SetCardCountInCollection_Call struct {
	*mock.Call
}

// SetCardCountInCollection is a helper method to define mock.On call
//   - collectionId
//   - card
func (_e *MockCardsServicer_Expecter) SetCardCountInCollection(collectionId interface{}, card interface{}) *MockCardsServicer_SetCardCountInCollection_Call {
	return &MockCardsServicer_SetCardCountInCollection_Call{Call: _e.mock.On("SetCardCountInCollection", collectionId, card)}
}

func (_c *MockCardsServicer_SetCardCountInCollection_Call) Run(run func(collectionId string, card *domain.Card)) *MockCardsServicer_SetCardCountInCollection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*domain.Card))
	})
	return _c
}

func (_c *MockCardsServicer_SetCardCountInCollection_Call) Return(responseErr *domain.ResponseErr) *MockCardsServicer_SetCardCountInCollection_Call {
	_c.Call.Return(responseErr)
	return _c
}

func (_c *MockCardsServicer_SetCardCountInCollection_Call) RunAndReturn(run func(collectionId string, card *domain.Card) *domain.ResponseErr) *MockCardsServicer_SetCardCountInCollection_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCollectionsServicer creates a new instance of MockCollectionsServicer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCollectionsServicer(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCollectionsServicer {
	mock := &MockCollectionsServicer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockCollectionsServicer is an autogenerated mock type for the CollectionsServicer type
type MockCollectionsServicer struct {
	mock.Mock
}

type MockCollectionsServicer_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCollectionsServicer) EXPECT() *MockCollectionsServicer_Expecter {
	return &MockCollectionsServicer_Expecter{mock: &_m.Mock}
}

// Create provides a mock function for the type MockCollectionsServicer
func (_mock *MockCollectionsServicer) Create(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr) {
	ret := _mock.Called(collection)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *domain.Collection
	var r1 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(*domain.Collection) (*domain.Collection, *domain.ResponseErr)); ok {
		return returnFunc(collection)
	}
	if returnFunc, ok := ret.Get(0).(func(*domain.Collection) *domain.Collection); ok {
		r0 = returnFunc(collection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Collection)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(*domain.Collection) *domain.ResponseErr); ok {
		r1 = returnFunc(collection)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.ResponseErr)
		}
	}
	return r0, r1
}

// MockCollectionsServicer_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockCollectionsServicer_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - collection
func (_e *MockCollectionsServicer_Expecter) Create(collection interface{}) *MockCollectionsServicer_Create_Call {
	return &MockCollectionsServicer_Create_Call{Call: _e.mock.On("Create", collection)}
}

func (_c *MockCollectionsServicer_Create_Call) Run(run func(collection *domain.Collection)) *MockCollectionsServicer_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.Collection))
	})
	return _c
}

func (_c *MockCollectionsServicer_Create_Call) Return(collection1 *domain.Collection, responseErr *domain.ResponseErr) *MockCollectionsServicer_Create_Call {
	_c.Call.Return(collection1, responseErr)
	return _c
}

func (_c *MockCollectionsServicer_Create_Call) RunAndReturn(run func(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr)) *MockCollectionsServicer_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function for the type MockCollectionsServicer
func (_mock *MockCollectionsServicer) Delete(userID string, collectionID string) *domain.ResponseErr {
	ret := _mock.Called(userID, collectionID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(string, string) *domain.ResponseErr); ok {
		r0 = returnFunc(userID, collectionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.ResponseErr)
		}
	}
	return r0
}

// MockCollectionsServicer_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockCollectionsServicer_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - userID
//   - collectionID
func (_e *MockCollectionsServicer_Expecter) Delete(userID interface{}, collectionID interface{}) *MockCollectionsServicer_Delete_Call {
	return &MockCollectionsServicer_Delete_Call{Call: _e.mock.On("Delete", userID, collectionID)}
}

func (_c *MockCollectionsServicer_Delete_Call) Run(run func(userID string, collectionID string)) *MockCollectionsServicer_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *MockCollectionsServicer_Delete_Call) Return(responseErr *domain.ResponseErr) *MockCollectionsServicer_Delete_Call {
	_c.Call.Return(responseErr)
	return _c
}

func (_c *MockCollectionsServicer_Delete_Call) RunAndReturn(run func(userID string, collectionID string) *domain.ResponseErr) *MockCollectionsServicer_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function for the type MockCollectionsServicer
func (_mock *MockCollectionsServicer) Get(collectionID string) (*domain.Collection, *domain.ResponseErr) {
	ret := _mock.Called(collectionID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *domain.Collection
	var r1 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(string) (*domain.Collection, *domain.ResponseErr)); ok {
		return returnFunc(collectionID)
	}
	if returnFunc, ok := ret.Get(0).(func(string) *domain.Collection); ok {
		r0 = returnFunc(collectionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Collection)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string) *domain.ResponseErr); ok {
		r1 = returnFunc(collectionID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.ResponseErr)
		}
	}
	return r0, r1
}

// MockCollectionsServicer_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockCollectionsServicer_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - collectionID
func (_e *MockCollectionsServicer_Expecter) Get(collectionID interface{}) *MockCollectionsServicer_Get_Call {
	return &MockCollectionsServicer_Get_Call{Call: _e.mock.On("Get", collectionID)}
}

func (_c *MockCollectionsServicer_Get_Call) Run(run func(collectionID string)) *MockCollectionsServicer_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCollectionsServicer_Get_Call) Return(collection *domain.Collection, responseErr *domain.ResponseErr) *MockCollectionsServicer_Get_Call {
	_c.Call.Return(collection, responseErr)
	return _c
}

func (_c *MockCollectionsServicer_Get_Call) RunAndReturn(run func(collectionID string) (*domain.Collection, *domain.ResponseErr)) *MockCollectionsServicer_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetAll provides a mock function for the type MockCollectionsServicer
func (_mock *MockCollectionsServicer) GetAll(userId string) ([]domain.UserCollectionRef, *domain.ResponseErr) {
	ret := _mock.Called(userId)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []domain.UserCollectionRef
	var r1 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(string) ([]domain.UserCollectionRef, *domain.ResponseErr)); ok {
		return returnFunc(userId)
	}
	if returnFunc, ok := ret.Get(0).(func(string) []domain.UserCollectionRef); ok {
		r0 = returnFunc(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.UserCollectionRef)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string) *domain.ResponseErr); ok {
		r1 = returnFunc(userId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.ResponseErr)
		}
	}
	return r0, r1
}

// MockCollectionsServicer_GetAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAll'
type MockCollectionsServicer_GetAll_Call struct {
	*mock.Call
}

// GetAll is a helper method to define mock.On call
//   - userId
func (_e *MockCollectionsServicer_Expecter) GetAll(userId interface{}) *MockCollectionsServicer_GetAll_Call {
	return &MockCollectionsServicer_GetAll_Call{Call: _e.mock.On("GetAll", userId)}
}

func (_c *MockCollectionsServicer_GetAll_Call) Run(run func(userId string)) *MockCollectionsServicer_GetAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockCollectionsServicer_GetAll_Call) Return(userCollectionRefs []domain.UserCollectionRef, responseErr *domain.ResponseErr) *MockCollectionsServicer_GetAll_Call {
	_c.Call.Return(userCollectionRefs, responseErr)
	return _c
}

func (_c *MockCollectionsServicer_GetAll_Call) RunAndReturn(run func(userId string) ([]domain.UserCollectionRef, *domain.ResponseErr)) *MockCollectionsServicer_GetAll_Call {
	_c.Call.Return(run)
	return _c
}

// Rename provides a mock function for the type MockCollectionsServicer
func (_mock *MockCollectionsServicer) Rename(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr) {
	ret := _mock.Called(collection)

	if len(ret) == 0 {
		panic("no return value specified for Rename")
	}

	var r0 *domain.Collection
	var r1 *domain.ResponseErr
	if returnFunc, ok := ret.Get(0).(func(*domain.Collection) (*domain.Collection, *domain.ResponseErr)); ok {
		return returnFunc(collection)
	}
	if returnFunc, ok := ret.Get(0).(func(*domain.Collection) *domain.Collection); ok {
		r0 = returnFunc(collection)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Collection)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(*domain.Collection) *domain.ResponseErr); ok {
		r1 = returnFunc(collection)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*domain.ResponseErr)
		}
	}
	return r0, r1
}

// MockCollectionsServicer_Rename_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rename'
type MockCollectionsServicer_Rename_Call struct {
	*mock.Call
}

// Rename is a helper method to define mock.On call
//   - collection
func (_e *MockCollectionsServicer_Expecter) Rename(collection interface{}) *MockCollectionsServicer_Rename_Call {
	return &MockCollectionsServicer_Rename_Call{Call: _e.mock.On("Rename", collection)}
}

func (_c *MockCollectionsServicer_Rename_Call) Run(run func(collection *domain.Collection)) *MockCollectionsServicer_Rename_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*domain.Collection))
	})
	return _c
}

func (_c *MockCollectionsServicer_Rename_Call) Return(collection1 *domain.Collection, responseErr *domain.ResponseErr) *MockCollectionsServicer_Rename_Call {
	_c.Call.Return(collection1, responseErr)
	return _c
}

func (_c *MockCollectionsServicer_Rename_Call) RunAndReturn(run func(collection *domain.Collection) (*domain.Collection, *domain.ResponseErr)) *MockCollectionsServicer_Rename_Call {
	_c.Call.Return(run)
	return _c
}

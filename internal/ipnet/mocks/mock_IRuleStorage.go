// Code generated by mockery v2.37.0. DO NOT EDIT.

package mocks

import (
	ipnet "github.com/petrenko-alex/api-rate-limiter/internal/ipnet"
	mock "github.com/stretchr/testify/mock"
)

// MockIRuleStorage is an autogenerated mock type for the IRuleStorage type
type MockIRuleStorage struct {
	mock.Mock
}

type MockIRuleStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIRuleStorage) EXPECT() *MockIRuleStorage_Expecter {
	return &MockIRuleStorage_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: rule
func (_m *MockIRuleStorage) Create(rule ipnet.Rule) (int, error) {
	ret := _m.Called(rule)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(ipnet.Rule) (int, error)); ok {
		return rf(rule)
	}
	if rf, ok := ret.Get(0).(func(ipnet.Rule) int); ok {
		r0 = rf(rule)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(ipnet.Rule) error); ok {
		r1 = rf(rule)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRuleStorage_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockIRuleStorage_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - rule ipnet.Rule
func (_e *MockIRuleStorage_Expecter) Create(rule interface{}) *MockIRuleStorage_Create_Call {
	return &MockIRuleStorage_Create_Call{Call: _e.mock.On("Create", rule)}
}

func (_c *MockIRuleStorage_Create_Call) Run(run func(rule ipnet.Rule)) *MockIRuleStorage_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ipnet.Rule))
	})
	return _c
}

func (_c *MockIRuleStorage_Create_Call) Return(_a0 int, _a1 error) *MockIRuleStorage_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRuleStorage_Create_Call) RunAndReturn(run func(ipnet.Rule) (int, error)) *MockIRuleStorage_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: id
func (_m *MockIRuleStorage) Delete(id int) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(int) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockIRuleStorage_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockIRuleStorage_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - id int
func (_e *MockIRuleStorage_Expecter) Delete(id interface{}) *MockIRuleStorage_Delete_Call {
	return &MockIRuleStorage_Delete_Call{Call: _e.mock.On("Delete", id)}
}

func (_c *MockIRuleStorage_Delete_Call) Run(run func(id int)) *MockIRuleStorage_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockIRuleStorage_Delete_Call) Return(_a0 error) *MockIRuleStorage_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockIRuleStorage_Delete_Call) RunAndReturn(run func(int) error) *MockIRuleStorage_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ip, ruleType
func (_m *MockIRuleStorage) Find(ip string, ruleType ipnet.RuleType) (*ipnet.Rules, error) {
	ret := _m.Called(ip, ruleType)

	var r0 *ipnet.Rules
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ipnet.RuleType) (*ipnet.Rules, error)); ok {
		return rf(ip, ruleType)
	}
	if rf, ok := ret.Get(0).(func(string, ipnet.RuleType) *ipnet.Rules); ok {
		r0 = rf(ip, ruleType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ipnet.Rules)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ipnet.RuleType) error); ok {
		r1 = rf(ip, ruleType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRuleStorage_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type MockIRuleStorage_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ip string
//   - ruleType ipnet.RuleType
func (_e *MockIRuleStorage_Expecter) Find(ip interface{}, ruleType interface{}) *MockIRuleStorage_Find_Call {
	return &MockIRuleStorage_Find_Call{Call: _e.mock.On("Find", ip, ruleType)}
}

func (_c *MockIRuleStorage_Find_Call) Run(run func(ip string, ruleType ipnet.RuleType)) *MockIRuleStorage_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(ipnet.RuleType))
	})
	return _c
}

func (_c *MockIRuleStorage_Find_Call) Return(_a0 *ipnet.Rules, _a1 error) *MockIRuleStorage_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRuleStorage_Find_Call) RunAndReturn(run func(string, ipnet.RuleType) (*ipnet.Rules, error)) *MockIRuleStorage_Find_Call {
	_c.Call.Return(run)
	return _c
}

// GetForIP provides a mock function with given fields: ip
func (_m *MockIRuleStorage) GetForIP(ip string) (*ipnet.Rules, error) {
	ret := _m.Called(ip)

	var r0 *ipnet.Rules
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*ipnet.Rules, error)); ok {
		return rf(ip)
	}
	if rf, ok := ret.Get(0).(func(string) *ipnet.Rules); ok {
		r0 = rf(ip)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ipnet.Rules)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ip)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRuleStorage_GetForIP_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetForIP'
type MockIRuleStorage_GetForIP_Call struct {
	*mock.Call
}

// GetForIP is a helper method to define mock.On call
//   - ip string
func (_e *MockIRuleStorage_Expecter) GetForIP(ip interface{}) *MockIRuleStorage_GetForIP_Call {
	return &MockIRuleStorage_GetForIP_Call{Call: _e.mock.On("GetForIP", ip)}
}

func (_c *MockIRuleStorage_GetForIP_Call) Run(run func(ip string)) *MockIRuleStorage_GetForIP_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockIRuleStorage_GetForIP_Call) Return(_a0 *ipnet.Rules, _a1 error) *MockIRuleStorage_GetForIP_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRuleStorage_GetForIP_Call) RunAndReturn(run func(string) (*ipnet.Rules, error)) *MockIRuleStorage_GetForIP_Call {
	_c.Call.Return(run)
	return _c
}

// GetForType provides a mock function with given fields: ruleType
func (_m *MockIRuleStorage) GetForType(ruleType ipnet.RuleType) (*ipnet.Rules, error) {
	ret := _m.Called(ruleType)

	var r0 *ipnet.Rules
	var r1 error
	if rf, ok := ret.Get(0).(func(ipnet.RuleType) (*ipnet.Rules, error)); ok {
		return rf(ruleType)
	}
	if rf, ok := ret.Get(0).(func(ipnet.RuleType) *ipnet.Rules); ok {
		r0 = rf(ruleType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ipnet.Rules)
		}
	}

	if rf, ok := ret.Get(1).(func(ipnet.RuleType) error); ok {
		r1 = rf(ruleType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIRuleStorage_GetForType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetForType'
type MockIRuleStorage_GetForType_Call struct {
	*mock.Call
}

// GetForType is a helper method to define mock.On call
//   - ruleType ipnet.RuleType
func (_e *MockIRuleStorage_Expecter) GetForType(ruleType interface{}) *MockIRuleStorage_GetForType_Call {
	return &MockIRuleStorage_GetForType_Call{Call: _e.mock.On("GetForType", ruleType)}
}

func (_c *MockIRuleStorage_GetForType_Call) Run(run func(ruleType ipnet.RuleType)) *MockIRuleStorage_GetForType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ipnet.RuleType))
	})
	return _c
}

func (_c *MockIRuleStorage_GetForType_Call) Return(_a0 *ipnet.Rules, _a1 error) *MockIRuleStorage_GetForType_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIRuleStorage_GetForType_Call) RunAndReturn(run func(ipnet.RuleType) (*ipnet.Rules, error)) *MockIRuleStorage_GetForType_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockIRuleStorage creates a new instance of MockIRuleStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockIRuleStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockIRuleStorage {
	mock := &MockIRuleStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

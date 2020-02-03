// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	flamingo "flamingo.me/flamingo/v3/framework/flamingo"
	mock "github.com/stretchr/testify/mock"
)

// eventSubscriber is an autogenerated mock type for the eventSubscriber type
type eventSubscriber struct {
	mock.Mock
}

// Notify provides a mock function with given fields: ctx, event
func (_m *eventSubscriber) Notify(ctx context.Context, event flamingo.Event) {
	_m.Called(ctx, event)
}
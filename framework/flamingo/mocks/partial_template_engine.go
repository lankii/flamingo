// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	io "io"

	mock "github.com/stretchr/testify/mock"
)

// PartialTemplateEngine is an autogenerated mock type for the PartialTemplateEngine type
type PartialTemplateEngine struct {
	mock.Mock
}

// RenderPartials provides a mock function with given fields: ctx, templateName, data, partials
func (_m *PartialTemplateEngine) RenderPartials(ctx context.Context, templateName string, data interface{}, partials []string) (map[string]io.Reader, error) {
	ret := _m.Called(ctx, templateName, data, partials)

	var r0 map[string]io.Reader
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, []string) map[string]io.Reader); ok {
		r0 = rf(ctx, templateName, data, partials)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]io.Reader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}, []string) error); ok {
		r1 = rf(ctx, templateName, data, partials)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
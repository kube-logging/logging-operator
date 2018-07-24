package sdkdecorator

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
)

// CallSdkFunctionWithLogging expands the sdk create/update/delete/get functions with error logging
func CallSdkFunctionWithLogging(function interface{}) func(sdk.Object, ...interface{}) {
	return func(object sdk.Object, opts ...interface{}) {
		var err error
		switch f := function.(type) {
		case func(sdk.Object) error:
			err = f(object)
		case func(sdk.Object, ...sdk.DeleteOption) error:
			var deleteOpts []sdk.DeleteOption
			for _, opt := range opts {
				deleteOpts = append(deleteOpts, opt.(sdk.DeleteOption))
			}
			err = f(object, deleteOpts...)
		default:
			logrus.Warnf("Unsupported SDK function call %T", f)
		}
		if err != nil {
			logrus.Error(err)
		}
	}
}

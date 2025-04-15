package inerr

import (
	"fmt"

	"github.com/AsaHero/whereismycity/pkg/logger"
	"github.com/AsaHero/whereismycity/pkg/utility"
	"github.com/sirupsen/logrus"
)

func Err(err error) error {
	if err == nil {
		err = fmt.Errorf("unknown error")
	}

	scope, caller, location := utility.GetFrameData(2)

	logger.Error(err.Error(), logrus.Fields{
		"scope":    scope,
		"caller":   caller,
		"location": location,
	})

	return err
}

func Newf(format string, msg ...any) error {
	scope, caller, location := utility.GetFrameData(2)

	err := fmt.Errorf(format, msg...)

	logger.Error(err.Error(), logrus.Fields{
		"scope":    scope,
		"caller":   caller,
		"location": location,
	})

	return err
}

func WithMessage(err error, format string, msg ...any) error {
	if err == nil {
		err = fmt.Errorf("empty")
	}

	scope, caller, location := utility.GetFrameData(2)

	message := fmt.Sprintf(format, msg...)
	wrappedErr := fmt.Errorf("%s: %w", message, err)

	logger.Error(wrappedErr.Error(), logrus.Fields{
		"scope":    scope,
		"caller":   caller,
		"location": location,
	})

	return wrappedErr
}

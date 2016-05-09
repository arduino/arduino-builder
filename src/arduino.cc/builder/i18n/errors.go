package i18n

import "arduino.cc/builder/constants"
import "github.com/go-errors/errors"
import "os"

func ErrorfWithLogger(logger Logger, format string, a ...interface{}) *errors.Error {
	if logger.Name() == "machine" {
		logger.Fprintln(os.Stderr, constants.LOG_LEVEL_ERROR, format, a...)
		return errors.Errorf("")
	}
	return errors.Errorf(Format(format, a...))
}

func WrapError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, 0)
}

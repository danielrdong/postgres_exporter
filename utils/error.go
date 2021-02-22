package utils

import "errors"

func CombineErr(errs ...error) error {
	var errStr string
	for _, err := range errs {
		if err != nil {
			if errStr == "" {
				errStr += err.Error()
			} else {
				errStr += "; " + err.Error()
			}
		}
	}
	if errStr == "" {
		return nil
	} else {
		return errors.New(errStr)
	}
}

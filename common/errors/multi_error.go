package errors

import (
	"strings"
)

type multiError []error // 错误的数组

func (e multiError) Error() string {
	var r strings.Builder
	r.WriteString("multierr: ")
	for _, err := range e {
		r.WriteString(err.Error())
		r.WriteString(" | ")
	}
	return r.String()
}

func Combine(maybeError ...error) error { // 错误数组为参数，输出合并后的错误
	var errs multiError
	for _, err := range maybeError {
		if err != nil {
			errs = append(errs, err) // 拼接所有错误
		}
	}
	if len(errs) == 0 { // 没有错误返回空
		return nil
	}
	return errs // 返回错误的集合
}

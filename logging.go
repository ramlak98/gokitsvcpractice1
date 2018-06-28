package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

type loggingMiddleware struct {
	logger log.Logger
	next StringService
}

func (mw loggingMiddleware) Uppercase(str string) (output string, err error){
	defer func(begin time.Time){
		_ = mw.logger.Log(
			"method", "uppercase",
			"input", str,
			"output", output,
			"err", err,
			"took", time.Since(begin).Nanoseconds(),
		)
	}(time.Now())

	output, err = mw.next.Uppercase(str)
	return
}

func (mw loggingMiddleware) Count(s string) (n int){
	defer func(begin time.Time){
		_ = mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin).Nanoseconds(),
		)
	}(time.Now())

	n = mw.next.Count(s)
	return
}
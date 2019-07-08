package gormungandr

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type Request struct {
	User     User          `form:"-"`
	Coverage string        `form:"-"` //requested coverage
	ID       uuid.UUID     `form:"-"`
	logger   *logrus.Entry `form:"-"` //logger associated to this request
}

func NewRequest() Request {
	id := uuid.NewV4()
	return Request{
		ID:     id,
		logger: logrus.WithField("request_id", id),
	}
}

func (r *Request) Logger() *logrus.Entry {
	return r.logger
}

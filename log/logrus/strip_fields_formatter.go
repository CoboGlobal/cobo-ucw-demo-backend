package logrus

import "github.com/sirupsen/logrus"

type stripFieldsFormatter struct {
	inner  logrus.Formatter
	fields []string
}

// NewStripFieldsFormatter creates a new stripFieldsFormatter with an inner formatter and list of keys.
// The list of keys are deleted from log entries, before they are formatted with the inner formatter.
func NewStripFieldsFormatter(inner logrus.Formatter, fields []string) logrus.Formatter {
	return &stripFieldsFormatter{inner, fields}
}

func (s *stripFieldsFormatter) Format(e *logrus.Entry) ([]byte, error) {
	for _, key := range s.fields {
		delete(e.Data, key)
	}
	return s.inner.Format(e)
}

package logrus

import (
	"io"
	"log"

	"github.com/sirupsen/logrus"
)

// WriterHook is a hook to handle writing to any io.Writer.
type WriterHook struct {
	levels        []logrus.Level
	formatter     logrus.Formatter
	defaultWriter io.Writer
}

// NewWriterHook returns new WriterHook.
func NewWriterHook(output io.Writer, formatter logrus.Formatter, level logrus.Level) *WriterHook {
	hook := &WriterHook{
		levels:        logrus.AllLevels[:level+1],
		formatter:     formatter,
		defaultWriter: output,
	}
	return hook
}

// Fire writes the log using the defined writer.
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	msg, err := hook.formatter.Format(entry)
	if err != nil {
		log.Println("failed to generate string for entry:", err)
		return err
	}
	_, err = hook.defaultWriter.Write(msg)
	return err
}

// Levels returns configured log levels.
func (hook *WriterHook) Levels() []logrus.Level {
	return hook.levels
}

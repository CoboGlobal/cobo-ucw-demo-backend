package logrus

import (
	"io"
	"os"
	"regexp"
	"sync"
	"time"

	"cobo-ucw-backend/internal/conf"
	"github.com/getsentry/sentry-go"
	"github.com/orandin/sentrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	DefaultLogger *logrus.Entry // the exposed logger (actually an logrus.Entry)

	defaultInnerLogger = logrus.New() // logger under the hood
	once               sync.Once
	exitHandler        func() = nil
)

// nolint:nolintlint,gochecknoinits
func init() {
	once.Do(func() {
		defaultInnerLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		})
		defaultInnerLogger.SetLevel(logrus.InfoLevel)
		DefaultLogger = logrus.NewEntry(defaultInnerLogger)
	})
}

func InitLogger(config *conf.Log, defaultFields logrus.Fields) *logrus.Logger {
	stripKeys := make([]string, 0, len(defaultFields))
	for k := range defaultFields {
		stripKeys = append(stripKeys, k)
	}

	if config == nil {
		return defaultInnerLogger
	}

	initDefaultLogger(config.Std, stripKeys)
	initFileLogger(config.JsonFile, defaultFields, true, nil)
	initFileLogger(config.StdFile, defaultFields, false, stripKeys)
	initSentryLogger(config.Sentry)
	return defaultInnerLogger
}

func initDefaultLogger(config *conf.Log_Std, stripKeys []string) {
	if config == nil {
		return
	}
	defaultInnerLogger = logrus.New()
	defaultInnerLogger.SetFormatter(
		// Strip some fields before print to stdout.
		// These fields are useful for ES log search, not useful for manual debugging
		// in a terminal.
		NewStripFieldsFormatter(
			&logrus.TextFormatter{
				ForceColors:            true,
				FullTimestamp:          true,
				DisableLevelTruncation: true,
			},
			// do not print "error" field to stdout, should be included in msg
			append([]string{"error"}, stripKeys...),
		),
	)
	// logger.SetReportCaller(true)
	defaultInnerLogger.SetLevel(ParseLogLevel(config.Level))
	if config.Enable {
		defaultInnerLogger.SetOutput(os.Stdout)
	} else {
		defaultInnerLogger.SetOutput(io.Discard)
	}
}

func fixFilename(in string, fields logrus.Fields) string {
	re := regexp.MustCompile(`\{[0-9A-Za-z_-]+\}`)
	sanitizer := regexp.MustCompile(`[^0-9A-Za-z_-]`)
	ret := re.ReplaceAllStringFunc(in, func(k string) string {
		k = k[1 : len(k)-1]
		if v, ok := fields[k].(string); ok {
			v = sanitizer.ReplaceAllString(v, "_")
			return v
		}
		return ""
	})
	return ret
}

func initFileLogger(config *conf.Log_File, filenameFields logrus.Fields, isJSON bool, stripKeys []string) {
	if config == nil || !config.Enable || config.FileName == "" {
		return
	}
	filename := fixFilename(config.FileName, filenameFields)
	// Debug("File log enable, base filename: ", filename)
	file := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    int(config.MaxSize),
		MaxBackups: int(config.MaxBackup),
		MaxAge:     int(config.MaxAge),
		Compress:   false,
	}

	var fmter logrus.Formatter
	if isJSON {
		fmter = &logrus.JSONFormatter{
			DataKey:         "extra",
			TimestampFormat: "2006-01-02T15:04:05.999999Z07:00",
		}
	} else {
		fmter = &logrus.TextFormatter{
			DisableColors:          true,
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		}
	}

	defaultInnerLogger.Hooks.Add(NewWriterHook(
		file,
		NewStripFieldsFormatter(
			fmter,
			// do not print "error" field to file, should be included in msg
			append([]string{"error"}, stripKeys...),
		),
		ParseLogLevel(config.Level),
	))
}

func initSentryLogger(config *conf.Log_Sentry) {
	if config == nil || !config.Enable || config.Dsn == "" {
		return
	}
	// Sentry init
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              config.Dsn,
		Environment:      config.Environment,
		Release:          config.Release,
		Debug:            config.Debug,
		TracesSampleRate: float64(config.TracesSampleRate),
	})
	if err != nil {
		return
	}

	levels := config.Levels
	var hookLevels []logrus.Level //nolint
	for _, level := range levels {
		hookLevels = append(hookLevels, ParseLogLevel(level))
	}
	// Add Sentrus hook
	defaultInnerLogger.AddHook(sentrus.NewHook(hookLevels))

	exitHandler := func() { sentry.Flush(5 * time.Second) }
	logrus.RegisterExitHandler(exitHandler)
}

func FreeLogger() {
	if exitHandler != nil {
		exitHandler()
	}
}

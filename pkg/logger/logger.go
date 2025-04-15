package logger

import (
	"io"
	"os"
	"sync"

	"github.com/AsaHero/whereismycity/pkg/config"
	"github.com/sirupsen/logrus"
)

// log is a private instance of logrus.Logger
var once sync.Once
var log *logrus.Logger

func Init(cfg *config.Config, logFileName string) *logrus.Logger {
	once.Do(func() {
		// init new logger
		log = logrus.New()

		// Set output format to JSON
		log.SetFormatter(&OrderedJSONFormatter{})

		// Setting up file output
		logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}

		// MultiWriter to write to stdout and file
		mw := io.MultiWriter(os.Stdout, logFile)

		// Set output sources
		log.SetOutput(mw)

		// Parsing and setting log level
		level, err := logrus.ParseLevel(cfg.LogLevel)
		if err != nil {
			log.Infof("failed to parse log level: %v \n Setting default log level (info)", err)
			log.SetLevel(logrus.InfoLevel)
		} else {
			log.SetLevel(level)
		}

	})

	return log
}

// Infoln logs a message at level Info with a new line.
func Info(args ...any) {
	fields, messages := formatMessageAndFields(args...)

	if len(fields) > 0 {
		log.WithFields(fields).Infoln(messages...)
	} else {
		log.Infoln(messages...)
	}
}

// Errorln logs a message at level Error with a new line.
func Error(args ...any) {
	fields, messages := formatMessageAndFields(args...)

	if len(fields) > 0 {
		log.WithFields(fields).Errorln(messages...)
	} else {
		log.Errorln(messages...)
	}
}

// Debugln logs a message at level Debug with a new line.
func Debug(args ...any) {
	fields, messages := formatMessageAndFields(args...)

	if len(fields) > 0 {
		log.WithFields(fields).Debugln(messages...)
	} else {
		log.Debugln(messages...)
	}
}

// Warnln logs a message at level Warn with a new line.
func Warn(args ...any) {
	fields, messages := formatMessageAndFields(args...)

	if len(fields) > 0 {
		log.WithFields(fields).Warnln(messages...)
	} else {
		log.Warnln(messages...)
	}
}

// Fatalln logs a message at level Fatal with a new line, then the process will exit with status set to 1.
func Fatal(args ...any) {
	fields, messages := formatMessageAndFields(args...)

	if len(fields) > 0 {
		log.WithFields(fields).Fatalln(messages...)
	} else {
		log.Fatalln(messages...)
	}
}

func AlertError(args ...any) {
	fields, messages := formatMessageAndFields(args...)

	fields["custom_level"] = "alert"

	log.WithFields(fields).Errorln(messages...)
}

func AlertWarn(args ...any) {
	fields, messages := formatMessageAndFields(args...)

	fields["custom_level"] = "alert"

	log.WithFields(fields).Warnln(messages...)
}

func AlertInfo(args ...any) {
	fields, messages := formatMessageAndFields(args...)

	fields["custom_level"] = "alert"

	log.WithFields(fields).Infoln(messages...)
}

func formatMessageAndFields(args ...any) (logrus.Fields, []any) {
	var (
		fields   logrus.Fields = make(logrus.Fields)
		messages []any         = make([]any, 0, len(args))
	)

	for _, arg := range args {
		switch v := arg.(type) {
		case logrus.Fields:
			for key, value := range v {
				fields[key] = value
			}
		case string:
			messages = append(messages, v)
		}
	}

	return fields, messages
}

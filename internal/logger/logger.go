package logger

import (
    "go.uber.org/zap"
)

var log *zap.Logger

func Init() {
    var err error
    log, err = zap.NewDevelopment()
    if err != nil {
        panic(err)
    }
}

func Logger() *zap.Logger {
    return log
}

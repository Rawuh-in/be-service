package logger

import (
	"context"
	"fmt"
	"os"
	"rawuh-service/internal/shared/constant"
	"rawuh-service/internal/shared/lib/utils"
	"reflect"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type LoggerConfig struct {
	Env           string
	ProductName   string
	ServiceName   string
	LogLevel      string
	LogOutput     string
	FluentbitHost string
	FluentbitPort int
	ProcessId     string
	FunctionName  string
	StartFunction *time.Time
}

type Logger struct {
	zapLog        *zap.Logger
	fluentBitHook *FluentBitHook
	loggerConfig  *LoggerConfig
}

var hostname string

func New(loggerConfig *LoggerConfig) *Logger {
	logLevel, levelErr := zap.ParseAtomicLevel(loggerConfig.LogLevel)
	if levelErr != nil {
		panic(levelErr)
	}

	hostname, _ = os.Hostname()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.MessageKey = "message"

	zapEncoder := zapcore.NewJSONEncoder(encoderCfg)

	var core zapcore.Core
	var fluentBitHook *FluentBitHook
	var err error

	if strings.EqualFold(loggerConfig.LogOutput, "elastic") {
		fluentBitHook, err = NewFluentBitHook(loggerConfig.ServiceName, loggerConfig.FluentbitHost, loggerConfig.FluentbitPort)
		if err != nil {
			panic(err)
		}

		if logLevel.Level() == zap.DebugLevel {
			core = zapcore.NewTee(
				zapcore.NewCore(zapEncoder, fluentBitHook, logLevel),
				zapcore.NewCore(zapEncoder, os.Stdout, logLevel),
			)
		} else {
			core = zapcore.NewTee(
				zapcore.NewCore(zapEncoder, os.Stdout, logLevel),
				zapcore.NewCore(zapEncoder, fluentBitHook, logLevel),
			)
		}
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(zapEncoder, os.Stdout, logLevel),
		)
	}

	zapLogger := zap.New(core)

	return &Logger{
		zapLogger,
		fluentBitHook,
		loggerConfig,
	}
}

func getProcessId(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if processId, ok := md[string(constant.ContextKeyProcessId)]; ok {
			return processId[0]
		}
	}

	return ""
}

func (l *Logger) GetZapLoggerTemplate() *zap.Logger {
	return l.zapLog
}

func (l *Logger) Sync() error {
	err := l.fluentBitHook.Sync()
	if err != nil {
		return err
	}
	return l.zapLog.Sync()
}

func (l *Logger) GetProcessIdFromLogger() string {
	return l.loggerConfig.ProcessId
}

func (l *Logger) Info(message string, params ...interface{}) {

	var metadata interface{}
	var responseTime string

	for _, param := range params {
		switch v := param.(type) {
		case time.Duration:
			responseTime = v.String()
		default:
			metadata = param
		}
	}

	l.zapLog.Info(message,
		zap.String("hostname", hostname),
		zap.String("product_name", l.loggerConfig.ProductName),
		zap.String("service_name", l.loggerConfig.ServiceName),
		zap.String("process_id", l.loggerConfig.ProcessId),
		zap.String("function_name", l.loggerConfig.FunctionName),
		zap.String("log_type", "application"),
		zap.String("response_time", responseTime),
		zap.Any("metadata", ParseMetadata(metadata)))

}

func (l *Logger) Warn(message string, params ...interface{}) {
	var metadata interface{}
	var responseTime string

	for _, param := range params {
		switch v := param.(type) {
		case time.Duration:
			responseTime = v.String()
		default:
			metadata = param
		}
	}

	l.zapLog.Warn(message,
		zap.String("hostname", hostname),
		zap.String("product_name", l.loggerConfig.ProductName),
		zap.String("service_name", l.loggerConfig.ServiceName),
		zap.String("process_id", l.loggerConfig.ProcessId),
		zap.String("function_name", l.loggerConfig.FunctionName),
		zap.String("log_type", "application"),
		zap.String("response_time", responseTime),
		zap.Any("metadata", ParseMetadata(metadata)))

}

func (l *Logger) Error(message string, errMsg error) {
	stsMsg := ""
	if errMsg != nil {
		stsMsg = "Error Message - " + status.Convert(errMsg).Message()
	}

	respTime := time.Duration(0)
	if l.loggerConfig.StartFunction != nil {
		respTime = time.Since(*l.loggerConfig.StartFunction)
	}

	l.zapLog.Error(message,
		zap.String("hostname", hostname),
		zap.String("product_name", l.loggerConfig.ProductName),
		zap.String("service_name", l.loggerConfig.ServiceName),
		zap.String("process_id", l.loggerConfig.ProcessId),
		zap.String("function_name", l.loggerConfig.FunctionName),
		zap.String("log_type", "application"),
		zap.String("response_time", respTime.String()),
		zap.Any("metadata", ParseMetadata(stsMsg)),
	)
}

func (l *Logger) Debug(message string, params ...interface{}) {
	var metadata interface{}
	var responseTime string

	for _, param := range params {
		switch v := param.(type) {
		case time.Duration:
			responseTime = v.String()
		default:
			metadata = param
		}
	}

	l.zapLog.Debug(message,
		zap.String("hostname", hostname),
		zap.String("product_name", l.loggerConfig.ProductName),
		zap.String("service_name", l.loggerConfig.ServiceName),
		zap.String("process_id", l.loggerConfig.ProcessId),
		zap.String("function_name", l.loggerConfig.FunctionName),
		zap.String("log_type", "application"),
		zap.String("response_time", responseTime),
		zap.Any("metadata", ParseMetadata(metadata)))
}

func (l *Logger) Fatal(message string, errMsg error) {
	stsMsg := "Fatal Message - " + status.Convert(errMsg).Message()

	l.zapLog.Fatal(message,
		zap.String("hostname", hostname),
		zap.String("product_name", l.loggerConfig.ProductName),
		zap.String("service_name", l.loggerConfig.ServiceName),
		zap.String("process_id", l.loggerConfig.ProcessId),
		zap.String("function_name", l.loggerConfig.FunctionName),
		zap.String("log_type", "application"),
		zap.String("response_time", ""),
		zap.Any("metadata", ParseMetadata(stsMsg)),
	)

}

// Used to write published/consumed queue message
func (l *Logger) QueueMessageInfo(queueMessage string, params ...interface{}) {
	var metadata interface{}
	var responseTime string

	for _, param := range params {
		switch v := param.(type) {
		case time.Duration:
			responseTime = v.String()
		default:
			metadata = param
		}
	}

	l.zapLog.Info(queueMessage,
		zap.String("hostname", hostname),
		zap.String("product_name", l.loggerConfig.ProductName),
		zap.String("service_name", l.loggerConfig.ServiceName),
		zap.String("process_id", l.loggerConfig.ProcessId),
		zap.String("function_name", l.loggerConfig.FunctionName),
		zap.String("log_type", "queue_message"),
		zap.String("response_time", responseTime),
		zap.Any("metadata", ParseMetadata(metadata)))

}

// Used to write published/consumed queue message
func (l *Logger) QueueMessageError(queueMessage string, params ...interface{}) {
	var metadata interface{}
	var responseTime string

	for _, param := range params {
		switch v := param.(type) {
		case time.Duration:
			responseTime = v.String()
		default:
			metadata = param
		}
	}

	l.zapLog.Error(queueMessage,
		zap.String("hostname", hostname),
		zap.String("product_name", l.loggerConfig.ProductName),
		zap.String("service_name", l.loggerConfig.ServiceName),
		zap.String("process_id", l.loggerConfig.ProcessId),
		zap.String("function_name", l.loggerConfig.FunctionName),
		zap.String("log_type", "queue_message"),
		zap.String("response_time", responseTime),
		zap.Any("metadata", ParseMetadata(metadata)))

}

func (l *Logger) StartLogger(ctx context.Context, funcName string, metadatas interface{}) (context.Context, *Logger) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
		l.zapLog.Debug("No incoming metadata found, initializing empty metadata")
	}

	processId := getProcessId(ctx)
	if processId == "" {
		processId = utils.GenerateProcessId()
	}

	md.Append(constant.ContextKeyProcessIdStr, processId)

	ctx = context.Background()

	newCtx := metadata.NewOutgoingContext(ctx, md)
	newCtx = metadata.NewIncomingContext(newCtx, md)

	t := time.Now()
	l.loggerConfig.StartFunction = &t
	l.loggerConfig.FunctionName = funcName
	l.loggerConfig.ProcessId = processId

	l.zapLog.Info("Start Function ...",
		zap.String("hostname", hostname),
		zap.String("product_name", l.loggerConfig.ProductName),
		zap.String("service_name", l.loggerConfig.ServiceName),
		zap.String("process_id", l.loggerConfig.ProcessId),
		zap.String("function_name", l.loggerConfig.FunctionName),
		zap.String("log_type", "application"),
		zap.Timep("response_time", nil),
		zap.Any("metadata", ParseMetadata(metadatas)))

	return newCtx, l
}

func ParseMetadata(metadata interface{}) interface{} {
	if metadata == nil {
		return nil
	}

	switch reflect.TypeOf(metadata).Kind() {
	case reflect.Map, reflect.Struct:
		return metadata
	default:
		return map[string]string{"default_value": fmt.Sprintf("%v", metadata)}
	}
}

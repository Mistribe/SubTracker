package logfx

import (
	"log/slog"
	"strings"

	"go.uber.org/fx/fxevent"
)

type slogFxLogger struct {
	sourceLogger *slog.Logger
}

func (l slogFxLogger) LogEvent(event fxevent.Event) {
	logger := l.sourceLogger.With("src", "fxevent")
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		logger.Debug("OnStart hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			logger.Error("OnStart hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"err", e.Err,
			)
		} else {
			logger.Debug("OnStart hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.OnStopExecuting:
		logger.Debug("OnStop hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			logger.Error("OnStop hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"err", e.Err,
			)
		} else {
			logger.Debug("OnStop hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			logger.Error("error encountered while applying options",
				"type", e.TypeName,
				"module", e.ModuleName,
				"err", e.Err,
			)
		} else {
			logger.Debug("supplied",
				"type", e.TypeName,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			logger.Debug("provided",
				"constructor", e.ConstructorName,
				"module", e.ModuleName,
				"type", rtype,
				"private", e.Private,
			)
		}
		if e.Err != nil {
			logger.Error("error encountered while providing",
				"constructor", e.ConstructorName,
				"module", e.ModuleName,
				"err", e.Err,
			)
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			logger.Debug("replaced",
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			logger.Error("error encountered while replacing",
				"module", e.ModuleName,
				"err", e.Err,
			)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			logger.Debug("decorated",
				"decorator", e.DecoratorName,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			logger.Error("error encountered while decorating",
				"decorator", e.DecoratorName,
				"module", e.ModuleName,
				"err", e.Err,
			)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		logger.Debug("invoking",
			"function", e.FunctionName,
			"module", e.ModuleName,
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			logger.Error("invoke failed",
				"function", e.FunctionName,
				"module", e.ModuleName,
				"err", e.Err,
			)
		}
	case *fxevent.Stopping:
		logger.Debug("received signal",
			"signal", strings.ToUpper(e.Signal.String()),
		)
	case *fxevent.Stopped:
		if e.Err != nil {
			logger.Error("stop failed",
				"err", e.Err,
			)
		}
	case *fxevent.RollingBack:
		logger.Error("start failed, rolling back",
			"err", e.StartErr,
		)
	case *fxevent.RolledBack:
		if e.Err != nil {
			logger.Error("rollback failed",
				"err", e.Err,
			)
		}
	case *fxevent.Started:
		if e.Err != nil {
			logger.Error("start failed",
				"err", e.Err,
			)
		} else {
			logger.Debug("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			logger.Error("custom logger initialization failed",
				"err", e.Err,
			)
		} else {
			logger.Debug("initialized custom logger",
				"function", e.ConstructorName,
			)
		}
	}
}

func NewFxLogger(logger *slog.Logger) fxevent.Logger {
	logger.Info("initializing fx logger")
	return &slogFxLogger{
		sourceLogger: logger,
	}
}

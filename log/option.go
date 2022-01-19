package log

/**
 * example
type Sentry struct{}

func (s Sentry) AddOption() zap.Option {
	err := sentry.Init(sentry.ClientOptions{Dsn: "", Environment: ""})
	if err != nil {
		fmt.Sprintf("sentry.Init: %s", err)
		return nil
	}
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(
			core,
			func(entry zapcore.Entry) error {
				if entry.Level > zapcore.InfoLevel {
					sentry.CaptureException(errors.New(entry.Message))
				}
				defer sentry.Flush(2 * time.Second)
				return nil
			})
	})
}
*/

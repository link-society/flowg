package forwarders

import (
	"context"
	"fmt"

	"log/syslog"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	"link-society.com/flowg/internal/models"
)

type syslogRuntime struct {
	config *models.ForwarderSyslogV2

	tagProg      *vm.Program
	severityProg *vm.Program
	facilityProg *vm.Program
	messageProg  *vm.Program
}

var _ Runtime = (*syslogRuntime)(nil)

var (
	severityMap = map[string]syslog.Priority{
		"emerg":   syslog.LOG_EMERG,
		"alert":   syslog.LOG_ALERT,
		"crit":    syslog.LOG_CRIT,
		"err":     syslog.LOG_ERR,
		"warning": syslog.LOG_WARNING,
		"notice":  syslog.LOG_NOTICE,
		"info":    syslog.LOG_INFO,
		"debug":   syslog.LOG_DEBUG,
	}

	facilityMap = map[string]syslog.Priority{
		"kern":     syslog.LOG_KERN,
		"user":     syslog.LOG_USER,
		"mail":     syslog.LOG_MAIL,
		"daemon":   syslog.LOG_DAEMON,
		"auth":     syslog.LOG_AUTH,
		"syslog":   syslog.LOG_SYSLOG,
		"lpr":      syslog.LOG_LPR,
		"news":     syslog.LOG_NEWS,
		"uucp":     syslog.LOG_UUCP,
		"cron":     syslog.LOG_CRON,
		"authpriv": syslog.LOG_AUTHPRIV,
		"ftp":      syslog.LOG_FTP,
		"local0":   syslog.LOG_LOCAL0,
		"local1":   syslog.LOG_LOCAL1,
		"local2":   syslog.LOG_LOCAL2,
		"local3":   syslog.LOG_LOCAL3,
		"local4":   syslog.LOG_LOCAL4,
		"local5":   syslog.LOG_LOCAL5,
		"local6":   syslog.LOG_LOCAL6,
		"local7":   syslog.LOG_LOCAL7,
	}
)

func (rt *syslogRuntime) Init(context.Context) error {
	var err error

	rt.tagProg, err = CompileDynamicField(string(rt.config.Tag))
	if err != nil {
		return fmt.Errorf("failed to compile tag field: %w", err)
	}

	rt.severityProg, err = CompileDynamicField(string(rt.config.Severity))
	if err != nil {
		return fmt.Errorf("failed to compile severity field: %w", err)
	}

	rt.facilityProg, err = CompileDynamicField(string(rt.config.Facility))
	if err != nil {
		return fmt.Errorf("failed to compile facility field: %w", err)
	}

	msg := string(rt.config.Message)
	if msg == "" {
		msg = "@expr:toJSON(log)"
	}
	rt.messageProg, err = CompileDynamicField(msg)
	if err != nil {
		return fmt.Errorf("failed to compile message field: %w", err)
	}

	return nil
}

func (rt *syslogRuntime) Close(context.Context) error {
	return nil
}

func (rt *syslogRuntime) Call(ctx context.Context, record *models.LogRecord) error {
	replyC := make(chan error, 1)
	defer close(replyC)

	go func() {
		env := map[string]any{
			"timestamp": record.Timestamp,
			"log":       record.Fields,
		}

		eval := func(prog *vm.Program, field string) (string, error) {
			out, err := expr.Run(prog, env)
			if err != nil {
				return "", fmt.Errorf("failed to evaluate %s expression: %w", field, err)
			}
			str, ok := out.(string)
			if !ok {
				return "", fmt.Errorf("%s expression did not evaluate to string", field)
			}
			return str, nil
		}

		tag, err := eval(rt.tagProg, "tag")
		if err != nil {
			replyC <- err
			return
		}

		severity, err := eval(rt.severityProg, "severity")
		if err != nil {
			replyC <- err
			return
		}

		facility, err := eval(rt.facilityProg, "facility")
		if err != nil {
			replyC <- err
			return
		}

		message, err := eval(rt.messageProg, "message")
		if err != nil {
			replyC <- err
			return
		}

		severityValue, ok := severityMap[severity]
		if !ok {
			severityValue = syslog.LOG_INFO
		}
		facilityValue, ok := facilityMap[facility]
		if !ok {
			facilityValue = syslog.LOG_USER
		}
		priority := severityValue | facilityValue

		writer, err := syslog.Dial(rt.config.Network, rt.config.Address, priority, tag)
		if err != nil {
			replyC <- fmt.Errorf("failed to dial syslog with evaluated parameters: %w", err)
			return
		}
		defer writer.Close()

		if _, err := writer.Write([]byte(message)); err != nil {
			replyC <- fmt.Errorf("failed to write syslog message: %w", err)
			return
		}

		replyC <- nil
	}()

	select {
	case <-ctx.Done():
		return nil

	case err := <-replyC:
		return err
	}
}

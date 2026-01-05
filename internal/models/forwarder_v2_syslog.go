package models

import (
	"context"
	"fmt"

	"log/syslog"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type forwarderStateSyslogV2 struct {
	tagProg      *vm.Program
	severityProg *vm.Program
	facilityProg *vm.Program
	messageProg  *vm.Program
}

type ForwarderSyslogV2 struct {
	Type     string                         `json:"type" enum:"syslog" required:"true"`
	Network  string                         `json:"network" enum:"tcp,udp" required:"true"`
	Address  string                         `json:"address" required:"true" pattern:"^(([a-zA-Z0-9.-]+)|(\\[[0-9A-Fa-f:]+\\])):[0-9]{1,5}$"`
	Tag      ForwarderSyslogV2TagField      `json:"tag" required:"true"`
	Severity ForwarderSyslogV2SeverityField `json:"severity" required:"true"`
	Facility ForwarderSyslogV2FacilityField `json:"facility" required:"true"`
	Message  ForwarderSyslogV2MessageField  `json:"message" required:"false"`

	state *forwarderStateSyslogV2
}

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

func (f *ForwarderSyslogV2) init(context.Context) error {
	var err error
	f.state = &forwarderStateSyslogV2{}

	f.state.tagProg, err = CompileDynamicField(string(f.Tag))
	if err != nil {
		return fmt.Errorf("failed to compile tag field: %w", err)
	}

	f.state.severityProg, err = CompileDynamicField(string(f.Severity))
	if err != nil {
		return fmt.Errorf("failed to compile severity field: %w", err)
	}

	f.state.facilityProg, err = CompileDynamicField(string(f.Facility))
	if err != nil {
		return fmt.Errorf("failed to compile facility field: %w", err)
	}

	msg := string(f.Message)
	if msg == "" {
		msg = "toJSON(log)"
	}
	f.state.messageProg, err = CompileDynamicField(msg)
	if err != nil {
		return fmt.Errorf("failed to compile message field: %w", err)
	}

	return nil
}

func (f *ForwarderSyslogV2) close(context.Context) error {
	return nil
}

func (f *ForwarderSyslogV2) call(ctx context.Context, record *LogRecord) error {
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

		tag, err := eval(f.state.tagProg, "tag")
		if err != nil {
			replyC <- err
			return
		}

		severity, err := eval(f.state.severityProg, "severity")
		if err != nil {
			replyC <- err
			return
		}

		facility, err := eval(f.state.facilityProg, "facility")
		if err != nil {
			replyC <- err
			return
		}

		message, err := eval(f.state.messageProg, "message")
		if err != nil {
			replyC <- err
			return
		}

		priority := severityMap[severity] | facilityMap[facility]
		writer, err := syslog.Dial(f.Network, f.Address, priority, tag)
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

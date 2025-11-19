package models

import (
	"context"
	"fmt"

	"encoding/json"
	"log/syslog"
)

type forwarderStateSyslogV2 struct {
	writer *syslog.Writer
}

type ForwarderSyslogV2 struct {
	Type     string `json:"type" enum:"syslog" required:"true"`
	Network  string `json:"network" enum:"tcp,udp" required:"true"`
	Address  string `json:"address" required:"true" pattern:"^(([a-zA-Z0-9.-]+)|(\\[[0-9A-Fa-f:]+\\])):[0-9]{1,5}$"`
	Tag      string `json:"tag" required:"true" minLength:"1"`
	Severity string `json:"severity" enum:"emerg,alert,crit,err,warning,notice,info,debug" required:"true"`
	Facility string `json:"facility" enum:"kern,user,mail,daemon,auth,syslog,lpr,news,uucp,cron,authpriv,ftp,local0,local1,local2,local3,local4,local5,local6,local7" required:"true"`

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

func (f *ForwarderSyslogV2) init(ctx context.Context) error {
	reply := make(chan error, 1)
	defer close(reply)

	go func() {
		severity := severityMap[f.Severity]
		facility := facilityMap[f.Facility]
		priority := severity | facility

		writer, err := syslog.Dial(f.Network, f.Address, priority, f.Tag)
		if err != nil {
			reply <- fmt.Errorf("failed to dial syslog: %w", err)
		}

		f.state = &forwarderStateSyslogV2{
			writer: writer,
		}

		reply <- nil
	}()

	select {
	case <-ctx.Done():
		return nil

	case err := <-reply:
		return err
	}
}

func (f *ForwarderSyslogV2) close(context.Context) error {
	if f.state != nil && f.state.writer != nil {
		return f.state.writer.Close()
	}
	return nil
}

func (f *ForwarderSyslogV2) call(ctx context.Context, record *LogRecord) error {
	reply := make(chan error, 1)
	defer close(reply)

	go func() {
		if err := json.NewEncoder(f.state.writer).Encode(record); err != nil {
			reply <- fmt.Errorf("failed to send log record to syslog: %w", err)
			return
		}

		reply <- nil
	}()

	select {
	case <-ctx.Done():
		return nil

	case err := <-reply:
		return err
	}
}

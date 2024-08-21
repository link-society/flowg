// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"time"

	"link-society.com/flowg/internal/storage"
)

type streamSearchbarProps struct {
	From        time.Time
	To          time.Time
	Filter      string
	AutoRefresh string
}

func streamSearchbar(props streamSearchbarProps) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form id=\"form_stream\" hx-get=\"\" hx-target=\"#stream_viewer\" hx-swap=\"outerHTML\" class=\"flex flex-row items-center gap-2 px-3 py-1 z-depth-1\"><div class=\"flex-grow\"><label for=\"data_stream_filter\">Filter:</label> <input id=\"data_stream_filter\" name=\"filter\" type=\"text\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(props.Filter)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/streams.templ`, Line: 30, Col: 27}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" placeholder=\"field = &#34;value&#34;\"></div><div><label for=\"data_stream_from\">From:</label> <input id=\"data_stream_from\" name=\"from\" type=\"datetime-local\" step=\"1\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(props.From.Format("2006-01-02T15:04:05"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/streams.templ`, Line: 41, Col: 55}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div><label for=\"data_stream_to\">To:</label> <input id=\"data_stream_to\" name=\"to\" type=\"datetime-local\" step=\"1\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 string
		templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(props.To.Format("2006-01-02T15:04:05"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/streams.templ`, Line: 51, Col: 53}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"></div><div><label for=\"data_stream_autorefresh\">Auto Refresh:</label> <select id=\"data_stream_autorefresh\" name=\"autorefresh\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		switch props.AutoRefresh {
		case "0":
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"0\" selected>No Auto Refresh</option> <option value=\"5\">Every 5s</option> <option value=\"10\">Every 10s</option> <option value=\"30\">Every 30s</option> <option value=\"60\">Every 1m</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case "5":
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"0\">No Auto Refresh</option> <option value=\"5\" selected>Every 5s</option> <option value=\"10\">Every 10s</option> <option value=\"30\">Every 30s</option> <option value=\"60\">Every 1m</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case "10":
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"0\">No Auto Refresh</option> <option value=\"5\">Every 5s</option> <option value=\"10\" selected>Every 10s</option> <option value=\"30\">Every 30s</option> <option value=\"60\">Every 1m</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case "30":
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"0\">No Auto Refresh</option> <option value=\"5\">Every 5s</option> <option value=\"10\">Every 10s</option> <option value=\"30\" selected>Every 30s</option> <option value=\"60\">Every 1m</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		case "60":
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"0\">No Auto Refresh</option> <option value=\"5\">Every 5s</option> <option value=\"10\">Every 10s</option> <option value=\"30\">Every 30s</option> <option value=\"60\" selected>Every 1m</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select></div><button type=\"submit\" class=\"btn waves-effect waves-light ml-5\"><i class=\"material-icons right\">search</i> Run Query</button></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func streamHistogram(data string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var5 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var5 == nil {
			templ_7745c5c3_Var5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"data_stream_histogram\" data-timeserie=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 string
		templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(data)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/streams.templ`, Line: 105, Col: 24}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"grey lighten-3 mt-1\"></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

type StreamViewerProps struct {
	LogEntries    []storage.LogEntry
	Fields        []string
	From          time.Time
	To            time.Time
	Filter        string
	AutoRefresh   string
	HistogramData string

	Notifications []string
}

func StreamViewer(props StreamViewerProps) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var7 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var7 == nil {
			templ_7745c5c3_Var7 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"stream_viewer\" class=\"col s10 h-full flex flex-col\"><div class=\"card-panel white p-0\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = streamSearchbar(streamSearchbarProps{
			From:        props.From,
			To:          props.To,
			Filter:      props.Filter,
			AutoRefresh: props.AutoRefresh,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = streamHistogram(props.HistogramData).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"\n        card-panel white\n        p-0 mb-0 h-0\n        flex-grow flex-shrink\n        overflow-auto\n      \"><table class=\"w-full table-responsive logs highlight\"><thead class=\"grey lighten-2 z-depth-1\"><tr><th class=\"text-center\">Ingested At</th>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, field := range props.Fields {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<th class=\"font-monospace\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var8 string
			templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(field)
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/streams.templ`, Line: 149, Col: 47}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</th>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</tr></thead> <tbody>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, entry := range props.LogEntries {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<tr><td class=\"font-monospace\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var9 string
			templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(entry.Timestamp.Format(time.RFC3339))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/streams.templ`, Line: 157, Col: 78}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</td>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, field := range props.Fields {
				if val, exists := entry.Fields[field]; exists {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<td class=\"font-monospace\"><p class=\"grey lighten-4 w-full px-1 m-0\">")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var10 string
					templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs(val)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/streams.templ`, Line: 162, Col: 66}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</p></td>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<td></td>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</tr>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</tbody></table></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if props.Notifications != nil {
			for _, notification := range props.Notifications {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script type=\"application/javascript\" data-message=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var11 string
				templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(notification)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/components/streams.templ`, Line: 179, Col: 34}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">\n        (() => {\n          const S = document.currentScript\n\n          document.addEventListener('htmx:load', () => {\n            setTimeout(() => {\n              M.toast({\n                html: S.dataset.message,\n                completeCallback: () => S.remove(),\n              });\n            }, 1000);\n          }, { once: true });\n        })();\n      </script>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		return templ_7745c5c3_Err
	})
}
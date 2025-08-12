import FormControl from '@mui/material/FormControl'
import InputAdornment from '@mui/material/InputAdornment'
import InputLabel from '@mui/material/InputLabel'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import TextField from '@mui/material/TextField'

import { SyslogIcon } from '@/components/icons/syslog'

import { ForwarderTypeLabelMap } from '@/lib/models/forwarder'
import {
  SyslogFacility,
  SyslogFacilityValues,
  SyslogForwarderModel,
  SyslogNetwork,
  SyslogNetworkValues,
  SyslogSeverity,
  SyslogSeverityValues,
} from '@/lib/models/forwarder/syslog'

type SyslogForwarderEditorProps = {
  config: SyslogForwarderModel
  onConfigChange: (config: SyslogForwarderModel) => void
}

export const SyslogForwarderEditor = ({
  config,
  onConfigChange,
}: SyslogForwarderEditorProps) => {
  return (
    <div
      id="container:editor.forwarders.syslog"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="mb-6 shadow">
        <TextField
          label="Forwarder Type"
          variant="outlined"
          className="w-full"
          type="text"
          value={ForwarderTypeLabelMap.syslog}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <SyslogIcon />
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <div className="flex flex-row items-center gap-3">
        <FormControl>
          <InputLabel id="label:editor.forwarders.syslog.network">
            Network
          </InputLabel>
          <Select<SyslogNetwork>
            labelId="label:editor.forwarders.syslog.network"
            id="select:editor.forwarders.syslog.network"
            value={config.network}
            label="Network"
            onChange={(e) => {
              onConfigChange({
                ...config,
                network: e.target.value as SyslogNetwork,
              })
            }}
          >
            {SyslogNetworkValues.map((t) => (
              <MenuItem
                id={`option:editor.forwarders.syslog.network.${t}`}
                key={t}
                value={t}
              >
                {t.toUpperCase()}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <TextField
          className="grow"
          id="input:editor.forwarders.syslog.address"
          label="Server Address"
          variant="outlined"
          type="text"
          value={config.address}
          onChange={(e) => {
            onConfigChange({
              ...config,
              address: e.target.value,
            })
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.syslog.tag"
        label="Tag"
        variant="outlined"
        type="text"
        value={config.tag}
        onChange={(e) => {
          onConfigChange({
            ...config,
            tag: e.target.value,
          })
        }}
      />

      <div className="flex flex-row gap-3">
        <FormControl className="grow">
          <InputLabel id="label:editor.forwarders.syslog.severity">
            Severity
          </InputLabel>
          <Select<SyslogSeverity>
            labelId="label:editor.forwarders.syslog.severity"
            id="select:editor.forwarders.syslog.severity"
            value={config.severity}
            label="Severity"
            onChange={(e) => {
              onConfigChange({
                ...config,
                severity: e.target.value as SyslogSeverity,
              })
            }}
          >
            {SyslogSeverityValues.map((t) => (
              <MenuItem
                id={`option:editor.forwarders.syslog.severity.${t}`}
                key={t}
                value={t}
              >
                {t.toUpperCase()}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <FormControl className="grow">
          <InputLabel id="label:editor.forwarders.syslog.facility">
            Facility
          </InputLabel>
          <Select<SyslogFacility>
            labelId="label:editor.forwarders.syslog.facility"
            id="select:editor.forwarders.syslog.facility"
            value={config.facility}
            label="Facility"
            onChange={(e) => {
              onConfigChange({
                ...config,
                facility: e.target.value as SyslogFacility,
              })
            }}
          >
            {SyslogFacilityValues.map((t) => (
              <MenuItem
                id={`option:editor.forwarders.syslog.facility.${t}`}
                key={t}
                value={t}
              >
                {t.toUpperCase()}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </div>
    </div>
  )
}

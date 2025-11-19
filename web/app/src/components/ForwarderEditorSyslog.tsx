import { useEffect } from 'react'

import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import * as validators from '@/lib/validators'

import ForwarderConfigSyslogModel, {
  SyslogFacility,
  SyslogFacilityValues,
  SyslogNetwork,
  SyslogNetworkValues,
  SyslogSeverity,
  SyslogSeverityValues,
} from '@/lib/models/ForwarderConfigSyslogModel'

type ForwarderEditorSyslogProps = {
  config: ForwarderConfigSyslogModel
  onConfigChange: (config: ForwarderConfigSyslogModel) => void
  onValidationChange: (valid: boolean) => void
}

const ForwarderEditorSyslog = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorSyslogProps) => {
  const [network, setNetwork] = useInput<SyslogNetwork>(config.network, [])
  const [address, setAddress] = useInput<string>(config.address, [
    validators.pattern(/^(([a-zA-Z0-9.-]+)|(\[[0-9A-Fa-f:]+\])):[0-9]{1,5}$/)
  ])
  const [tag, setTag] = useInput<string>(config.tag, [
    validators.minLength(1),
  ])
  const [severity, setSeverity] = useInput<SyslogSeverity>(config.severity, [])
  const [facility, setFacility] = useInput<SyslogFacility>(config.facility, [])

  useEffect(() => {
    const valid =
      network.valid &&
      address.valid &&
      tag.valid &&
      severity.valid &&
      facility.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'syslog',
        network: network.value,
        address: address.value,
        tag: tag.value,
        severity: severity.value,
        facility: facility.value,
      })
    }
  }, [network, address, tag, severity, facility, onValidationChange, onConfigChange])

  return (
    <div
      id="container:editor.forwarders.syslog"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="flex flex-row items-center gap-3">
        <FormControl>
          <InputLabel id="label:editor.forwarders.syslog.network">
            Network
          </InputLabel>
          <Select<SyslogNetwork>
            labelId="label:editor.forwarders.syslog.network"
            id="select:editor.forwarders.syslog.network"
            value={network.value}
            label="Network"
            onChange={(e) => { setNetwork(e.target.value as SyslogNetwork) }}
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
          error={!address.valid}
          value={address.value}
          onChange={(e) => { setAddress(e.target.value) }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.syslog.tag"
        label="Tag"
        variant="outlined"
        type="text"
        error={!tag.valid}
        value={tag.value}
        onChange={(e) => { setTag(e.target.value) }}
      />

      <div className="flex flex-row gap-3">
        <FormControl className="grow">
          <InputLabel id="label:editor.forwarders.syslog.severity">
            Severity
          </InputLabel>
          <Select<SyslogSeverity>
            labelId="label:editor.forwarders.syslog.severity"
            id="select:editor.forwarders.syslog.severity"
            value={severity.value}
            label="Severity"
            onChange={(e) => { setSeverity(e.target.value as SyslogSeverity) }}
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
            value={facility.value}
            label="Facility"
            onChange={(e) => { setFacility(e.target.value as SyslogFacility) }}
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

export default ForwarderEditorSyslog

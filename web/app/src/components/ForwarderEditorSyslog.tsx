import { useEffect } from 'react'

import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import { DynamicField } from '@/lib/models/DynamicField'
import ForwarderConfigSyslogModel, {
  SyslogFacility,
  SyslogFacilityValues,
  SyslogNetwork,
  SyslogNetworkValues,
  SyslogSeverity,
  SyslogSeverityValues,
} from '@/lib/models/ForwarderConfigSyslogModel'

import * as validators from '@/lib/validators'

import DynamicFieldControl from '@/components/DynamicFieldControl'

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
    validators.pattern(/^(([a-zA-Z0-9.-]+)|(\[[0-9A-Fa-f:]+\])):[0-9]{1,5}$/),
  ])
  const [tag, setTag] = useInput<DynamicField<string>>(config.tag, [
    validators.dynamicField([validators.minLength(1)]),
  ])
  const [severity, setSeverity] = useInput<SyslogSeverity>(config.severity, [
    validators.dynamicField([]),
  ])
  const [facility, setFacility] = useInput<SyslogFacility>(config.facility, [
    validators.dynamicField([]),
  ])
  const [message, setMessage] = useInput<DynamicField<string>>(config.message, [
    validators.dynamicField([]),
  ])

  useEffect(() => {
    const valid =
      network.valid &&
      address.valid &&
      tag.valid &&
      severity.valid &&
      facility.valid &&
      message.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'syslog',
        network: network.value,
        address: address.value,
        tag: tag.value,
        severity: severity.value,
        facility: facility.value,
        message: message.value,
      })
    }
  }, [network, address, tag, severity, facility])

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
            onChange={(e) => {
              setNetwork(e.target.value as SyslogNetwork)
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
          error={!address.valid}
          value={address.value}
          onChange={(e) => {
            setAddress(e.target.value)
          }}
        />
      </div>

      <DynamicFieldControl
        id="input:editor.forwarders.syslog.tag"
        label="Tag"
        variant="outlined"
        type="text"
        error={!tag.valid}
        value={tag.value}
        onChange={(value) => {
          setTag(value)
        }}
      />

      <div className="flex flex-row gap-3">
        <FormControl className="grow">
          <DynamicFieldControl
            id="select:editor.forwarders.syslog.severity"
            label="Severity"
            variant="outlined"
            select
            error={!severity.valid}
            value={severity.value}
            onChange={(value) => {
              setSeverity(value)
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
          </DynamicFieldControl>
        </FormControl>

        <FormControl className="grow">
          <DynamicFieldControl
            id="select:editor.forwarders.syslog.facility"
            label="Facility"
            variant="outlined"
            select
            error={!facility.valid}
            value={facility.value}
            onChange={(value) => {
              setFacility(value)
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
          </DynamicFieldControl>
        </FormControl>
      </div>

      <DynamicFieldControl
        id="input:editor.forwarders.syslog.message"
        label="Message"
        variant="outlined"
        error={!message.valid}
        value={message.value}
        onChange={(value) => {
          setMessage(value)
        }}
      />
    </div>
  )
}

export default ForwarderEditorSyslog

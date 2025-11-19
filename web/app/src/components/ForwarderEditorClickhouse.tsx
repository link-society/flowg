import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormGroup from '@mui/material/FormGroup'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import ForwarderConfigClickhouseModel from '@/lib/models/ForwarderConfigClickhouseModel'
import { ForwarderConfigTypeLabelMap } from '@/lib/models/ForwarderConfigModel'

import ForwarderIconClickhouse from '@/components/ForwarderIconClickhouse'

type ForwarderEditorClickhouseProps = {
  config: ForwarderConfigClickhouseModel
  onConfigChange: (config: ForwarderConfigClickhouseModel) => void
}

const ForwarderEditorClickhouse = ({
  config,
  onConfigChange,
}: ForwarderEditorClickhouseProps) => {
  return (
    <div
      id="container:editor.forwarders.clickhouse"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="mb-6 shadow">
        <TextField
          label="Forwarder Type"
          variant="outlined"
          className="w-full"
          type="text"
          value={ForwarderConfigTypeLabelMap.clickhouse}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <ForwarderIconClickhouse />
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.clickhouse.address"
        label="Clickhouse Connection Address"
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

      <TextField
        id="input:editor.forwarders.clickhouse.db"
        label="Database Name"
        variant="outlined"
        type="text"
        value={config.db}
        onChange={(e) => {
          onConfigChange({
            ...config,
            db: e.target.value,
          })
        }}
      />

      <TextField
        id="input:editor.forwarders.clickhouse.table"
        label="Table Name"
        variant="outlined"
        type="text"
        value={config.table}
        onChange={(e) => {
          onConfigChange({
            ...config,
            table: e.target.value,
          })
        }}
      />

      <TextField
        id="input:editor.forwarders.clickhouse.user"
        label="Database Username"
        variant="outlined"
        type="text"
        value={config.user}
        onChange={(e) => {
          onConfigChange({
            ...config,
            user: e.target.value,
          })
        }}
      />

      <TextField
        id="input:editor.forwarders.clickhouse.pass"
        label="Database Password"
        variant="outlined"
        type="password"
        value={config.pass}
        onChange={(e) => {
          onConfigChange({
            ...config,
            pass: e.target.value,
          })
        }}
      />

      <FormGroup>
        <FormControlLabel
          control={
            <Checkbox
              id="input:editor.forwarders.clickhouse.tsl"
              checked={config.tls}
              onChange={(e) => {
                onConfigChange({
                  ...config,
                  tls: e.target.checked,
                })
              }}
            />
          }
          label="Use TLS"
        />
      </FormGroup>
    </div>
  )
}

export default ForwarderEditorClickhouse

import Checkbox from '@mui/material/Checkbox'
import Divider from '@mui/material/Divider'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormGroup from '@mui/material/FormGroup'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import { ForwarderTypeLabelMap } from '@/lib/models/forwarder'
import { ClickhouseForwarderModel } from '@/lib/models/forwarder/clickhouse'

// import { ClickhouseIcon } from '@/components/icons/clickhouse'

type ClickhouseForwarderEditorProps = {
  config: ClickhouseForwarderModel
  onConfigChange: (config: ClickhouseForwarderModel) => void
}

export const ClickhouseForwarderEditor = ({
  config,
  onConfigChange,
}: ClickhouseForwarderEditorProps) => {
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
          value={ForwarderTypeLabelMap.clickhouse}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  {/* <ClickhouseIcon /> */}
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.clickhouse.url"
        label="Clickhouse connection URL"
        variant="outlined"
        type="text"
        value={config.url}
        onChange={(e) => {
          onConfigChange({
            ...config,
            url: e.target.value,
          })
        }}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.clickhouse.db"
        label="Database name"
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

      <Divider />

      <TextField
        id="input:editor.forwarders.clickhouse.table"
        label="Table name"
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

      <Divider />

      <TextField
        id="input:editor.forwarders.clickhouse.user"
        label="Database username"
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

      <Divider />

      <TextField
        id="input:editor.forwarders.clickhouse.pass"
        label="Database password"
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

      <Divider />

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

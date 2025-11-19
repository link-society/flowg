import { useEffect } from 'react'

import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormGroup from '@mui/material/FormGroup'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import ForwarderConfigClickhouseModel from '@/lib/models/ForwarderConfigClickhouseModel'

import * as validators from '@/lib/validators'

type ForwarderEditorClickhouseProps = {
  config: ForwarderConfigClickhouseModel
  onConfigChange: (config: ForwarderConfigClickhouseModel) => void
  onValidationChange: (valid: boolean) => void
}

const ForwarderEditorClickhouse = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorClickhouseProps) => {
  const [address, setAddress] = useInput(config.address, [
    validators.pattern(/^(([a-zA-Z0-9.-]+)|(\[[0-9A-Fa-f:]+\])):[0-9]{1,5}$/),
  ])
  const [db, setDb] = useInput(config.db, [validators.minLength(1)])
  const [table, setTable] = useInput(config.table, [
    validators.minLength(1),
    validators.maxLength(64),
    validators.pattern(/^[a-zA-Z_][a-zA-Z0-9_]*$/),
  ])
  const [user, setUser] = useInput(config.user, [validators.minLength(1)])
  const [pass, setPass] = useInput(config.pass, [validators.minLength(1)])
  const [tls, setTls] = useInput(config.tls, [])

  useEffect(() => {
    const valid =
      address.valid && db.valid && table.valid && user.valid && pass.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'clickhouse',
        address: address.value,
        db: db.value,
        table: table.value,
        user: user.value,
        pass: pass.value,
        tls: tls.value,
      })
    }
  }, [address, db, table, user, pass, tls, onValidationChange, onConfigChange])

  return (
    <div
      id="container:editor.forwarders.clickhouse"
      className="flex flex-col items-stretch gap-3"
    >
      <TextField
        id="input:editor.forwarders.clickhouse.address"
        label="Clickhouse Connection Address"
        variant="outlined"
        type="text"
        error={!address.valid}
        value={address.value}
        onChange={(e) => {
          setAddress(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.clickhouse.db"
        label="Database Name"
        variant="outlined"
        type="text"
        error={!db.valid}
        value={db.value}
        onChange={(e) => {
          setDb(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.clickhouse.table"
        label="Table Name"
        variant="outlined"
        type="text"
        error={!table.valid}
        value={table.value}
        onChange={(e) => {
          setTable(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.clickhouse.user"
        label="Database Username"
        variant="outlined"
        type="text"
        error={!user.valid}
        value={user.value}
        onChange={(e) => {
          setUser(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.clickhouse.pass"
        label="Database Password"
        variant="outlined"
        type="password"
        error={!pass.valid}
        value={pass.value}
        onChange={(e) => {
          setPass(e.target.value)
        }}
      />

      <FormGroup>
        <FormControlLabel
          control={
            <Checkbox
              id="input:editor.forwarders.clickhouse.tsl"
              checked={tls.value}
              onChange={(e) => {
                setTls(e.target.checked)
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

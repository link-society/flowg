import { useEffect } from 'react'
import { useTranslation } from 'react-i18next'

import Checkbox from '@mui/material/Checkbox'
import FormControlLabel from '@mui/material/FormControlLabel'
import FormGroup from '@mui/material/FormGroup'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import * as validators from '@/lib/validators'

import { ForwarderEditorClickhouseRoot } from './styles'
import { ForwarderEditorClickhouseProps } from './types'

const ForwarderEditorClickhouse = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorClickhouseProps) => {
  const { t } = useTranslation()
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
  }, [address, db, table, user, pass, tls])

  return (
    <ForwarderEditorClickhouseRoot id="container:editor.forwarders.clickhouse">
      <TextField
        id="input:editor.forwarders.clickhouse.address"
        label={t('components.forwarderEditorClickhouse.addressLabel')}
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
        label={t('components.forwarderEditorClickhouse.dbLabel')}
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
        label={t('components.forwarderEditorClickhouse.tableLabel')}
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
        label={t('components.forwarderEditorClickhouse.userLabel')}
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
        label={t('components.forwarderEditorClickhouse.passLabel')}
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
          label={t('components.forwarderEditorClickhouse.tlsLabel')}
        />
      </FormGroup>
    </ForwarderEditorClickhouseRoot>
  )
}

export default ForwarderEditorClickhouse

import { useEffect } from 'react'

import Divider from '@mui/material/Divider'
import TextField from '@mui/material/TextField'

import { useInput } from '@/lib/hooks/input'

import ForwarderConfigElasticModel from '@/lib/models/ForwarderConfigElasticModel'

import * as validators from '@/lib/validators'

import InputList from '@/components/InputList'

type ForwarderEditorElasticProps = {
  config: ForwarderConfigElasticModel
  onConfigChange: (config: ForwarderConfigElasticModel) => void
  onValidationChange: (valid: boolean) => void
}

const ForwarderEditorElastic = ({
  config,
  onConfigChange,
  onValidationChange,
}: ForwarderEditorElasticProps) => {
  const [index, setIndex] = useInput(config.index, [validators.minLength(1)])
  const [username, setUsername] = useInput(config.username, [
    validators.minLength(1),
  ])
  const [password, setPassword] = useInput(config.password, [
    validators.minLength(1),
  ])
  const [addresses, setAddresses] = useInput(config.addresses, [
    validators.minItems(1),
    validators.items([validators.minLength(1), validators.formatUri]),
  ])
  const [ca, setCa] = useInput<string | undefined>(config.ca, [
    (value) => {
      if (value === undefined || value === '') {
        return true
      }

      return validators.minLength(1)(value)
    },
  ])

  useEffect(() => {
    const valid =
      index.valid &&
      username.valid &&
      password.valid &&
      addresses.valid &&
      ca.valid
    onValidationChange(valid)

    if (valid) {
      onConfigChange({
        type: 'elastic',
        index: index.value,
        username: username.value,
        password: password.value,
        addresses: addresses.value,
        ca: ca.value || undefined,
      })
    }
  }, [index, username, password, addresses, ca])

  return (
    <div
      id="container:editor.forwarders.elastic"
      className="flex flex-col items-stretch gap-3"
    >
      <TextField
        id="input:editor.forwarders.elastic.index"
        label="Index"
        variant="outlined"
        type="text"
        error={!index.valid}
        value={index.value}
        onChange={(e) => {
          setIndex(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.elastic.username"
        label="Username"
        variant="outlined"
        type="text"
        error={!username.valid}
        value={username.value}
        onChange={(e) => {
          setUsername(e.target.value)
        }}
      />

      <TextField
        id="input:editor.forwarders.elastic.password"
        label="Password"
        variant="outlined"
        type="password"
        error={!password.valid}
        value={password.value}
        onChange={(e) => {
          setPassword(e.target.value)
        }}
      />

      <Divider />

      <InputList
        id="editor.forwarders.elastic.addresses"
        itemLabel="Address"
        items={config.addresses}
        itemValidators={[validators.minLength(1), validators.formatUri]}
        onChange={setAddresses}
      />

      <Divider />

      <TextField
        id="input:editor.forwarders.elastic.ca"
        label="CA Certificate"
        variant="outlined"
        type="text"
        multiline
        maxRows={8}
        rows={8}
        error={!ca.valid}
        value={ca.value}
        onChange={(e) => {
          setCa(e.target.value)
        }}
      />
    </div>
  )
}

export default ForwarderEditorElastic

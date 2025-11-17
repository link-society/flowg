import Divider from '@mui/material/Divider'
import InputAdornment from '@mui/material/InputAdornment'
import TextField from '@mui/material/TextField'

import ForwarderConfigElasticModel from '@/lib/models/ForwarderConfigElasticModel'
import { ForwarderConfigTypeLabelMap } from '@/lib/models/ForwarderConfigModel'

import ForwarderIconElastic from '@/components/ForwarderIconElastic'
import InputList from '@/components/InputList'

type ForwarderEditorElasticProps = {
  config: ForwarderConfigElasticModel
  onConfigChange: (config: ForwarderConfigElasticModel) => void
}

const ForwarderEditorElastic = ({
  config,
  onConfigChange,
}: ForwarderEditorElasticProps) => {
  return (
    <div
      id="container:editor.forwarders.elastic"
      className="flex flex-col items-stretch gap-3"
    >
      <div className="mb-6 shadow">
        <TextField
          label="Forwarder Type"
          variant="outlined"
          className="w-full"
          type="text"
          value={ForwarderConfigTypeLabelMap.elastic}
          disabled
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <ForwarderIconElastic />
                </InputAdornment>
              ),
            },
          }}
        />
      </div>

      <TextField
        id="input:editor.forwarders.elastic.index"
        label="Index"
        variant="outlined"
        type="text"
        value={config.index}
        onChange={(e) => {
          onConfigChange({
            ...config,
            index: e.target.value,
          })
        }}
      />

      <TextField
        id="input:editor.forwarders.elastic.username"
        label="Username"
        variant="outlined"
        type="text"
        value={config.username}
        onChange={(e) => {
          onConfigChange({
            ...config,
            username: e.target.value,
          })
        }}
      />

      <TextField
        id="input:editor.forwarders.elastic.password"
        label="Password"
        variant="outlined"
        type="password"
        value={config.password}
        onChange={(e) => {
          onConfigChange({
            ...config,
            password: e.target.value,
          })
        }}
      />

      <Divider />

      <InputList
        id="editor.forwarders.elastic.addresses"
        itemLabel="Address"
        items={config.addresses}
        onChange={(addresses) => {
          onConfigChange({
            ...config,
            addresses,
          })
        }}
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
        value={config.ca ?? ''}
        onChange={(e) => {
          onConfigChange({
            ...config,
            ca: e.target.value || undefined,
          })
        }}
      />
    </div>
  )
}

export default ForwarderEditorElastic

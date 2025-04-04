import { useState } from 'react'

import Button from '@mui/material/Button'
import Divider from '@mui/material/Divider'
import Paper from '@mui/material/Paper'
import TextField from '@mui/material/TextField'

import AddIcon from '@mui/icons-material/Add'
import DeleteIcon from '@mui/icons-material/Delete'

import { StreamConfigModel } from '@/lib/models/storage'

type StreamEditorProps = {
  streamConfig: StreamConfigModel
  onStreamConfigChange: (config: StreamConfigModel) => void
}

export const StreamEditor = ({
  streamConfig,
  onStreamConfigChange,
}: StreamEditorProps) => {
  const [newField, setNewField] = useState('')

  return (
    <div className="h-full flex flex-row items-stretch gap-3">
      <Paper className="h-full flex-1 flex flex-col items-stretch">
        <h1 className="p-3 bg-gray-100 text-xl text-center font-semibold">
          Retention
        </h1>
        <Divider />

        <div
          className="
            p-3 grow shrink h-0 overflow-auto
            flex flex-col items-stretch gap-3
          "
        >
          <TextField
            id="input:editor.streams.retention-size"
            label="Retention size (in MB)"
            variant="outlined"
            type="number"
            value={streamConfig.size}
            onChange={(e) => {
              onStreamConfigChange({
                ...streamConfig,
                size: Number(e.target.value),
              })
            }}
          />

          <TextField
            id="input:editor.streams.retention-ttl"
            label="Retention time (in seconds)"
            variant="outlined"
            type="number"
            value={streamConfig.ttl}
            onChange={(e) => {
              onStreamConfigChange({
                ...streamConfig,
                ttl: Number(e.target.value),
              })
            }}
          />

          <p className="italic">
            Use{' '}
            <code className="font-mono bg-gray-200 text-red-500 px-1">0</code>{' '}
            to disable
          </p>
        </div>
      </Paper>

      <Paper className="h-full flex-1 flex flex-col items-stretch">
        <h1 className="p-3 bg-gray-100 text-xl text-center font-semibold">
          Indexes
        </h1>
        <Divider />

        <div
          className="
            p-3 grow shrink h-0 overflow-auto
            flex flex-col items-stretch gap-3
          "
        >
          {streamConfig.indexed_fields.map((field) => (
            <div key={field} className="flex flex-row items-stretch gap-3">
              <TextField
                id={`input:editor.streams.indexed-field.item.${field}.name`}
                label="Field"
                variant="outlined"
                type="text"
                value={field}
                onChange={(e) => {
                  onStreamConfigChange({
                    ...streamConfig,
                    indexed_fields: streamConfig.indexed_fields.map((f) =>
                      f === field ? e.target.value : f
                    ),
                  })
                }}
                className="grow"
              />

              <Button
                id={`btn:editor.streams.indexed-field.item.${field}.delete`}
                variant="contained"
                color="error"
                size="small"
                onClick={() => {
                  onStreamConfigChange({
                    ...streamConfig,
                    indexed_fields: streamConfig.indexed_fields.filter(
                      (f) => f !== field
                    ),
                  })
                }}
              >
                <DeleteIcon />
              </Button>
            </div>
          ))}

          <div className="flex flex-row items-stretch gap-3">
            <TextField
              id="input:editor.streams.indexed-field.new.name"
              label="Field"
              variant="outlined"
              type="text"
              value={newField}
              onChange={(e) => {
                setNewField(e.target.value)
              }}
              className="grow"
            />

            <Button
              id="btn:editor.streams.indexed-field.new.add"
              variant="contained"
              color="primary"
              size="small"
              onClick={() => {
                onStreamConfigChange({
                  ...streamConfig,
                  indexed_fields: [...streamConfig.indexed_fields, newField],
                })
                setNewField('')
              }}
            >
              <AddIcon />
            </Button>
          </div>
        </div>
      </Paper>
    </div>
  )
}

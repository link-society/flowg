import Divider from '@mui/material/Divider'
import Paper from '@mui/material/Paper'
import TextField from '@mui/material/TextField'

import StreamConfigModel from '@/lib/models/StreamConfigModel'

import ListEdit from '@/components/ListEdit.tsx'

type StreamEditorProps = {
  streamConfig: StreamConfigModel
  onStreamConfigChange: (config: StreamConfigModel) => void
}

const StreamEditor = ({
  streamConfig,
  onStreamConfigChange,
}: StreamEditorProps) => (
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
          Use <code className="font-mono bg-gray-200 text-red-500 px-1">0</code>{' '}
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
        <ListEdit
          id="editor.streams.indexed-field"
          list={streamConfig.indexed_fields}
          setList={(list) =>
            onStreamConfigChange({ ...streamConfig, indexed_fields: list })
          }
        />
      </div>
    </Paper>
  </div>
)

export default StreamEditor

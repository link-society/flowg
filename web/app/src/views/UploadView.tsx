import { useState } from 'react'
import { LoaderFunction, useLoaderData } from 'react-router'

import Button from '@mui/material/Button'
import Card from '@mui/material/Card'
import CircularProgress from '@mui/material/CircularProgress'
import Divider from '@mui/material/Divider'
import FormControl from '@mui/material/FormControl'
import InputLabel from '@mui/material/InputLabel'
import MenuItem from '@mui/material/MenuItem'
import Select from '@mui/material/Select'
import { styled } from '@mui/material/styles'

import AttachFileIcon from '@mui/icons-material/AttachFile'
import UploadFileIcon from '@mui/icons-material/UploadFile'

import * as configApi from '@/lib/api/operations/config'
import * as logsApi from '@/lib/api/operations/logs'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import { loginRequired } from '@/lib/decorators/loaders'

type LoaderData = {
  pipelines: string[]
}

export const loader: LoaderFunction = loginRequired(
  async (): Promise<LoaderData> => {
    const pipelines = await configApi.listPipelines()

    return { pipelines }
  }
)

const VisuallyHiddenInput = styled('input')({
  clip: 'rect(0 0 0 0)',
  clipPath: 'inset(50%)',
  height: 1,
  overflow: 'hidden',
  position: 'absolute',
  bottom: 0,
  left: 0,
  whiteSpace: 'nowrap',
  width: 1,
})

const UploadView = () => {
  const notify = useNotify()

  const { pipelines } = useLoaderData() as LoaderData
  const [pipeline, setPipeline] = useState<string | null>(null)
  const [file, setFile] = useState<File | null>(null)

  const onSelectPipeline = (event: any) => {
    setPipeline(event.target.value)
  }

  const [uploadLogs, loading] = useApiOperation(async () => {
    if (pipeline === null || file === null) {
      return
    }

    await logsApi.uploadTextLogs(pipeline, file)
    notify.success('Logs uploaded successfully')

    setPipeline(null)
    setFile(null)
  }, [pipeline, file])

  return (
    <div className="w-1/3 py-6 m-auto">
      <header className="mb-6 flex flex-row items-center justify-center gap-2">
        <UploadFileIcon />
        <h1 className="text-3xl text-center font-bold">Upload Logs</h1>
      </header>
      <Card className="p-3">
        <form
          className="flex flex-col items-stretch gap-3"
          onSubmit={(e) => {
            e.preventDefault()
            uploadLogs()
          }}
        >
          <FormControl fullWidth>
            <InputLabel id="label:upload.field.pipeline">
              Select Pipeline
            </InputLabel>
            <Select
              labelId="label:upload.field.pipeline"
              id="select:upload.field.pipeline"
              value={pipeline}
              label="Pipeline"
              onChange={onSelectPipeline}
            >
              {pipelines.map((p) => (
                <MenuItem key={p} value={p}>
                  {p}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <Button
            id="btn:upload.field.file"
            component="label"
            role={undefined}
            variant="contained"
            startIcon={<AttachFileIcon />}
          >
            {file === null ? <>Select log file</> : <>{file.name}</>}

            <VisuallyHiddenInput
              id="input:upload.field.file"
              type="file"
              onChange={(e) =>
                setFile(e.target.files ? e.target.files[0] : null)
              }
            />
          </Button>

          <Divider />

          <Button
            id="btn:upload.submit"
            variant="contained"
            color="secondary"
            className="w-full"
            type="submit"
            startIcon={!loading && <UploadFileIcon />}
            disabled={pipeline === null || file === null}
          >
            {loading ? (
              <CircularProgress color="inherit" size={24} />
            ) : (
              <>Upload</>
            )}
          </Button>
        </form>
      </Card>
    </div>
  )
}

export default UploadView

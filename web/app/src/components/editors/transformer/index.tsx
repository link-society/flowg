import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@toolpad/core/useNotifications'
import { useConfig } from '@/lib/context/config'

import PlayArrowIcon from '@mui/icons-material/PlayArrow'

import Grid from '@mui/material/Grid2'
import Paper from '@mui/material/Paper'
import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'

import { KeyValueEditor } from '@/components/form/kv-editor'
import { CodeEditor } from './code-editor'

import { UnauthenticatedError, PermissionDeniedError } from '@/lib/api/errors'
import * as testApi from '@/lib/api/operations/tests'

type TransformerEditorProps = {
  code: string
  onCodeChange: (value: string) => void
}

export const TransformerEditor = (props: TransformerEditorProps) => {
  const [code, setCode] = useState(props.code)
  const [testRecord, setTestRecord] = useState<[string, string][]>([])

  const [testLoading, setTestLoading] = useState(false)
  const [testResult, setTestResult] = useState<Record<string, string>>({})

  const navigate = useNavigate()
  const notifications = useNotifications()
  const config = useConfig()

  useEffect(
    () => { props.onCodeChange(code) },
    [code],
  )

  const onTest = async () => {
    setTestLoading(true)

    try {
      const input = Object.fromEntries(testRecord)
      const output = await testApi.testTransformer(code, input)
      setTestResult(output)
    }
    catch (error) {
      if (error instanceof UnauthenticatedError) {
        notifications.show('Session expired', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
        navigate('/web/login')
      }
      else if (error instanceof PermissionDeniedError) {
        notifications.show('Permission denied', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }
      else {
        notifications.show('Unknown error', {
          severity: 'error',
          autoHideDuration: config.notifications?.autoHideDuration,
        })
      }

      console.error(error)
    }

    setTestLoading(false)
  }

  return (
    <Grid container spacing={1} className="w-full h-full">
      <Grid size={8}>
        <Paper className="w-full h-full">
          <CodeEditor code={code} onCodeChange={setCode} />
        </Paper>
      </Grid>
      <Grid size={4}>
        <div className="h-full flex flex-col items-stretch gap-2">
          <div className="flex-1 h-0">
            <Paper className="p-2 h-full overflow-auto">
              <p className="text-sm text-gray-700 font-semibold mb-2">Input Record:</p>
              <KeyValueEditor
                keyLabel="Field"
                keyValues={testRecord}
                onChange={setTestRecord}
              />
            </Paper>
          </div>
          <div className="flex flex-col items-center">
            <Button
              variant="contained"
              color="primary"
              endIcon={!testLoading && <PlayArrowIcon />}
              disabled={testLoading}
              onClick={() => onTest()}
            >
              {testLoading
                ? <CircularProgress size={24} />
                : <>Run</>
              }
            </Button>
          </div>
          <div className="flex-1 h-0">
            <Paper className="p-2 h-full flex flex-col items-stretch">
              <p className="text-sm text-gray-700 font-semibold mb-2">Output Record:</p>

              <Paper
                variant="outlined"
                className="
                  p-2 flex-grow flex-shrink h-0 overflow-auto
                  font-mono !bg-gray-100
                "
                component="pre"
              >
                {JSON.stringify(testResult, null, 2)}
              </Paper>
            </Paper>
          </div>
        </div>
      </Grid>
    </Grid>
  )
}
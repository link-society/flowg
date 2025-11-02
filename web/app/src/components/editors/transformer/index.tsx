import { useEffect, useState } from 'react'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import Grid from '@mui/material/Grid'
import Paper from '@mui/material/Paper'

import PlayArrowIcon from '@mui/icons-material/PlayArrow'

import * as testApi from '@/lib/api/operations/tests'
import { useApiOperation } from '@/lib/hooks/api'

import { KeyValueEditor } from '@/components/form/kv-editor'

import { CodeEditor } from './code-editor'

type TransformerEditorProps = Readonly<{
  code: string
  onCodeChange: (value: string) => void
}>

export const TransformerEditor = (props: TransformerEditorProps) => {
  const [code, setCode] = useState(props.code)
  const [testRecord, setTestRecord] = useState<[string, string][]>([])
  const [testResult, setTestResult] = useState<string>('')

  useEffect(() => {
    props.onCodeChange(code)
  }, [code])

  const [onTest, testLoading] = useApiOperation(async () => {
    const input = Object.fromEntries(testRecord)
    const output = await testApi.testTransformer(code, input)

    if (output.success) {
      setTestResult(JSON.stringify(output.record, null, 2))
    } else {
      setTestResult(output.error)
    }
  }, [code, testRecord])

  return (
    <Grid container spacing={1} className="w-full h-full">
      <Grid size={8}>
        <Paper className="w-full h-full">
          <CodeEditor
            id="monaco:transformers.editor"
            code={code}
            onCodeChange={setCode}
          />
        </Paper>
      </Grid>
      <Grid size={4}>
        <div className="h-full flex flex-col items-stretch gap-2">
          <div className="flex-1 h-0">
            <Paper className="p-2 h-full overflow-auto">
              <p className="text-sm text-gray-700 font-semibold mb-2">
                Input Record:
              </p>
              <KeyValueEditor
                id="kv:transformers.test.record"
                keyLabel="Field"
                keyValues={testRecord}
                onChange={setTestRecord}
              />
            </Paper>
          </div>
          <div className="flex flex-col items-center">
            <Button
              id="btn:transformers.test.run"
              variant="contained"
              color="primary"
              endIcon={!testLoading && <PlayArrowIcon />}
              disabled={testLoading}
              onClick={() => onTest()}
            >
              {testLoading ? <CircularProgress size={24} /> : <>Run</>}
            </Button>
          </div>
          <div className="flex-1 h-0">
            <Paper className="p-2 h-full flex flex-col items-stretch">
              <p className="text-sm text-gray-700 font-semibold mb-2">
                Output Record:
              </p>

              <Paper
                id="container:transformers.test.result"
                variant="outlined"
                className="
                  p-2 grow shrink h-0 overflow-auto
                  font-mono bg-gray-100!
                "
                component="pre"
              >
                {testResult}
              </Paper>
            </Paper>
          </div>
        </div>
      </Grid>
    </Grid>
  )
}

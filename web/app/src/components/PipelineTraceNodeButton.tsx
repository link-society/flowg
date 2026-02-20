import { Visibility } from '@mui/icons-material'

import { useState } from 'react'

import Button from '@mui/material/Button'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Paper from '@mui/material/Paper'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

type PipelineTraceNodeButtonProps = {
  trace: NodeTrace
}

const PipelineTraceNodeButton = ({ trace }: PipelineTraceNodeButtonProps) => {
  const [open, setOpen] = useState<boolean>(false)

  return (
    <>
      <Button
        variant="contained"
        size="small"
        color="primary"
        startIcon={<Visibility />}
        onClick={() => setOpen(true)}
      >
        Inspect
      </Button>

      <Dialog open={open} onClose={() => setOpen(false)}>
        <DialogTitle>Node trace</DialogTitle>
        <DialogContent>
          <div className="flex flex-col gap-5">
            {trace.error && (
              <div className="flex flex-col gap-2">
                <p className="text-sm text-gray-700 font-semibold mb-2">
                  Error:
                </p>
                <Paper
                  id="container:transformers.test.result"
                  variant="outlined"
                  className="
                      p-2 grow shrink overflow-auto
                      font-mono bg-gray-100! min-w-50
                    "
                  component="pre"
                >
                  {trace.error}
                </Paper>
              </div>
            )}
            <div className="flex gap-5">
              {trace.input && (
                <div className="flex flex-col gap-2">
                  <p className="text-sm text-gray-700 font-semibold mb-2">
                    Input Record:
                  </p>
                  <Paper
                    id="container:transformers.test.result"
                    variant="outlined"
                    className="
                      p-2 grow shrink overflow-auto
                      font-mono bg-gray-100! min-w-50
                    "
                    component="pre"
                  >
                    {JSON.stringify(trace.input, null, 2)}
                  </Paper>
                </div>
              )}

              {trace.output && (
                <div className="flex flex-col gap-2">
                  <p className="text-sm text-gray-700 font-semibold mb-2">
                    Output Record:
                  </p>
                  <Paper
                    id="container:transformers.test.result"
                    variant="outlined"
                    className="
                      p-2 grow shrink overflow-auto
                      font-mono bg-gray-100! min-w-50
                    "
                    component="pre"
                  >
                    {JSON.stringify(trace.output, null, 2)}
                  </Paper>
                </div>
              )}
            </div>
          </div>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </>
  )
}

export default PipelineTraceNodeButton

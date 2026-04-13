import { Visibility } from '@mui/icons-material'

import { useState } from 'react'

import Button from '@mui/material/Button'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'

import {
  TraceColumns,
  TraceDialogContent,
  TraceLabel,
  TracePaper,
  TraceSection,
} from './styles'
import { PipelineTraceNodeButtonProps } from './types'

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
          <TraceDialogContent>
            {trace.error && (
              <TraceSection>
                <TraceLabel>Error:</TraceLabel>
                <TracePaper variant="outlined" component="pre">
                  {trace.error}
                </TracePaper>
              </TraceSection>
            )}
            <TraceColumns>
              {trace.input && (
                <TraceSection>
                  <TraceLabel>Input Record:</TraceLabel>
                  <TracePaper variant="outlined" component="pre">
                    {JSON.stringify(trace.input, null, 2)}
                  </TracePaper>
                </TraceSection>
              )}

              {trace.output && (
                <TraceSection>
                  <TraceLabel>Output Record:</TraceLabel>
                  <TracePaper variant="outlined" component="pre">
                    {JSON.stringify(trace.output, null, 2)}
                  </TracePaper>
                </TraceSection>
              )}
            </TraceColumns>
          </TraceDialogContent>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </>
  )
}

export default PipelineTraceNodeButton

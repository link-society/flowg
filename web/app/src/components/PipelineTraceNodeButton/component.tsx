import { Visibility } from '@mui/icons-material'
import { Tab, Tabs } from '@mui/material'

import { useState } from 'react'

import Button from '@mui/material/Button'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'

import { TraceDialogContent } from './styles'
import { PipelineTraceNodeButtonProps } from './types'

import NodeTraceTabPanel from '../NodeTraceTabPanel/component'

const PipelineTraceNodeButton = ({ traces }: PipelineTraceNodeButtonProps) => {
  const [open, setOpen] = useState<boolean>(false)
  const [tab, setTab] = useState<number>(0)

  const onTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTab(newValue)
  }
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

      <Dialog
        open={open}
        onClose={() => setOpen(false)}
        maxWidth={false}
        slotProps={{
          paper: {
            sx: {
              width: '80%',
              height: '90%',
            },
          },
        }}
      >
        <DialogTitle>Node traces</DialogTitle>
        <DialogContent>
          <TraceDialogContent>
            <Tabs
              variant="scrollable"
              sx={{ borderBottom: 1, borderColor: 'divider' }}
              value={tab}
              onChange={onTabChange}
            >
              {traces.map((_trace, index) => (
                <Tab key={index} label={`Event #${index + 1}`} />
              ))}
            </Tabs>

            {traces.map((trace, index) => (
              <NodeTraceTabPanel
                trace={trace}
                key={index}
                value={tab}
                index={index}
              />
            ))}
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

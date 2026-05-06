import { Visibility } from '@mui/icons-material'

import { useState } from 'react'

import Button from '@mui/material/Button'
import Dialog from '@mui/material/Dialog'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Tab from '@mui/material/Tab'
import Tabs from '@mui/material/Tabs'

import { NodeTrace } from '@/lib/models/PipelineTrace.ts'

import NodeTraceTabPanel from '@/components/NodeTraceTabPanel'

type PipelineTraceNodeButtonProps = {
  traces: NodeTrace[]
}

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
          <div className="flex flex-col gap-5">
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

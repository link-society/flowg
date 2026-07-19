import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'
import DialogActions from '@mui/material/DialogActions'
import DialogContent from '@mui/material/DialogContent'
import DialogTitle from '@mui/material/DialogTitle'
import Tab from '@mui/material/Tab'

import VisibilityIcon from '@mui/icons-material/Visibility'

import { TraceDialog, TraceDialogContent, TraceTabs } from './styles'
import { PipelineTraceNodeButtonProps } from './types'

import NodeTraceTabPanel from '../NodeTraceTabPanel/component'

const PipelineTraceNodeButton = ({ traces }: PipelineTraceNodeButtonProps) => {
  const { t } = useTranslation()
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
        startIcon={<VisibilityIcon />}
        onClick={() => setOpen(true)}
      >
        {t('components.pipelineTraceNodeButton.inspect')}
      </Button>

      <TraceDialog open={open} onClose={() => setOpen(false)} maxWidth={false}>
        <DialogTitle>
          {t('components.pipelineTraceNodeButton.title')}
        </DialogTitle>
        <DialogContent>
          <TraceDialogContent>
            <TraceTabs variant="scrollable" value={tab} onChange={onTabChange}>
              {traces.map((_trace, index) => (
                <Tab
                  key={`tab-${index + 1}`}
                  label={t('components.pipelineTraceNodeButton.eventLabel', {
                    n: index + 1,
                  })}
                />
              ))}
            </TraceTabs>

            {traces.map((trace, index) => (
              <NodeTraceTabPanel
                trace={trace}
                key={`panel-${index + 1}`}
                value={tab}
                index={index}
              />
            ))}
          </TraceDialogContent>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>
            {t('common.actions.close')}
          </Button>
        </DialogActions>
      </TraceDialog>
    </>
  )
}

export default PipelineTraceNodeButton

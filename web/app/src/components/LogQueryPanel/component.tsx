import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import Button from '@mui/material/Button'
import CircularProgress from '@mui/material/CircularProgress'
import TextField from '@mui/material/TextField'

import SearchIcon from '@mui/icons-material/Search'

import TimeWindowSelector, {
  TimeWindowFactory,
} from '@/components/TimeWindowSelector/component'

import {
  LogQueryPanelActions,
  LogQueryPanelContainer,
  LogQueryPanelFilterForm,
  LogQueryPanelTimeWindow,
} from './styles'
import { LogQueryPanelProps } from './types'

const LogQueryPanel = (props: LogQueryPanelProps) => {
  const { t } = useTranslation()
  const [filter, setFilter] = useState('')
  const [timeWindowFactory, setTimeWindowFactory] =
    useState<TimeWindowFactory | null>(null)

  const requestFetch = useCallback(() => {
    if (timeWindowFactory !== null) {
      const { from, to, live } = timeWindowFactory.make()
      props.onFetchRequested(filter, from, to, live)
    }
  }, [timeWindowFactory, filter, props.onFetchRequested])

  useEffect(() => {
    requestFetch()
  }, [timeWindowFactory])

  return (
    <LogQueryPanelContainer>
      <LogQueryPanelFilterForm
        onSubmit={(e) => {
          e.preventDefault()
          requestFetch()
        }}
      >
        <TextField
          id="input:streams.filter"
          label={t('components.logQueryPanel.filterLabel')}
          variant="outlined"
          size="small"
          value={filter}
          onChange={(e) => {
            setFilter(e.target.value)
          }}
          disabled={props.loading}
          fullWidth
        />
      </LogQueryPanelFilterForm>

      <LogQueryPanelTimeWindow>
        <TimeWindowSelector onTimeWindowChanged={setTimeWindowFactory} />
      </LogQueryPanelTimeWindow>

      <LogQueryPanelActions>
        <Button
          id="btn:streams.query"
          variant="contained"
          size="small"
          color="secondary"
          onClick={() => requestFetch()}
          endIcon={!props.loading && <SearchIcon />}
          disabled={props.loading}
          fullWidth
        >
          {props.loading ? (
            <CircularProgress color="inherit" size={24} />
          ) : (
            <>{t('components.logQueryPanel.submit')}</>
          )}
        </Button>
      </LogQueryPanelActions>
    </LogQueryPanelContainer>
  )
}

export default LogQueryPanel

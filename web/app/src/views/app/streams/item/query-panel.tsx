import { useCallback, useEffect, useState } from 'react'

import SearchIcon from '@mui/icons-material/Search'

import Grid from '@mui/material/Grid2'
import CircularProgress from '@mui/material/CircularProgress'
import Button from '@mui/material/Button'
import TextField from '@mui/material/TextField'

import { TimeWindowSelector, TimeWindowFactory } from './timewindow-selector'

type QueryPanelProps = Readonly<{
  loading: boolean
  onFetchRequested: (filter: string, from: Date, to: Date, live: boolean) => void
}>

export const QueryPanel = (props: QueryPanelProps) => {
  const [filter, setFilter] = useState('')
  const [timeWindowFactory, setTimeWindowFactory] = useState<TimeWindowFactory | null>(null)

  const requestFetch = useCallback(
    () => {
      if (timeWindowFactory !== null) {
        const { from, to, live } = timeWindowFactory.make()
        props.onFetchRequested(filter, from, to, live)
      }
    },
    [props.onFetchRequested, timeWindowFactory, filter],
  )

  useEffect(
    () => {
      requestFetch()
    },
    [timeWindowFactory],
  )

  return (
    <Grid container spacing={2} className="p-3 items-center">
      <Grid
        size={{ xs: 6 }}
        component="form"
        onSubmit={(e) => {
          e.preventDefault()
          requestFetch()
        }}
      >
        <TextField
          label="Filter"
          variant="outlined"
          size="small"
          value={filter}
          onChange={(e) => {
            setFilter(e.target.value)
          }}
          disabled={props.loading}
          className="w-full"
        />
      </Grid>

      <Grid size={{ xs: 4 }}>
        <TimeWindowSelector
          loading={props.loading}
          onTimeWindowChanged={setTimeWindowFactory}
        />
      </Grid>

      <Grid size={{ xs: 2 }}>
        <Button
          className="w-full"
          variant="contained"
          size="small"
          color="secondary"
          onClick={() => requestFetch()}
          endIcon={!props.loading && <SearchIcon />}
          disabled={props.loading}
        >
          {props.loading
            ? <CircularProgress color="inherit" size={24} />
            : <>Query Logs</>
          }
        </Button>
      </Grid>
    </Grid>
  )
}

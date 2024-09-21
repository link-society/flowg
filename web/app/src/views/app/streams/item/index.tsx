import { useEffect, useState } from 'react'
import { useLoaderData, useParams } from 'react-router-dom'
import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import Grid from '@mui/material/Grid2'
import Paper from '@mui/material/Paper'
import Divider from '@mui/material/Divider'

import { StreamList } from './stream-list'
import { QueryPanel } from './query-panel'
import { Chart } from './chart'
import { LogTable } from './log-table'

import * as logApi from '@/lib/api/operations/logs'

import { LogEntryModel } from '@/lib/models'
import { ColDef } from 'ag-grid-community'

import { LoaderData } from '../loader'

export const StreamView = () => {
  const notify = useNotify()

  const { stream: currentStream } = useParams()
  const { streams, fields } = useLoaderData() as LoaderData

  const timestampColumnDef = (): ColDef<LogEntryModel> => ({
    headerName: 'Ingested At',
    headerClass: 'font-semibold',
    field: 'timestamp',
    valueFormatter: ({ value }) => {
      return (value as Date).toLocaleString()
    },
    suppressMovable: true,
    sort: 'desc',
  })

  const fieldToColumnDef = (field: string): ColDef<LogEntryModel> => ({
    headerName: field,
    headerClass: 'font-semibold',
    field: `fields.${field}`,
    sortable: false,
    cellClass: 'font-mono',
  })

  const [watcher, setWatcher] = useState<{ enabled: Boolean, filter: string }>({
    enabled: false,
    filter: '',
  })

  const [timeWindow, setTimeWindow] = useState<{ from: Date, to: Date }>({
    from: new Date(Date.now() - 15 * 60 * 1000),
    to: new Date(),
  })

  const [rowData, setRowData] = useState<LogEntryModel[]>([])
  const [columnDefs, setColumnDefs] = useState<ColDef<LogEntryModel>[]>([
    timestampColumnDef(),
    ...fields!.map(fieldToColumnDef),
  ])

  const [fetchLogs, loading] = useApiOperation(
    async (filter: string, from: Date, to: Date, live: boolean) => {
      const logs = await logApi.queryLogs(
        currentStream!,
        from, to,
        filter ? filter : undefined,
      )
      setRowData(logs)
      setTimeWindow({ from, to })
      setWatcher({ enabled: live, filter })
    },
    [currentStream, setRowData],
  )

  const [handleLiveError] = useApiOperation(
    async (err: Error) => { throw err },
    [],
  )

  useEffect(
    () => {
      if (watcher.enabled) {
        const bus = logApi.watchLogs(currentStream!, watcher.filter)

        bus.control.addEventListener('error', (event) => {
          const evt = event as CustomEvent
          handleLiveError(evt.detail)
        })

        bus.messages.addEventListener('log', (event) => {
          const evt = event as CustomEvent
          type RawLogEntry = {
            timestamp: string
            fields: Record<string, string>
          }
          const rawlogEntry: RawLogEntry = JSON.parse(evt.detail.data)
          const logEntry: LogEntryModel = {
            timestamp: new Date(rawlogEntry.timestamp),
            fields: rawlogEntry.fields,
          }
          setRowData((prev) => [...prev, logEntry])

          const allFields = [...fields!]

          for (const field of Object.keys(logEntry.fields)) {
            if (!allFields!.includes(field)) {
              allFields.push(field)
            }
          }

          allFields.sort()

          setColumnDefs([
            timestampColumnDef(),
            ...allFields.map(fieldToColumnDef),
          ])
        })

        bus.messages.addEventListener('exception', (event) => {
          const evt = event as CustomEvent
          notify.error('An error occured while watching logs')
          console.error(evt.detail.data)
        })

        return () => bus.close()
      }
    },
    [currentStream, watcher],
  )

  return (
    <Grid container spacing={1} className="p-2 h-full">
      <Grid size={{ xs: 2 }} className="h-full">
        <Paper className="h-full overflow-auto">
          <StreamList streams={streams} currentStream={currentStream!} />
        </Paper>
      </Grid>
      <Grid size={{ xs: 10 }} className="flex flex-col items-stretch gap-2">
        <Paper>
          <QueryPanel
            loading={loading}
            onFetchRequested={fetchLogs}
          />
          <Divider />
          <Chart
            rowData={rowData}
            from={timeWindow.from}
            to={timeWindow.to}
          />
        </Paper>

        <LogTable rowData={rowData} columnDefs={columnDefs} />
      </Grid>
    </Grid>
  )
}

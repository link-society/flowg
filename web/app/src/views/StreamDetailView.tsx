import { useEffect, useState } from 'react'
import { LoaderFunction, useLoaderData } from 'react-router'

import Divider from '@mui/material/Divider'
import Grid from '@mui/material/Grid'
import Paper from '@mui/material/Paper'

import { ColDef } from 'ag-grid-community'

import * as configApi from '@/lib/api/operations/config'
import * as logApi from '@/lib/api/operations/logs'

import { useApiOperation } from '@/lib/hooks/api'
import { useNotify } from '@/lib/hooks/notify'

import LogEntryModel from '@/lib/models/LogEntryModel'

import { loginRequired } from '@/lib/decorators/loaders'

import LogChart from '@/components/LogChart'
import LogQueryPanel from '@/components/LogQueryPanel'
import LogTable from '@/components/LogTable'
import SideNavList from '@/components/SideNavList'
import StreamIndexSelector from '@/components/StreamIndexSelector'

type LoaderData = {
  streams: string[]
  currentStream: string
  fields: string[]
  indices: Record<string, Array<string>>
}

export const loader: LoaderFunction = loginRequired(async ({ params }) => {
  const [streamConfigs, fields, indices] = await Promise.all([
    configApi.listStreams(),
    configApi.listStreamFields(params.stream!),
    logApi.getStreamIndices(params.stream!),
  ])

  const streams = Object.keys(streamConfigs)
  streams.sort((a, b) => a.localeCompare(b))

  return { streams, currentStream: params.stream!, fields, indices }
})

const StreamDetailView = () => {
  const notify = useNotify()

  const { streams, currentStream, fields, indices } =
    useLoaderData() as LoaderData

  const [selectedIndices, setSelectedIndices] = useState<
    Record<string, Array<string>>
  >({})

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
    valueGetter: ({ data }) => data?.fields[field],
  })

  const [watcher, setWatcher] = useState<{ enabled: boolean; filter: string }>({
    enabled: false,
    filter: '',
  })

  const [timeWindow, setTimeWindow] = useState<{ from: Date; to: Date }>({
    from: new Date(),
    to: new Date(),
  })

  const [rowData, setRowData] = useState<LogEntryModel[]>([])
  const [columnDefs, setColumnDefs] = useState<ColDef<LogEntryModel>[]>([
    timestampColumnDef(),
    ...fields.map(fieldToColumnDef),
  ])

  const [fetchLogs, loading] = useApiOperation(
    async (filter: string, from: Date, to: Date, live: boolean) => {
      console.log(selectedIndices)
      const logs = await logApi.queryLogs(
        currentStream,
        from,
        to,
        filter === '' ? undefined : filter,
        selectedIndices
      )
      setRowData(logs)
      setTimeWindow({ from, to })
      setWatcher({ enabled: live, filter })
    },
    [currentStream, setRowData, selectedIndices]
  )

  const [handleLiveError] = useApiOperation(async (err: Error) => {
    throw err
  }, [])

  useEffect(() => {
    if (watcher.enabled) {
      const bus = logApi.watchLogs(
        currentStream,
        watcher.filter,
        selectedIndices
      )

      const incomingState = {
        rowData: [] as LogEntryModel[],
        columnDefs,
      }

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
        incomingState.rowData.push(logEntry)

        const allFields = [...fields]

        for (const field of Object.keys(logEntry.fields)) {
          if (!allFields.includes(field)) {
            allFields.push(field)
          }
        }

        allFields.sort((a, b) => a.localeCompare(b))

        incomingState.columnDefs = [
          timestampColumnDef(),
          ...allFields.map(fieldToColumnDef),
        ]
      })

      bus.messages.addEventListener('exception', (event) => {
        const evt = event as CustomEvent
        notify.error('An error occured while watching logs')
        console.error(evt.detail.data)
      })

      const token = setInterval(() => {
        setRowData((prev) => {
          return [...prev, ...incomingState.rowData]
        })
        setColumnDefs(incomingState.columnDefs)
        setTimeWindow((prev) => ({
          from: prev.from,
          to: new Date(),
        }))
        incomingState.rowData = []
      }, 1000)

      return () => {
        bus.close()
        clearInterval(token)
      }
    }
  }, [currentStream, watcher, selectedIndices])

  return (
    <Grid container spacing={1} className="p-2 h-full">
      <Grid size={{ xs: 2 }} className="h-full">
        <SideNavList
          namespace="streams"
          urlPrefix="/web/streams"
          items={streams}
          currentItem={currentStream}
        />
      </Grid>
      <Grid size={{ xs: 8 }} className="flex flex-col items-stretch gap-2">
        <Paper>
          <LogQueryPanel loading={loading} onFetchRequested={fetchLogs} />
          <Divider />
          <LogChart
            rowData={rowData}
            from={timeWindow.from}
            to={timeWindow.to}
          />
        </Paper>

        <LogTable rowData={rowData} columnDefs={columnDefs} />
      </Grid>
      <Grid size={{ xs: 2 }} className="h-full">
        <StreamIndexSelector
          indices={indices}
          selection={selectedIndices}
          onSelectionChange={setSelectedIndices}
        />
      </Grid>
    </Grid>
  )
}

export default StreamDetailView

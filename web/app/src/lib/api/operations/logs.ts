import * as request from '@/lib/api/request'

import { LogEntryModel } from '@/lib/models/log'

export const queryLogs = async (
  stream: string,
  from: Date,
  to: Date,
  filter?: string,
): Promise<LogEntryModel[]> => {
  type QueryLogsResponse = {
    success: boolean
    records: Array<{
      timestamp: string
      fields: Record<string, string>
    }>
  }

  const { body } = await request.GET<QueryLogsResponse>({
    path: `/api/v1/streams/${stream}/logs`,
    searchParams: new URLSearchParams(filter === undefined
      ? {
        from: from.toISOString(),
        to: to.toISOString(),
      }
      : {
        from: from.toISOString(),
        to: to.toISOString(),
        filter,
      }
    ),
  })

  return body.records.map(
    ({ timestamp, fields }) => ({
      timestamp: new Date(timestamp),
      fields,
    }),
  )
}

export const watchLogs = (stream: string, filter: string) => {
  return request.SSE({
    path: `/api/v1/streams/${stream}/logs/watch`,
    searchParams: new URLSearchParams({ filter }),
  })
}

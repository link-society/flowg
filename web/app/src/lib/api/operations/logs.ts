import * as request from '@/lib/api/request'

import LogEntryModel from '@/lib/models/LogEntryModel'

export const queryLogs = async (
  stream: string,
  from: Date,
  to: Date,
  filter?: string,
  indexing?: Record<string, Array<string>>
): Promise<LogEntryModel[]> => {
  type QueryLogsResponse = {
    success: boolean
    records: Array<{
      timestamp: string
      fields: Record<string, string>
    }>
  }

  const searchParams: {
    from: string
    to: string
    filter?: string
    indexing?: string
  } = {
    from: from.toISOString(),
    to: to.toISOString(),
  }

  if (filter !== undefined) {
    searchParams.filter = filter
  }

  if (indexing !== undefined) {
    searchParams.indexing = JSON.stringify(indexing)
  }

  const { body } = await request.GET<QueryLogsResponse>({
    path: `/api/v1/streams/${stream}/logs`,
    searchParams: new URLSearchParams(searchParams),
  })

  return body.records.map(({ timestamp, fields }) => ({
    timestamp: new Date(timestamp),
    fields,
  }))
}

export const watchLogs = (
  stream: string,
  filter: string,
  indexing: Record<string, Array<string>>
) => {
  return request.SSE({
    path: `/api/v1/streams/${stream}/logs/watch`,
    searchParams: new URLSearchParams({
      filter,
      indexing: JSON.stringify(indexing),
    }),
  })
}

export const uploadTextLogs = async (
  pipeline: string,
  content: BodyInit
): Promise<void> => {
  await request.POST({
    path: `/api/v1/pipelines/${pipeline}/logs/text`,
    body: content,
    contentType: 'text/plain',
  })
}

export const getStreamUsage = async (stream: string): Promise<number> => {
  type GetStreamUsageResponse = {
    success: boolean
    usage: number
  }

  const { body } = await request.GET<GetStreamUsageResponse>({
    path: `/api/v1/streams/${stream}/usage`,
  })

  return body.usage
}

export const getStreamIndices = async (
  stream: string
): Promise<Record<string, Array<string>>> => {
  type GetStreamIndicesResponse = {
    success: boolean
    indices: Record<string, Array<string>>
  }

  const { body } = await request.GET<GetStreamIndicesResponse>({
    path: `/api/v1/streams/${stream}/indices`,
  })

  return body.indices
}

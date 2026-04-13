export type LogQueryPanelProps = Readonly<{
  loading: boolean
  onFetchRequested: (
    filter: string,
    from: Date,
    to: Date,
    live: boolean
  ) => void
}>

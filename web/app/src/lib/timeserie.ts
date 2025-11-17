import LogEntryModel from '@/lib/models/LogEntryModel'

const INTERVAL_OPTIONS = [
  5 * 1000, // 5 seconds
  5 * 60 * 1000, // 5 minutes
  60 * 60 * 1000, // 1 hour
  24 * 60 * 60 * 1000, // 1 day
]

const TARGET_POINT_COUNT = 100

const bestInterval = (from: Date, to: Date) => {
  const timewindowLength = to.getTime() - from.getTime()
  let bestInterval = INTERVAL_OPTIONS[0]
  let minDiff = Infinity

  for (const interval of INTERVAL_OPTIONS) {
    const pointCount = timewindowLength / interval
    const diff = Math.abs(pointCount - TARGET_POINT_COUNT)

    if (diff < minDiff) {
      minDiff = diff
      bestInterval = interval
    }
  }

  return bestInterval
}

const floorToInterval = (date: Date, interval: number) => {
  const timestamp = date.getTime()
  const floored = Math.floor(timestamp / interval) * interval
  return new Date(floored)
}

const ceilToInterval = (date: Date, interval: number) => {
  const timestamp = date.getTime()
  const ceiled = Math.ceil(timestamp / interval) * interval
  return new Date(ceiled)
}

export const aggregateLogs = (logs: LogEntryModel[], from: Date, to: Date) => {
  const interval = bestInterval(from, to)
  const flooredFrom = floorToInterval(from, interval)
  const ceiledTo = ceilToInterval(to, interval)

  const timeserie: [number, number][] = []
  const index: { [key: number]: number } = {}

  for (
    let current = flooredFrom;
    current < ceiledTo;
    current = new Date(current.getTime() + interval)
  ) {
    const idx = timeserie.length
    timeserie.push([current.getTime(), 0])
    index[current.getTime()] = idx
  }

  for (const log of logs) {
    if (log.timestamp < from || log.timestamp >= to) {
      continue
    }

    const floored = floorToInterval(log.timestamp, interval).getTime()
    timeserie[index[floored]][1] += 1
  }

  return timeserie
}

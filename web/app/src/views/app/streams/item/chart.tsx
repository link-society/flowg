import { useMemo } from 'react'

import Box from '@mui/material/Box'

import ApexChart from 'react-apexcharts'
import type { ApexOptions } from 'apexcharts'

import { LogEntryModel } from '@/lib/models/log'
import { aggregateLogs } from '@/lib/timeserie'

type ChartProps = Readonly<{
  rowData: LogEntryModel[]
  from: Date,
  to: Date,
}>

const CHART_OPTIONS: ApexOptions = {
  chart: {
    animations: {
      enabled: false,
    },
  },
  dataLabels: {
    enabled: false,
  },
  xaxis: {
    type: 'datetime',
  },
}

export const Chart = ({ rowData, from, to }: ChartProps) => {
  const series = useMemo(
    () => [
      {
        name: 'Logs',
        data: aggregateLogs(rowData, from, to),
      }
    ],
    [rowData, from, to],
  )

  return (
    <Box className="bg-gray-100 min-h-[150px]">
      <ApexChart
        options={CHART_OPTIONS}
        series={series}
        type="bar"
        width="100%"
        height={150}
      />
    </Box>
  )
}

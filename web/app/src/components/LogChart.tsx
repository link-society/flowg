import type { ApexOptions } from 'apexcharts'

import { useMemo } from 'react'
import ApexChart from 'react-apexcharts'

import Box from '@mui/material/Box'

import LogEntryModel from '@/lib/models/LogEntryModel'

import { aggregateLogs } from '@/lib/timeserie'

type LogChartProps = Readonly<{
  rowData: LogEntryModel[]
  from: Date
  to: Date
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

const LogChart = ({ rowData, from, to }: LogChartProps) => {
  const series = useMemo(
    () => [
      {
        name: 'Logs',
        data: aggregateLogs(rowData, from, to),
      },
    ],
    [rowData, from, to]
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

export default LogChart

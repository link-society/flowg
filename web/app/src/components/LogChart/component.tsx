import type { ApexOptions } from 'apexcharts'

import { useMemo } from 'react'
import ApexChart from 'react-apexcharts'

import { aggregateLogs } from '@/lib/timeserie'

import { LogChartContainer } from './styles'
import { LogChartProps } from './types'

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
    <LogChartContainer>
      <ApexChart
        options={CHART_OPTIONS}
        series={series}
        type="bar"
        width="100%"
        height={150}
      />
    </LogChartContainer>
  )
}

export default LogChart

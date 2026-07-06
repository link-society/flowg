import { useColorMode } from '@/theme'
import type { ApexOptions } from 'apexcharts'

import { useMemo } from 'react'
import ApexChart from 'react-apexcharts'

import { aggregateLogs } from '@/lib/timeserie'

import { LogChartContainer } from './styles'
import { LogChartProps } from './types'

const LogChart = ({ rowData, from, to }: LogChartProps) => {
  const { mode } = useColorMode()

  const options: ApexOptions = useMemo(
    () => ({
      chart: {
        animations: { enabled: false },
        foreColor: mode === 'dark' ? '#ffffff' : undefined,
      },
      dataLabels: { enabled: false },
      xaxis: { type: 'datetime' },
      tooltip: { theme: mode, x: { format: 'dd MMM HH:mm:ss' } },
    }),
    [mode]
  )

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
        options={options}
        series={series}
        type="bar"
        width="100%"
        height={150}
      />
    </LogChartContainer>
  )
}

export default LogChart

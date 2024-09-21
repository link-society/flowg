import { useCallback, useEffect, useMemo, useState } from 'react'

import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'

import Box from '@mui/material/Box'
import Divider from '@mui/material/Divider'
import Button from '@mui/material/Button'
import Menu from '@mui/material/Menu'
import ToggleButtonGroup from '@mui/material/ToggleButtonGroup'
import ToggleButton from '@mui/material/ToggleButton'

import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker'
import dayjs from 'dayjs'

type RelTimeWindowOption = {
  label: string
  value: number
}

const RELATIVE_TIMEWINDOW_OPTIONS: RelTimeWindowOption[] = [
  { label: 'Last 15 Minutes', value: 15 * 60 * 1000 },
  { label: 'Last Hour', value: 60 * 60 * 1000 },
  { label: 'Last Day', value: 24 * 60 * 60 * 1000 },
  { label: 'Last Week', value: 7 * 24 * 60 * 60 * 1000 },
]

const RELATIVE_TIMEWINDOW_LABELS_BY_VALUE = RELATIVE_TIMEWINDOW_OPTIONS.reduce(
  (acc, option) => {
    acc[option.value] = option.label
    return acc
  },
  {} as Record<number, string>,
)

export const DEFAULT_TIMEWINDOW_VALUE = RELATIVE_TIMEWINDOW_OPTIONS[0].value

type LabelRendererProps = Readonly<{
  timewindowType: 'relative' | 'absolute'
  relativeTimewindow: number
  from: Date
  to: Date
  live: boolean
}>

const LabelRenderer = (props: LabelRendererProps) => (
  <Box className="flex flex-row items-center justify-center gap-1">
    {props.live
      ? (
        <>
          <span className="font-semibold">From</span>
          <span>{props.from.toLocaleString()}</span>
        </>
      )
      : (
        props.timewindowType === 'relative'
          ? (
            <span>
              {RELATIVE_TIMEWINDOW_LABELS_BY_VALUE[props.relativeTimewindow] ?? '#ERR#'}
            </span>
          )
          : (
            <>
              <span className="font-semibold">From</span>
              <span>{props.from.toLocaleString()}</span>
              <span className="font-semibold">to</span>
              <span>{props.to.toLocaleString()}</span>
            </>
          )
      )
    }
  </Box>
)

type TimeWindow = {
  from: Date,
  to: Date,
  live: boolean,
}

export type TimeWindowFactory = {
  make: () => TimeWindow,
}

type TimeWindowSelectorProps = Readonly<{
  onTimeWindowChanged: (factory: TimeWindowFactory) => void
}>

export const TimeWindowSelector = ({ onTimeWindowChanged }: TimeWindowSelectorProps) => {
  const now = useMemo(
    () => new Date(),
    [],
  )

  const [menu, setMenu] = useState<HTMLElement | null>(null)

  const [from, setFrom] = useState(new Date(now.getTime() - DEFAULT_TIMEWINDOW_VALUE))
  const [to, setTo] = useState(now)
  const [live, setLive] = useState(false)

  const [timeWindowType, setTimeWindowType] = useState<'relative' | 'absolute'>('relative')
  const [relativeTimeWindow, setRelativeTimeWindow] = useState(DEFAULT_TIMEWINDOW_VALUE)

  const timeWindowFactory = useCallback(
    () => {
      switch (timeWindowType) {
        case 'relative':
          const now = new Date()
          return {
            from: new Date(now.getTime() - relativeTimeWindow),
            to: now,
            live,
          }

        case 'absolute':
          return { from, to, live }

        default:
          throw new Error('Invalid timewindow type')
      }
    },
    [timeWindowType, relativeTimeWindow, from, to, live],
  )

  useEffect(
    () => {
      onTimeWindowChanged({ make: timeWindowFactory })
    },
    [],
  )

  return (
    <>
      <Button
        variant="outlined"
        onClick={(evt) => {
          setMenu(evt.currentTarget)
        }}
        className="w-full h-full"
        endIcon={<ArrowDropDownIcon />}
        sx={{
          '& .MuiButton-icon': {
            marginLeft: 'auto',
          }
        }}
      >
        <LabelRenderer
          timewindowType={timeWindowType}
          relativeTimewindow={relativeTimeWindow}
          from={from}
          to={to}
          live={live}
        />
      </Button>
      <Menu
        anchorEl={menu}
        open={Boolean(menu)}
        onClose={() => {
          setMenu(null)
          onTimeWindowChanged({make: timeWindowFactory})
        }}
        MenuListProps={{
          sx: { width: menu?.offsetWidth },
        }}
      >
        <Box>
          <ToggleButtonGroup
            color="primary"
            value={timeWindowType}
            exclusive
            onChange={(_, value) => {
              if (value !== null) {
                setTimeWindowType(value)
              }
            }}
            className="p-3 w-full"
          >
            <ToggleButton value="relative" className="flex-grow">Relative</ToggleButton>
            <ToggleButton value="absolute" className="flex-grow">Absolute</ToggleButton>
          </ToggleButtonGroup>
          <Divider />

          {timeWindowType === 'relative' && (
            <ToggleButtonGroup
              color="secondary"
              orientation="vertical"
              exclusive
              value={relativeTimeWindow}
              onChange={(_, value) => {
                if (value !== null) {
                  setRelativeTimeWindow(value)
                }
              }}
              className="p-3 w-full"
            >
              {RELATIVE_TIMEWINDOW_OPTIONS.map((option) => (
                <ToggleButton
                  key={option.value}
                  value={option.value}
                >
                  {option.label}
                </ToggleButton>
              ))}
            </ToggleButtonGroup>
          )}
          {timeWindowType === 'absolute' && (
            <Box className="p-3 w-full flex flex-col items-stretch gap-3">
              <DateTimePicker
                label="From"
                value={dayjs(from)}
                onChange={(date) => {
                  if (date !== null) {
                    setFrom(date.toDate())
                  }
                }}
              />

              <DateTimePicker
                label="To"
                value={dayjs(to)}
                onChange={(date) => {
                  if (date !== null) {
                    setTo(date.toDate())
                  }
                }}
              />
            </Box>
          )}

          <Divider />

          <ToggleButtonGroup
            color="info"
            orientation="vertical"
            exclusive
            value={live}
            onChange={(_, value) => {
              setLive(value !== null)
            }}
            className="p-3 w-full"
          >
            <ToggleButton value={true}>
              Watch Logs
            </ToggleButton>
          </ToggleButtonGroup>
        </Box>
      </Menu>
    </>
  )
}

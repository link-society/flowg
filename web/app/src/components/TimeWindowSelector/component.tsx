import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker'
import dayjs from 'dayjs'

import { useCallback, useEffect, useMemo, useState } from 'react'

import Button from '@mui/material/Button'
import Divider from '@mui/material/Divider'
import Menu from '@mui/material/Menu'
import ToggleButton from '@mui/material/ToggleButton'
import Typography from '@mui/material/Typography'

import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import CheckIcon from '@mui/icons-material/Check'

import {
  LabelRow,
  LabelStrong,
  MenuBody,
  MenuPad,
  MenuSection,
  TimeWindowButton,
  TimeWindowToggleButton,
  TimeWindowToggleGroup,
} from './styles'
import {
  LabelRendererProps,
  RelTimeWindowOption,
  TimeWindowSelectorProps,
} from './types'

export type { TimeWindowFactory } from './types'

const RELATIVE_TIMEWINDOW_OPTIONS: RelTimeWindowOption[] = [
  { key: 'last-15min', label: 'Last 15 Minutes', value: 15 * 60 * 1000 },
  { key: 'last-1h', label: 'Last Hour', value: 60 * 60 * 1000 },
  { key: 'last-day', label: 'Last Day', value: 24 * 60 * 60 * 1000 },
  { key: 'last-week', label: 'Last Week', value: 7 * 24 * 60 * 60 * 1000 },
]

const RELATIVE_TIMEWINDOW_LABELS_BY_VALUE = RELATIVE_TIMEWINDOW_OPTIONS.reduce(
  (acc, option) => {
    acc[option.value] = option.label
    return acc
  },
  {} as Record<number, string>
)

export const DEFAULT_TIMEWINDOW_VALUE = RELATIVE_TIMEWINDOW_OPTIONS[0].value

const LabelRenderer = (props: LabelRendererProps) => (
  <LabelRow>
    {props.live ? (
      <>
        <LabelStrong variant="text">From</LabelStrong>
        <Typography variant="text">{props.from.toLocaleString()}</Typography>
      </>
    ) : props.timewindowType === 'relative' ? (
      <Typography variant="text">
        {RELATIVE_TIMEWINDOW_LABELS_BY_VALUE[props.relativeTimewindow] ??
          '#ERR#'}
      </Typography>
    ) : (
      <>
        <LabelStrong variant="text">From</LabelStrong>
        <Typography variant="text">{props.from.toLocaleString()}</Typography>
        <LabelStrong variant="text">to</LabelStrong>
        <Typography variant="text">{props.to.toLocaleString()}</Typography>
      </>
    )}
  </LabelRow>
)

const TimeWindowSelector = ({
  onTimeWindowChanged,
}: TimeWindowSelectorProps) => {
  const now = useMemo(() => new Date(), [])

  const [menu, setMenu] = useState<HTMLElement | null>(null)
  const open = Boolean(menu)

  const [from, setFrom] = useState(
    new Date(now.getTime() - DEFAULT_TIMEWINDOW_VALUE)
  )
  const [to, setTo] = useState(now)
  const [live, setLive] = useState(false)

  const [timeWindowType, setTimeWindowType] = useState<'relative' | 'absolute'>(
    'relative'
  )
  const [relativeTimeWindow, setRelativeTimeWindow] = useState(
    DEFAULT_TIMEWINDOW_VALUE
  )

  const timeWindowFactory = useCallback((): {
    from: Date
    to: Date
    live: boolean
  } => {
    switch (timeWindowType) {
      case 'relative': {
        const now = new Date()
        return {
          from: new Date(now.getTime() - relativeTimeWindow),
          to: now,
          live,
        }
      }

      case 'absolute':
        return { from, to, live }

      default:
        throw new Error('Invalid timewindow type')
    }
  }, [timeWindowType, relativeTimeWindow, from, to, live])

  const handleClose = () => {
    setMenu(null)
    onTimeWindowChanged({ make: timeWindowFactory })
  }

  useEffect(() => {
    onTimeWindowChanged({ make: timeWindowFactory })
  }, [])

  return (
    <>
      <TimeWindowButton
        id="btn:streams.timewindow-selector.open"
        variant="outlined"
        onClick={(evt) => {
          setMenu(evt.currentTarget)
        }}
        endIcon={<ArrowDropDownIcon />}
      >
        <LabelRenderer
          timewindowType={timeWindowType}
          relativeTimewindow={relativeTimeWindow}
          from={from}
          to={to}
          live={live}
        />
      </TimeWindowButton>
      <Menu
        id="menu:streams.timewindow-selector"
        anchorEl={menu}
        open={open}
        onClose={handleClose}
        slotProps={{
          list: {
            style: { width: menu?.offsetWidth },
          },
        }}
      >
        <MenuBody>
          <TimeWindowToggleGroup
            color="primary"
            value={timeWindowType}
            exclusive
            onChange={(_, value) => {
              if (value !== null) {
                setTimeWindowType(value)
              }
            }}
          >
            <TimeWindowToggleButton
              id="btn:streams.timewindow-selector.type.relative"
              value="relative"
            >
              Relative
            </TimeWindowToggleButton>
            <TimeWindowToggleButton
              id="btn:streams.timewindow-selector.type.absolute"
              value="absolute"
            >
              Absolute
            </TimeWindowToggleButton>
          </TimeWindowToggleGroup>
          <Divider />

          {timeWindowType === 'relative' && (
            <TimeWindowToggleGroup
              color="secondary"
              orientation="vertical"
              exclusive
              value={relativeTimeWindow}
              onChange={(_, value) => {
                if (value !== null) {
                  setRelativeTimeWindow(value)
                }
              }}
            >
              {RELATIVE_TIMEWINDOW_OPTIONS.map((option) => (
                <ToggleButton
                  id={`btn:streams.timewindow-selector.relative.${option.key}`}
                  key={option.key}
                  value={option.value}
                >
                  {option.label}
                </ToggleButton>
              ))}
            </TimeWindowToggleGroup>
          )}
          {timeWindowType === 'absolute' && (
            <MenuSection>
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
            </MenuSection>
          )}

          <Divider />

          <TimeWindowToggleGroup
            color="info"
            orientation="vertical"
            exclusive
            value={live}
            onChange={(_, value) => {
              setLive(value !== null)
            }}
          >
            <ToggleButton
              id="btn:streams.timewindow-selector.live"
              value={true}
            >
              Watch Logs
            </ToggleButton>
          </TimeWindowToggleGroup>

          <MenuPad>
            <Button
              id="btn:streams.timewindow-selector.apply"
              variant="contained"
              color="primary"
              fullWidth
              endIcon={<CheckIcon />}
              onClick={handleClose}
            >
              Apply
            </Button>
          </MenuPad>
        </MenuBody>
      </Menu>
    </>
  )
}

export default TimeWindowSelector

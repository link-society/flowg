import InputBase from '@mui/material/InputBase'
import Select from '@mui/material/Select'
import { alpha, styled } from '@mui/material/styles'

import { BorderState } from './types'

export const DynamicFieldGroupBorder = styled('fieldset')<{
  ownerState: BorderState
}>(({ theme, ownerState }) => ({
  margin: 0,
  padding: '0 8px',
  position: 'absolute',
  inset: '-5px 0 0',
  border: `${ownerState.focused ? 2 : 1}px solid ${
    ownerState.error
      ? theme.palette.error.main
      : ownerState.focused
        ? theme.palette.primary.main
        : ownerState.disabled
          ? alpha(theme.palette.text.primary, 0.26)
          : alpha(theme.palette.text.primary, 0.23)
  }`,
  borderRadius: `${theme.shape.borderRadius}px`,
  pointerEvents: 'none',

  '& legend': {
    display: 'block',
    float: 'unset',
    width: 'auto',
    padding: 0,
    height: 11,
    fontSize: '0.75em',
    visibility: 'hidden',
    overflow: 'hidden',
    maxWidth: ownerState.shrunk ? '100%' : '0.01px',
    transition: 'max-width 100ms cubic-bezier(0.0, 0, 0.2, 1) 50ms',
    whiteSpace: 'nowrap',

    '& > span': {
      paddingLeft: '5px',
      paddingRight: '5px',
      display: 'inline-block',
      opacity: 0,
      visibility: 'visible',
    },
  },
}))

export const DynamicFieldGroupRoot = styled('div')({
  position: 'relative',
  display: 'flex',
  alignItems: 'stretch',
})

const borderlessSelectSx = {
  '& .MuiOutlinedInput-notchedOutline': { border: 'none' },
  '&:hover .MuiOutlinedInput-notchedOutline': { border: 'none' },
  '&.Mui-focused .MuiOutlinedInput-notchedOutline': { border: 'none' },
} as const

export const DynamicFieldModeSelect = styled(Select)(({ theme }) => ({
  ...borderlessSelectSx,
  '& .MuiSelect-select': {
    display: 'flex',
    alignItems: 'center',
    paddingLeft: theme.spacing(1.5),
    paddingRight: `${theme.spacing(3.5)} !important`,
    paddingTop: theme.spacing(1),
    paddingBottom: theme.spacing(1),
  },
}))

export const DynamicFieldInput = styled(InputBase)(({ theme }) => ({
  flex: 1,
  '& .MuiInputBase-input': {
    padding: theme.spacing(1, 1.5),
  },
  '& .MuiInputBase-inputMultiline': {
    padding: 0,
  },
}))

export const DynamicFieldValueSelect = styled(Select)(({ theme }) => ({
  ...borderlessSelectSx,
  flex: 1,
  '& .MuiSelect-select': {
    padding: theme.spacing(1, 1.5),
  },
}))

import Box from '@mui/material/Box'
import Button from '@mui/material/Button'
import ToggleButton from '@mui/material/ToggleButton'
import ToggleButtonGroup from '@mui/material/ToggleButtonGroup'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const LabelRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'center',
  gap: theme.spacing(0.5),
}))

export const LabelStrong = styled(Typography)({
  fontWeight: 600,
})

export const MenuSection = styled(Box)(({ theme }) => ({
  padding: theme.spacing(1.5),
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
}))

export const MenuPad = styled(Box)(({ theme }) => ({
  padding: theme.spacing(1.5),
  width: '100%',
}))

export const MenuBody = styled(Box)({
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const TimeWindowButton = styled(Button)({
  width: '100%',
  height: '100%',
  '& .MuiButton-icon': { marginLeft: 'auto' },
})

export const TimeWindowToggleGroup = styled(ToggleButtonGroup)(({ theme }) => ({
  padding: theme.spacing(1.5),
  width: '100%',
}))

export const TimeWindowToggleButton = styled(ToggleButton)({
  flexGrow: 1,
})

import Box from '@mui/material/Box'
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

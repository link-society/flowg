import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const LabelRow = styled(Box)({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'center',
  gap: '0.25rem',
})

export const LabelStrong = styled(Typography)({
  fontWeight: 600,
})

export const MenuSection = styled(Box)({
  padding: '0.75rem',
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.75rem',
})

export const MenuPad = styled(Box)({
  padding: '0.75rem',
  width: '100%',
})

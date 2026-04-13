import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const FormStack = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.75rem',
})

export const FieldRow = styled(Box)({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'flex-end',
})

export const FieldStack = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.5rem',
})

export const FieldLabel = styled(Typography)({
  fontWeight: 600,
})

import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const TraceSection = styled('div')(({ theme }) => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1),
}))

export const TraceRow = styled('div')(({ theme }) => ({
  display: 'flex',
  gap: theme.spacing(2.5),
}))

export const TraceLabel = styled(Typography)(({ theme }) => ({
  color: theme.tokens.colors.labelText,
  fontWeight: 600,
}))

export const TraceCode = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(1),
  flex: 1,
  flexShrink: 1,
  overflow: 'auto',
  fontFamily: 'monospace',
  backgroundColor: theme.tokens.colors.codeBg,
  minWidth: 200,
}))

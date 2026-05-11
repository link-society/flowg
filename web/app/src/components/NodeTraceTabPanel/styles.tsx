import Paper from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const TraceSection = styled('div')({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  gap: 8,
})

export const TraceRow = styled('div')({
  display: 'flex',
  gap: 20,
})

export const TraceLabel = styled(Typography)({
  fontSize: '0.875rem',
  color: '#374151',
  fontWeight: 600,
  marginBottom: 8,
})

export const TraceCode = styled(Paper)({
  padding: 8,
  flex: 1,
  flexShrink: 1,
  overflow: 'auto',
  fontFamily: 'monospace',
  backgroundColor: '#f3f4f6',
  minWidth: 200,
})

import Box from '@mui/material/Box'
import Paper, { PaperProps } from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const TraceDialogContent = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  gap: '1.25rem',
})

export const TraceSection = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  gap: '0.5rem',
})

export const TraceColumns = styled(Box)({
  display: 'flex',
  gap: '1.25rem',
})

export const TraceLabel = styled(Typography)(({ theme }) => ({
  fontSize: '0.875rem',
  color: theme.palette.text.secondary,
  fontWeight: 600,
  marginBottom: '0.5rem',
}))

export const TracePaper = styled(Paper)<PaperProps<'pre'>>(({ theme }) => ({
  padding: '0.5rem',
  flexGrow: 1,
  flexShrink: 1,
  overflow: 'auto',
  fontFamily: 'monospace',
  backgroundColor: theme.palette.grey[100],
  minWidth: '12.5rem',
}))

import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const TraceDialogContent = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2.5),
}))

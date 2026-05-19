import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const Root = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1.5),
}))

export const Row = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
  '& .MuiTextField-root': { flexGrow: 1 },
}))

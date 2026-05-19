import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const InputListRoot = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1),
}))

export const InputListRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  gap: theme.spacing(1),
  '& .MuiTextField-root': { flexGrow: 1 },
}))

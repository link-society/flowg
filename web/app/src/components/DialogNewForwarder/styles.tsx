import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const DialogFormBody = styled(Box)(({ theme }) => ({
  paddingTop: theme.spacing(1.5),
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
}))

export const TypeOption = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1),
}))

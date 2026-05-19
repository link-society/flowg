import Box from '@mui/material/Box'
import Dialog from '@mui/material/Dialog'
import Tabs from '@mui/material/Tabs'
import { styled } from '@mui/material/styles'

export const TraceDialogContent = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(2.5),
}))

export const TraceDialog = styled(Dialog)({
  '& .MuiDialog-paper': {
    width: '80%',
    height: '90%',
  },
})

export const TraceTabs = styled(Tabs)(({ theme }) => ({
  borderBottom: `1px solid ${theme.palette.divider}`,
}))

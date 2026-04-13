import Box from '@mui/material/Box'
import Paper from '@mui/material/Paper'
import { styled } from '@mui/material/styles'

export const TransferRoot = styled(Box)({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'center',
  gap: '1rem',
})

export const TransferColumn = styled(Box)({
  flex: 1,
  minWidth: 0,
})

export const TransferControls = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  flexShrink: 0,
})

export const TransferListPaper = styled(Paper)({
  minHeight: '15rem',
  maxHeight: '15rem',
  overflow: 'auto',
})

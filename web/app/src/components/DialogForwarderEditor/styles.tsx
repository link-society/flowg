import { Box, Toolbar, styled } from '@mui/material'

export const DialogToolbar = styled(Toolbar)({
  gap: '0.75rem',
})

export const DialogToolbarName = styled('div')({
  flexGrow: 1,
})

export const DialogBody = styled(Box)({
  flex: 1,
  padding: '1.5rem',
  overflow: 'auto',
})

export const DialogLoading = styled(Box)({
  width: '100%',
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
})

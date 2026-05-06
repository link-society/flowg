import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import { styled } from '@mui/material/styles'

import AppContainer from '@/components/AppContainer/component'

export const UploadViewRoot = styled(AppContainer)({
  width: 'calc(100% / 3)',
  margin: 'auto',
  display: 'flex',
  flexDirection: 'column',
})

export const UploadViewHeader = styled(Box)({
  marginBottom: '1.5rem',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'center',
  gap: '0.5rem',
})

export const UploadViewCard = styled(Card)({
  padding: '0.75rem',
})

export const UploadViewForm = styled('form')({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.75rem',
})

export const VisuallyHiddenInput = styled('input')({
  clip: 'rect(0 0 0 0)',
  clipPath: 'inset(50%)',
  height: 1,
  overflow: 'hidden',
  position: 'absolute',
  bottom: 0,
  left: 0,
  whiteSpace: 'nowrap',
  width: 1,
})

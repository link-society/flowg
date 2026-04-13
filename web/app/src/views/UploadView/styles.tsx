import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import { styled } from '@mui/material/styles'

import AppContainer from '@/components/AppContainer/component'

export const UploadViewRoot = styled(AppContainer)({
  flexDirection: 'column',
})

export const UploadViewHeader = styled(Box)({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'center',
  gap: '0.5rem',
})

export const UploadViewCard = styled(Card)({
  padding: '0.75rem',
  maxWidth: 400,
  width: '100%',
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

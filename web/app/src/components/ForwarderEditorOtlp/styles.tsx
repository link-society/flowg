import { TextField, styled } from '@mui/material'

export const ForwarderEditorOtlpRoot = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.75rem',
})

export const ForwarderEditorOtlpEndpointField = styled(TextField)({
  width: '100%',
})

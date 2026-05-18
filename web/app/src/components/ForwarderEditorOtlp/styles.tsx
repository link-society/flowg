import { TextField, styled } from '@mui/material'

export const ForwarderEditorOtlpRoot = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
}))

export const ForwarderEditorOtlpEndpointField = styled(TextField)({
  width: '100%',
})

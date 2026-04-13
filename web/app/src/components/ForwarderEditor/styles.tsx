import { TextField, styled } from '@mui/material'

export const ForwarderEditorTypeField = styled('div')(({ theme }) => ({
  marginBottom: '1.5rem',
  boxShadow: theme.shadows[1],
}))

export const ForwarderEditorTextField = styled(TextField)({
  width: '100%',
})

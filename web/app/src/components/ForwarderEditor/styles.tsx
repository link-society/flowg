import { TextField, styled } from '@mui/material'

export const ForwarderEditorTypeField = styled('div')`
  margin-bottom: 1.5rem;
  box-shadow: ${({ theme }) => theme.shadows[1]};
`

export const ForwarderEditorTextField = styled(TextField)`
  width: 100%;
`

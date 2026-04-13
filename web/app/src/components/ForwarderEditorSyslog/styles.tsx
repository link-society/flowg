import { FormControl, TextField, styled } from '@mui/material'

export const ForwarderEditorSyslogRoot = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.75rem',
})

export const ForwarderEditorSyslogRow = styled('div')({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: '0.75rem',
})

export const ForwarderEditorSyslogColumn = styled('div')({
  display: 'flex',
  flexDirection: 'row',
  gap: '0.75rem',
})

export const ForwarderEditorSyslogAddressField = styled(TextField)({
  flexGrow: 1,
})

export const ForwarderEditorSyslogGrowFormControl = styled(FormControl)({
  flexGrow: 1,
})

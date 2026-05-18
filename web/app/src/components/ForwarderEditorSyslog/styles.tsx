import { FormControl, TextField, styled } from '@mui/material'

export const ForwarderEditorSyslogRoot = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
}))

export const ForwarderEditorSyslogRow = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const ForwarderEditorSyslogColumn = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  gap: theme.spacing(1.5),
}))

export const ForwarderEditorSyslogAddressField = styled(TextField)({
  flexGrow: 1,
})

export const ForwarderEditorSyslogGrowFormControl = styled(FormControl)({
  flexGrow: 1,
})

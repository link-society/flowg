import { FormControl, TextField, styled } from '@mui/material'

export const ForwarderEditorSyslogRoot = styled('div')`
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 0.75rem;
`

export const ForwarderEditorSyslogRow = styled('div')`
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 0.75rem;
`

export const ForwarderEditorSyslogColumn = styled('div')`
  display: flex;
  flex-direction: row;
  gap: 0.75rem;
`

export const ForwarderEditorSyslogAddressField = styled(TextField)`
  flex-grow: 1;
`

export const ForwarderEditorSyslogGrowFormControl = styled(FormControl)`
  flex-grow: 1;
`

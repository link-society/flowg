import { TextField, styled } from '@mui/material'

export const ForwarderEditorGoogleCloudLoggingRoot = styled('div')(
  ({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'stretch',
    gap: theme.spacing(1.5),
  })
)

export const ForwarderEditorGoogleCloudLoggingRow = styled('div')(
  ({ theme }) => ({
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    gap: theme.spacing(1.5),
  })
)

export const ForwarderEditorGoogleCloudLoggingField = styled(TextField)({
  flexGrow: 1,
})

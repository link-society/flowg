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

export const ForwarderEditorGoogleCloudLoggingJsonField = styled('fieldset')(
  ({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    margin: 0,
    padding: theme.spacing(0.75, 1.5, 1.5),
    border: `1px solid ${
      theme.palette.mode === 'dark'
        ? 'rgba(255, 255, 255, 0.23)'
        : 'rgba(0, 0, 0, 0.23)'
    }`,
    borderRadius: theme.shape.borderRadius,
    overflow: 'hidden',
    transition: theme.transitions.create('border-color'),

    '&:hover': {
      borderColor: theme.palette.text.primary,
    },

    '&:focus-within': {
      borderColor: theme.palette.primary.main,

      '& legend': {
        color: theme.palette.primary.main,
      },
    },
  })
)

export const ForwarderEditorGoogleCloudLoggingJsonLabel = styled('legend')(
  ({ theme }) => ({
    padding: theme.spacing(0, 0.75),
    fontSize: '0.75rem',
    color:
      theme.palette.mode === 'dark'
        ? 'rgba(255, 255, 255, 0.7)'
        : 'rgba(0, 0, 0, 0.6)',
  })
)

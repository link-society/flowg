import Box from '@mui/material/Box'
import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const FormStack = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
}))

export const FieldStack = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1),
}))

export const FieldLabel = styled(Typography)({
  fontWeight: 600,
})

export const FormTextField = styled(TextField)(({ theme }) => ({
  marginTop: theme.spacing(2),
}))

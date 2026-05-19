import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const FormStack = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
}))

export const FieldRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'flex-end',
  '& > .MuiSvgIcon-root': {
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(0.5),
    marginBottom: theme.spacing(0.5),
  },
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

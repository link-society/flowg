import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const Label = styled(Typography)({
  fontWeight: 600,
  marginBottom: 8,
})

export const FormRow = styled('form')({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: 8,
})

export const IconBox = styled(Box)({
  flexGrow: 0,
})

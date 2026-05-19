import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const FieldRow = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'flex-end',
  marginBottom: theme.spacing(1.5),
  '& > .MuiSvgIcon-root': {
    marginRight: theme.spacing(1),
    marginTop: theme.spacing(0.5),
    marginBottom: theme.spacing(0.5),
  },
}))

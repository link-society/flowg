import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const Label = styled(Typography)(({ theme }) => ({
  fontWeight: 600,
  marginBottom: theme.spacing(0.5),
  display: 'block',
}))

export const PermissionGrid = styled(Box)(({ theme }) => ({
  padding: theme.spacing(0.5),
  display: 'block',
  [theme.breakpoints.up('md')]: {
    display: 'grid',
    gridTemplateColumns: 'repeat(4, 1fr)',
    gap: theme.spacing(0.5),
  },
}))

export const PermissionLabel = styled(Typography)({
  fontSize: '0.875rem',
})

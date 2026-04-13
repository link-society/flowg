import Box from '@mui/material/Box'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const Label = styled(Typography)({
  fontWeight: 600,
  marginBottom: 4,
  display: 'block',
})

export const PermissionGrid = styled(Box)(({ theme }) => ({
  padding: 4,
  display: 'block',
  [theme.breakpoints.up('md')]: {
    display: 'grid',
    gridTemplateColumns: 'repeat(4, 1fr)',
    gap: 4,
  },
}))

export const PermissionLabel = styled(Typography)({
  fontSize: '0.875rem',
})

import Paper, { PaperProps } from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const Label = styled(Typography)(({ theme }) => ({
  fontWeight: 600,
  marginBottom: theme.spacing(0.5),
  display: 'block',
}))

export const RolesPaper = styled(Paper)<PaperProps<'ul'>>(({ theme }) => ({
  padding: theme.spacing(0.5),
  display: 'flex',
  flexDirection: 'row',
  justifyContent: 'center',
  flexWrap: 'wrap',
  listStyle: 'none',
}))

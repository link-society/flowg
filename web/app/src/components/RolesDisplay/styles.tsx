import Paper, { PaperProps } from '@mui/material/Paper'
import Typography from '@mui/material/Typography'
import { styled } from '@mui/material/styles'

export const Label = styled(Typography)({
  fontWeight: 600,
  marginBottom: 4,
  display: 'block',
})

export const RolesPaper = styled(Paper)<PaperProps<'ul'>>({
  padding: 4,
  display: 'flex',
  flexDirection: 'row',
  justifyContent: 'center',
  flexWrap: 'wrap',
  listStyle: 'none',
})

import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const DragHandleRoot = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  alignSelf: 'stretch',
  padding: theme.spacing(0, 0.25),
  color: theme.tokens.colors.mutedText,
  cursor: 'grab',
  touchAction: 'none',
  '&:hover': {
    color: theme.tokens.colors.labelText,
  },
  '&:active': {
    cursor: 'grabbing',
  },
}))

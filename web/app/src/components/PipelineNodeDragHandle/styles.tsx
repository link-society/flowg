import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const DragHandleRoot = styled(Box)(({ theme }) => ({
  position: 'absolute',
  right: 'calc(100% + 8px)',
  top: '50%',
  transform: 'translateY(-50%)',
  width: 28,
  height: 28,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  backgroundColor: theme.palette.background.paper,
  border: '1px solid',
  borderRadius: theme.spacing(0.75),
  cursor: 'grab',
  touchAction: 'none',
  opacity: 0,
  pointerEvents: 'none',
  transition: 'opacity 0.15s ease',
  '&:active': {
    cursor: 'grabbing',
  },
}))

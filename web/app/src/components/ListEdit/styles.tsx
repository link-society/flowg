import Box from '@mui/material/Box'
import { styled } from '@mui/material/styles'

export const Root = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  gap: 12,
})

export const Row = styled(Box)({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  gap: 12,
})

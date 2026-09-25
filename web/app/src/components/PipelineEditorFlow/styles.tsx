import { Paper, styled } from '@mui/material'

export const FlowRoot = styled(Paper)({
  width: '100%',
  height: '100%',
  position: 'relative',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'stretch',
  overflow: 'hidden',

  '.react-flow__panel': {
    backgroundColor: 'transparent',
  },
})

export const FlowCanvasWrap = styled('div')({
  flex: 1,
  minWidth: 0,
  height: '100%',
  position: 'relative',
  backgroundColor: 'transparent',
})

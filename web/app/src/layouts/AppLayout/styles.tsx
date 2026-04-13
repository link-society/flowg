import { styled } from '@mui/material'

export const AppLayoutContainer = styled('div')({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  overflow: 'hidden',
  '& > main': {
    flexGrow: 1,
    flexShrink: 1,
    display: 'flex',
    flexDirection: 'column',
    overflowY: 'auto',
  },
})

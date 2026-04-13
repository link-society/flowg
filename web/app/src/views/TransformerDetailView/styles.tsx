import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const TransformerDetailViewContainer = styled('div')({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
})

export const TransformerDetailViewToolbar = styled(AppContainer)(
  ({ theme }) => ({
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'stretch',
    color: 'white',
    backgroundColor: theme.tokens.colors.editorToolbarBkg,
    zIndex: 10,
    flex: 0,
    boxShadow:
      '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
  })
)

export const TransformerDetailViewToolbarLeft = styled('div')({
  flexGrow: 1,
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: 12,
})

export const TransformerDetailViewToolbarRight = styled('div')({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: 12,
})

export const TransformerDetailViewBody = styled(AppContainer)({
  flexGrow: 1,
  flexShrink: 1,
  height: 0,
  display: 'flex',
  flexDirection: 'column',
  gap: 8,
  '@media (min-width: 990px)': {
    flexDirection: 'row',
  },
})

export const TransformerDetailViewSidebar = styled('div')({
  flex: '0 0 16.67%',
  height: '100%',
  width: '100%',
  '@media (min-width: 990px)': {
    flexDirection: 'row',
  },
})

export const TransformerDetailViewEditor = styled('div')({
  flex: 1,
  height: '100%',
  flexDirection: 'column',
  '@media (min-width: 990px)': {
    flexDirection: 'row',
  },
  '& > div': {
    display: 'flex',
    flexDirection: 'column',
    '@media (min-width: 990px)': {
      flexDirection: 'row',
    },
  },
})

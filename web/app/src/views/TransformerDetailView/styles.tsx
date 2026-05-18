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
    color: theme.tokens.colors.primaryContrast,
    backgroundColor: theme.tokens.colors.editorToolbarBkg,
    zIndex: 10,
    flex: 0,
    boxShadow: theme.tokens.shadows.md,
  })
)

export const TransformerDetailViewToolbarLeft = styled('div')(({ theme }) => ({
  flexGrow: 1,
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const TransformerDetailViewToolbarRight = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  gap: theme.spacing(1.5),
}))

export const TransformerDetailViewBody = styled(AppContainer)(({ theme }) => ({
  flexGrow: 1,
  flexShrink: 1,
  height: 0,
  display: 'flex',
  flexDirection: 'column',
  gap: theme.spacing(1),
  [theme.breakpoints.up('md')]: {
    flexDirection: 'row',
  },
}))

export const TransformerDetailViewSidebar = styled('div')(({ theme }) => ({
  flex: '0 0 16.67%',
  height: '100%',
  width: '100%',
  [theme.breakpoints.up('md')]: {
    flexDirection: 'row',
  },
}))

export const TransformerDetailViewEditor = styled('div')(({ theme }) => ({
  flex: 1,
  height: '100%',
  flexDirection: 'column',
  [theme.breakpoints.up('md')]: {
    flexDirection: 'row',
  },
  '& > div': {
    display: 'flex',
    flexDirection: 'column',
    [theme.breakpoints.up('md')]: {
      flexDirection: 'row',
    },
  },
}))

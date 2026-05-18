import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const HomeViewContainer = styled(AppContainer)(({ theme }) => ({
  placeContent: 'space-evenly',
  flexDirection: 'column',
  '& h1': {
    display: 'flex',
    gap: theme.spacing(1),
    placeItems: 'center',
    img: {
      height: '2rem',
    },
  },
}))

export const HomeViewPermissionsWrapper = styled('div')(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: theme.spacing(2),
  width: '100%',
  [theme.breakpoints.up('md')]: {
    gridAutoColumns: '1fr',
    gridAutoFlow: 'column',
    gridTemplateColumns: 'unset',
    width: 'auto',
  },
}))

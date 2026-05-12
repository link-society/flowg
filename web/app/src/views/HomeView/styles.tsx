import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const HomeViewContainer = styled(AppContainer)({
  placeContent: 'space-evenly',
  flexDirection: 'column',
  '& h1': {
    display: 'flex',
    gap: 8,
    placeItems: 'center',
    img: {
      height: '2rem',
    },
  },
})

export const HomeViewPermissionsWrapper = styled('div')({
  display: 'grid',
  gridTemplateColumns: '1fr',
  gap: 16,
  width: '100%',
  '@media (min-width: 990px)': {
    gridAutoColumns: '1fr',
    gridAutoFlow: 'column',
    gridTemplateColumns: 'unset',
    width: 'auto',
  },
})

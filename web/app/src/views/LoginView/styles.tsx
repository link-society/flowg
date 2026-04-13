import { Card, styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const LoginViewContainer = styled(AppContainer)({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  height: '100%',
  '& > header': {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: '0.5rem',
    '& img': {
      height: '4rem',
    },
    '& h1': {
      fontSize: '3rem',
      lineHeight: 1,
      fontWeight: 700,
      textAlign: 'center',
    },
  },
})

export const LoginViewCard = styled(Card)({
  width: '100%',
  maxWidth: '28rem',
  padding: '0.75rem',
  minWidth: 550,
  '& > form': {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'stretch',
    gap: '0.75rem',
    '& header h2': {
      fontSize: '1.5rem',
      lineHeight: '2rem',
      textAlign: 'center',
      fontWeight: 400,
    },
  },
  '@media (max-width: 990px)': {
    minWidth: '100%',
  },
})

export const LoginViewCardFields = styled('section')({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: '0.75rem',
  '& > div': {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'flex-end',
    '& .icon': {
      marginRight: 8,
      marginTop: 4,
      marginBottom: 4,
    },
  },
})

import { Card, styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const LoginViewContainer = styled(AppContainer)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  height: '100%',
  '& > header': {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: theme.spacing(1),
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
}))

export const LoginViewCard = styled(Card)(({ theme }) => ({
  width: '100%',
  maxWidth: '28rem',
  padding: theme.spacing(1.5),
  minWidth: 550,
  '& > form': {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'stretch',
    gap: theme.spacing(1.5),
    '& header h2': {
      fontSize: '1.5rem',
      lineHeight: '2rem',
      textAlign: 'center',
      fontWeight: 400,
    },
  },
  [theme.breakpoints.down('md')]: {
    minWidth: '100%',
  },
}))

export const LoginViewCardFields = styled('section')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'stretch',
  gap: theme.spacing(1.5),
  '& > div': {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'flex-end',
    '& .icon': {
      marginRight: theme.spacing(1),
      marginTop: theme.spacing(0.5),
      marginBottom: theme.spacing(0.5),
    },
  },
}))

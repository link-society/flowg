import { styled } from '@mui/material'

import AppContainer from '@/components/AppContainer/component'

export const ForwarderSectionViewRoot = styled(AppContainer)({
  placeContent: 'center',
  flexDirection: 'column',
  '& h1': {
    fontWeight: 700,
    textAlign: 'center',
  },
})

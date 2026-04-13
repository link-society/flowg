import { styled } from '@mui/material'

import List from '@mui/material/List'
import Paper from '@mui/material/Paper'

export const SideNavListContainer = styled(Paper)({
  height: '100%',
  overflow: 'auto',
})

export const SideNavListNav = styled(List)(({ theme }) => ({
  padding: '0 !important',
  '& .MuiListItemButton-root': {
    color: theme.palette.secondary.main,
  },
  '& .MuiListItemButton-root.active': {
    backgroundColor: theme.palette.secondary.main,
    '&:hover': {
      backgroundColor: theme.palette.secondary.main,
    },
  },
  '& .MuiButtonBase-root.active span': {
    color: theme.tokens.colors.white,
    fontWeight: 500,
  },
}))

import { styled } from '@mui/material'

import List from '@mui/material/List'
import Paper from '@mui/material/Paper'

export const SideNavListContainer = styled(Paper)`
  height: 100%;
  overflow: auto;
`

export const SideNavListNav = styled(List)`
  padding: 0 !important;

  .MuiListItemButton-root {
    color: ${({ theme }) => theme.palette.secondary.main};
  }

  .MuiListItemButton-root.active {
    background-color: ${({ theme }) => theme.palette.secondary.main};
    color: white;

    &:hover {
      background-color: ${({ theme }) => theme.palette.secondary.main};
    }
  }
`

import { styled } from '@mui/material'

import Drawer from '@mui/material/Drawer'
import Paper from '@mui/material/Paper'

export const LogTableContainer = styled(Paper)`
  flex-grow: 1;
`

export const LogTableDrawer = styled(Drawer)`
  & .MuiDrawer-paper {
    width: 33vw;
    padding: 0.75rem;
  }
`

export const LogTableDetailPre = styled('pre')`
  padding: 0.5rem;
  width: 100%;
  overflow: auto;
  font-family: monospace;
  background-color: ${({ theme }) => theme.tokens.colors.codeBg};
  border: 1px solid rgba(0, 0, 0, 0.12);
  border-radius: 4px;
`

import { styled } from '@mui/material'

export const StyledAppLayout = styled('div')`
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  > main {
    flex-grow: 1;
    flex-shrink: 1;
    height: 0px;
  }
`

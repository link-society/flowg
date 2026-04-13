import { styled } from '@mui/material'

export const AppLayoutContainer = styled('div')`
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  > main {
    flex-grow: 1;
    flex-shrink: 1;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    /* padding: 1.5rem; */
  }
`

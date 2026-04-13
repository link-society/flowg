import { styled } from '@mui/material'

export const HomeViewContainer = styled('section')`
  flex: 1;
  display: flex;
  flex-direction: column;
  place-content: space-evenly;
  place-items: center;
  gap: 16px;

  h1 {
    display: flex;
    gap: 8px;
  }
`

export const HomeViewPermissionsWrapper = styled('div')`
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;

  @media (min-width: 990px) {
    grid-auto-columns: 1fr;
    grid-auto-flow: column;
    grid-template-columns: unset;
  }
`

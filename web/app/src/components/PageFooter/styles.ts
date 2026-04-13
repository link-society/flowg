import { styled } from '@mui/material'

export const PageFooterContainer = styled('footer')`
  display: flex;
  padding: 0.75rem;
  background-color: ${({ theme }) => theme.tokens.colors.borderLight};

  > div {
    margin-left: auto;
  }
`

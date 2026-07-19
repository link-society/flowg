import { Typography, styled } from '@mui/material'

export const PageFooterContainer = styled('footer')`
  display: flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing(1)};
  padding: ${({ theme }) => theme.spacing(1.5)};
  background-color: ${({ theme }) => theme.tokens.colors.borderLight};
  place-content: space-between;
`

export const PageFooterActions = styled('div')`
  display: flex;
  > div:nth-child(2) {
    margin: 0 5px 0 10px;
    background-color: ${({ theme }) => theme.tokens.colors.borderLight};
  }
`

export const PageFooterText = styled(Typography)({
  fontWeight: 600,
})

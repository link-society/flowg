import { Typography, styled } from '@mui/material'

export const PageFooterContainer = styled('footer')`
  display: flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing(1)};
  padding: ${({ theme }) => theme.spacing(1.5)};
  background-color: ${({ theme }) => theme.tokens.colors.borderLight};

  > div {
    margin-left: auto;
  }
`

export const PageFooterText = styled(Typography)({
  fontWeight: 600,
})

import { CardContent, styled } from '@mui/material'

export const StatCardHeaderWrapper = styled('div')`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  font-size: 1.5rem;
  line-height: 2rem;
  font-weight: 600;
`

export const StatCardContent = styled(CardContent)`
  padding: 0;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 10px;
`

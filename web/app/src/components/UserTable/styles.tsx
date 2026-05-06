import { Card, CardContent, CardHeader, styled } from '@mui/material'

export const UserTableCard = styled(Card)`
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: stretch;

  @media (max-width: 1200px) {
    min-height: 24rem;
  }
`

export const UserTableCardHeader = styled(CardHeader)`
  background-color: ${({ theme }) => theme.tokens.colors.headerCardBkg};
  color: ${({ theme }) => theme.tokens.colors.primaryContrast};
  box-shadow: ${({ theme }) => theme.shadows[4]};
  z-index: 20;
`

export const UserTableCardHeaderTitle = styled('div')`
  display: flex;
  align-items: center;
  gap: 0.75rem;
`

export const UserTableCardHeaderTitleText = styled('span')`
  flex: 1;
`

export const UserTableCardContent = styled(CardContent)`
  padding: 0 !important;
  flex: 1 1 0;
  overflow: hidden;

  .ag-cell-wrapper {
    height: auto !important;
  }

  .ag-cell {
    display: flex;
    align-items: center;
  }

  .flowg-actions-cell {
    justify-content: center;
  }
`

export const RolesCellRoot = styled('div')`
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  padding: 8px 0;
  align-content: center;
  width: 100%;
`

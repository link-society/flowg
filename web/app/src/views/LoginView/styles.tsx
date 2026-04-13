import { Card, styled } from '@mui/material'

export const LoginViewContainer = styled('div')`
  padding-top: 1.5rem;
  padding-bottom: 1.5rem;
  flex: 1;
  place-items: center;

  > header {
    margin-bottom: 1.5rem;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;

    img {
      height: 4rem;
    }
    h1 {
      font-size: 3rem;
      line-height: 1;
      font-weight: 700;
      text-align: center;
    }
  }
`

export const LoginViewCard = styled(Card)`
  width: 100%;
  max-width: 28rem;
  padding: 0.75rem;

  > form {
    display: flex;
    flex-direction: column;
    align-items: stretch;
    gap: 0.75rem;

    header {
      h2 {
        font-size: 1.5rem;
        line-height: 2rem;
        text-align: center;
        font-weight: 400;
      }
    }
  }
`
export const LoginViewCardFields = styled('section')`
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 0.75rem;

  > div {
    display: flex;
    flex-direction: row;
    align-items: flex-end;
    .icon {
      margin-right: 8px;
      margin-top: 4px;
      margin-bottom: 4px;
    }
  }
`

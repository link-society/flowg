import { Button, IconButton, Link, styled } from '@mui/material'

export const NavBarToolbar = styled('div')`
  background-color: ${({ theme }) => theme.tokens.colors.black};
  display: flex;
  align-items: center;
  padding: 0 16px;
  min-height: 64px;
  width: 100%;
`

export const NavBarLeftSection = styled('section')`
  height: 100%;
  display: flex;
  flex-direction: row;
  align-items: stretch;
  gap: 0.75rem;
  flex-grow: 1;

  img {
    height: 2rem;
  }
`

export const NavBarRightSection = styled('section')`
  height: 100%;
  display: flex;
  flex-direction: row-reverse;
  align-items: stretch;
  gap: 0.75rem;
`

export const NavBarDesktopMenus = styled('div')`
  display: contents;
`

export const NavBarDrawerClose = styled(IconButton)`
  align-self: flex-end;
  margin: 8px;
`

export const NavBarDrawerList = styled('div')`
  width: 240px;
`

export const NavBarButton = styled(Button)`
  display: flex;
  gap: 8px;
  width: auto;
  text-transform: none;

  @media (max-width: 990px) {
    padding: 8px;
    min-width: fit-content;
  }

  .nav-text {
    font-size: 14px;
    display: none;
    @media (min-width: 990px) {
      display: block;
    }
  }
`

export const NavBarLink = styled(Link)`
  display: flex;
  gap: 8px;
  width: auto;
  text-transform: none;
  place-items: center;

  @media (max-width: 990px) {
    padding: 8px;
    min-width: fit-content;
  }

  .nav-text {
    font-size: 14px;
    display: none;
    @media (min-width: 990px) {
      display: block;
    }
  }
`

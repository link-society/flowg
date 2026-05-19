import { Button, IconButton, Link, styled } from '@mui/material'

export const NavToolbar = styled('div')`
  background-color: ${({ theme }) => theme.tokens.colors.black};
  display: flex;
  align-items: center;
  padding: 0 ${({ theme }) => theme.spacing(2)};
  min-height: 64px;
  width: 100%;
`

export const NavBarLeftSection = styled('section')`
  height: 100%;
  display: flex;
  flex-direction: row;
  align-items: stretch;
  gap: ${({ theme }) => theme.spacing(1.5)};
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
  gap: ${({ theme }) => theme.spacing(1.5)};
`

export const NavBarDesktopMenus = styled('div')`
  display: contents;
`

export const NavBarDrawerClose = styled(IconButton)`
  align-self: flex-end;
  margin: ${({ theme }) => theme.spacing(1)};
`

export const NavBarDrawerList = styled('div')`
  width: 240px;
`

export const NavBarButtonLogo = styled(Button)`
  display: flex;
  gap: ${({ theme }) => theme.spacing(1)};
  width: auto;
  text-transform: none;

  ${({ theme }) => theme.breakpoints.down('md')} {
    padding: ${({ theme }) => theme.spacing(1)};
    min-width: fit-content;
  }

  .nav-text {
    display: none;
    ${({ theme }) => theme.breakpoints.up('md')} {
      display: block;
    }
  }
`

export const NavBarButton = styled(Button)`
  display: flex;
  gap: ${({ theme }) => theme.spacing(1)};
  width: auto;
  text-transform: none;

  ${({ theme }) => theme.breakpoints.down('md')} {
    padding: ${({ theme }) => theme.spacing(1)};
    min-width: fit-content;
  }

  .nav-text {
    display: none;
    font-size: 14px;
    ${({ theme }) => theme.breakpoints.up('md')} {
      display: block;
    }
  }
`

export const NavBarLink = styled(Link)`
  display: flex;
  gap: ${({ theme }) => theme.spacing(1)};
  width: auto;
  text-transform: none;
  place-items: center;

  ${({ theme }) => theme.breakpoints.down('md')} {
    padding: ${({ theme }) => theme.spacing(1)};
    min-width: fit-content;
  }

  .nav-text {
    font-size: 14px;
    display: none;
    ${({ theme }) => theme.breakpoints.up('md')} {
      display: block;
    }
  }
`

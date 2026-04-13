import { useState } from 'react'
import { useNavigate } from 'react-router'

import Divider from '@mui/material/Divider'
import IconButton from '@mui/material/IconButton'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'
import useMediaQuery from '@mui/material/useMediaQuery'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import DashboardIcon from '@mui/icons-material/Dashboard'
import LogoutIcon from '@mui/icons-material/Logout'

import { useProfile } from '@/lib/hooks/profile'

import { NavMenuProfileButton } from './styles'

const NavMenuProfile = () => {
  const { user, permissions } = useProfile()
  const isMobile = useMediaQuery('(max-width: 990px)')
  const navigate = useNavigate()

  const [anchor, setAnchor] = useState<null | HTMLElement>(null)
  const handleOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchor(event.currentTarget)
  }
  const handleClose = () => {
    setAnchor(null)
  }
  const handleNavigate = (path: string) => {
    handleClose()
    navigate(path)
  }

  return (
    <>
      {isMobile ? (
        <IconButton
          id="menu:navbar.profile"
          color="inherit"
          onClick={handleOpen}
        >
          <AccountCircleIcon fontSize="small" />
        </IconButton>
      ) : (
        <NavMenuProfileButton
          id="menu:navbar.profile"
          color="inherit"
          onClick={handleOpen}
          startIcon={<AccountCircleIcon />}
          endIcon={<ArrowDropDownIcon />}
        >
          {user.name}
        </NavMenuProfileButton>
      )}

      <Menu
        anchorEl={anchor}
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
        keepMounted
        open={Boolean(anchor)}
        onClose={handleClose}
      >
        <MenuItem
          variant="navLink"
          id="link:navbar.profile.account"
          onClick={() => handleNavigate('/web/account')}
        >
          <AccountCircleIcon fontSize="small" />
          Account
        </MenuItem>

        {permissions.can_view_acls && (
          <MenuItem
            variant="navLink"
            id="link:navbar.profile.admin"
            onClick={() => handleNavigate('/web/admin')}
          >
            <DashboardIcon fontSize="small" />
            Admin
          </MenuItem>
        )}

        <Divider />

        <MenuItem
          variant="navLink"
          id="link:navbar.profile.logout"
          onClick={() => handleNavigate('/web/logout')}
        >
          <LogoutIcon fontSize="small" />
          Logout
        </MenuItem>
      </Menu>
    </>
  )
}

export default NavMenuProfile

import { useState } from 'react'

import Button from '@mui/material/Button'
import Divider from '@mui/material/Divider'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import DashboardIcon from '@mui/icons-material/Dashboard'
import LogoutIcon from '@mui/icons-material/Logout'

import { useProfile } from '@/lib/hooks/profile'

const NavMenuProfile = () => {
  const { user, permissions } = useProfile()

  const [anchor, setAnchor] = useState<null | HTMLElement>(null)
  const handleOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchor(event.currentTarget)
  }
  const handleClose = () => {
    setAnchor(null)
  }

  return (
    <>
      <Button
        id="menu:navbar.profile"
        color="inherit"
        onClick={handleOpen}
        startIcon={<AccountCircleIcon />}
        endIcon={<ArrowDropDownIcon />}
        sx={{ textTransform: 'none' }}
      >
        {user.name}
      </Button>

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
        <MenuItem onClick={handleClose}>
          <Button
            id="link:navbar.profile.account"
            href="/web/account"
            color="secondary"
            startIcon={<AccountCircleIcon />}
          >
            Account
          </Button>
        </MenuItem>

        {permissions.can_view_acls && (
          <MenuItem onClick={handleClose}>
            <Button
              id="link:navbar.profile.admin"
              href="/web/admin"
              color="secondary"
              startIcon={<DashboardIcon />}
            >
              Admin
            </Button>
          </MenuItem>
        )}

        <Divider />

        <MenuItem onClick={handleClose}>
          <Button
            id="link:navbar.profile.logout"
            href="/web/logout"
            color="secondary"
            startIcon={<LogoutIcon />}
          >
            Logout
          </Button>
        </MenuItem>
      </Menu>
    </>
  )
}

export default NavMenuProfile

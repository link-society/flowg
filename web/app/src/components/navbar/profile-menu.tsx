import { useState } from 'react'

import AccountCircleIcon from '@mui/icons-material/AccountCircle'
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import DashboardIcon from '@mui/icons-material/Dashboard'
import LogoutIcon from '@mui/icons-material/Logout'

import Divider from '@mui/material/Divider'
import Button from '@mui/material/Button'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'

import { useProfile } from '@/lib/context/profile'

export const ProfileMenu = () => {
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

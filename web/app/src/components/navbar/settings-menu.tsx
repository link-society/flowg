import { useState } from 'react'

import AccountTreeIcon from '@mui/icons-material/AccountTree'
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import NotificationsActiveIcon from '@mui/icons-material/NotificationsActive'
import SdStorageIcon from '@mui/icons-material/SdStorage'
import SettingsIcon from '@mui/icons-material/Settings'

import Button from '@mui/material/Button'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'

import { useProfile } from '@/lib/context/profile'

export const SettingsMenu = () => {
  const { permissions } = useProfile()

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
        id="menu:navbar.settings"
        color="inherit"
        onClick={handleOpen}
        startIcon={<SettingsIcon />}
        endIcon={<ArrowDropDownIcon />}
        sx={{ textTransform: 'none' }}
      >
        Settings
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
        {permissions.can_view_transformers && (
          <MenuItem onClick={handleClose}>
            <Button
              id="link:navbar.settings.transformers"
              href="/web/transformers"
              color="secondary"
              startIcon={<FilterAltIcon />}
            >
              Transformers
            </Button>
          </MenuItem>
        )}

        {permissions.can_view_alerts && (
          <MenuItem onClick={handleClose}>
            <Button
              id="link:navbar.settings.alerts"
              href="/web/alerts"
              color="secondary"
              startIcon={<NotificationsActiveIcon />}
            >
              Alerts
            </Button>
          </MenuItem>
        )}

        {permissions.can_view_streams && (
          <MenuItem onClick={handleClose}>
            <Button
              id="link:navbar.settings.storage"
              href="/web/storage"
              color="secondary"
              startIcon={<SdStorageIcon />}
            >
              Storage
            </Button>
          </MenuItem>
        )}

        {permissions.can_view_pipelines && (
          <MenuItem onClick={handleClose}>
            <Button
              id="link:navbar.settings.pipelines"
              href="/web/pipelines"
              color="secondary"
              startIcon={<AccountTreeIcon />}
            >
              Pipelines
            </Button>
          </MenuItem>
        )}

      </Menu>
    </>
  )
}

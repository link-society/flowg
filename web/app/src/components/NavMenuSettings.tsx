import { useState } from 'react'

import Button from '@mui/material/Button'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'

import AccountTreeIcon from '@mui/icons-material/AccountTree'
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'
import SdStorageIcon from '@mui/icons-material/SdStorage'
import SettingsIcon from '@mui/icons-material/Settings'

import { useProfile } from '@/lib/hooks/profile'

const NavMenuSettings = () => {
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

        {permissions.can_view_forwarders && (
          <MenuItem onClick={handleClose}>
            <Button
              id="link:navbar.settings.forwarders"
              href="/web/forwarders"
              color="secondary"
              startIcon={<ForwardToInboxIcon />}
            >
              Forwarders
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

        {permissions.can_read_system_configuration && (
          <MenuItem onClick={handleClose}>
            <Button
              id="link:navbar.settings.configuration"
              href="/web/system-configuration"
              color="secondary"
              startIcon={<SettingsIcon />}
            >
              System configuration
            </Button>
          </MenuItem>
        )}
      </Menu>
    </>
  )
}

export default NavMenuSettings

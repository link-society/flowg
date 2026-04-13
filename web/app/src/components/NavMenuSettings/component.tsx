import { useState } from 'react'
import { useNavigate } from 'react-router'

import IconButton from '@mui/material/IconButton'
import Menu from '@mui/material/Menu'
import MenuItem from '@mui/material/MenuItem'
import useMediaQuery from '@mui/material/useMediaQuery'

import AccountTreeIcon from '@mui/icons-material/AccountTree'
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown'
import FilterAltIcon from '@mui/icons-material/FilterAlt'
import ForwardToInboxIcon from '@mui/icons-material/ForwardToInbox'
import SdStorageIcon from '@mui/icons-material/SdStorage'
import SettingsIcon from '@mui/icons-material/Settings'

import { useProfile } from '@/lib/hooks/profile'

import { NavMenuSettingsButton } from './styles'

const NavMenuSettings = () => {
  const { permissions } = useProfile()
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
          id="menu:navbar.settings"
          color="inherit"
          onClick={handleOpen}
        >
          <SettingsIcon fontSize="small" />
        </IconButton>
      ) : (
        <NavMenuSettingsButton
          id="menu:navbar.settings"
          color="inherit"
          onClick={handleOpen}
          startIcon={<SettingsIcon />}
          endIcon={<ArrowDropDownIcon />}
        >
          Settings
        </NavMenuSettingsButton>
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
        {permissions.can_view_transformers && (
          <MenuItem
            variant="navLink"
            id="link:navbar.settings.transformers"
            onClick={() => handleNavigate('/web/transformers')}
          >
            <FilterAltIcon fontSize="small" />
            Transformers
          </MenuItem>
        )}

        {permissions.can_view_forwarders && (
          <MenuItem
            variant="navLink"
            id="link:navbar.settings.forwarders"
            onClick={() => handleNavigate('/web/forwarders')}
          >
            <ForwardToInboxIcon fontSize="small" />
            Forwarders
          </MenuItem>
        )}

        {permissions.can_view_streams && (
          <MenuItem
            variant="navLink"
            id="link:navbar.settings.storage"
            onClick={() => handleNavigate('/web/storage')}
          >
            <SdStorageIcon fontSize="small" />
            Storage
          </MenuItem>
        )}

        {permissions.can_view_pipelines && (
          <MenuItem
            variant="navLink"
            id="link:navbar.settings.pipelines"
            onClick={() => handleNavigate('/web/pipelines')}
          >
            <AccountTreeIcon fontSize="small" />
            Pipelines
          </MenuItem>
        )}

        {permissions.can_read_system_configuration && (
          <MenuItem
            variant="navLink"
            id="link:navbar.settings.configuration"
            onClick={() => handleNavigate('/web/system-configuration')}
          >
            <SettingsIcon fontSize="small" />
            System configuration
          </MenuItem>
        )}
      </Menu>
    </>
  )
}

export default NavMenuSettings

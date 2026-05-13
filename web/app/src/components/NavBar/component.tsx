import { useNavigate } from 'react-router'

import AppBar from '@mui/material/AppBar'
import Toolbar from '@mui/material/Toolbar'
import Typography from '@mui/material/Typography'

import ApiIcon from '@mui/icons-material/Api'
import GitHubIcon from '@mui/icons-material/GitHub'
import StorageIcon from '@mui/icons-material/Storage'
import UploadFileIcon from '@mui/icons-material/UploadFile'

import { useProfile } from '@/lib/hooks/profile'

import NavMenuProfile from '@/components/NavMenuProfile/component'
import NavMenuSettings from '@/components/NavMenuSettings/component'

import { buildUrl } from '@/router'

import {
  NavBarButton,
  NavBarButtonLogo,
  NavBarLeftSection,
  NavBarLink,
  NavBarRightSection,
} from './styles'

const NavBar = () => {
  const { permissions } = useProfile()
  const navigate = useNavigate()

  const handleNavigate = (path: string) => () => {
    navigate(path)
  }

  return (
    <AppBar position="static">
      <Toolbar>
        <NavBarLeftSection>
          <NavBarButtonLogo
            onClick={handleNavigate(buildUrl('/'))}
            color="inherit"
          >
            <img src="./assets/logo.png" alt="Logo FlowG" />

            <Typography variant="titleSm" className="nav-text">
              FlowG
            </Typography>
          </NavBarButtonLogo>

          <NavBarLink
            href="https://github.com/link-society/flowg"
            color="inherit"
            target="_blank"
            underline="none"
          >
            <GitHubIcon fontSize="small" />
            <Typography variant="titleSm" className="nav-text">
              Github
            </Typography>
          </NavBarLink>

          <NavBarLink
            href="/api/docs"
            target="_blank"
            color="inherit"
            underline="none"
          >
            <ApiIcon fontSize="small" />
            <Typography variant="titleSm" className="nav-text">
              API Docs
            </Typography>
          </NavBarLink>
        </NavBarLeftSection>

        <NavBarRightSection>
          <NavMenuProfile />
          <NavMenuSettings />

          {permissions.can_view_streams && (
            <NavBarButton
              id="link:navbar.streams"
              onClick={handleNavigate(buildUrl('/streams'))}
              color="inherit"
            >
              <StorageIcon fontSize="small" />

              <Typography variant="titleSm" className="nav-text">
                Streams
              </Typography>
            </NavBarButton>
          )}

          {permissions.can_send_logs && (
            <NavBarButton
              id="link:navbar.upload"
              onClick={handleNavigate(buildUrl('/upload'))}
              color="inherit"
            >
              <UploadFileIcon fontSize="small" />

              <Typography variant="titleSm" className="nav-text">
                Upload
              </Typography>
            </NavBarButton>
          )}
        </NavBarRightSection>
      </Toolbar>
    </AppBar>
  )
}

export default NavBar

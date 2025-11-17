import AppBar from '@mui/material/AppBar'
import Button from '@mui/material/Button'
import Toolbar from '@mui/material/Toolbar'
import Typography from '@mui/material/Typography'

import ApiIcon from '@mui/icons-material/Api'
import GitHubIcon from '@mui/icons-material/GitHub'
import StorageIcon from '@mui/icons-material/Storage'
import UploadFileIcon from '@mui/icons-material/UploadFile'

import { useProfile } from '@/lib/hooks/profile'

import NavMenuProfile from '@/components/NavMenuProfile'
import NavMenuSettings from '@/components/NavMenuSettings'

const NavBar = () => {
  const { permissions } = useProfile()

  return (
    <AppBar position="static">
      <Toolbar>
        <section className="h-full flex flex-row items-stretch gap-3 grow">
          <Button
            href="/web/"
            color="inherit"
            startIcon={
              <img src="/web/assets/logo.png" alt="Logo" className="h-8" />
            }
            sx={{ textTransform: 'none' }}
          >
            <Typography variant="h6">FlowG</Typography>
          </Button>
          <Button
            href="https://github.com/link-society/flowg"
            target="_blank"
            color="inherit"
            startIcon={<GitHubIcon />}
            sx={{ textTransform: 'none' }}
          >
            Github
          </Button>
          <Button
            href="/api/docs"
            target="_blank"
            color="inherit"
            startIcon={<ApiIcon />}
            sx={{ textTransform: 'none' }}
          >
            API Docs
          </Button>
        </section>

        <section className="h-full flex flex-row-reverse items-stretch gap-3">
          <NavMenuProfile />
          <NavMenuSettings />
          {permissions.can_view_streams && (
            <Button
              id="link:navbar.streams"
              href="/web/streams"
              color="inherit"
              startIcon={<StorageIcon />}
              sx={{ textTransform: 'none' }}
            >
              Streams
            </Button>
          )}
          {permissions.can_send_logs && (
            <Button
              id="link:navbar.upload"
              href="/web/upload"
              color="inherit"
              startIcon={<UploadFileIcon />}
              sx={{ textTransform: 'none' }}
            >
              Upload
            </Button>
          )}
        </section>
      </Toolbar>
    </AppBar>
  )
}

export default NavBar

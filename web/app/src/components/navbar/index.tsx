import AppBar from '@mui/material/AppBar'
import Button from '@mui/material/Button'
import Toolbar from '@mui/material/Toolbar'
import Typography from '@mui/material/Typography'

import ApiIcon from '@mui/icons-material/Api'
import GitHubIcon from '@mui/icons-material/GitHub'
import StorageIcon from '@mui/icons-material/Storage'

import { useProfile } from '@/lib/context/profile'

import { ProfileMenu } from './profile-menu'
import { SettingsMenu } from './settings-menu'

export const NavBar = () => {
  const { permissions } = useProfile()

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" sx={{ mr: 2 }}>
          FlowG
        </Typography>
        <section className="h-full flex flex-row items-stretch gap-3 grow">
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
          <ProfileMenu />
          <SettingsMenu />
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
        </section>
      </Toolbar>
    </AppBar>
  )
}

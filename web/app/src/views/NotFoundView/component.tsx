import { Typography } from '@mui/material'

import { useNavigate } from 'react-router'

import Button from '@mui/material/Button'

import SearchOffIcon from '@mui/icons-material/SearchOff'

import { buildUrl } from '@/router'

import { NotFoundViewContainer } from './styles'

const NotFoundView = () => {
  const navigate = useNavigate()

  return (
    <NotFoundViewContainer>
      <SearchOffIcon sx={{ fontSize: '5rem', opacity: 0.3 }} />

      <Typography variant="titleLg" component="h1" sx={{ fontWeight: 700 }}>
        404
      </Typography>

      <Typography variant="titleMd" component="h2">
        Page not found
      </Typography>

      <Typography
        variant="text"
        sx={{
          color: (theme) => `rgba(0, 0, 0, ${theme.tokens.opacity.overlay})`,
        }}
      >
        The page you are looking for does not exist or has been moved.
      </Typography>

      <Button
        variant="contained"
        color="secondary"
        onClick={() => navigate(buildUrl('/'))}
      >
        Back to home
      </Button>
    </NotFoundViewContainer>
  )
}

export default NotFoundView

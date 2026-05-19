import { useNavigate } from 'react-router'

import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'

import { buildUrl } from '@/router'

import {
  NotFoundHint,
  NotFoundIcon,
  NotFoundTitle,
  NotFoundViewContainer,
} from './styles'

const NotFoundView = () => {
  const navigate = useNavigate()

  return (
    <NotFoundViewContainer>
      <NotFoundIcon />

      <NotFoundTitle variant="titleLg" component="h1">
        404
      </NotFoundTitle>

      <Typography variant="titleMd" component="h2">
        Page not found
      </Typography>

      <NotFoundHint variant="text">
        The page you are looking for does not exist or has been moved.
      </NotFoundHint>

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

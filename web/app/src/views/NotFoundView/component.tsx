import { useTranslation } from 'react-i18next'
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
  const { t } = useTranslation()
  const navigate = useNavigate()

  return (
    <NotFoundViewContainer>
      <NotFoundIcon />

      <NotFoundTitle variant="titleLg">404</NotFoundTitle>

      <Typography variant="titleMd" component="h2">
        {t('pages.notFound.title')}
      </Typography>

      <NotFoundHint variant="text">
        {t('pages.notFound.description')}
      </NotFoundHint>

      <Button
        variant="contained"
        color="secondary"
        onClick={() => navigate(buildUrl('/'))}
      >
        {t('pages.notFound.backHome')}
      </Button>
    </NotFoundViewContainer>
  )
}

export default NotFoundView

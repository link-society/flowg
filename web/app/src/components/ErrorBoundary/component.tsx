import { useTranslation } from 'react-i18next'
import { useRouteError } from 'react-router'

import TextField from '@mui/material/TextField'
import Typography from '@mui/material/Typography'

import SentimentVeryDissatisfiedIcon from '@mui/icons-material/SentimentVeryDissatisfied'

import * as errors from '@/lib/api/errors'

import { CodeBlock, ErrorHeading, ErrorHeadingLabel, ErrorRoot } from './styles'

const Heading = ({ title }: { title: string }) => (
  <ErrorHeading variant="titleLg">
    <SentimentVeryDissatisfiedIcon fontSize="large" />
    <ErrorHeadingLabel variant="titleLg">{title}</ErrorHeadingLabel>
  </ErrorHeading>
)

const ErrorBoundary = () => {
  const { t } = useTranslation()
  const error = useRouteError()

  if (error instanceof errors.InvalidResponseError) {
    return (
      <ErrorRoot>
        <Heading title={error.message} />

        <TextField
          label={t('components.errorBoundary.statusCodeLabel')}
          value={error.statusCode}
          variant="outlined"
          disabled
          error={!error.ok}
        />

        <CodeBlock component="pre">
          <code>{error.body}</code>
        </CodeBlock>
      </ErrorRoot>
    )
  } else if (error instanceof Error) {
    return (
      <ErrorRoot>
        <Heading title={error.message} />

        {error.stack ? (
          <CodeBlock component="pre">
            <code>{error.stack}</code>
          </CodeBlock>
        ) : (
          <Typography variant="text">
            {t('components.errorBoundary.noStacktrace')}
          </Typography>
        )}
      </ErrorRoot>
    )
  }

  console.error(error)

  return (
    <ErrorRoot>
      <Heading title={t('components.errorBoundary.unknownError')} />
      <Typography variant="text">
        {t('components.errorBoundary.unknownErrorDetail')}
      </Typography>
    </ErrorRoot>
  )
}

export default ErrorBoundary

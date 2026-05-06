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
  const error = useRouteError()

  if (error instanceof errors.InvalidResponseError) {
    return (
      <ErrorRoot>
        <Heading title={error.message} />

        <TextField
          label="Status Code"
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
          <Typography variant="text">No stacktrace available</Typography>
        )}
      </ErrorRoot>
    )
  }

  console.error(error)

  return (
    <ErrorRoot>
      <Heading title="An unknown exception has been thrown." />
      <Typography variant="text">
        The exception has been logged to the console.
      </Typography>
    </ErrorRoot>
  )
}

export default ErrorBoundary

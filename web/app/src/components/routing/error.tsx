import { useRouteError } from 'react-router'

import TextField from '@mui/material/TextField'

import SentimentVeryDissatisfiedIcon from '@mui/icons-material/SentimentVeryDissatisfied'

import * as errors from '@/lib/api/errors'

const Heading = ({ title }: { title: string }) => (
  <h1 className="text-2xl text-red-600 flex flex-row items-center">
    <SentimentVeryDissatisfiedIcon fontSize="large" />

    <span className="ml-2">{title}</span>
  </h1>
)

const CodeBlock = ({ content }: { content: string }) => (
  <pre className="p-2 bg-black text-gray-400 shadow">
    <code>{content}</code>
  </pre>
)

export const ErrorBoundary = () => {
  const error = useRouteError()

  if (error instanceof errors.InvalidResponseError) {
    return (
      <div className="p-3 flex flex-col gap-3">
        <Heading title={error.message} />

        <TextField
          label="Status Code"
          value={error.statusCode}
          variant="outlined"
          disabled
          error={!error.ok}
        />

        <CodeBlock content={error.body} />
      </div>
    )
  } else if (error instanceof Error) {
    return (
      <div className="p-3 flex flex-col gap-3">
        <Heading title={error.message} />

        {error.stack ? (
          <CodeBlock content={error.stack} />
        ) : (
          <p>No stacktrace available</p>
        )}
      </div>
    )
  }

  return (
    <div className="p-3 flex flex-col gap-3">
      <Heading title="An unknown exception has been thrown." />

      <p>The exception has been logged to the console.</p>
    </div>
  )
}

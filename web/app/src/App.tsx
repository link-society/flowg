import { useMemo } from 'react'
import { RouterProvider } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

import router from '@/router'

export default function() {
  const queryClient = useMemo(
    () => new QueryClient(),
    [],
  )

  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  )
}

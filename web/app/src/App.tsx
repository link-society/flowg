import { ThemeRegistry } from '@/theme'

import { Suspense } from 'react'
import { RouterProvider } from 'react-router'

import LinearProgress from '@mui/material/LinearProgress'

import router from '@/router'

const App = () => (
  <ThemeRegistry>
    <Suspense fallback={<LinearProgress />}>
      <RouterProvider router={router} />
    </Suspense>
  </ThemeRegistry>
)

export default App

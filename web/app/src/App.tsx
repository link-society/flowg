import { ThemeRegistry } from '@/theme'

import { RouterProvider } from 'react-router'

import router from '@/router'

const App = () => (
  <ThemeRegistry>
    <RouterProvider router={router} />
  </ThemeRegistry>
)

export default App

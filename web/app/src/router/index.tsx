import { createBrowserRouter } from 'react-router-dom'

import BaseLayout from '@/layouts/base'

export default createBrowserRouter([
  {
    path: '/web/',
    element: <BaseLayout />,
    children: [
      {
        path: 'login',
        element: <div>Login</div>
      },
      {
        path: 'logout',
        element: <div>Logout</div>
      },
      {
        path: '',
        element: <div>Home</div>
      }
    ]
  }
])

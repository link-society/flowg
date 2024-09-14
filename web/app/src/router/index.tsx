import { createBrowserRouter } from 'react-router-dom'

import BaseLayout from '@/layouts/base'
import AppLayout, { loader as AppLoader } from '@/layouts/app'

import LoginView from '@/views/onboarding/login'
import LogoutView, { loader as LogoutLoader } from '@/views/onboarding/logout'

export default createBrowserRouter([
  {
    path: '/web/',
    element: <BaseLayout />,
    children: [
      {
        path: 'login',
        element: <LoginView />,
      },
      {
        path: 'logout',
        element: <LogoutView />,
        loader: LogoutLoader,
      },
      {
        path: '',
        element: <AppLayout />,
        loader: AppLoader,
        children: [
          {
            path: '',
            element: <div>Home</div>
          },
        ],
      },
    ],
  },
])

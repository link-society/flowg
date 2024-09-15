import { createBrowserRouter } from 'react-router-dom'

export default createBrowserRouter([
  {
    path: '/web/',
    lazy: async () => {
      const { BaseLayout: Component } = await import('@/layouts/base')
      return { Component }
    },
    children: [
      {
        path: 'login',
        lazy: async () => {
          const { LoginView: Component } = await import('@/views/onboarding/login')
          return { Component }
        },
      },
      {
        path: 'logout',
        lazy: async () => {
          const { LogoutView: Component, loader } = await import('@/views/onboarding/logout')
          return { Component, loader }
        },
      },
      {
        path: '',
        lazy: async () => {
          const { AppLayout: Component, loader } = await import('@/layouts/app')
          return { Component, loader }
        },
        children: [
          {
            path: '',
            lazy: async () => {
              const { HomeView: Component } = await import('@/views/app/home')
              return { Component }
            },
          },
        ],
      },
    ],
  },
])

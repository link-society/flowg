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
            path: 'account',
            lazy: async () => {
              const { AccountView: Component, loader } = await import('@/views/app/account')
              return { Component, loader }
            },
          },
          {
            path: 'admin',
            lazy: async () => {
              const { AdminView: Component, loader } = await import('@/views/app/admin')
              return { Component, loader }
            },
          },
          {
            path: 'transformers',
            lazy: async () => {
              const { TransformerView: Component } = await import('@/views/app/transformers/section')
              const { loader } = await import('@/views/app/transformers/loader')
              return { Component, loader }
            },
          },
          {
            path: 'transformers/:transformer',
            lazy: async () => {
              const { TransformerView: Component } = await import('@/views/app/transformers/item')
              const { loader } = await import('@/views/app/transformers/loader')
              return { Component, loader }
            },
          },
          {
            path: 'pipelines',
            lazy: async () => {
              const { PipelineView: Component } = await import('@/views/app/pipelines/section')
              const { loader } = await import('@/views/app/pipelines/loader')
              return { Component, loader }
            },
          },
          {
            path: 'pipelines/:pipeline',
            lazy: async () => {
              const { PipelineView: Component } = await import('@/views/app/pipelines/item')
              const { loader } = await import('@/views/app/pipelines/loader')
              return { Component, loader }
            },
          },
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

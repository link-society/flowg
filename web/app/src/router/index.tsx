import { createBrowserRouter } from 'react-router'

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
              const { AccountView: Component } = await import('@/views/app/account')
              const { loader } = await import('@/views/app/account/loader')
              return { Component, loader }
            },
          },
          {
            path: 'admin',
            lazy: async () => {
              const { AdminView: Component } = await import('@/views/app/admin')
              const { loader } = await import('@/views/app/admin/loader')
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
            path: 'storage',
            lazy: async () => {
              const { StreamView: Component } = await import('@/views/app/storage/section')
              const { loader } = await import('@/views/app/storage/loader')
              return { Component, loader }
            },
          },
          {
            path: 'storage/:stream',
            lazy: async () => {
              const { StreamView: Component } = await import('@/views/app/storage/item')
              const { loader } = await import('@/views/app/storage/loader')
              return { Component, loader }
            },
          },
          {
            path: 'alerts',
            lazy: async () => {
              const { AlertView: Component } = await import('@/views/app/alerts/section')
              const { loader } = await import('@/views/app/alerts/loader')
              return { Component, loader }
            },
          },
          {
            path: 'alerts/:alert',
            lazy: async () => {
              const { AlertView: Component } = await import('@/views/app/alerts/item')
              const { loader } = await import('@/views/app/alerts/loader')
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
            path: 'streams',
            lazy: async () => {
              const { StreamView: Component } = await import('@/views/app/streams/section')
              const { loader } = await import('@/views/app/streams/loader')
              return { Component, loader }
            },
          },
          {
            path: 'streams/:stream',
            lazy: async () => {
              const { StreamView: Component } = await import('@/views/app/streams/item')
              const { loader } = await import('@/views/app/streams/loader')
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

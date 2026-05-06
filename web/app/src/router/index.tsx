import { createBrowserRouter } from 'react-router'

import LinearProgress from '@mui/material/LinearProgress'

import ErrorBoundary from '@/components/ErrorBoundary/component'

export default createBrowserRouter([
  {
    path: '/web/',
    lazy: async () => {
      const { default: Component } =
        await import('@/layouts/BaseLayout/component')
      return {
        Component,
        HydrateFallback: () => <LinearProgress />,
        ErrorBoundary: () => <ErrorBoundary />,
      }
    },
    children: [
      {
        path: 'login',
        lazy: async () => {
          const { default: Component } =
            await import('@/views/LoginView/component')
          return { Component }
        },
      },
      {
        path: 'logout',
        lazy: async () => {
          const { default: Component, loader } =
            await import('@/views/LogoutView/component')
          return { Component, loader }
        },
      },
      {
        path: '',
        lazy: async () => {
          const { default: Component, loader } =
            await import('@/layouts/AppLayout/component')
          return { Component, loader }
        },
        children: [
          {
            path: 'account',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/AccountView/component')
              return { Component, loader }
            },
          },
          {
            path: 'admin',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/AdminView/component')
              return { Component, loader }
            },
          },
          {
            path: 'transformers',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/TransformerSectionView/component')
              return { Component, loader }
            },
          },
          {
            path: 'transformers/:transformer',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/TransformerDetailView/component')
              return { Component, loader }
            },
          },
          {
            path: 'storage',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/StorageSectionView/component')
              return { Component, loader }
            },
          },
          {
            path: 'storage/:stream',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/StorageDetailView/component')
              return { Component, loader }
            },
          },
          {
            path: 'forwarders',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/ForwarderSectionView/component')
              return { Component, loader }
            },
          },
          {
            path: 'forwarders/:forwarder',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/ForwarderDetailView/component')
              return { Component, loader }
            },
          },
          {
            path: 'pipelines',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/PipelineSectionView/component')
              return { Component, loader }
            },
          },
          {
            path: 'pipelines/:pipeline',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/PipelineDetailView/component')
              return { Component, loader }
            },
          },
          {
            path: 'system-configuration',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/SystemConfiguration/component')
              return { Component, loader }
            },
          },
          {
            path: 'streams',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/StreamSectionView/component')
              return { Component, loader }
            },
          },
          {
            path: 'streams/:stream',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/StreamDetailView/component')
              return { Component, loader }
            },
          },
          {
            path: 'upload',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/UploadView/component')
              return { Component, loader }
            },
          },
          {
            path: '',
            lazy: async () => {
              const { default: Component } =
                await import('@/views/HomeView/component')
              return { Component }
            },
          },
        ],
      },
    ],
  },
])

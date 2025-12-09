import { createBrowserRouter } from 'react-router'

import LinearProgress from '@mui/material/LinearProgress'

import ErrorBoundary from '@/components/ErrorBoundary'

export default createBrowserRouter([
  {
    path: '/web/',
    lazy: async () => {
      const { default: Component } = await import('@/layouts/base')
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
          const { default: Component } = await import('@/views/LoginView')
          return { Component }
        },
      },
      {
        path: 'logout',
        lazy: async () => {
          const { default: Component, loader } =
            await import('@/views/LogoutView')
          return { Component, loader }
        },
      },
      {
        path: '',
        lazy: async () => {
          const { default: Component, loader } = await import('@/layouts/app')
          return { Component, loader }
        },
        children: [
          {
            path: 'account',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/AccountView')
              return { Component, loader }
            },
          },
          {
            path: 'admin',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/AdminView')
              return { Component, loader }
            },
          },
          {
            path: 'transformers',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/TransformerSectionView')
              return { Component, loader }
            },
          },
          {
            path: 'transformers/:transformer',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/TransformerDetailView')
              return { Component, loader }
            },
          },
          {
            path: 'storage',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/StorageSectionView')
              return { Component, loader }
            },
          },
          {
            path: 'storage/:stream',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/StorageDetailView')
              return { Component, loader }
            },
          },
          {
            path: 'forwarders',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/ForwarderSectionView')
              return { Component, loader }
            },
          },
          {
            path: 'forwarders/:forwarder',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/ForwarderDetailView')
              return { Component, loader }
            },
          },
          {
            path: 'pipelines',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/PipelineSectionView')
              return { Component, loader }
            },
          },
          {
            path: 'pipelines/:pipeline',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/PipelineDetailView')
              return { Component, loader }
            },
          },
          {
            path: 'system-configuration',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/SystemConfiguration.tsx')
              return { Component, loader }
            },
          },
          {
            path: 'streams',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/StreamSectionView')
              return { Component, loader }
            },
          },
          {
            path: 'streams/:stream',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/StreamDetailView')
              return { Component, loader }
            },
          },
          {
            path: 'upload',
            lazy: async () => {
              const { default: Component, loader } =
                await import('@/views/UploadView')
              return { Component, loader }
            },
          },
          {
            path: '',
            lazy: async () => {
              const { default: Component } = await import('@/views/HomeView')
              return { Component }
            },
          },
        ],
      },
    ],
  },
])

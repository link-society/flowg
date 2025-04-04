import '@xyflow/react/dist/style.css'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-balham.css'
import 'ag-grid-community/styles/ag-theme-material.css'

import './styles/main.css'
import './styles/table.css'

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import {
  AllCommunityModule,
  ModuleRegistry,
  provideGlobalGridOptions,
} from 'ag-grid-community'

import App from '@/App'

ModuleRegistry.registerModules([AllCommunityModule])
provideGlobalGridOptions({ theme: 'legacy' })

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>
)

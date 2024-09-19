import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import './styles/main.css'
import '@xyflow/react/dist/style.css'
import 'ag-grid-community/styles/ag-grid.css'
import 'ag-grid-community/styles/ag-theme-material.css'
import './styles/table.css'

import App from '@/App'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>
)

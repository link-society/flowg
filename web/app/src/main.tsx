import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'

import * as M from '@materializecss/materialize'

import App from '@/App'

import './styles/main.css'
import '@materializecss/materialize/style'

document.documentElement.setAttribute('theme', 'light')

const rootElement = document.getElementById('root') as HTMLElement
const observer = new MutationObserver((mutations) => {
  for (const mutation of mutations) {
    if (mutation.type === 'childList') {
      for (const node of mutation.addedNodes) {
        if (node instanceof HTMLElement) {
          M.AutoInit(node, {
            Dropdown: {
              constrainWidth: false,
            }
          })
        }
      }
    }
  }
})
observer.observe(document.body, { childList: true, subtree: true })

const root = createRoot(rootElement)
root.render(
  <StrictMode>
    <App />
  </StrictMode>
)

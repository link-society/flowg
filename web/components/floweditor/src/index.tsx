import ReactDOM from 'react-dom/client'
import { ReactFlowProvider } from '@xyflow/react'

import FlowEditor from './flow/FlowEditor'
import { DnDProvider } from './dnd/context'


class FlowEditorElement extends HTMLElement {
  private root: ReactDOM.Root

  constructor() {
    super()

    this.root = ReactDOM.createRoot(this)
  }

  connectedCallback() {
    this.render()
  }

  static get observedAttributes() {
    return ['flow']
  }

  attributeChangedCallback(name: string, oldValue: string, newValue: string) {
    if (name === 'flow' && oldValue !== newValue) {
      this.render()
    }
  }

  private render() {
    const flow = this.getAttribute('flow') ?? ''

    const handleChange = (value: string) => {
      this.setAttribute('flow', value)
    }

    this.root.render(
      <ReactFlowProvider>
        <DnDProvider>
          <FlowEditor
            flow={flow}
            onFlowChange={handleChange}
          />
        </DnDProvider>
      </ReactFlowProvider>
    )
  }
}

customElements.define('flow-editor', FlowEditorElement)

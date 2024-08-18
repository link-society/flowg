import ReactDOM from 'react-dom/client'

import FlowEditor from './FlowEditor'


class FlowEditorElement extends HTMLElement {
  connectedCallback() {
    const root = ReactDOM.createRoot(this)
    root.render(<FlowEditor />)
  }
}

customElements.define('flow-editor', FlowEditorElement)

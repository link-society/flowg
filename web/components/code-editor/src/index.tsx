import ReactDOM from 'react-dom/client'

import CodeEditor from './CodeEditor'


class CodeEditorElement extends HTMLElement {
  private root: ReactDOM.Root

  constructor() {
    super()

    this.root = ReactDOM.createRoot(this)
  }

  connectedCallback() {
    this.render()
  }

  static get observedAttributes() {
    return ['code']
  }

  attributeChangedCallback(name: string, oldValue: string, newValue: string) {
    if (name === 'code' && oldValue !== newValue) {
      this.render()
    }
  }

  private render() {
    const code = this.getAttribute('code') ?? ''

    const handleChange = (value: string) => {
      this.setAttribute('code', value)
    }

    this.root.render(<CodeEditor code={code} onCodeChange={handleChange} />)
  }
}

customElements.define('code-editor', CodeEditorElement)

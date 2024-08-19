import ReactDOM from 'react-dom/client'

import FlowEditor from './FlowEditor'
import { AddNodeEvent } from './event'


class FlowEditorElement extends HTMLElement {
  private root: ReactDOM.Root
  private eventTarget: EventTarget

  constructor() {
    super()

    this.root = ReactDOM.createRoot(this)
    this.eventTarget = new EventTarget()
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
      <FlowEditor
        flow={flow}
        onFlowChange={handleChange}
        eventTarget={this.eventTarget}
      />
    )
  }

  addNode(type: string) {
    const event = new AddNodeEvent(type)
    this.eventTarget.dispatchEvent(event)
  }
}

customElements.define('flow-editor', FlowEditorElement)

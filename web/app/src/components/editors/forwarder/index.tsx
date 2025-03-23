import { ForwarderModel } from '@/lib/models/forwarder'

import { ForwarderConfigHttpEditor } from './http'

type ForwarderEditorProps = {
  forwarder: ForwarderModel
  onForwarderChange: (forwarder: ForwarderModel) => void
}

export const ForwarderEditor = ({ forwarder, onForwarderChange }: ForwarderEditorProps) => {
  switch (forwarder.config.type) {
    case 'http':
      return (
        <ForwarderConfigHttpEditor
          config={forwarder.config}
          onConfigChange={(config) => {
            onForwarderChange({
              ...forwarder,
              config,
            })
          }}
        />
      )

    default:
      throw new Error(`Unknown forwarder type: ${forwarder.config.type}`)
  }
}

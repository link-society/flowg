import ForwarderModel from '@/lib/models/ForwarderModel'

export type LoaderData = {
  forwarders: string[]
  currentForwarder: {
    name: string
    forwarder: ForwarderModel
  }
}

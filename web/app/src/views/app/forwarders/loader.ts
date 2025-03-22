import { LoaderFunction } from 'react-router'

import { loginRequired } from '@/lib/decorators/loaders'
import * as configApi from '@/lib/api/operations/config'
import { ForwarderModel } from '@/lib/models'

export type LoaderData = {
  forwarders: string[]
  currentForwarder?: {
    name: string
    forwarder: ForwarderModel
  }
}

export const loader: LoaderFunction = loginRequired(
  async ({ params }): Promise<LoaderData> => {
    const forwarders = await configApi.listForwarders()

    if (params.forwarder !== undefined) {
      if (!forwarders.includes(params.forwarder)) {
        throw new Response(
          `Forwarder ${params.forwarder} not found`,
          { status: 404 },
        )
      }

      const forwarder = await configApi.getForwarder(params.forwarder)
      return {
        forwarders: forwarders,
        currentForwarder: {
          name: params.forwarder,
          forwarder,
        },
      }
    }

    return { forwarders: forwarders }
  },
)

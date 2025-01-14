import { LoaderFunction } from 'react-router'

import { loginRequired } from '@/lib/decorators/loaders'
import * as configApi from '@/lib/api/operations/config'
import { WebhookModel } from '@/lib/models'

export type LoaderData = {
  alerts: string[]
  currentAlert?: {
    name: string
    webhook: WebhookModel
  }
}

export const loader: LoaderFunction = loginRequired(
  async ({ params }): Promise<LoaderData> => {
    const alerts = await configApi.listAlerts()

    if (params.alert !== undefined) {
      if (!alerts.includes(params.alert)) {
        throw new Response(
          `Alert ${params.alert} not found`,
          { status: 404 },
        )
      }

      const webhook = await configApi.getAlert(params.alert)
      return {
        alerts,
        currentAlert: {
          name: params.alert,
          webhook,
        },
      }
    }

    return { alerts }
  },
)

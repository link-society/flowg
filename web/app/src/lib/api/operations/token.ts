import * as request from '@/lib/api/request'

import { TokenModel } from '@/lib/models'

export const listTokens = async (): Promise<string[]> => {
  type ListTokensResponse = {
    success: boolean
    token_uuids: string[]
  }
  const { body } = await request.GET<ListTokensResponse>('/api/v1/tokens')
  return body.token_uuids
}

export const createToken = async (): Promise<TokenModel> => {
  type CreateTokenResponse = {
    success: boolean
    token: string
    token_uuid: string
  }

  const { body } = await request.POST<{}, CreateTokenResponse>('/api/v1/token', {})
  return { token: body.token, token_uuid: body.token_uuid }
}

export const deleteToken = async (tokenUuid: string): Promise<void> => {
  await request.DELETE(`/api/v1/tokens/${tokenUuid}`)
}

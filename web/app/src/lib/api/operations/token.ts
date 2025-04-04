import * as request from '@/lib/api/request'
import { TokenModel } from '@/lib/models/auth'

export const listTokens = async (): Promise<string[]> => {
  type ListTokensResponse = {
    success: boolean
    token_uuids: string[]
  }
  const { body } = await request.GET<ListTokensResponse>({
    path: '/api/v1/tokens',
  })
  return body.token_uuids
}

export const createToken = async (): Promise<TokenModel> => {
  type CreateTokenResponse = {
    success: boolean
    token: string
    token_uuid: string
  }

  const { body } = await request.POST<unknown, CreateTokenResponse>({
    path: '/api/v1/token',
    body: {},
  })
  return { token: body.token, token_uuid: body.token_uuid }
}

export const deleteToken = async (tokenUuid: string): Promise<void> => {
  await request.DELETE({
    path: `/api/v1/tokens/${tokenUuid}`,
  })
}

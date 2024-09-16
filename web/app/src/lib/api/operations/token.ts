import * as request from '@/lib/api/request'

export const listTokens = async (): Promise<string[]> => {
  type ListTokensResponse = {
    success: boolean
    token_uuids: string[]
  }
  const { body } = await request.GET<ListTokensResponse>('/api/v1/tokens')
  return body.token_uuids
}

export const deleteToken = async (tokenUuid: string): Promise<void> => {
  await request.DELETE(`/api/v1/tokens/${tokenUuid}`)
}

import { UserModel } from '@/lib/models'

export type WhoamiResponse = {
  success: boolean
  user: UserModel
}

export const whoami = async (): Promise<WhoamiResponse> => {
  const response = await fetch('/api/auth/whoami')
  return response.json()
}

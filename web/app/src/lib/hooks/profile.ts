import { useContext } from 'react'

import ProfileContext from '@/lib/context/profile'

export const useProfile = () => {
  return useContext(ProfileContext)
}

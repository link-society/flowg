import { createContext, useContext } from 'react'
import { ProfileModel } from '@/lib/models'

const ProfileContext = createContext<ProfileModel>(null!)

export const ProfileProvider = ProfileContext.Provider

export const useProfile = () => {
  return useContext(ProfileContext)
}

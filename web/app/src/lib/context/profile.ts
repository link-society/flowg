import { createContext } from 'react'

import ProfileModel from '@/lib/models/ProfileModel'

const ProfileContext = createContext<ProfileModel>(null!)

export default ProfileContext

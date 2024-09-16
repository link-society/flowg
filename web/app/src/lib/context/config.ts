import { createContext, useContext } from 'react'

type Config = {
  notifications?: {
    autoHideDuration?: number
  }
}

const ConfigContext = createContext<Config>(null!)

export const ConfigProvider = ConfigContext.Provider

export const useConfig = () => {
  return useContext(ConfigContext)
}

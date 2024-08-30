import React, { createContext, useContext, useState } from 'react'

type DnDContextType = ReturnType<typeof useState<string | undefined>>

const DnDContext = createContext<DnDContextType>(['', () => {}])

export const useDnD = () => {
  return useContext(DnDContext)
}

export const DnDProvider: React.FC<{children: React.ReactNode}> = ({ children }) => {
  const [type, setType] = useState<string | undefined>('')

  return (
    <DnDContext.Provider value={[type, setType]}>
      {children}
    </DnDContext.Provider>
  )
}

export default DnDContext

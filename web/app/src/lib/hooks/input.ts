import { Dispatch, SetStateAction, useCallback, useMemo, useState } from 'react'

import { Validator } from '@/lib/validators'

type InputData<T> = {
  value: T
  valid: boolean
}

export const useInput = <T>(
  initialState: T | (() => T),
  validators?: Array<Validator<T>>
): [InputData<T>, Dispatch<SetStateAction<T>>] => {
  const validate = useCallback(
    (value: T): boolean => {
      for (const validator of validators ?? []) {
        if (!validator(value)) {
          return false
        }
      }
      return true
    },
    [validators]
  )

  const [value, setValue] = useState<T>(initialState)
  const valid = useMemo(() => validate(value), [value, validate])
  const container = useMemo(() => ({ value, valid }), [value, valid])

  return [container, setValue]
}

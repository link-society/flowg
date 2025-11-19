import { Validator } from '@/lib/validators'

import { Dispatch, SetStateAction, useCallback, useMemo, useState } from 'react'

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

  return [{ value, valid }, setValue]
}

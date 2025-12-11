export type Validator<T> = (value: T) => boolean

export const minLength =
  (length: number) =>
  (value: string): boolean => {
    return value.length >= length
  }

export const maxLength =
  (length: number) =>
  (value: string): boolean => {
    return value.length <= length
  }

export const minItems =
  <T>(count: number) =>
  (value: T[]): boolean => {
    return value.length >= count
  }

export const items =
  <T>(validators: Array<Validator<T>>) =>
  (value: T[]): boolean => {
    for (const item of value) {
      for (const validator of validators) {
        if (!validator(item)) {
          return false
        }
      }
    }

    return true
  }

export const pattern =
  (regex: RegExp) =>
  (value: string): boolean => {
    return regex.test(value)
  }

export const formatUri = (value: string): boolean => {
  try {
    new URL(value)
    return true
  } catch {
    return false
  }
}

export const useDirty = <T>(initial: T, current: T): boolean =>
  JSON.stringify(initial) !== JSON.stringify(current)

export type DynamicField<T extends string> = T | `@expr:${string}`

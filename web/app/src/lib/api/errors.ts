export class ApiError extends Error {}

export class UnauthenticatedError extends ApiError {}

export class PermissionDeniedError extends ApiError {}

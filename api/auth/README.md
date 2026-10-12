# auth

The package at `api/auth` provides the authentication and
authorization building blocks shared by FlowG's REST API handlers.

It exists to keep identity and permission concerns out of the individual
endpoint use-cases. Endpoints only declare _who_ may call them and _which_
permissions are required; this package takes care of establishing the caller's
identity from the incoming request and enforcing the declared permissions.

## Responsibilities

- **Authentication** — `ApiMiddleware` is the HTTP boundary that turns an
  anonymous request into an authenticated one, accepting personal access tokens
  and JWTs as bearer credentials.
- **Session tokens** — `Keys.NewJWT` mints the time-limited token handed back
  on successful login, and `Keys.VerifyJWT` validates it on subsequent requests
  so the password never has to be replayed.
- **Identity propagation** — `ContextWithUser` and `GetContextUser` carry the
  authenticated user across the request pipeline through its context, so
  handlers and decorators never need to re-authenticate.
- **Authorization** — `RequireScopeApiDecorator` and
  `RequireScopesApiDecorator` wrap use-case interactors to enforce that the
  caller holds the required permission scopes before any business logic runs.

## Usage shape

A typical endpoint is composed by chaining these pieces: the middleware
authenticates the request and binds the user to the context, then a scope
decorator wraps the interactor so it only runs for authorized callers. Both
are given the logger they log through (the operation's API logger).

```text
request ──▶ ApiMiddleware ──▶ RequireScope(s)ApiDecorator ──▶ use-case
            (authenticate)     (authorize)                     (business logic)
```

## Configuration

JWTs are signed and verified with the `Keys` provided to the dependency graph,
built from the `security.secret_key` setting (or `FLOWG_SECRET_KEY`). When it is
unset a random key is generated, which is fine for local development but
invalidates all issued tokens whenever the process restarts — set it explicitly
in production so tokens survive restarts.

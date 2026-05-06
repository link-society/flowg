# Frontend Architecture

This document covers the structure and conventions of the FlowG frontend. It is meant as a reference for anyone working on the codebase.

---

## Project Architecture

The frontend is a React SPA. Source code lives under `src/`, split into folders by concern.

```
src/
 ├─ components/
 ├─ layouts/
 ├─ lib/
 │   ├─ api/
 │   ├─ context/
 │   ├─ decorators/
 │   ├─ hooks/
 │   └─ models/
 ├─ router/
 ├─ theme/
 ├─ views/
 └─ App.tsx
```

### `components/`

Generic, reusable UI components. They should not depend on any specific view or feature. If a component can be used in more than one place, it belongs here.

### `layouts/`

Page templates (sidebars, headers, shells). They define the frame around view content but hold no business logic.

### `lib/`

Shared logic used across the app:

- **`api/`** — All backend communication. HTTP calls should not leak outside this folder.
- **`context/`** — React contexts and their providers for shared application state (auth, notifications, etc.).
- **`decorators/`** — Higher-order components for cross-cutting concerns like access control.
- **`hooks/`** — Custom hooks that encapsulate reusable stateful logic (data fetching, form state, etc.).
- **`models/`** — TypeScript types and interfaces for API payloads and domain objects.

### `router/`

Route definitions. Maps URLs to views and handles navigation guards and lazy loading.

### `theme/`

Design system configuration: MUI theme, global styles, and design tokens (colors, typography, spacing). All application-wide visual customization goes here.

### `views/`

One component per page/screen. Views fetch data (via hooks or context), compose layouts and components, and handle screen-specific interactions. Keep business logic out of views — delegate to `lib/`.

### `App.tsx`

The app entry point. Sets up the router, context providers, and theme.

---

## Component Structure

Each component is a folder rather than a single file. This keeps concerns separated and avoids large, hard-to-navigate files.

```
ComponentName/
 ├─ component.tsx   # JSX and rendering logic
 ├─ styles.tsx      # Styled components using MUI's `styled` utility
 └─ types.ts        # Props and any types local to this component
```

Additional files can be added as needed, for example, `hooks.ts` for component-local state logic, or `utils.ts` for pure helpers — but only when the component genuinely requires them.

This convention applies to `components/`, `layouts/`, and `views/`.

---

## Component Declaration Convention

All components are declared as typed constants using React's `FC` type, rather than standard function declarations.

```tsx
import type { FC } from 'react'

import type { BannerProps } from './types'

export const Banner: FC<BannerProps> = ({
  dataTestId,
  profileImageUrl,
  variant,
}) => {
  // component logic
  return <div>...</div>
}

export default Banner
```

This style is preferred because it makes the props type explicit at the declaration site and keeps component definitions consistent across the codebase. The named export is used for direct imports; the default export is provided for compatibility with lazy-loaded routes.

---

## Styling

All styles are written using MUI's `styled` utility and live exclusively in the colocated `styles.tsx` file. `component.tsx` must not contain any `sx` prop or inline `style` attribute — if something needs styling, it belongs in `styles.tsx`.

**Always use theme tokens instead of hard-coded values.** This ensures consistency and makes theme changes propagate automatically.

```tsx
// ComponentName/styles.tsx
import { styled } from '@mui/material/styles'

export const Container = styled('div')(({ theme }) => ({
  position: 'absolute',
  left: theme.spacing(2),
  display: 'flex',
  alignItems: 'center',
  width: 150,
  height: 150,
  backgroundColor: theme.palette.grey[200],
  overflow: 'hidden',
}))
```

```tsx
// ComponentName/component.tsx
import { Container } from './styles'

export const ComponentName: FC<ComponentNameProps> = () => (
  <Container>...</Container>
)
```

**Guidelines:**

- **Colors** — always access colors via `theme.tokens.colors` inside `styled`. Never use hard-coded hex values, direct imports of `colors.ts`, or `theme.palette` inline in `styles.tsx`. If a color is missing from `theme/tokens/colors.ts`, add it there first.
- Use `theme.spacing` for spacing and `theme.typography` for text styles.
- Use `theme.shadows` for elevation/shadow values.
- Import `styled` from `@mui/material/styles` to get automatic theme access.
- Avoid inline `style={{}}` or `sx={{}}`

```tsx
// ✅ correct — color via theme.tokens.colors
export const Header = styled('div')`
  background-color: ${({ theme }) => theme.tokens.colors.primary};
  color: ${({ theme }) => theme.tokens.colors.primaryContrast};
`

// ❌ wrong — direct import of colors
import { colors } from '@/theme/tokens'
export const Header = styled('div')`
  background-color: ${colors.primary};
`

// ❌ wrong — hard-coded hex
export const Header = styled('div')`
  background-color: #1565c0;
`
```

---

## Typography

All visible text — headings, labels, descriptions, captions, button labels when not handled by a dedicated MUI component — must be wrapped in MUI's `<Typography>` component. Never use bare `<span>`, `<p>`, `<h1>`–`<h6>`, or string literals directly in JSX.

Use the custom variants defined in `theme/theme.ts`:

| Variant   | Use case                         |
| --------- | -------------------------------- |
| `titleLg` | Page titles                      |
| `titleMd` | Section headings                 |
| `titleSm` | Card headers, sub-section titles |
| `text`    | Body text, labels                |

```tsx
// ✅ correct
<Typography variant="titleSm">Account Information</Typography>

// ❌ wrong — bare span
<span>Account Information</span>

// ❌ wrong — bare text node
Account Information
```

---

## Extending the Design System

New colors, typography scales, spacing values, or any other design token must be defined in `theme/tokens/`. Never hard-coded in a component.

```
theme/
 ├─ tokens/
 │   ├─ colors.ts
 │   ├─ typography.ts
 │   └─ spacing.ts
 └─ index.ts   # builds and exports the MUI theme
```

Define the raw value in the relevant token file, then consume it in `theme/index.ts` when building the MUI theme object.

If a token introduces a property that MUI does not have by default (custom palette keys, new typography variants, etc.), declare it in `styled.d.ts` via MUI's theme augmentation so TypeScript is aware of it:

```ts
// styled.d.ts
import '@mui/material/styles'

declare module '@mui/material/styles' {
  interface Palette {
    brand: {
      primary: string
      secondary: string
    }
  }
  interface PaletteOptions {
    brand?: {
      primary?: string
      secondary?: string
    }
  }
}
```

This keeps the type system in sync and prevents silent `any` access on custom theme properties.

---

## Global Styles

Application-wide CSS resets or baseline styles are defined using MUI's `GlobalStyles` component, placed in `theme/`. This keeps global style concerns co-located with the rest of the design system rather than scattered across the app.

```tsx
// theme/global-styles.tsx
import GlobalStyles from '@mui/material/GlobalStyles'

export const AppGlobalStyles = () => (
  <GlobalStyles
    styles={(theme) => ({
      '*, *::before, *::after': {
        boxSizing: 'border-box',
      },
      body: {
        margin: 0,
        backgroundColor: theme.palette.background.default,
        color: theme.palette.text.primary,
        fontFamily: theme.typography.fontFamily,
      },
    })}
  />
)
```

Render `<AppGlobalStyles />` once at the root of the app, inside the `ThemeProvider`, so theme tokens are available:

```tsx
// App.tsx
<ThemeProvider theme={theme}>
  <AppGlobalStyles />
  <RouterProvider router={router} />
</ThemeProvider>
```

**Guidelines:**

- Only put truly global rules here (resets, `body`, `html`, custom scrollbar styles, etc.).
- Never use `GlobalStyles` to style a specific component — that belongs in the component's `styles.tsx`.
- Always use theme tokens; avoid hard-coded values.

---

## i18n

The app uses [i18next](https://www.i18next.com/) for internationalization. Translation strings are stored in `.po` files (gettext format), which is the most portable format and is natively supported by most professional translation tools.

### File format

We use `.po` / `.mo` files (gettext) rather than plain JSON. This makes it straightforward to hand off translation files to external translators without requiring any knowledge of the codebase.

The [`i18next-gettext-converter`](https://github.com/i18next/i18next-gettext-converter) plugin handles conversion between gettext and the JSON format i18next consumes at runtime.

### Workflow

1. Strings are written in source code using i18next's `t()` function.
2. `.po` files are extracted and sent to translators.
3. Translators work with their preferred tool and return the translated `.po` files.
4. The converter produces the JSON files consumed by the app at build time.

### Usage in components

Use the `useTranslation` hook and alias `t` to `translate` for clarity. Keys follow a `<scope>.<screen>.<element>` naming convention.

```tsx
import type { FC } from 'react'
import { useTranslation } from 'react-i18next'

import { Typography } from '@/components'

export const HomePage: FC = () => {
  const { t } = useTranslation()

  return (
    <>
      <Typography as="h1" variant="titleXl">
        {t('pages.home.title')}
      </Typography>

      <Typography as="p" variant="textMd" color="textTertiary">
        {t('pages.home.description')}
      </Typography>
    </>
  )
}
```

### `.po` file example

```po
# src/locales/en.po
msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Language: en\n"

msgid "pages.home.title"
msgstr "Welcome to FlowG"

msgid "pages.home.description"
msgstr "A log management platform."
```

### Guidelines

**1. Never hardcode text**

All visible text must come from `t()`. No string literals in JSX.

**2. Always wrap translated text in `Typography`**

This applies to headings, labels, descriptions, and button text (when not handled by a dedicated component).

**3. Use a clear, hierarchical key structure**

| Prefix                         | Usage                               |
| ------------------------------ | ----------------------------------- |
| `common.*`                     | Shared labels reused across the app |
| `pages.<pageName>.*`           | Content specific to a page          |
| `components.<componentName>.*` | Text owned by a specific component  |
| `error.*`                      | Error messages                      |

> **`common.*`** is reserved for truly generic UI labels reused across the entire application (e.g. _Save_, _Cancel_, _Back_). Feature-specific labels must stay scoped to their own domain and must not be placed under `common.*`.

**4. Provide context for ambiguous keys**

Include a `description` whenever possible — it is propagated to the generated JSON and `.po` files and helps translators understand the intended context.

```ts
t('some.key', { description: 'Shown when the user saves their profile' })
```

Identical keys cannot have different descriptions unless a distinct `context` is also provided.

```ts
// invalid — same key, different descriptions, no context
t('some.key', { description: 'a' })
t('some.key', { description: 'b' })

// valid — same key, different descriptions, distinguished by context
t('some.key', { description: 'a', context: 'foo' })
t('some.key', { description: 'b', context: 'bar' })
```

**5. Prefer full, explicit paths**

Use `pages.home.title` rather than `home.title`. Explicit paths make keys searchable and avoid ambiguity when the same word appears in different contexts.

**6. Use predictable key names**

Stick to a consistent vocabulary for leaf keys:

- `title`, `subtitle`, `description`
- `action`, `cta`, `successMessage`
- `emptyState.title`, `emptyState.description`

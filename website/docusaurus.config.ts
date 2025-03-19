import prismMonokaiTheme from './prism.theme.monokai'
import type { Config } from '@docusaurus/types'
import type * as Preset from '@docusaurus/preset-classic'
import type * as Redocusaurus from 'redocusaurus'

import path from 'path'

const config: Config = {
  title: 'FlowG',
  tagline: 'Free and Open-Source Low-Code log processing solution',
  favicon: 'img/favicon.ico',

  url: 'https://link-society.github.io/',
  baseUrl: '/flowg/',
  trailingSlash: false,

  organizationName: 'link-society',
  projectName: 'flowg',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
    [
      'redocusaurus',
      {
        config: path.join(__dirname, 'redocly.yaml'),
        specs: [
          { spec: './src/openapi.json' },
        ]
      },
    ] satisfies Redocusaurus.PresetEntry,
  ],

  markdown: {
    mermaid: true,
  },

  themes: [
    '@docusaurus/theme-mermaid',
  ],

  themeConfig: {
    colorMode: {
      defaultMode: 'light',
      disableSwitch: true,
      respectPrefersColorScheme: false,
    },
    navbar: {
      title: 'FlowG',
      style: 'primary',
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docSidebar',
          position: 'left',
          label: 'Documentation',
        },
        {
          href: 'https://github.com/link-society/flowg',
          html: `
            <div style="display: flex; align-items: center;">
              <img
                alt="GitHub Release"
                src="https://img.shields.io/github/v/release/link-society/flowg?style=social"
              />
            </div>
          `,
          position: 'right',
        },
        {
          href: 'https://github.com/link-society/flowg',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'light',
      copyright: `
        Website content is distributed under the terms of the
        <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC-BY-SA 4.0</a>
        license.
      `,
    },
    prism: {
      theme: prismMonokaiTheme as any,
      darkTheme: prismMonokaiTheme as any,
      additionalLanguages: ['bash', 'ini', 'apacheconf', 'nginx'],
    },
  } satisfies Preset.ThemeConfig,
}

export default config

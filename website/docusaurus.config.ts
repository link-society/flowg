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
        blog: {
          blogTitle: 'FlowG Blog',
          blogDescription: 'Updates and news about FlowG',
          postsPerPage: 'ALL',
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
    hooks: {
      onBrokenMarkdownImages: 'warn',
      onBrokenMarkdownLinks: 'warn',
    },
  },

  plugins: [
    'plugin-image-zoom',
  ],

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
          to: 'blog',
          label: 'Blog',
          position: 'left',
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
    imageZoom: {
      selector: '.markdown div.with-zoom img',
      options: {
        margin: 24,
        background: 'rgba(0, 0, 0, 0.9)',
      },
    }
  } satisfies Preset.ThemeConfig,
}

export default config

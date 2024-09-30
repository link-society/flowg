import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

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
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;

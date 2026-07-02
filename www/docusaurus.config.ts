import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Kobo Docs',
  tagline: 'Virtual Account Infrastructure',
  favicon: 'img/favicon.png',

  url: 'https://kobo.triumphsystems.tech',
  baseUrl: '/',

  organizationName: 'triumphsystems',
  projectName: 'kobo',

  onBrokenLinks: 'throw',

  markdown: {
    hooks: {
      onBrokenMarkdownLinks: 'warn',
    },
  },

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
          docItemComponent: "@theme/ApiItem",
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  plugins: [
    [
      'docusaurus-plugin-openapi-docs',
      {
        id: "api", 
        docsPluginId: "classic",
        config: {
          kobo: {
            specPath: "../core/openapi.yaml",
            outputDir: "docs/api",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "tag",
            },
          },
        }
      },
    ],
    './src/plugins/webpack-fallback/index.js',
  ],

  themes: ["docusaurus-theme-openapi-docs"],

  themeConfig: {
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: true,
      respectPrefersColorScheme: false,
    },
    navbar: {
      title: '',
      logo: {
        alt: 'Kobo Logo',
        src: 'img/logo.png',
      },
      items: [
        {
          type: 'html',
          position: 'left',
          value: '<div class="navbar-brand-extras"><span class="navbar-docs-pill">DOCS</span><div class="navbar-separator"></div></div>',
        },
        {
          type: 'html',
          position: 'left',
          value: `
            <a href="/docs/intro" class="navbar-version-dropdown">
              <span>v1.x</span>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 9l6 6 6-6"/></svg>
            </a>
          `,
        },
        {
          type: 'html',
          position: 'left',
          className: 'navbar-search-wrapper',
          value: `
            <div class="navbar-fake-search">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
              <span class="search-text">Search</span>
              <div class="search-keys"><kbd>CTRL</kbd><kbd>K</kbd></div>
            </div>
          `,
        },
        {
          type: 'html',
          position: 'right',
          value: `
            <a href="/docs/intro" class="navbar-assistant-link">
              <span style="color: #FACC15;">✦</span> <span>Documentation</span>
            </a>
          `,
        },
        {
          type: 'html',
          position: 'right',
          value: `
            <a href="/docs/api/kobo-api" class="navbar-assistant-link" style="margin-left: 1rem;">
              <span>API Reference</span>
            </a>
          `,
        },
        {
          href: 'https://github.com/leoemaxie/kobo',
          label: 'GitHub',
          position: 'right',
          className: 'navbar-simple-link',
        },
        {
          type: 'html',
          position: 'right',
          value: `
            <a href="#" class="navbar-menu-dots">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="#059669"><circle cx="5" cy="12" r="2"/><circle cx="12" cy="12" r="2"/><circle cx="19" cy="12" r="2"/></svg>
            </a>
          `,
        },
        {
          href: 'https://app.kobo.triumphsystems.tech',
          position: 'right',
          className: 'button button--primary nav-login-btn',
          label: 'Log In' /* Using standard text, arrow can be styled */
        },
      ],
    },
    footer: {
      style: 'dark',
      copyright: `
        <div class="kobo-custom-footer">
          <div class="kobo-footer-brand">
            <img src="/img/logo.png" alt="Kobo Logo" />
          </div>
          <div class="kobo-footer-copyright">
            © ${new Date().getFullYear()} Triumph Systems. All rights reserved
          </div>
          <div class="kobo-footer-socials">
            <a href="https://github.com/leoemaxie/kobo" target="_blank" rel="noopener noreferrer" aria-label="GitHub">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.166 6.839 9.489.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.462-1.11-1.462-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.379.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.161 22 16.416 22 12c0-5.523-4.477-10-10-10z"/></svg>
            </a>
            <a href="#" target="_blank" rel="noopener noreferrer" aria-label="Twitter">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M22.46 6c-.77.35-1.6.58-2.46.69.88-.53 1.56-1.37 1.88-2.38-.83.5-1.75.85-2.72 1.05C18.37 4.5 17.26 4 16 4c-2.35 0-4.27 1.92-4.27 4.29 0 .34.04.67.11.98C8.28 9.09 5.11 7.38 3 4.79c-.37.63-.58 1.37-.58 2.15 0 1.49.75 2.81 1.91 3.56-.71 0-1.37-.2-1.95-.5v.05c0 2.08 1.48 3.82 3.44 4.21a4.22 4.22 0 0 1-1.93.07 4.28 4.28 0 0 0 4 2.98 8.521 8.521 0 0 1-5.33 1.84c-.34 0-.68-.02-1.02-.06C3.44 20.29 5.7 21 8.12 21 16 21 20.33 14.46 20.33 8.79c0-.19 0-.37-.01-.56.84-.6 1.56-1.36 2.14-2.23z"/></svg>
            </a>
            <a href="#" target="_blank" rel="noopener noreferrer" aria-label="Discord">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor"><path d="M20.317 4.3698a19.7913 19.7913 0 00-4.8851-1.5152.0741.0741 0 00-.0785.0371c-.211.3753-.4447.8648-.6083 1.2495-1.8447-.2762-3.68-.2762-5.4868 0-.1636-.3933-.4058-.8742-.6177-1.2495a.077.077 0 00-.0785-.037 19.7363 19.7363 0 00-4.8852 1.515.0699.0699 0 00-.0321.0277C.5334 9.0458-.319 13.5799.0992 18.0578a.0824.0824 0 00.0312.0561c2.0528 1.5076 4.0413 2.4228 5.9929 3.0294a.0777.0777 0 00.0842-.0276c.4616-.6304.8731-1.2952 1.226-1.9942a.076.076 0 00-.0416-.1057c-.6528-.2476-1.2743-.5495-1.8722-.8923a.077.077 0 01-.0076-.1277c.1258-.0943.2517-.1923.3718-.2914a.0743.0743 0 01.0776-.0105c3.9278 1.7933 8.18 1.7933 12.0614 0a.0739.0739 0 01.0785.0095c.1202.099.246.1981.3728.2924a.077.077 0 01-.0066.1276 12.2986 12.2986 0 01-1.873.8914.0766.0766 0 00-.0407.1067c.3604.698.7719 1.3628 1.225 1.9932a.076.076 0 00.0842.0286c1.961-.6067 3.9495-1.5219 6.0023-3.0294a.077.077 0 00.0313-.0552c.5004-5.177-.8382-9.6739-3.5485-13.6604a.061.061 0 00-.0312-.0286zM8.02 15.3312c-1.1825 0-2.1569-1.0857-2.1569-2.419 0-1.3332.9555-2.4189 2.157-2.4189 1.2108 0 2.1757 1.0952 2.1568 2.419 0 1.3332-.9555 2.4189-2.1569 2.4189zm7.9748 0c-1.1825 0-2.1569-1.0857-2.1569-2.419 0-1.3332.9554-2.4189 2.1569-2.4189 1.2108 0 2.1757 1.0952 2.1568 2.419 0 1.3332-.946 2.4189-2.1568 2.4189Z"/></svg>
            </a>
          </div>
        </div>
      `,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['java', 'go', 'json', 'bash'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;

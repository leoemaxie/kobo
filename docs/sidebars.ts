import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */
const sidebars: SidebarsConfig = {
  tutorialSidebar: [
    {
      type: 'category',
      label: 'How to build with Kobo',
      collapsed: false,
      items: [
        'intro',
        'concepts/primitives',
        'concepts/lifecycle',
        'concepts/reconciliation'
      ],
    },
    {
      type: 'category',
      label: 'SDKs',
      collapsed: false,
      items: [
        'sdks/typescript',
        'sdks/go',
        'sdks/java'
      ],
    }
  ],
  apiSidebar: [
    {
      type: "category",
      label: "API Reference",
      link: {
        type: "generated-index",
        title: "Kobo API",
        description: "Explore the Kobo REST API endpoints.",
        slug: "/api"
      },
      items: require("./docs/api/sidebar.ts")
    }
  ]
};

export default sidebars;

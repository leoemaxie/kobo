import React, {type ReactNode} from 'react';
import Layout from '@theme-original/ApiItem/Layout';
import type LayoutType from '@theme/ApiItem/Layout';
import type {WrapperProps} from '@docusaurus/types';

type Props = WrapperProps<typeof LayoutType>;

import { McpInstallButton } from 'docusaurus-plugin-mcp-server/theme';

export default function LayoutWrapper(props: Props): ReactNode {
  return (
    <>
      <div style={{ padding: '0 0 16px 0', display: 'flex', alignItems: 'center' }}>
        <McpInstallButton 
          serverUrl="https://kobo.triumphsystems.tech/mcp" 
          serverName="kobo-docs" 
          label="Open in Claude"
          className="nomba-mcp-btn"
        />
      </div>
      <Layout {...props} />
    </>
  );
}

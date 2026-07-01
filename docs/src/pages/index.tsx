import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Heading from '@theme/Heading';
import HomepageFeatures from '../components/HomepageFeatures';

import styles from './index.module.css';
import { JSX } from 'react';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className="container">
        <Heading as="h1" className="hero__title" style={{ fontFamily: 'Outfit, sans-serif', fontWeight: 800, fontSize: '4rem', letterSpacing: '-0.02em', color: '#FFFFFF' }}>
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle" style={{ color: '#A1A1AA', fontSize: '1.5rem', maxWidth: '600px', margin: '0 auto 2rem auto' }}>
          {siteConfig.tagline}
        </p>
        <div className={styles.buttons}>
          <Link
            className="button button--primary button--lg"
            style={{ borderRadius: '8px', padding: '12px 24px', fontWeight: 600, color: '#09090B' }}
            to="/docs/intro">
            Read the Docs
          </Link>
          <Link
            className="button button--secondary button--lg"
            style={{ borderRadius: '8px', padding: '12px 24px', fontWeight: 600, marginLeft: '16px', backgroundColor: '#27272A', color: '#FFFFFF', border: 'none' }}
            to="https://github.com/leoemaxie/kobo">
            View on GitHub
          </Link>
        </div>
      </div>
    </header>
  );
}

export default function Home(): JSX.Element {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`Welcome to ${siteConfig.title}`}
      description="The definitive B2B ledger and reconciliation engine.">
      <main style={{ backgroundColor: '#09090B' }}>
        <HomepageHeader />
        <HomepageFeatures />
      </main>
    </Layout>
  );
}

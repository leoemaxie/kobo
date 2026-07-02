import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '../components/HomepageFeatures';

import styles from './index.module.css';
import { JSX } from 'react';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <header className={clsx('hero', styles.heroBanner)}>
      <div className="container">
        <h1 className={styles.heroTitle}>
          Identity-anchored virtual account infrastructure
        </h1>
        <p className={styles.heroSubtitle}>
          Nomba-backed, per-identity dedicated accounts, automatic reconciliation.
        </p>
        <div className={styles.buttons}>
          <Link
            className={styles.primaryButton}
            to="/docs/intro">
            Read the docs
          </Link>
          <Link
            className={styles.ghostButton}
            to="/docs/api/kobo-api">
            View API reference
          </Link>
        </div>
        
        {/* Terminal Code Block */}
        <div className={styles.codeBlockWrapper}>
          <div className={styles.codeBlock}>
            <div className={styles.codeBlockHeader}>
              <span className={styles.codeDot} style={{backgroundColor: '#ff5f56'}}></span>
              <span className={styles.codeDot} style={{backgroundColor: '#ffbd2e'}}></span>
              <span className={styles.codeDot} style={{backgroundColor: '#27c93f'}}></span>
            </div>
            <pre className={styles.codePre}>
              <code>
<span className={styles.codePrompt}>$</span> curl -X POST https://api.kobo.triumphsystems.tech/v1/identities \
  -H "Authorization: Bearer sk_test_..." \
  -H "Content-Type: application/json" \
  -d '&#123;
    "type": "individual",
    "bvn": "22345678901",
    "first_name": "Satoshi",
    "last_name": "Nakamoto"
  &#125;'
              </code>
            </pre>
          </div>
        </div>
      </div>
    </header>
  );
}

function HowItWorks() {
  return (
    <section className={styles.howItWorksSection}>
      <div className="container">
        <h2 className={styles.sectionEyebrow}>CORE PRIMITIVES</h2>
        <h3 className={styles.sectionTitle}>How it works</h3>
        <div className={styles.stepsContainer}>
          <div className={styles.stepItem}>
            <div className={styles.stepNumber}>1</div>
            <h4 className={styles.stepTitle}>Register an identity</h4>
            <p className={styles.stepDesc}>Kobo provisions a dedicated Nomba virtual account.</p>
          </div>
          <div className={styles.stepConnector} />
          <div className={styles.stepItem}>
            <div className={styles.stepNumber}>2</div>
            <h4 className={styles.stepTitle}>Funds arrive</h4>
            <p className={styles.stepDesc}>Automatically reconciled against that identity.</p>
          </div>
          <div className={styles.stepConnector} />
          <div className={styles.stepItem}>
            <div className={styles.stepNumber}>3</div>
            <h4 className={styles.stepTitle}>Lifecycle events</h4>
            <p className={styles.stepDesc}>Handled by a defined state machine, not ad hoc code.</p>
          </div>
          <div className={styles.stepConnector} />
          <div className={styles.stepItem}>
            <div className={styles.stepNumber}>4</div>
            <h4 className={styles.stepTitle}>Your product</h4>
            <p className={styles.stepDesc}>Reads statements and transaction history through one API.</p>
          </div>
        </div>
      </div>
    </section>
  );
}

function ReferenceImplementation() {
  return (
    <section className={styles.referenceSection}>
      <div className="container">
        <div className={styles.referenceContent}>
          <p>
            Built and demonstrated with a school-fee collection app as a reference integrator — proving a second product team could build against the same API with zero special access.
          </p>
          <a href="https://fees.kobo.triumphsystems.tech" target="_blank" rel="noopener noreferrer" className={styles.referenceLink}>
            View reference application →
          </a>
        </div>
      </div>
    </section>
  );
}

function FinalCTA() {
  return (
    <section className={styles.finalCtaSection}>
      <div className="container">
        <div className={styles.finalCtaCard}>
          <h2 className={styles.finalCtaTitle}>
            Start building with <span className={styles.inlineHighlight}>Kobo</span> today
          </h2>
          <div className={styles.buttons} style={{ marginTop: '2rem' }}>
            <Link className={styles.primaryButton} to="/docs/intro">
              Read the docs
            </Link>
            <Link className={styles.ghostButton} to="/docs/api/kobo-api">
              View API reference
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}

export default function Home(): JSX.Element {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`Welcome to ${siteConfig.title}`}
      description="Identity-anchored virtual account infrastructure for Nigerian fintech">
      <main className={styles.mainCanvas}>
        <HomepageHeader />
        <HomepageFeatures />
        <HowItWorks />
        <ReferenceImplementation />
        <FinalCTA />
      </main>
    </Layout>
  );
}

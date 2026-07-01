import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import { JSX } from 'react';

type FeatureItem = {
  title: string;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Virtual Accounts',
    description: (
      <>
        Provision dedicated NUBANs instantly via the Nomba platform. Perfect for user wallets, payment collections, and reconciliation.
      </>
    ),
  },
  {
    title: 'Automated Ledger',
    description: (
      <>
        Never calculate balances manually. Every inbound transfer automatically creates an immutable, cryptographically verifiable ledger entry.
      </>
    ),
  },
  {
    title: 'Type-Safe SDKs',
    description: (
      <>
        Integrate seamlessly with our zero-dependency, generated SDKs for TypeScript, Go, and Java. Built for microservices.
      </>
    ),
  },
];

function Feature({title, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className={styles.featureCard}>
        <Heading as="h3" className={styles.featureTitle}>{title}</Heading>
        <p className={styles.featureDescription}>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}

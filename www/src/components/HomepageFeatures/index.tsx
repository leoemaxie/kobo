import clsx from 'clsx';
import styles from './styles.module.css';
import { JSX } from 'react';

type FeatureItem = {
  title: string;
  description: JSX.Element;
  icon: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Dedicated identity, not a shared pool',
    description: (
      <>
        Every customer gets their own NUBAN-style account. No more manual matching or relying on shared collection pools.
      </>
    ),
    icon: (
      <svg viewBox="0 0 24 24"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
    ),
  },
  {
    title: 'Reconciliation that\'s actually reliable',
    description: (
      <>
        Idempotent design. Handles duplicate webhooks, delayed webhook delivery, partial payments, and overpayments out of the box.
      </>
    ),
    icon: (
      <svg viewBox="0 0 24 24"><path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path><polyline points="22,6 12,13 2,6"></polyline></svg>
    ),
  },
  {
    title: 'A lifecycle that\'s actually modeled',
    description: (
      <>
        Renames, closures, KYC tier changes, and misdirected payments have defined, tested behavior, not just a happy path.
      </>
    ),
    icon: (
      <svg viewBox="0 0 24 24"><polyline points="16 18 22 12 16 6"></polyline><polyline points="8 6 2 12 8 18"></polyline></svg>
    ),
  },
];

function Feature({title, description, icon}: FeatureItem) {
  return (
    <div className={clsx('col col--4')}>
      <div className={styles.featureCard}>
        <div className={styles.featureIconWrapper}>
          {icon}
        </div>
        <h3 className={styles.featureTitle}>{title}</h3>
        <p className={styles.featureDescription}>{description}</p>
        <a href="/docs/intro" className={styles.featureLink}>Explore →</a>
      </div>
    </div>
  );
}

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <h2 className={styles.sectionEyebrow}>CAPABILITIES</h2>
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}

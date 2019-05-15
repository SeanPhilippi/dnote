import React from 'react';

import Plan from './internal';
import ServerIcon from '../../Icons/Server';
import GlobeIcon from '../../Icons/Globe';

import styles from './Plan.module.scss';

const proFeatures = [
  {
    id: 'core',
    label: <div className={styles['feature-bold']}>Everything in core</div>
  },
  {
    id: 'host',
    label: <div>Hosting</div>
  },
  {
    id: 'auto',
    label: <div>Automatic update and migration</div>
  },
  {
    id: 'email-support',
    label: <div>Email support</div>
  }
];

const proPerks = [
  {
    id: 'hosted',
    icon: <ServerIcon width="16" height="16" fill="#245fc5" />,
    value: 'Fully hosted and managed'
  },
  {
    id: 'support',
    icon: <GlobeIcon width="16" height="16" fill="#245fc5" />,
    value: 'Support the Dnote community and development'
  }
];

function ProPlan({ wrapperClassName, ctaContent }) {
  return (
    <Plan
      name="Pro"
      price="$3"
      interval="month"
      features={proFeatures}
      perks={proPerks}
      ctaContent={ctaContent}
      wrapperClassName={wrapperClassName}
    />
  );
}

export default ProPlan;

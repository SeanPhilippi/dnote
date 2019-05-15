import React from 'react';

import BoxIcon from '../../Icons/Box';
import Plan from './internal';

const selfHostedPerks = [
  {
    id: 'own-machine',
    icon: <BoxIcon width="16" height="16" fill="#6e6e6e" />,
    value: 'Host on your own machine'
  }
];

const baseFeatures = [
  {
    id: 'backup',
    label: <div>Encrypted backup using AES256</div>
  },
  {
    id: 'sync',
    label: <div>Multi-device sync</div>
  },
  {
    id: 'cli',
    label: <div>Command line interface</div>
  },
  {
    id: 'atom',
    label: <div>Atom plugin</div>
  },
  {
    id: 'web',
    label: <div>Web client</div>
  },
  {
    id: 'digest',
    label: <div>Automated email digest</div>
  },
  {
    id: 'ext',
    label: <div>Firefox/Chrome extension</div>
  },
  {
    id: 'foss',
    label: <div>Free and open source</div>
  },
  {
    id: 'forum-support',
    label: <div>Forum support</div>
  }
];

function Core({ wrapperClassName, ctaContent }) {
  return (
    <Plan
      name="Core"
      price="Free"
      features={baseFeatures}
      perks={selfHostedPerks}
      wrapperClassName={wrapperClassName}
      ctaContent={ctaContent}
    />
  );
}

export default Core;

import React from 'react';
import classnames from 'classnames';

import Button from '../../Common/Button';

import styles from './Sidebar.module.scss';

function Sidebar({ isReady, transacting }) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.header}>
        <div className={styles['plan-name']}>Pro</div>

        <div className={styles['price-wrapper']}>
          <strong className={styles.price}>$3</strong>
          <div className={styles.interval}>/ month</div>
        </div>

        <Button
          id="T-unlock-pro-btn"
          type="submit"
          className={classnames(
            'button button-large button-third button-stretch'
          )}
          disabled={transacting}
          isBusy={transacting || !isReady}
        >
          Unlock
        </Button>
      </div>
    </div>
  );
}

export default Sidebar;

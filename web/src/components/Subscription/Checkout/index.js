import React from 'react';
import classnames from 'classnames';

import * as paymentService from '../../../services/payment';
import { useScript } from '../../../libs/hooks';
import ProPlan from '../Plan/Pro';

import styles from './Checkout.module.scss';

function Checkout({}) {
  const [stripeLoaded, stripeLoadError] = useScript(
    'https://checkout.stripe.com/checkout.js'
  );

  return (
    <div className={classnames('container', styles.wrapper)}>
      <div className="row">
        <div className="col-12 col-md-8">
          <h1>You are almost there.</h1>
        </div>
        <div className="col-12 col-md-4">
          <ProPlan />
        </div>
      </div>
    </div>
  );
}

export default Checkout;

import React, { useState } from 'react';
import classnames from 'classnames';
import { injectStripe, CardElement } from 'react-stripe-elements';

import Sidebar from './Sidebar';
import Flash from '../../Common/Flash';
import * as paymentService from '../../../services/payment';

import styles from './Form.module.scss';

function Form({ isReady, stripe }) {
  const [nameOnCard, setNameOnCard] = useState('');
  const [transacting, setTransacting] = useState(false);
  const [errMessage, setErrMessage] = useState('');

  async function handleSubmit(e) {
    e.preventDefault();

    if (!isReady) {
      return;
    }

    setTransacting(true);

    try {
      const source = await stripe.createSource({
        type: 'card',
        owner: {
          name: nameOnCard
        }
      });

      console.log('source', source);

      await paymentService.createSubscription({ source });
    } catch (err) {
      console.log('error subscribing', err);
      setErrMessage(err.message);
    }

    setTransacting(false);
  }

  return (
    <form
      className={classnames('container-narrow', styles.wrapper)}
      onSubmit={handleSubmit}
    >
      {errMessage && (
        <Flash type="danger">Failed to subscribe. Error: {errMessage}</Flash>
      )}

      <div className="row">
        <div className="col-12 col-md-8">
          <h1>You are almost there.</h1>

          <label htmlFor="name">
            <input
              type="text"
              value={nameOnCard}
              onChange={e => {
                const val = e.target.value;
                setNameOnCard(val);
              }}
            />
          </label>

          {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
          <label htmlFor="card" className={styles.card}>
            Card detail
            <CardElement
              id="card"
              onReady={el => {
                el.focus();
              }}
            />
          </label>
        </div>

        <div className="col-12 col-md-4">
          <Sidebar isReady={isReady} transacting={transacting} />
        </div>
      </div>
    </form>
  );
}

export default injectStripe(Form);

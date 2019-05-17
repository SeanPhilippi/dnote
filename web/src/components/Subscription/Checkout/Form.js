import React, { useState } from 'react';
import classnames from 'classnames';
import Helmet from 'react-helmet';
import { injectStripe, CardElement } from 'react-stripe-elements';

import Sidebar from './Sidebar';
import CountrySelect from './CountrySelect';
import Flash from '../../Common/Flash';
import * as paymentService from '../../../services/payment';

import styles from './Form.module.scss';

const elementStyles = {
  base: {
    color: '#32325D',
    fontFamily: 'Source Code Pro, Consolas, Menlo, monospace',
    fontSize: '16px',
    fontSmoothing: 'antialiased',

    '::placeholder': {
      color: '#CFD7DF'
    },
    ':-webkit-autofill': {
      color: '#e39f48'
    }
  },
  invalid: {
    color: '#E25950',

    '::placeholder': {
      color: '#FFCCA5'
    }
  }
};

function Form({ isReady, stripe, stripeLoadError }) {
  const [nameOnCard, setNameOnCard] = useState('');
  const [cardElementFocused, setCardElementFocused] = useState(false);
  const [billingCountry, setBillingCountry] = useState('');
  const [transacting, setTransacting] = useState(false);
  const [errMessage, setErrMessage] = useState('');

  async function handleSubmit(e) {
    e.preventDefault();

    if (!isReady) {
      return;
    }

    setTransacting(true);

    try {
      const { source } = await stripe.createSource({
        type: 'card',
        owner: {
          name: nameOnCard
        }
      });

      await paymentService.createSubscription({
        source,
        country: billingCountry
      });
    } catch (err) {
      console.log('error subscribing', err);
      setErrMessage(err.message);
    }

    setTransacting(false);
  }

  return (
    <form
      className={classnames('container', styles.wrapper)}
      onSubmit={handleSubmit}
    >
      <Helmet>
        <title>Subscriptions</title>
      </Helmet>

      {errMessage && (
        <Flash type="danger" className={styles.flash}>
          Failed to subscribe. Error: {errMessage}
        </Flash>
      )}
      {stripeLoadError && (
        <Flash type="danger" className={styles.flash}>
          Failed to load stripe. {stripeLoadError}
        </Flash>
      )}

      <div className="row">
        <div className="col-12 col-lg-7 col-xl-8">
          <div className={styles['content-wrapper']}>
            <h1 className={styles.heading}>You are almost there.</h1>

            <div className={styles.content}>
              <div className={styles['input-row']}>
                <label htmlFor="name" className="label-full">
                  <span className={styles.label}>Name on Card</span>
                  <input
                    autoFocus
                    id="name"
                    className={classnames(
                      'text-input text-input-stretch text-input-medium',
                      styles.input
                    )}
                    type="text"
                    value={nameOnCard}
                    onChange={e => {
                      const val = e.target.value;
                      setNameOnCard(val);
                    }}
                  />
                </label>
              </div>

              <div
                className={classnames(styles['input-row'], styles['card-row'])}
              >
                {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
                <label htmlFor="card-number" className={styles.number}>
                  <span className={styles.label}>Card Number</span>

                  <CardElement
                    id="card"
                    className={classnames(styles['card-number'], styles.input, {
                      [styles['card-number-active']]: cardElementFocused
                    })}
                    onFocus={() => {
                      setCardElementFocused(true);
                    }}
                    onBlur={() => {
                      setCardElementFocused(false);
                    }}
                    style={elementStyles}
                  />
                </label>
              </div>

              <div className={styles['input-row']}>
                {/* eslint-disable-next-line jsx-a11y/label-has-associated-control */}
                <label htmlFor="billing-country" className="label-full">
                  <span className={styles.label}>Country</span>
                  <CountrySelect
                    id="billing-country"
                    className={classnames(
                      styles['countries-select'],
                      styles.input
                    )}
                    value={billingCountry}
                    onChange={e => {
                      const val = e.target.value;
                      setBillingCountry(val);
                    }}
                  />
                </label>
              </div>
            </div>
          </div>
        </div>

        <div className="col-12 col-lg-5 col-xl-4">
          <Sidebar isReady={isReady} transacting={transacting} />
        </div>
      </div>
    </form>
  );
}

export default injectStripe(Form);

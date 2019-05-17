import React from 'react';
import { StripeProvider, Elements } from 'react-stripe-elements';

import { useScript } from '../../../libs/hooks';
import CheckoutForm from './Form';

function Checkout() {
  const [stripeLoaded, stripeLoadError] = useScript('https://js.stripe.com/v3');

  let key;
  if (__PRODUCTION__) {
    key = 'pk_live_xvouPZFPDDBSIyMUSLZwkXfR';
  } else {
    key = 'pk_test_5926f65DQoIilZeNOiKydfoN';
  }

  let stripe = null;
  if (stripeLoaded) {
    stripe = window.Stripe(key);
  }

  return (
    <StripeProvider stripe={stripe}>
      <Elements>
        <CheckoutForm stripeLoadError={stripeLoadError} />
      </Elements>
    </StripeProvider>
  );
}

export default Checkout;

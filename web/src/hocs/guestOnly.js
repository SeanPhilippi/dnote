/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

import React from 'react';
import { connect } from 'react-redux';
import { Redirect } from 'react-router-dom';
import hoistNonReactStatics from 'hoist-non-react-statics';

import { isEmptyObj } from '../libs/obj';
import { getReferrer } from '../libs/url';

function renderFallback(referrer) {
  let destination;
  if (referrer) {
    destination = referrer;
  } else {
    destination = '/';
  }

  return <Redirect to={{ pathname: destination }} />;
}

// guestOnly returns a HOC that renders the given component only if user is not
// logged in
export default function(Component) {
  const HOC = props => {
    const { userData, location } = props;

    const loggedIn = Boolean(userData.isFetched && !isEmptyObj(userData.data));

    if (loggedIn) {
      const referrer = getReferrer(location);
      return renderFallback(referrer);
    }

    return <Component {...props} />;
  };

  // Copy over static methods
  hoistNonReactStatics(HOC, Component);

  function mapStateToProps(state) {
    return {
      userData: state.auth.user
    };
  }

  return connect(mapStateToProps)(HOC);
}

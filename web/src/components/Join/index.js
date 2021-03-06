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

import React, { useState } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';

import JoinForm from './JoinForm';
import Logo from '../Icons/Logo';
import Flash from '../Common/Flash';

import { getReferrer } from '../../libs/url';
import { updateAuthEmail } from '../../actions/form';
import { register } from '../../services/users';
import { getCurrentUser } from '../../actions/auth';
import { registerHelper } from '../../crypto';
import { DEFAULT_KDF_ITERATION } from '../../crypto/consts';

import authStyles from '../Common/Auth.module.scss';

function Join({ doGetCurrentUser, formData, doUpdateAuthEmail, location }) {
  const [errMsg, setErrMsg] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const referrer = getReferrer(location);

  async function handleJoin(email, password, passwordConfirmation) {
    if (!email) {
      setErrMsg('Please enter an email.');
      return;
    }
    if (!password) {
      setErrMsg('Please enter a pasasword.');
      return;
    }
    if (!passwordConfirmation) {
      setErrMsg('The passwords do not match.');
      return;
    }

    setErrMsg('');
    setSubmitting(true);

    try {
      const { cipherKey, cipherKeyEnc, authKey } = await registerHelper({
        email,
        password,
        iteration: DEFAULT_KDF_ITERATION
      });
      await register({
        email,
        authKey,
        cipherKeyEnc,
        iteration: DEFAULT_KDF_ITERATION
      });
      localStorage.setItem('cipherKey', cipherKey);

      // guestOnly HOC will redirect the user accordingly after the current user is fetched
      await doGetCurrentUser();
      doUpdateAuthEmail('');
    } catch (err) {
      console.log(err);
      setErrMsg(err.message);
      setSubmitting(false);
    }
  }

  return (
    <div className={authStyles.page}>
      <Helmet>
        <title>Join</title>
      </Helmet>
      <div className="container">
        <a href="/">
          <Logo fill="#252833" width="60" height="60" />
        </a>
        <h1 className={authStyles.heading}>Join Dnote</h1>

        <div className={authStyles.body}>
          {referrer && (
            <Flash type="info" wrapperClassName={authStyles['referrer-flash']}>
              Please join to continue.
            </Flash>
          )}

          <div className={authStyles.panel}>
            {errMsg && (
              <Flash type="danger" wrapperClassName={authStyles['error-flash']}>
                {errMsg}
              </Flash>
            )}

            <JoinForm
              onJoin={handleJoin}
              submitting={submitting}
              email={formData.auth.email}
              onUpdateEmail={doUpdateAuthEmail}
            />
          </div>

          <div className={authStyles.footer}>
            <div className={authStyles.callout}>Already have an account?</div>
            <Link to="/login" className={authStyles.cta}>
              Sign in
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

function mapStateToProps(state) {
  return {
    formData: state.form
  };
}

const mapDispatchToProps = {
  doGetCurrentUser: getCurrentUser,
  doUpdateAuthEmail: updateAuthEmail
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Join);

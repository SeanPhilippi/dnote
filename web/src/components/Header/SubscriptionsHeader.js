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
import { Link } from 'react-router-dom';

import Logo from '../Icons/LogoWithText';
import { homePath } from '../../libs/paths';
import styles from './SubscriptionsHeader.module.scss';

function SubscriptionsHeader({ userData }) {
  const user = userData.data;

  return (
    <header className={styles.wrapper}>
      <div className={styles.content}>
        <Link to={homePath({})} className={styles.brand}>
          <Logo width={32} height={32} fill="black" className={styles.logo} />
        </Link>

        <div className={styles.email}>{user.email}</div>
      </div>
    </header>
  );
}

function mapStateToProps(state) {
  return {
    userData: state.auth.user
  };
}

export default connect(mapStateToProps)(SubscriptionsHeader);

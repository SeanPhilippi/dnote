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
import classnames from 'classnames';
import { Link } from 'react-router-dom';

import LockIcon from '../Icons/Lock';
import { subscriptionsPath, homePath } from '../../libs/paths';

import styles from './SubscriberWall.module.scss';

function SubscriberWall({ wrapperClassName }) {
  return (
    <div className={classnames(styles.wrapper, wrapperClassName)}>
      <LockIcon width="64" height="64" />

      <div className={styles.lead}>Unlock Dnote Pro to get started.</div>

      <div className={styles.actions}>
        <Link to={subscriptionsPath()} className="button button-first">
          Get started
        </Link>
        <Link
          to={homePath({}, { demo: true })}
          className="button button-first-outline "
        >
          Live demo
        </Link>
      </div>
    </div>
  );
}

export default SubscriberWall;

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
import classnames from 'classnames';
import { Link } from 'react-router-dom';

import { homePath } from '../../../libs/paths';
import Actions from './Actions';
import MobileActions from './MobileActions';

import styles from './BookItem.module.scss';

export default ({ book, demo, isFocused, setFocusedOptEl, onDeleteBook }) => {
  const [isHovered, setIsHovered] = useState(false);
  const isActive = isFocused || isHovered;

  return (
    <li
      className={classnames(styles.item, `book-item-${book.uuid}`, {
        [styles.active]: isActive
      })}
      ref={el => {
        if (isFocused) {
          // eslint-disable-next-line no-param-reassign
          setFocusedOptEl(el);
        }
      }}
      onMouseEnter={() => {
        setIsHovered(true);
      }}
      onMouseLeave={() => {
        setIsHovered(false);
      }}
    >
      <Link
        to={homePath({ book: book.uuid }, { demo })}
        className={styles.link}
        tabIndex="-1"
      >
        <h2 className={styles.label}>{book.label}</h2>
      </Link>

      <MobileActions bookUUID={book.uuid} onDeleteBook={onDeleteBook} />

      <Actions
        bookUUID={book.uuid}
        onDeleteBook={onDeleteBook}
        shown={isActive}
      />
    </li>
  );
};

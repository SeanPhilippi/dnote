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

import { apiClient } from '../libs/http';
import { getPath } from '../libs/url';

export function fetch(digestUUID, { demo }) {
  let endpoint;
  if (demo) {
    endpoint = `/demo/digests/${digestUUID}`;
  } else {
    endpoint = `/digests/${digestUUID}`;
  }

  return apiClient.get(endpoint);
}

export function fetchAll({ page, demo }) {
  let path;
  if (demo) {
    path = `/demo/digests`;
  } else {
    path = '/digests';
  }

  const endpoint = getPath(path, { page });

  return apiClient.get(endpoint);
}

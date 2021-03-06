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

import qs from 'qs';
import { apiClient } from '../libs/http';

export function fetch(queryObj = {}, options = {}) {
  const { demo } = options;

  let baseURL;
  if (demo) {
    baseURL = '/v1/demo/books';
  } else {
    baseURL = '/v1/books';
  }

  const queryStr = qs.stringify(queryObj);

  let endpoint;
  if (queryStr) {
    endpoint = `${baseURL}?${queryStr}`;
  } else {
    endpoint = baseURL;
  }

  return apiClient.get(endpoint);
}

export function create({ name }) {
  const payload = { name };

  return apiClient.post('/v2/books', payload);
}

export function remove(uuid) {
  return apiClient.del(`/v1/books/${uuid}`);
}

export function update(uuid, payload) {
  return apiClient.patch(`/v1/books/${uuid}`, payload);
}

export function get(bookUUID) {
  const endpoint = `/v1/books/${bookUUID}`;

  return apiClient.get(endpoint);
}

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

export function updateUser({ name }) {
  const payload = { name };

  return apiClient.patch('/account/profile', payload);
}

export function updateEmail({
  oldAuthKey,
  newAuthKey,
  newEmail,
  newCipherKeyEnc
}) {
  const payload = {
    old_auth_key: oldAuthKey,
    new_auth_key: newAuthKey,
    new_email: newEmail,
    new_cipher_key_enc: newCipherKeyEnc
  };

  return apiClient.patch('/account/email', payload);
}

export function updatePassword({
  oldAuthKey,
  newAuthKey,
  newCipherKeyEnc,
  newKdfIteration
}) {
  const payload = {
    old_auth_key: oldAuthKey,
    new_auth_key: newAuthKey,
    new_cipher_key_enc: newCipherKeyEnc,
    new_kdf_iteration: newKdfIteration
  };

  return apiClient.patch('/account/password', payload);
}

export function register({ email, authKey, cipherKeyEnc, iteration }) {
  const payload = {
    email,
    auth_key: authKey,
    iteration,
    cipher_key_enc: cipherKeyEnc
  };

  return apiClient.post('/v1/register', payload);
}

export function presignin({ email }) {
  return apiClient.get(`/v1/presignin?email=${email}`);
}

export function signin({ email, authKey }) {
  const payload = { email, auth_key: authKey };

  return apiClient.post('/v1/signin', payload);
}

export function signout() {
  return apiClient.post('/v1/signout');
}

export function sendResetPasswordEmail({ email }) {
  const payload = { email };

  return apiClient.post('/reset-token', payload);
}

export function sendEmailVerificationEmail() {
  return apiClient.post('/verification-token');
}

export function verifyEmail({ token }) {
  const payload = { token };

  return apiClient.patch('/verify-email', payload);
}

export function updateEmailPreference({ token, digestFrequency }) {
  const payload = { digest_weekly: digestFrequency === 'weekly' };

  let endpoint = '/account/email-preference';
  if (token) {
    endpoint = `${endpoint}?token=${token}`;
  }
  return apiClient.patch(endpoint, payload);
}

export function getEmailPreference({ token }) {
  let endpoint = '/account/email-preference';
  if (token) {
    endpoint = `${endpoint}?token=${token}`;
  }

  return apiClient.get(endpoint);
}

export function getMe() {
  return apiClient.get('/me').then(res => {
    return res.user;
  });
}

export function fetchCalendar(options = {}) {
  const { demo } = options;

  let endpoint = '/calendar';
  if (demo) {
    endpoint = `/demo${endpoint}`;
  }

  return apiClient.get(endpoint);
}

export function legacySignin({ email, password }) {
  const payload = { email, password };

  return apiClient.post('/legacy/signin', payload);
}

export function legacyGetMe() {
  return apiClient.get('/legacy/me').then(res => {
    return res.user;
  });
}

export function legacyRegister({ email, authKey, cipherKeyEnc, iteration }) {
  const payload = {
    email,
    auth_key: authKey,
    iteration,
    cipher_key_enc: cipherKeyEnc
  };

  return apiClient.post('/legacy/register', payload);
}

export function legacyMigrate() {
  return apiClient.patch('/legacy/migrate', {});
}

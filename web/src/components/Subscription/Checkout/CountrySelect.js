import React from 'react';
import { countries } from '../../../libs/countries';

function CountrySelect({ id, className, onChange, value }) {
  return (
    <select id={id} className={className} value={value} onChange={onChange}>
      <option value="" />

      {countries.map(country => {
        return (
          <option key={country.code} value={country.code}>
            {country.name}
          </option>
        );
      })}
    </select>
  );
}

export default CountrySelect;

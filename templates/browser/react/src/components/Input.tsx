import React from 'react';

interface Props extends HTMLInputElement {
  label: string;
}

export function Input({
  label,
  id,
  onChange: handleChange,
  type,
  value,
}: Props) {
  return (
    <label
      className="block text-gray-700 text-sm font-bold mb-2"
      htmlFor={id}
    >
      {label}
      <input
        type={type}
        id={id}
        value={value}
        onChange={handleChange}
      />
    </label>
  );
}

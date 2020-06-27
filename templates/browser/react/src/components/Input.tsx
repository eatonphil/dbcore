import React from 'react';

interface Props extends React.HTMLProps<HTMLInputElement> {
  label: string;
}

export function Input({
  disabled,
  label,
  id,
  onChange: handleChange,
  type,
  value,
}: Props) {
  return (
    <label
      className="block text-gray-700 text-sm font-bold mb-2 uppercase text-sm"
      htmlFor={id}
    >
      {label}
      <input
        disabled={disabled}
        className="border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight"
        type={type}
        id={id}
        value={value}
        onChange={handleChange}
      />
    </label>
  );
}

import React from 'react';

export function Form({
  children,
  onSubmit,
}: HTMLFormElement) {
  return (
    <form
      className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4"
      onSubmit={onSubmit}
      children={children}
    />
  );
}

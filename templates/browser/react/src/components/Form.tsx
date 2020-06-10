import React from 'react';

import { Button } from '../components/Button';

interface Props extends HTMLFormElement {
  buttonText: string;
  error?: string;
}

export function Form({
  buttonText,
  children,
  error,
  onSubmit,
}: Props) {
  return (
    <form
      className="border px-8 pt-6 pb-8 mb-4"
      onSubmit={onSubmit}
    >
      {children}
      <Button type="submit">{buttonText}</Button>
      {error && <div className="text-red-600 text-sm mt-4">{error}</div>}
    </form>
  );
}

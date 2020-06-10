import React from 'react';

interface Props {
  actions?: React.ReactNode;
  size: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  children: React.ReactNode;
}

function H(size) {
  
}

export function Heading({
  actions,
  size,
  children,
}: Props) {
  const h = (() => {
    switch (size) {
      case 'xs':
        return <h6 className="text-xs font-semibold">{children}</h6>;
      case 'sm':
        return <h5 className="text-sm font-semibold">{children}</h5>;
      case 'md':
        return <h4 className="text-md font-semibold">{children}</h4>;
      case 'lg':
        return <h3 className="text-lg font-semibold">{children}</h3>;
      case 'xl':
        return <h2 className="text-xl font-semibold mb-4">{children}</h2>;
    }
  })();

  return (
    <div className="flex justify-between items-center">
      {h}
      {actions}
    </div>
  )
}

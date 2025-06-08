import React from 'react';

// Placeholder component to resolve import errors
interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'success' | 'warning' | 'outline' | 'destructive';
}

export const Badge: React.FC<BadgeProps> = ({ children, variant='default', className='', ...props }) => {
  const base = 'px-2 py-1 text-xs rounded';
  const styles: Record<string,string> = {
    default: 'bg-gray-200 text-gray-800',
    success: 'bg-green-100 text-green-800',
    warning: 'bg-yellow-100 text-yellow-800',
    outline: 'border border-gray-300 text-gray-800',
    destructive: 'bg-red-100 text-red-800',
  };
  return (
    <span className={`${base} ${styles[variant]} ${className}`} {...props}>
      {children}
    </span>
  );
};

export default Badge;

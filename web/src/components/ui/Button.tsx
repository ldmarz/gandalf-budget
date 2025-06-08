import React from 'react';
import { buttonClasses, destructiveButtonClasses, secondaryButtonClasses, disabledClasses } from '../../styles/commonClasses';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
  onClick?: () => void;
  variant?: 'primary' | 'destructive' | 'secondary' | 'ghost';
  size?: 'sm' | 'md';
}

export const Button: React.FC<ButtonProps> = ({
  children,
  onClick,
  variant = 'primary',
  size = 'md',
  disabled = false,
  className = '',
  type = 'button',
  ...props // Pass down any other native button props
}) => {
  let baseClasses: string;

  switch (variant) {
    case 'destructive':
      baseClasses = destructiveButtonClasses;
      break;
    case 'secondary':
      baseClasses = secondaryButtonClasses;
      break;
    case 'ghost':
      baseClasses = 'text-gray-700 hover:bg-gray-100';
      break;
    case 'primary':
    default:
      baseClasses = buttonClasses;
      break;
  }

  const sizeClasses = size === 'sm' ? 'px-2 py-1 text-sm' : 'px-3 py-2';

  const combinedClasses = `
    ${baseClasses}
    ${sizeClasses}
    ${disabled ? disabledClasses : ''}
    ${className}
  `.replace(/\s+/g, ' ').trim(); // Replace multiple spaces with a single space and trim

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={combinedClasses}
      {...props}
    >
      {children}
    </button>
  );
};

export default Button;

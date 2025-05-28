import React from 'react';
import { buttonClasses, destructiveButtonClasses, secondaryButtonClasses, disabledClasses } from '../../styles/commonClasses';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
  onClick?: () => void;
  variant?: 'primary' | 'destructive' | 'secondary';
  // disabled is already part of ButtonHTMLAttributes
  // type is already part of ButtonHTMLAttributes
}

export const Button: React.FC<ButtonProps> = ({
  children,
  onClick,
  variant = 'primary',
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
    case 'primary':
    default:
      baseClasses = buttonClasses;
      break;
  }

  const combinedClasses = `
    ${baseClasses}
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

import React, { ChangeEvent, FocusEvent } from 'react';
import { inputClasses } from '../../styles/commonClasses';
import { disabledClasses } from '../../styles/commonClasses'; // Also import disabledClasses

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  type?: string;
  value?: string | number;
  onChange?: (event: ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (event: FocusEvent<HTMLInputElement>) => void;
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  className?: string;
  id?: string;
  min?: string | number;
  step?: string | number;
  defaultValue?: string | number;
}

export const Input: React.FC<InputProps> = ({
  type = 'text',
  value,
  onChange,
  onBlur,
  placeholder,
  required,
  disabled = false,
  className = '',
  id,
  min,
  step,
  defaultValue,
  ...props // Pass down any other native input props
}) => {
  const combinedClasses = `
    ${inputClasses}
    ${disabled ? disabledClasses : ''}
    ${className}
  `.replace(/\s+/g, ' ').trim();

  return (
    <input
      type={type}
      id={id}
      value={value}
      onChange={onChange}
      onBlur={onBlur}
      placeholder={placeholder}
      required={required}
      disabled={disabled}
      className={combinedClasses}
      min={min}
      step={step}
      defaultValue={defaultValue}
      {...props}
    />
  );
};

export default Input;

import React, { ChangeEvent } from 'react';
import { inputClasses, disabledClasses } from '../../styles/commonClasses';

interface Option {
  value: string | number;
  label: string;
  disabled?: boolean;
}

interface SelectProps {
  value?: string | number;
  onChange?: (event: ChangeEvent<HTMLSelectElement>) => void;
  options: Option[];
  required?: boolean;
  disabled?: boolean;
  className?: string;
  id?: string;
}

export const Select: React.FC<SelectProps> = ({
  value,
  onChange,
  options,
  required,
  disabled = false,
  className = '',
  id,
  ...props // Pass down any other native select props though less common for select
}) => {
  const combinedClasses = `
    ${inputClasses}
    ${disabled ? disabledClasses : ''}
    ${className}
  `.replace(/\s+/g, ' ').trim();

  return (
    <select
      id={id}
      value={value}
      onChange={onChange}
      required={required}
      disabled={disabled}
      className={combinedClasses}
      {...props}
    >
      {options.map((option) => (
        <option key={option.value} value={option.value} disabled={option.disabled}>
          {option.label}
        </option>
      ))}
    </select>
  );
};

export default Select;

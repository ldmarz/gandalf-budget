import React from 'react';

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg';
  text?: string;
  className?: string; // Allow for additional custom styling
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  text,
  className = '',
}) => {
  let sizeClasses = 'h-8 w-8'; // Default to md
  let borderClasses = 'border-t-2 border-b-2'; // Default border thickness

  switch (size) {
    case 'sm':
      sizeClasses = 'h-5 w-5';
      borderClasses = 'border-t border-b'; // Thinner border for smaller size
      break;
    case 'lg':
      sizeClasses = 'h-12 w-12';
      borderClasses = 'border-t-4 border-b-4'; // Thicker border for larger size
      break;
    case 'md':
    default:
      // Default classes are already set
      break;
  }

  const combinedClasses = `
    flex flex-col justify-center items-center 
    ${className}
  `.replace(/\s+/g, ' ').trim();

  return (
    <div className={combinedClasses}>
      <div
        className={`animate-spin rounded-full ${sizeClasses} ${borderClasses} border-blue-500`}
      ></div>
      {text && <span className="mt-2 text-sm text-gray-300">{text}</span>}
    </div>
  );
};

export default LoadingSpinner;

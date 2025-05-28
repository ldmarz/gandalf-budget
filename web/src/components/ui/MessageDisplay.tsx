import React from 'react';

interface MessageDisplayProps {
  message: React.ReactNode | string | null | undefined;
  type?: 'error' | 'success' | 'info';
  className?: string;
}

export const MessageDisplay: React.FC<MessageDisplayProps> = ({
  message,
  type = 'info',
  className = '',
}) => {
  if (!message) {
    return null;
  }

  let typeClasses = '';
  switch (type) {
    case 'error':
      typeClasses = 'bg-red-800 border border-red-700 text-white';
      break;
    case 'success':
      typeClasses = 'bg-green-800 border border-green-700 text-white';
      break;
    case 'info':
    default:
      typeClasses = 'bg-blue-800 border border-blue-700 text-white';
      break;
  }

  const combinedClasses = `
    px-4 py-3 rounded relative 
    ${typeClasses} 
    ${className}
  `.replace(/\s+/g, ' ').trim();

  return (
    <div className={combinedClasses} role="alert">
      {message}
    </div>
  );
};

export default MessageDisplay;

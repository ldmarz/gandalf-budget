import React from 'react';
import { cardClasses } from '../../styles/commonClasses';

interface CardProps {
  children: React.ReactNode;
  className?: string;
}

export const Card: React.FC<CardProps> = ({ children, className }) => {
  const combinedClasses = `${cardClasses} ${className || ''}`.trim();
  return (
    <div className={combinedClasses}>
      {children}
    </div>
  );
};

export default Card;

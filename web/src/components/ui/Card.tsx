import React from 'react';
import { cardClasses } from '../../styles/commonClasses'; // Assuming this path is correct and file exists

interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  className?: string;
}

export const Card: React.FC<CardProps> = ({ children, className, ...props }) => {
  const combinedClasses = `${cardClasses} ${className || ''}`.trim();
  return (
    <div className={combinedClasses} {...props}>
      {children}
    </div>
  );
};

// Placeholder exports for Card sub-components
export const CardHeader: React.FC<CardProps> = ({ children, className, ...props }) => (
  <div className={`card-header ${className || ''}`.trim()} {...props}>
    {children}
  </div>
);

export const CardTitle: React.FC<CardProps> = ({ children, className, ...props }) => (
  <h3 className={`card-title ${className || ''}`.trim()} {...props}>
    {children}
  </h3>
);

export const CardDescription: React.FC<CardProps> = ({ children, className, ...props }) => (
  <p className={`card-description ${className || ''}`.trim()} {...props}>
    {children}
  </p>
);

export const CardContent: React.FC<CardProps> = ({ children, className, ...props }) => (
  <div className={`card-content ${className || ''}`.trim()} {...props}>
    {children}
  </div>
);

export default Card;

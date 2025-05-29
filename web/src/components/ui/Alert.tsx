import React from 'react';

// Placeholder components to resolve import errors
export const Alert: React.FC<any> = ({ children, ...props }) => <div {...props}>{children}</div>;
export const AlertDescription: React.FC<any> = ({ children, ...props }) => <p {...props}>{children}</p>;
export const AlertTitle: React.FC<any> = ({ children, ...props }) => <h5 {...props}>{children}</h5>;

export default Alert;

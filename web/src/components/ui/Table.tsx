import React from 'react';

// Placeholder components to resolve import errors
export const Table: React.FC<any> = ({ children, ...props }) => <table {...props}>{children}</table>;
export const TableBody: React.FC<any> = ({ children, ...props }) => <tbody {...props}>{children}</tbody>;
export const TableCell: React.FC<any> = ({ children, ...props }) => <td {...props}>{children}</td>;
export const TableHead: React.FC<any> = ({ children, ...props }) => <th {...props}>{children}</th>;
export const TableHeader: React.FC<any> = ({ children, ...props }) => <thead {...props}>{children}</thead>;
export const TableRow: React.FC<any> = ({ children, ...props }) => <tr {...props}>{children}</tr>;

export default Table;

import React from 'react';

interface Props extends React.ImgHTMLAttributes<HTMLImageElement> {
  src: string;
}

export const Avatar: React.FC<Props> = ({ src, className='', ...props }) => (
  <img src={src} className={`rounded-full ${className}`} {...props} />
);

export default Avatar;

import React from 'react';

interface Category {
  name: string;
  color: string;
}

interface CategoryBadgeProps {
  category: Category | undefined | null;
  className?: string;
}

const defaultColor = 'bg-gray-500';

export const CategoryBadge: React.FC<CategoryBadgeProps> = ({ category, className = '' }) => {
  const categoryName = category?.name || 'Unknown Category';
  const categoryColor = category?.color || defaultColor;

  const combinedClasses = `flex items-center ${className}`.trim();

  return (
    <div className={combinedClasses}>
      <span className={`inline-block w-4 h-4 rounded-full mr-2 ${categoryColor}`}></span>
      <span>{categoryName}</span>
    </div>
  );
};

export default CategoryBadge;

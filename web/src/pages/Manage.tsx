import React, { useState, useEffect, FormEvent } from 'react';
import { get, post, put, del } from '../lib/api'; // Adjust path as necessary

interface Category {
  id: number;
  name: string;
  color: string; // Tailwind color class e.g., 'bg-red-500'
}

// Basic Tailwind classes - can be expanded
const inputClasses = "border border-gray-300 rounded px-2 py-1 text-black"; // Added text-black
const buttonClasses = "bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded";
const destructiveButtonClasses = "bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-2 rounded";
const cardClasses = "bg-gray-800 p-4 rounded shadow-md mb-4"; // Changed to dark theme card
const textMutedClasses = "text-gray-400"; // For placeholder text or muted info

export default function ManagePage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Form state for adding/editing
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [currentCategory, setCurrentCategory] = useState<Category | null>(null);
  const [formName, setFormName] = useState<string>('');
  const [formColor, setFormColor] = useState<string>(''); // e.g., 'bg-blue-500'

  const tailwindColors = [
    'bg-red-500', 'bg-orange-500', 'bg-amber-500', 'bg-yellow-500', 'bg-lime-500',
    'bg-green-500', 'bg-emerald-500', 'bg-teal-500', 'bg-cyan-500', 'bg-sky-500',
    'bg-blue-500', 'bg-indigo-500', 'bg-violet-500', 'bg-purple-500', 'bg-fuchsia-500',
    'bg-pink-500', 'bg-rose-500', 'bg-slate-500', 'bg-gray-500',
  ];


  // Fetch categories
  const fetchCategories = async () => {
    setIsLoading(true);
    try {
      const data = await get<Category[]>('/categories');
      setCategories(data || []); // Handle null response from API if no categories
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch categories');
      setCategories([]); // Clear categories on error
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchCategories();
  }, []);

  // Form handling
  const handleFormSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!formName.trim() || !formColor.trim()) {
      alert('Name and color are required.');
      return;
    }

    const categoryData = { name: formName, color: formColor };

    try {
      if (isEditing && currentCategory) {
        await put<Category, Category>(`/categories/${currentCategory.id}`, { ...categoryData, id: currentCategory.id });
      } else {
        await post<Category, Omit<Category, 'id'>>('/categories', categoryData);
      }
      await fetchCategories(); // Re-fetch all categories
      resetForm();
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : (isEditing ? 'Failed to update category' : 'Failed to create category');
      alert(errorMsg); // Simple alert for now
      setError(errorMsg);
    }
  };

  const handleEdit = (category: Category) => {
    setIsEditing(true);
    setCurrentCategory(category);
    setFormName(category.name);
    setFormColor(category.color);
    window.scrollTo(0, 0); // Scroll to top to see the form
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this category?')) {
      try {
        await del<null>(`/categories/${id}`);
        await fetchCategories(); // Re-fetch
      } catch (err) {
        const errorMsg = err instanceof Error ? err.message : 'Failed to delete category';
        alert(errorMsg);
        setError(errorMsg);
      }
    }
  };

  const resetForm = () => {
    setIsEditing(false);
    setCurrentCategory(null);
    setFormName('');
    setFormColor('');
  };

  if (isLoading) return <div className="p-4 text-white">Loading categories...</div>; // Added text-white
  // Removed redundant error display here, will display below list

  return (
    <div className="p-4 bg-gray-900 min-h-screen text-white"> {/* Changed to dark theme */}
      <h1 className="text-2xl font-bold mb-6 text-center">Manage Categories</h1>

      {/* Add/Edit Form Card */}
      <div className={cardClasses}>
        <h2 className="text-xl font-semibold mb-3">{isEditing ? 'Edit Category' : 'Add New Category'}</h2>
        <form onSubmit={handleFormSubmit} className="space-y-3">
          <div>
            <label htmlFor="name" className="block text-sm font-medium mb-1">Name:</label>
            <input
              id="name"
              type="text"
              value={formName}
              onChange={(e) => setFormName(e.target.value)}
              className={`${inputClasses} w-full`}
              placeholder="e.g., Groceries"
              required
            />
          </div>
          <div>
            <label htmlFor="color" className="block text-sm font-medium mb-1">Color (Tailwind Class):</label>
            <select
              id="color"
              value={formColor}
              onChange={(e) => setFormColor(e.target.value)}
              className={`${inputClasses} w-full`}
              required
            >
              <option value="" disabled className={textMutedClasses}>Select a color</option>
              {tailwindColors.map(colorClass => (
                <option key={colorClass} value={colorClass}>
                  {colorClass}
                </option>
              ))}
            </select>
            {formColor && (
              <div className="mt-2 p-2 flex items-center">
                <span className={`inline-block w-6 h-6 rounded mr-2 ${formColor}`}></span>
                <span>Selected color preview</span>
              </div>
            )}
          </div>
          <div className="flex space-x-2">
            <button type="submit" className={buttonClasses}>
              {isEditing ? 'Update Category' : 'Add Category'}
            </button>
            {isEditing && (
              <button type="button" onClick={resetForm} className={`${buttonClasses} bg-gray-500 hover:bg-gray-600`}>
                Cancel Edit
              </button>
            )}
          </div>
        </form>
      </div>
      
      {error && (
        <div className="mt-4 p-3 bg-red-800 border border-red-700 text-white rounded text-center"> {/* Darker error for dark theme */}
          <p>Error: {error}</p>
        </div>
      )}

      {/* Categories List Card */}
      <div className={`${cardClasses} mt-8`}>
        <h2 className="text-xl font-semibold mb-4">Existing Categories</h2>
        {categories.length === 0 && !isLoading && !error && (
          <p className={textMutedClasses}>No categories found. Add some using the form above.</p>
        )}
        {categories.length > 0 && (
          <ul className="space-y-3">
            {categories.map(cat => (
              <li key={cat.id} className="flex items-center justify-between p-3 bg-gray-700 rounded"> {/* Darker list items */}
                <div className="flex items-center">
                  <span className={`inline-block w-5 h-5 rounded mr-3 ${cat.color}`}></span>
                  <span className="font-medium">{cat.name}</span>
                </div>
                <div className="space-x-2">
                  <button onClick={() => handleEdit(cat)} className={`${buttonClasses} bg-yellow-500 hover:bg-yellow-600`}>
                    Edit
                  </button>
                  <button onClick={() => handleDelete(cat.id)} className={destructiveButtonClasses}>
                    Delete
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}

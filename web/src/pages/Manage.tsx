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

  // --- Budget Line State ---
  const [budgetLines, setBudgetLines] = useState<api.BudgetLine[]>([]);
  const [currentMonthId, setCurrentMonthId] = useState<number>(1); // Default to month 1 for now
  const [isLoadingBL, setIsLoadingBL] = useState(true);
  const [errorBL, setErrorBL] = useState<string | null>(null);

  // Form state for budget lines
  const [isEditingBL, setIsEditingBL] = useState<boolean>(false);
  const [currentBL, setCurrentBL] = useState<api.BudgetLine | null>(null);
  const [formBLLabel, setFormBLLabel] = useState<string>('');
  const [formBLExpected, setFormBLExpected] = useState<string>(''); // Store as string for form input
  const [formBLCategoryID, setFormBLCategoryID] = useState<string>(''); // Store as string for form input

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
    // Fetch budget lines for the default/current month
    if (currentMonthId) {
      fetchBudgetLines(currentMonthId);
    }
  }, [currentMonthId]); // Re-fetch if currentMonthId changes

  // --- Category Functions ---
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

  const resetForm = () => {
    setIsEditing(false);
    setCurrentCategory(null);
    setFormName('');
    setFormColor('');
  };
  
  // --- Budget Line Functions ---
  const fetchBudgetLines = async (monthId: number) => {
    setIsLoadingBL(true);
    try {
      const data = await api.getBudgetLinesByMonth(monthId);
      // Map category details to budget lines
      const linesWithCategoryDetails = data.map(line => {
        const category = categories.find(c => c.id === line.category_id);
        return {
          ...line,
          category_name: category?.name || 'Unknown',
          category_color: category?.color || 'bg-gray-500',
        };
      });
      setBudgetLines(linesWithCategoryDetails || []);
      setErrorBL(null);
    } catch (err) {
      setErrorBL(err instanceof Error ? err.message : 'Failed to fetch budget lines');
      setBudgetLines([]);
    } finally {
      setIsLoadingBL(false);
    }
  };

  useEffect(() => {
    if (categories.length > 0 && currentMonthId) {
        fetchBudgetLines(currentMonthId);
    }
  }, [categories, currentMonthId]); // Re-fetch budget lines if categories array or monthId changes

  const resetBLForm = () => {
    setIsEditingBL(false);
    setCurrentBL(null);
    setFormBLLabel('');
    setFormBLExpected('');
    setFormBLCategoryID('');
  };

  const handleBLEdit = (bl: api.BudgetLine) => {
    setIsEditingBL(true);
    setCurrentBL(bl);
    setFormBLLabel(bl.label);
    setFormBLExpected(bl.expected.toString());
    setFormBLCategoryID(bl.category_id.toString());
    window.scrollTo(0, document.getElementById('budget-lines-section')?.offsetTop || 0);
  };
  
  const handleBLFormSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!formBLLabel.trim() || !formBLExpected.trim() || !formBLCategoryID.trim() || !currentMonthId) {
      alert('Label, expected amount, category, and a selected month are required.');
      return;
    }

    const expectedAmount = parseFloat(formBLExpected);
    const categoryId = parseInt(formBLCategoryID, 10);

    if (isNaN(expectedAmount) || isNaN(categoryId)) {
      alert('Invalid expected amount or category ID.');
      return;
    }

    const budgetLineData = { 
      month_id: currentMonthId, 
      category_id: categoryId, 
      label: formBLLabel, 
      expected: expectedAmount 
    };

    try {
      if (isEditingBL && currentBL) {
        // For PUT, only send fields that can be updated (label, expected)
        await api.updateBudgetLine(currentBL.id, { label: formBLLabel, expected: expectedAmount });
      } else {
        await api.createBudgetLine(budgetLineData);
      }
      await fetchBudgetLines(currentMonthId); // Re-fetch budget lines for the current month
      resetBLForm();
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : (isEditingBL ? 'Failed to update budget line' : 'Failed to create budget line');
      alert(errorMsg); 
      setErrorBL(errorMsg);
    }
  };

  const handleBLDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this budget line?')) {
      try {
        await api.deleteBudgetLine(id);
        await fetchBudgetLines(currentMonthId); // Re-fetch
      } catch (err) {
        const errorMsg = err instanceof Error ? err.message : 'Failed to delete budget line';
        alert(errorMsg);
        setErrorBL(errorMsg);
      }
    }
  };


  if (isLoading) return <div className="p-4 text-white">Loading page...</div>;

  return (
    <div className="p-4 bg-gray-900 min-h-screen text-white">
      <h1 className="text-2xl font-bold mb-6 text-center">Manage Data</h1>

      {/* Month Selector (Simplified) */}
      <div className={`${cardClasses} mb-6`}>
        <h2 className="text-xl font-semibold mb-3">Select Month</h2>
        <div className="flex items-center space-x-2">
          <label htmlFor="month_id_selector" className="block text-sm font-medium">Month ID:</label>
          <input
            id="month_id_selector"
            type="number"
            value={currentMonthId}
            onChange={(e) => setCurrentMonthId(parseInt(e.target.value, 10) || 1)}
            className={`${inputClasses} w-24`}
            min="1"
          />
          <button onClick={() => fetchBudgetLines(currentMonthId)} className={buttonClasses}>Load Budget Lines</button>
        </div>
        <p className={textMutedClasses}>Note: This is a simplified month selector. Use a valid Month ID from your database.</p>
      </div>
      
      {/* Categories Section */}
      <section id="categories-section">
        <h2 className="text-2xl font-semibold mb-4 text-center">Categories</h2>
        {/* Add/Edit Category Form Card */}
        <div className={cardClasses}>
          <h3 className="text-xl font-semibold mb-3">{isEditing ? 'Edit Category' : 'Add New Category'}</h3>
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
        <div className="mt-4 p-3 bg-red-800 border border-red-700 text-white rounded text-center">
          <p>Category Error: {error}</p>
        </div>
      )}

      {/* Categories List Card */}
      <div className={`${cardClasses} mt-8`}>
        <h3 className="text-xl font-semibold mb-4">Existing Categories</h3>
        {isLoading && <p>Loading categories...</p>}
        {!isLoading && categories.length === 0 && !error && (
          <p className={textMutedClasses}>No categories found. Add some using the form above.</p>
        )}
        {categories.length > 0 && (
          <ul className="space-y-3">
            {categories.map(cat => (
              <li key={cat.id} className="flex items-center justify-between p-3 bg-gray-700 rounded">
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
      </section>

      {/* Budget Lines Section - Placeholder for now */}
      <section id="budget-lines-section" className="mt-12">
        <h2 className="text-2xl font-semibold mb-4 text-center">Budget Lines for Month ID: {currentMonthId}</h2>
        
        {/* Add/Edit Budget Line Form Card */}
        <div className={cardClasses}>
          <h3 className="text-xl font-semibold mb-3">{isEditingBL ? 'Edit Budget Line' : 'Add New Budget Line'}</h3>
          <form onSubmit={handleBLFormSubmit} className="space-y-3">
            <div>
              <label htmlFor="bl_label" className="block text-sm font-medium mb-1">Label:</label>
              <input
                id="bl_label"
                type="text"
                value={formBLLabel}
                onChange={(e) => setFormBLLabel(e.target.value)}
                className={`${inputClasses} w-full`}
                placeholder="e.g., Coffee Supplies"
                required
              />
            </div>
            <div>
              <label htmlFor="bl_expected" className="block text-sm font-medium mb-1">Expected Amount:</label>
              <input
                id="bl_expected"
                type="number"
                value={formBLExpected}
                onChange={(e) => setFormBLExpected(e.target.value)}
                className={`${inputClasses} w-full`}
                placeholder="e.g., 50.00"
                step="0.01"
                required
              />
            </div>
            <div>
              <label htmlFor="bl_category" className="block text-sm font-medium mb-1">Category:</label>
              <select
                id="bl_category"
                value={formBLCategoryID}
                onChange={(e) => setFormBLCategoryID(e.target.value)}
                className={`${inputClasses} w-full`}
                required
                disabled={categories.length === 0}
              >
                <option value="" disabled className={textMutedClasses}>
                  {categories.length === 0 ? 'Please add categories first' : 'Select a category'}
                </option>
                {categories.map(cat => (
                  <option key={cat.id} value={cat.id.toString()}>
                    {cat.name}
                  </option>
                ))}
              </select>
            </div>
            <div className="flex space-x-2">
              <button type="submit" className={buttonClasses} disabled={categories.length === 0 || !currentMonthId}>
                {isEditingBL ? 'Update Budget Line' : 'Add Budget Line'}
              </button>
              {isEditingBL && (
                <button type="button" onClick={resetBLForm} className={`${buttonClasses} bg-gray-500 hover:bg-gray-600`}>
                  Cancel Edit
                </button>
              )}
            </div>
          </form>
        </div>

        {errorBL && (
          <div className="mt-4 p-3 bg-red-800 border border-red-700 text-white rounded text-center">
            <p>Budget Line Error: {errorBL}</p>
          </div>
        )}

        {/* Budget Lines List Card */}
        <div className={`${cardClasses} mt-8`}>
          <h3 className="text-xl font-semibold mb-4">Existing Budget Lines</h3>
          {isLoadingBL && <p>Loading budget lines...</p>}
          {!isLoadingBL && budgetLines.length === 0 && !errorBL && (
             <p className={textMutedClasses}>
             {currentMonthId ? 'No budget lines found for this month.' : 'Select a month ID to see budget lines.'}
           </p>
          )}
          {budgetLines.length > 0 && (
            <ul className="space-y-3">
              {budgetLines.map(bl => (
                <li key={bl.id} className="flex items-center justify-between p-3 bg-gray-700 rounded">
                  <div className="flex items-center">
                    <span className={`inline-block w-5 h-5 rounded mr-3 ${bl.category_color || 'bg-gray-500'}`}></span>
                    <span className="font-medium mr-2">{bl.label}</span>
                    <span className="text-sm text-gray-400">({bl.category_name})</span>
                  </div>
                  <div className="flex items-center space-x-3">
                    <span className="font-mono text-sm">${bl.expected.toFixed(2)}</span>
                    <button onClick={() => handleBLEdit(bl)} className={`${buttonClasses} bg-yellow-500 hover:bg-yellow-600 text-xs`}>
                      Edit
                    </button>
                    <button onClick={() => handleBLDelete(bl.id)} className={`${destructiveButtonClasses} text-xs`}>
                      Delete
                    </button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
      </section>
    </div>
  );
}
// Add api import
import * as api from '../lib/api';

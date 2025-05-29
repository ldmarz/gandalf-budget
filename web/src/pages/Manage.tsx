import React, { useState, useEffect, FormEvent } from 'react';
import { get, post, put, del } from '../lib/api'; // Adjust path as necessary
import { textMutedClasses } from '../styles/commonClasses'; // Removed inputClasses, buttonClasses, destructiveButtonClasses
import Card from '../components/ui/Card'; // Import Card component
import Button from '../components/ui/Button'; // Import Button component
import Input from '../components/ui/Input'; // Import Input component
import Select from '../components/ui/Select'; // Import Select component
import LoadingSpinner from '../components/ui/LoadingSpinner'; // Import LoadingSpinner
import MessageDisplay from '../components/ui/MessageDisplay'; // Import MessageDisplay
import CategoryBadge from '../components/CategoryBadge'; // Import CategoryBadge

interface Category {
  id: number;
  name: string;
  color: string;
}

export default function ManagePage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [currentCategory, setCurrentCategory] = useState<Category | null>(null);
  const [formName, setFormName] = useState<string>('');
  const [formColor, setFormColor] = useState<string>('');

  const [budgetLines, setBudgetLines] = useState<api.BudgetLine[]>([]);
  const [currentMonthId, setCurrentMonthId] = useState<number>(1);
  const [isLoadingBL, setIsLoadingBL] = useState(true);
  const [errorBL, setErrorBL] = useState<string | null>(null);

  const [isEditingBL, setIsEditingBL] = useState<boolean>(false);
  const [currentBL, setCurrentBL] = useState<api.BudgetLine | null>(null);
  const [formBLLabel, setFormBLLabel] = useState<string>('');
  const [formBLExpected, setFormBLExpected] = useState<string>('');
  const [formBLCategoryID, setFormBLCategoryID] = useState<string>('');

  const tailwindColors = [
    'bg-red-500', 'bg-orange-500', 'bg-amber-500', 'bg-yellow-500', 'bg-lime-500',
    'bg-green-500', 'bg-emerald-500', 'bg-teal-500', 'bg-cyan-500', 'bg-sky-500',
    'bg-blue-500', 'bg-indigo-500', 'bg-violet-500', 'bg-purple-500', 'bg-fuchsia-500',
    'bg-pink-500', 'bg-rose-500', 'bg-slate-500', 'bg-gray-500',
  ];

  const fetchCategories = async () => {
    setIsLoading(true);
    try {
      const data = await get<Category[]>('/categories');
      setCategories(data || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch categories');
      setCategories([]);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchCategories();
    if (currentMonthId) {
      fetchBudgetLines(currentMonthId);
    }
  }, [currentMonthId]);

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
      await fetchCategories();
      resetForm();
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : (isEditing ? 'Failed to update category' : 'Failed to create category');
      alert(errorMsg);
      setError(errorMsg);
    }
  };

  const handleEdit = (category: Category) => {
    setIsEditing(true);
    setCurrentCategory(category);
    setFormName(category.name);
    setFormColor(category.color);
    window.scrollTo(0, 0);
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this category?')) {
      try {
        await del<null>(`/categories/${id}`);
        await fetchCategories();
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

  const fetchBudgetLines = async (monthId: number) => {
    setIsLoadingBL(true);
    try {
      const data = await api.getBudgetLinesByMonth(monthId);
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
  }, [categories, currentMonthId]);

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
        await api.updateBudgetLine(currentBL.id, { label: formBLLabel, expected: expectedAmount });
      } else {
        await api.createBudgetLine(budgetLineData);
      }
      await fetchBudgetLines(currentMonthId);
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
        await fetchBudgetLines(currentMonthId);
      } catch (err) {
        const errorMsg = err instanceof Error ? err.message : 'Failed to delete budget line';
        alert(errorMsg);
        setErrorBL(errorMsg);
      }
    }
  };

  if (isLoading) return <LoadingSpinner text="Loading page..." />;

  return (
    <div className="p-4 bg-gray-900 min-h-screen text-white">
      <h1 className="text-2xl font-bold mb-6 text-center">Manage Data</h1>

      <Card className="mb-6">
        <h2 className="text-xl font-semibold mb-3">Select Month</h2>
        <div className="flex items-center space-x-2">
          <label htmlFor="month_id_selector" className="block text-sm font-medium">Month ID:</label>
          <Input
            id="month_id_selector"
            type="number"
            value={currentMonthId}
            onChange={(e) => setCurrentMonthId(parseInt(e.target.value, 10) || 1)}
            className="w-24"
            min="1"
          />
          <Button onClick={() => fetchBudgetLines(currentMonthId)}>Load Budget Lines</Button>
        </div>
        <p className={textMutedClasses}>Note: This is a simplified month selector. Use a valid Month ID from your database.</p>
      </Card>
      
      <section id="categories-section">
        <h2 className="text-2xl font-semibold mb-4 text-center">Categories</h2>
        <Card>
          <h3 className="text-xl font-semibold mb-3">{isEditing ? 'Edit Category' : 'Add New Category'}</h3>
          <form onSubmit={handleFormSubmit} className="space-y-3">
            <div>
              <label htmlFor="name" className="block text-sm font-medium mb-1">Name:</label>
            <Input
              id="name"
              type="text"
              value={formName}
              onChange={(e) => setFormName(e.target.value)}
              className="w-full"
              placeholder="e.g., Groceries"
              required
            />
          </div>
          <div>
            <label htmlFor="color" className="block text-sm font-medium mb-1">Color (Tailwind Class):</label>
            <Select
              id="color"
              value={formColor}
              onChange={(e) => setFormColor(e.target.value)}
              className="w-full"
              required
              options={[
                { value: '', label: 'Select a color', disabled: true },
                ...tailwindColors.map(colorClass => ({ value: colorClass, label: colorClass }))
              ]}
            />
            {formColor && (
              <div className="mt-2 p-2 flex items-center">
                <span className={`inline-block w-6 h-6 rounded mr-2 ${formColor}`}></span>
                <span>Selected color preview</span>
              </div>
            )}
          </div>
          <div className="flex space-x-2">
            <Button type="submit">
              {isEditing ? 'Update Category' : 'Add Category'}
            </Button>
            {isEditing && (
              <Button type="button" variant="secondary" onClick={resetForm}>
                Cancel Edit
              </Button>
            )}
          </div>
        </form>
      </Card>
      
      <MessageDisplay message={error ? `Category Error: ${error}` : null} type="error" className="mt-4 text-center" />

      <Card className="mt-8">
        <h3 className="text-xl font-semibold mb-4">Existing Categories</h3>
        {isLoading && <LoadingSpinner text="Loading categories..." />}
        {!isLoading && categories.length === 0 && !error && (
          <p className={textMutedClasses}>No categories found. Add some using the form above.</p>
        )}
        {categories.length > 0 && (
          <ul className="space-y-3">
            {categories.map(cat => (
              <li key={cat.id} className="flex items-center justify-between p-3 bg-gray-700 rounded">
                <CategoryBadge category={cat} />
                <div className="space-x-2">
                  <Button onClick={() => handleEdit(cat)} variant="secondary" className="bg-yellow-500 hover:bg-yellow-600">
                    Edit
                  </Button>
                  <Button onClick={() => handleDelete(cat.id)} variant="destructive">
                    Delete
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </Card>
      </section>

      <section id="budget-lines-section" className="mt-12">
        <h2 className="text-2xl font-semibold mb-4 text-center">Budget Lines for Month ID: {currentMonthId}</h2>
        
        <Card>
          <h3 className="text-xl font-semibold mb-3">{isEditingBL ? 'Edit Budget Line' : 'Add New Budget Line'}</h3>
          <form onSubmit={handleBLFormSubmit} className="space-y-3">
            <div>
              <label htmlFor="bl_label" className="block text-sm font-medium mb-1">Label:</label>
              <Input
                id="bl_label"
                type="text"
                value={formBLLabel}
                onChange={(e) => setFormBLLabel(e.target.value)}
                className="w-full"
                placeholder="e.g., Coffee Supplies"
                required
              />
            </div>
            <div>
              <label htmlFor="bl_expected" className="block text-sm font-medium mb-1">Expected Amount:</label>
              <Input
                id="bl_expected"
                type="number"
                value={formBLExpected}
                onChange={(e) => setFormBLExpected(e.target.value)}
                className="w-full"
                placeholder="e.g., 50.00"
                step="0.01"
                required
              />
            </div>
            <div>
              <label htmlFor="bl_category" className="block text-sm font-medium mb-1">Category:</label>
              <Select
                id="bl_category"
                value={formBLCategoryID}
                onChange={(e) => setFormBLCategoryID(e.target.value)}
                className="w-full"
                required
                disabled={categories.length === 0}
                options={[
                  { value: '', label: categories.length === 0 ? 'Please add categories first' : 'Select a category', disabled: true },
                  ...categories.map(cat => ({ value: cat.id.toString(), label: cat.name }))
                ]}
              />
            </div>
            <div className="flex space-x-2">
              <Button type="submit" disabled={categories.length === 0 || !currentMonthId}>
                {isEditingBL ? 'Update Budget Line' : 'Add Budget Line'}
              </Button>
              {isEditingBL && (
                <Button type="button" variant="secondary" onClick={resetBLForm}>
                  Cancel Edit
                </Button>
              )}
            </div>
          </form>
        </Card>

      <MessageDisplay message={errorBL ? `Budget Line Error: ${errorBL}` : null} type="error" className="mt-4 text-center" />

        <Card className="mt-8">
          <h3 className="text-xl font-semibold mb-4">Existing Budget Lines</h3>
        {isLoadingBL && <LoadingSpinner text="Loading budget lines..." />}
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
                    <CategoryBadge category={{ name: bl.category_name, color: bl.category_color || 'bg-gray-500' }} className="mr-2"/>
                    <span className="font-medium">{bl.label}</span>
                  </div>
                  <div className="flex items-center space-x-3">
                    <span className="font-mono text-sm">${bl.expected.toFixed(2)}</span>
                    <Button onClick={() => handleBLEdit(bl)} variant="secondary" className="bg-yellow-500 hover:bg-yellow-600 text-xs">
                      Edit
                    </Button>
                    <Button onClick={() => handleBLDelete(bl.id)} variant="destructive" className="text-xs">
                      Delete
                    </Button>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </Card>
      </section>
    </div>
  );
}

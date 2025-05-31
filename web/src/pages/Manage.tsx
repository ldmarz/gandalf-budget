import React, { useState, useEffect, FormEvent } from 'react';
import * as api from '../lib/api'; // Assuming api.get, api.post etc. are defined
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/Card';
import { Button } from '../components/ui/Button';
import { Input } from '../components/ui/Input';
import { Select as UiSelect } from '../components/ui/Select'; // Renamed to avoid conflict with HTMLSelectElement
import { Loader2 } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '../components/ui/Alert';
import { CategoryBadge } from '../components/CategoryBadge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../components/ui/Table';
import { formatCurrency } from '../lib/utils';

interface Category {
  id: number;
  name: string;
  color: string;
}

// Define BudgetLine interface if not already available from api module
interface BudgetLine extends api.BudgetLine { // Assuming api.BudgetLine exists
    category_name?: string;
    category_color?: string;
}


export default function ManagePage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoadingCategories, setIsLoadingCategories] = useState(true);
  const [errorCategories, setErrorCategories] = useState<string | null>(null);

  const [isEditingCategory, setIsEditingCategory] = useState<boolean>(false);
  const [currentCategory, setCurrentCategory] = useState<Category | null>(null);
  const [categoryFormName, setCategoryFormName] = useState<string>('');
  const [categoryFormColor, setCategoryFormColor] = useState<string>('');

  const [budgetLines, setBudgetLines] = useState<BudgetLine[]>([]);
  const [currentMonthId, setCurrentMonthId] = useState<number>(() => {
    const params = new URLSearchParams(window.location.search);
    return parseInt(params.get('month_id') || '1', 10);
  });
  const [isLoadingBudgetLines, setIsLoadingBudgetLines] = useState(true);
  const [errorBudgetLines, setErrorBudgetLines] = useState<string | null>(null);

  const [isEditingBudgetLine, setIsEditingBudgetLine] = useState<boolean>(false);
  const [currentBudgetLine, setCurrentBudgetLine] = useState<BudgetLine | null>(null);
  const [budgetLineFormLabel, setBudgetLineFormLabel] = useState<string>('');
  const [budgetLineFormExpected, setBudgetLineFormExpected] = useState<string>('');
  const [budgetLineFormCategoryID, setBudgetLineFormCategoryID] = useState<string>('');

  // Tailwind CSS color palette (ensure these are full class names for JIT compilation)
  const tailwindColorOptions = [
    { value: 'bg-red-500', label: 'Red' }, { value: 'bg-orange-500', label: 'Orange' },
    { value: 'bg-amber-500', label: 'Amber' }, { value: 'bg-yellow-500', label: 'Yellow' },
    { value: 'bg-lime-500', label: 'Lime' }, { value: 'bg-green-500', label: 'Green' },
    { value: 'bg-emerald-500', label: 'Emerald' }, { value: 'bg-teal-500', label: 'Teal' },
    { value: 'bg-cyan-500', label: 'Cyan' }, { value: 'bg-sky-500', label: 'Sky' },
    { value: 'bg-blue-500', label: 'Blue' }, { value: 'bg-indigo-500', label: 'Indigo' },
    { value: 'bg-violet-500', label: 'Violet' }, { value: 'bg-purple-500', label: 'Purple' },
    { value: 'bg-fuchsia-500', label: 'Fuchsia' }, { value: 'bg-pink-500', label: 'Pink' },
    { value: 'bg-rose-500', label: 'Rose' }, { value: 'bg-slate-500', label: 'Slate' },
    { value: 'bg-gray-500', label: 'Gray' },
  ];

  const fetchCategories = async () => {
    setIsLoadingCategories(true);
    setErrorCategories(null);
    try {
      const data = await api.get<Category[]>('/categories');
      setCategories(data || []);
    } catch (err) {
      setErrorCategories(err instanceof Error ? err.message : 'Failed to fetch categories');
      setCategories([]);
    } finally {
      setIsLoadingCategories(false);
    }
  };

  const fetchBudgetLines = async (monthId: number) => {
    setIsLoadingBudgetLines(true);
    setErrorBudgetLines(null);
    try {
      const data = await api.getBudgetLinesByMonth(monthId);
      // Map category details to budget lines
      const linesWithCategoryDetails: BudgetLine[] = data.map(line => {
        const category = categories.find(c => c.id === line.category_id);
        return {
          ...line,
          category_name: category?.name || 'N/A',
          category_color: category?.color || 'bg-gray-300', // Default color
        };
      });
      setBudgetLines(linesWithCategoryDetails);
    } catch (err) {
      setErrorBudgetLines(err instanceof Error ? err.message : 'Failed to fetch budget lines');
      setBudgetLines([]);
    } finally {
      setIsLoadingBudgetLines(false);
    }
  };

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(() => {
    if (currentMonthId && categories.length > 0) { // Fetch budget lines only if categories are loaded
      fetchBudgetLines(currentMonthId);
    }
     const newSearch = new URLSearchParams(window.location.search);
    newSearch.set('month_id', currentMonthId.toString());
    window.history.replaceState({}, '', `${window.location.pathname}?${newSearch}`);
  }, [currentMonthId, categories]);


  const handleCategoryFormSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!categoryFormName.trim() || !categoryFormColor.trim()) {
      setErrorCategories('Name and color are required.');
      return;
    }
    setErrorCategories(null);

    const categoryData = { name: categoryFormName, color: categoryFormColor };

    try {
      if (isEditingCategory && currentCategory) {
        await api.put<Category, Category>(`/categories/${currentCategory.id}`, { ...categoryData, id: currentCategory.id });
      } else {
        await api.post<Category, Omit<Category, 'id'>>('/categories', categoryData);
      }
      await fetchCategories(); // Refresh categories list
      resetCategoryForm();
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : (isEditingCategory ? 'Failed to update category' : 'Failed to create category');
      setErrorCategories(errorMsg);
    }
  };

  const handleEditCategory = (category: Category) => {
    setIsEditingCategory(true);
    setCurrentCategory(category);
    setCategoryFormName(category.name);
    setCategoryFormColor(category.color);
    document.getElementById('category-form-section')?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleDeleteCategory = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this category? This might affect existing budget lines.')) {
      setErrorCategories(null);
      try {
        await api.del<null>(`/categories/${id}`);
        await fetchCategories(); // Refresh categories list
      } catch (err) {
        const errorMsg = err instanceof Error ? err.message : 'Failed to delete category. It might be in use.';
        setErrorCategories(errorMsg);
      }
    }
  };

  const resetCategoryForm = () => {
    setIsEditingCategory(false);
    setCurrentCategory(null);
    setCategoryFormName('');
    setCategoryFormColor('');
    setErrorCategories(null);
  };


  const resetBudgetLineForm = () => {
    setIsEditingBudgetLine(false);
    setCurrentBudgetLine(null);
    setBudgetLineFormLabel('');
    setBudgetLineFormExpected('');
    setBudgetLineFormCategoryID('');
    setErrorBudgetLines(null);
  };

  const handleEditBudgetLine = (bl: BudgetLine) => {
    setIsEditingBudgetLine(true);
    setCurrentBudgetLine(bl);
    setBudgetLineFormLabel(bl.label);
    setBudgetLineFormExpected(bl.expected.toString());
    setBudgetLineFormCategoryID(bl.category_id.toString());
    document.getElementById('budget-line-form-section')?.scrollIntoView({ behavior: 'smooth' });
  };
  
  const handleBudgetLineFormSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!budgetLineFormLabel.trim() || !budgetLineFormExpected.trim() || !budgetLineFormCategoryID.trim() || !currentMonthId) {
      setErrorBudgetLines('Label, expected amount, category, and a selected month are required.');
      return;
    }
    setErrorBudgetLines(null);

    const expectedAmount = parseFloat(budgetLineFormExpected);
    const categoryId = parseInt(budgetLineFormCategoryID, 10);

    if (isNaN(expectedAmount) || isNaN(categoryId)) {
      setErrorBudgetLines('Invalid expected amount or category ID.');
      return;
    }

    const budgetLineData = {
      month_id: currentMonthId,
      category_id: categoryId,
      label: budgetLineFormLabel,
      expected: expectedAmount
    };

    try {
      if (isEditingBudgetLine && currentBudgetLine) {
        // Ensure all necessary fields are passed for update if your API expects the full object
        await api.updateBudgetLine(currentBudgetLine.id, {
            ...currentBudgetLine, // spread existing fields if API needs them
            label: budgetLineFormLabel,
            expected: expectedAmount,
            category_id: categoryId, // ensure category_id can be updated if needed
            month_id: currentMonthId // ensure month_id can be updated if needed
        });
      } else {
        await api.createBudgetLine(budgetLineData);
      }
      await fetchBudgetLines(currentMonthId); // Refresh budget lines
      resetBudgetLineForm();
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : (isEditingBudgetLine ? 'Failed to update budget line' : 'Failed to create budget line');
      setErrorBudgetLines(errorMsg);
    }
  };

  const handleDeleteBudgetLine = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this budget line?')) {
      setErrorBudgetLines(null);
      try {
        await api.deleteBudgetLine(id);
        await fetchBudgetLines(currentMonthId); // Refresh budget lines
      } catch (err) {
        const errorMsg = err instanceof Error ? err.message : 'Failed to delete budget line';
        setErrorBudgetLines(errorMsg);
      }
    }
  };


  if (isLoadingCategories && isLoadingBudgetLines) { // Initial full page load
    return (
      <div className="flex flex-col justify-center items-center h-screen bg-gray-100 p-4 text-center">
        <Loader2 className="h-12 w-12 animate-spin text-blue-500 mb-4" />
        <p className="text-xl text-gray-700">Loading management data...</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4 sm:p-6 lg:p-8 bg-gray-100 min-h-screen">
      <header className="text-center mb-10">
        <h1 className="text-4xl font-bold text-gray-800">Manage Data</h1>
      </header>

      <Card className="mb-8 shadow-lg">
        <CardHeader>
          <CardTitle className="text-2xl font-semibold text-gray-700">Select Month for Budget Lines</CardTitle>
        </CardHeader>
        <CardContent className="pt-4 flex flex-col sm:flex-row items-start sm:items-center space-y-3 sm:space-y-0 sm:space-x-3">
          <div className="flex-grow sm:flex-grow-0">
            <label htmlFor="month_id_selector_manage" className="block text-sm font-medium text-gray-700 mb-1">
              Month ID:
            </label>
            <Input
              id="month_id_selector_manage"
              type="number"
              value={currentMonthId}
              onChange={(e) => setCurrentMonthId(parseInt(e.target.value, 10) || 1)}
              className="w-full sm:w-32 border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
              min="1"
            />
          </div>
          {/* Removed the explicit "Load Budget Lines" button as it fetches on monthId change */}
          <p className="text-sm text-gray-500 pt-1 sm:pt-6">
            Budget lines for the selected month will load automatically. This is a simplified month selector.
          </p>
        </CardContent>
      </Card>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 lg:gap-x-8 space-y-10 lg:space-y-0">
        {/* Categories Section */}
        <section id="categories-section" className="space-y-6">
          <h2 className="text-2xl font-semibold text-gray-700 text-center lg:text-left">Manage Categories</h2>

          {/* Category Form Card */}
          <Card className="shadow-lg" id="category-form-section">
            <CardHeader>
              <CardTitle className="text-xl font-semibold text-gray-700">
                {isEditingCategory ? 'Edit Category' : 'Add New Category'}
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-2">
              <form onSubmit={handleCategoryFormSubmit} className="space-y-4">
                <div>
                  <label htmlFor="categoryName" className="block text-sm font-medium text-gray-700 mb-1">Name:</label>
                  <Input
                    id="categoryName"
                    type="text"
                    value={categoryFormName}
                    onChange={(e) => setCategoryFormName(e.target.value)}
                    className="w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 p-2"
                    placeholder="e.g., Groceries"
                    required
                  />
                </div>
                <div>
                  <label htmlFor="categoryColor" className="block text-sm font-medium text-gray-700 mb-1">Color:</label>
                  <UiSelect
                    id="categoryColor"
                    value={categoryFormColor}
                    onChange={(e) => setCategoryFormColor(e.target.value)}
                    className="w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 p-2"
                    required
                  >
                    <option value="" disabled>Select a color</option>
                    {tailwindColorOptions.map(opt => (
                      <option key={opt.value} value={opt.value}>{opt.label}</option>
                    ))}
                  </UiSelect>
                  {categoryFormColor && (
                    <div className="mt-2 p-2 flex items-center space-x-2">
                      <span className={`inline-block w-5 h-5 rounded-full ${categoryFormColor}`}></span>
                      <span className="text-sm text-gray-600">Selected color preview</span>
                    </div>
                  )}
                </div>
                <div className="flex items-center space-x-3 pt-2">
                  <Button
                    type="submit"
                    className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-md shadow-sm disabled:opacity-50"
                    disabled={isLoadingCategories}
                  >
                    {isLoadingCategories && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                    {isEditingCategory ? 'Update Category' : 'Add Category'}
                  </Button>
                  {isEditingCategory && (
                    <Button
                      type="button"
                      variant="outline"
                      onClick={resetCategoryForm}
                      className="font-medium py-2 px-4 rounded-md shadow-sm"
                    >
                      Cancel Edit
                    </Button>
                  )}
                </div>
              </form>
            </CardContent>
          </Card>

          {errorCategories && (
            <Alert variant="destructive" className="max-w-md mx-auto">
              <AlertTitle className="font-semibold">Category Error</AlertTitle>
              <AlertDescription>{errorCategories}</AlertDescription>
            </Alert>
          )}

          {/* Existing Categories List/Table Card */}
          <Card className="shadow-lg">
            <CardHeader>
              <CardTitle className="text-xl font-semibold text-gray-700">Existing Categories</CardTitle>
            </CardHeader>
            <CardContent>
              {isLoadingCategories && (
                <div className="flex flex-col items-center justify-center py-10">
                  <Loader2 className="h-10 w-10 animate-spin text-blue-500 mb-3" />
                  <p className="text-gray-600">Loading categories...</p>
                </div>
              )}
              {!isLoadingCategories && categories.length === 0 && !errorCategories && (
                <p className="text-center text-gray-500 py-4">No categories found. Add some using the form above.</p>
              )}
              {!isLoadingCategories && categories.length > 0 && (
                <div className="overflow-x-auto">
                  <Table>
                    <TableHeader className="bg-gray-100">
                      <TableRow>
                        <TableHead className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Preview</TableHead>
                        <TableHead className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</TableHead>
                        <TableHead className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody className="bg-white divide-y divide-gray-200">
                      {categories.map(cat => (
                        <TableRow key={cat.id} className="hover:bg-gray-50">
                          <TableCell className="px-4 py-3 whitespace-nowrap">
                             <CategoryBadge category={cat} />
                          </TableCell>
                          <TableCell className="px-4 py-3 whitespace-nowrap text-sm text-gray-700">{cat.name}</TableCell>
                          <TableCell className="px-4 py-3 whitespace-nowrap text-right space-x-2">
                            <Button
                              onClick={() => handleEditCategory(cat)}
                              variant="outline"
                              size="sm"
                              className="text-xs py-1 px-2 rounded-md border-yellow-500 text-yellow-600 hover:bg-yellow-50"
                            >
                              Edit
                            </Button>
                            <Button
                              onClick={() => handleDeleteCategory(cat.id)}
                              variant="destructive"
                              size="sm"
                              className="text-xs py-1 px-2 rounded-md bg-red-500 hover:bg-red-600 text-white"
                            >
                              Delete
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              )}
            </CardContent>
          </Card>
        </section>

        {/* Budget Lines Section */}
        <section id="budget-lines-section" className="space-y-6">
           <h2 className="text-2xl font-semibold text-gray-700 text-center lg:text-left">Manage Budget Lines (Month ID: {currentMonthId})</h2>

          {/* Budget Line Form Card */}
          <Card className="shadow-lg" id="budget-line-form-section">
            <CardHeader>
              <CardTitle className="text-xl font-semibold text-gray-700">
                {isEditingBudgetLine ? 'Edit Budget Line' : 'Add New Budget Line'}
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-2">
              <form onSubmit={handleBudgetLineFormSubmit} className="space-y-4">
                <div>
                  <label htmlFor="bl_label" className="block text-sm font-medium text-gray-700 mb-1">Label:</label>
                  <Input
                    id="bl_label"
                    type="text"
                    value={budgetLineFormLabel}
                    onChange={(e) => setBudgetLineFormLabel(e.target.value)}
                    className="w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 p-2"
                    placeholder="e.g., Coffee Supplies"
                    required
                  />
                </div>
                <div>
                  <label htmlFor="bl_expected" className="block text-sm font-medium text-gray-700 mb-1">Expected Amount:</label>
                  <Input
                    id="bl_expected"
                    type="number"
                    value={budgetLineFormExpected}
                    onChange={(e) => setBudgetLineFormExpected(e.target.value)}
                    className="w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 p-2"
                    placeholder="e.g., 50.00"
                    step="0.01"
                    required
                  />
                </div>
                <div>
                  <label htmlFor="bl_category" className="block text-sm font-medium text-gray-700 mb-1">Category:</label>
                  <UiSelect
                    id="bl_category"
                    value={budgetLineFormCategoryID}
                    onChange={(e) => setBudgetLineFormCategoryID(e.target.value)}
                    className="w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 p-2"
                    required
                    disabled={categories.length === 0}
                  >
                    <option value="" disabled>
                      {categories.length === 0 ? 'Create categories first' : 'Select a category'}
                    </option>
                    {categories.map(cat => ({ value: cat.id.toString(), label: cat.name }))
                     .map(opt => <option key={opt.value} value={opt.value}>{opt.label}</option>)}
                  </UiSelect>
                </div>
                <div className="flex items-center space-x-3 pt-2">
                  <Button
                    type="submit"
                    disabled={categories.length === 0 || !currentMonthId || isLoadingBudgetLines}
                    className="bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-md shadow-sm disabled:opacity-50"
                  >
                    {isLoadingBudgetLines && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                    {isEditingBudgetLine ? 'Update Budget Line' : 'Add Budget Line'}
                  </Button>
                  {isEditingBudgetLine && (
                    <Button
                      type="button"
                      variant="outline"
                      onClick={resetBudgetLineForm}
                      className="font-medium py-2 px-4 rounded-md shadow-sm"
                    >
                      Cancel Edit
                    </Button>
                  )}
                </div>
              </form>
            </CardContent>
          </Card>

          {errorBudgetLines && (
            <Alert variant="destructive" className="max-w-md mx-auto">
              <AlertTitle className="font-semibold">Budget Line Error</AlertTitle>
              <AlertDescription>{errorBudgetLines}</AlertDescription>
            </Alert>
          )}

          {/* Existing Budget Lines List/Table Card */}
          <Card className="shadow-lg">
            <CardHeader>
              <CardTitle className="text-xl font-semibold text-gray-700">Existing Budget Lines</CardTitle>
               <CardDescription className="text-sm text-gray-500">For Month ID: {currentMonthId}</CardDescription>
            </CardHeader>
            <CardContent>
              {isLoadingBudgetLines && (
                 <div className="flex flex-col items-center justify-center py-10">
                  <Loader2 className="h-10 w-10 animate-spin text-blue-500 mb-3" />
                  <p className="text-gray-600">Loading budget lines...</p>
                </div>
              )}
              {!isLoadingBudgetLines && budgetLines.length === 0 && !errorBudgetLines && (
                <p className="text-center text-gray-500 py-4">
                  {currentMonthId ? `No budget lines found for Month ID ${currentMonthId}. Add some above.` : 'Select a month ID to see budget lines.'}
                </p>
              )}
              {!isLoadingBudgetLines && budgetLines.length > 0 && (
                 <div className="overflow-x-auto">
                  <Table>
                    <TableHeader className="bg-gray-100">
                      <TableRow>
                        <TableHead className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Label</TableHead>
                        <TableHead className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</TableHead>
                        <TableHead className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Expected</TableHead>
                        <TableHead className="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody className="bg-white divide-y divide-gray-200">
                      {budgetLines.map(bl => (
                        <TableRow key={bl.id} className="hover:bg-gray-50">
                          <TableCell className="px-4 py-3 whitespace-nowrap text-sm text-gray-700">{bl.label}</TableCell>
                          <TableCell className="px-4 py-3 whitespace-nowrap">
                             <CategoryBadge category={{ name: bl.category_name || 'N/A', color: bl.category_color || 'bg-gray-300' }} />
                          </TableCell>
                          <TableCell className="px-4 py-3 whitespace-nowrap text-sm text-gray-700 text-right">{formatCurrency(bl.expected)}</TableCell>
                          <TableCell className="px-4 py-3 whitespace-nowrap text-right space-x-2">
                            <Button
                              onClick={() => handleEditBudgetLine(bl)}
                              variant="outline"
                              size="sm"
                              className="text-xs py-1 px-2 rounded-md border-yellow-500 text-yellow-600 hover:bg-yellow-50"
                            >
                              Edit
                            </Button>
                            <Button
                              onClick={() => handleDeleteBudgetLine(bl.id)}
                              variant="destructive"
                              size="sm"
                              className="text-xs py-1 px-2 rounded-md bg-red-500 hover:bg-red-600 text-white"
                            >
                              Delete
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </div>
              )}
            </CardContent>
          </Card>
        </section>
      </div>
    </div>
  );
// No changes to this part of the file
}

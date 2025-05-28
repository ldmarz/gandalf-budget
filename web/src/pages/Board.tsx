import React, { useState, useEffect, ChangeEvent } from 'react';
import * as api from '../lib/api'; // Assuming api.ts is in ../lib

// Using existing Tailwind classes from Manage.tsx for consistency
const inputClasses = "border border-gray-300 rounded px-2 py-1 text-black text-sm w-24"; // Adjusted for board
const cardClasses = "bg-gray-800 p-4 rounded shadow-md mb-4";
const textMutedClasses = "text-gray-400";
const buttonClasses = "bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded text-sm";

export default function BoardPage() {
  const [boardLines, setBoardLines] = useState<api.BudgetLine[]>([]);
  const [categories, setCategories] = useState<api.Category[]>([]);
  const [currentMonthId, setCurrentMonthId] = useState<number>(1); // Default to month 1
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isLoadingCategories, setIsLoadingCategories] = useState(true);

  // Fetch categories
  const fetchCategories = async () => {
    setIsLoadingCategories(true);
    try {
      const data = await api.getAllCategories();
      setCategories(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch categories');
      setCategories([]);
    } finally {
      setIsLoadingCategories(false);
    }
  };

  // Fetch budget lines for the board
  const fetchBoardData = async (monthId: number) => {
    if (categories.length === 0) { // Ensure categories are loaded first
        // console.log("Categories not loaded yet, deferring board data fetch");
        return;
    }
    setIsLoading(true);
    setError(null);
    try {
      const data = await api.getBudgetLinesByMonth(monthId);
      const enrichedData = data.map(line => {
        const category = categories.find(c => c.id === line.category_id);
        return {
          ...line,
          category_name: category?.name || 'Unknown',
          category_color: category?.color || 'bg-gray-500',
          actual_amount: line.actual_amount || 0, // Ensure actual_amount is initialized
          actual_id: line.actual_id // Will be undefined if not provided by API
        };
      });
      setBoardLines(enrichedData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch board data');
      setBoardLines([]);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(() => {
    if (categories.length > 0) { // Only fetch board data if categories are available
        fetchBoardData(currentMonthId);
    }
  }, [currentMonthId, categories]);


  const handleActualAmountChange = async (
    budgetLineId: number, // Keep for finding the line in state
    actualLineId: number | undefined, 
    newActualString: string
  ) => {
    const newActual = parseFloat(newActualString);
    if (isNaN(newActual) || newActual < 0) {
      alert("Please enter a valid positive number for the actual amount.");
      // Optionally, revert the input to its previous state if needed
      fetchBoardData(currentMonthId); // Or update specific line from a backup
      return;
    }

    if (actualLineId === undefined) {
        setError(`Cannot update actual amount: ActualLine ID is missing for budget line ${budgetLineId}.`);
        // This indicates an issue with data from the backend (actual_id not provided)
        console.error("ActualLine ID is undefined for budget line:", budgetLineId);
        return;
    }

    // Optimistically update UI - or you can wait for API response
    setBoardLines(prevLines =>
      prevLines.map(line =>
        line.id === budgetLineId ? { ...line, actual_amount: newActual } : line
      )
    );

    try {
      await api.updateActualLine(actualLineId, { actual: newActual });
      // Optionally re-fetch to confirm or if backend does more logic
      // await fetchBoardData(currentMonthId); 
    } catch (err) {
      alert(`Failed to update actual amount: ${err instanceof Error ? err.message : 'Unknown error'}`);
      // Revert optimistic update on error
      setError(err instanceof Error ? `Failed to update: ${err.message}` : 'Failed to update actual amount.');
      fetchBoardData(currentMonthId); // Re-fetch to get the source of truth
    }
  };
  
  const getRowColor = (line: api.BudgetLine): string => {
    if (line.actual_amount && line.actual_amount > 0) {
      return 'bg-green-700 hover:bg-green-600'; // Darker green for dark theme
    }
    return 'bg-yellow-700 hover:bg-yellow-600'; // Darker yellow for dark theme
  };


  if (isLoadingCategories) return <div className="p-4 text-white">Loading categories...</div>;

  return (
    <div className="p-4 bg-gray-900 min-h-screen text-white">
      <h1 className="text-2xl font-bold mb-6 text-center">Monthly Budget Board</h1>

      <div className={`${cardClasses} mb-6`}>
        <h2 className="text-xl font-semibold mb-3">Select Month</h2>
        <div className="flex items-center space-x-2">
          <label htmlFor="month_id_selector_board" className="block text-sm font-medium">Month ID:</label>
          <input
            id="month_id_selector_board"
            type="number"
            value={currentMonthId}
            onChange={(e) => setCurrentMonthId(parseInt(e.target.value, 10) || 1)}
            className={`${inputClasses} w-24 !text-black`} // Ensure text is visible
            min="1"
          />
           {/* Removed the "Load Budget Lines" button, will auto-load on month change */}
        </div>
      </div>

      {error && (
        <div className="my-4 p-3 bg-red-800 border border-red-700 text-white rounded text-center">
          <p>Error: {error}</p>
        </div>
      )}

      {isLoading && <div className="text-center py-4">Loading board data...</div>}

      {!isLoading && !error && boardLines.length === 0 && (
        <div className={`${cardClasses} text-center`}>
          <p className={textMutedClasses}>No budget lines found for Month ID: {currentMonthId}.</p>
          <p className={textMutedClasses}>You can add budget lines in the 'Manage' page.</p>
        </div>
      )}

      {!isLoading && boardLines.length > 0 && (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-700">
            <thead className="bg-gray-700">
              <tr>
                <th scope="col" className="px-4 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Category</th>
                <th scope="col" className="px-4 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Budget Line Item</th>
                <th scope="col" className="px-4 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Expected (CLP)</th>
                <th scope="col" className="px-4 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Actual (CLP)</th>
              </tr>
            </thead>
            <tbody className="bg-gray-800 divide-y divide-gray-700">
              {boardLines.map(line => (
                <tr key={line.id} className={`${getRowColor(line)} transition-colors duration-150`}>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <div className="flex items-center">
                      <span className={`inline-block w-4 h-4 rounded mr-2 ${line.category_color || 'bg-gray-500'}`}></span>
                      {line.category_name}
                    </div>
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">{line.label}</td>
                  <td className="px-4 py-3 whitespace-nowrap text-right">{line.expected.toFixed(0)}</td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <input
                      type="number"
                      defaultValue={line.actual_amount || 0} // Use defaultValue for onBlur updates
                      onBlur={(e: ChangeEvent<HTMLInputElement>) => 
                        handleActualAmountChange(line.id, line.actual_id, e.target.value)
                      }
                      className={`${inputClasses} !text-black`} // Ensure text is visible
                      placeholder="0"
                      min="0"
                      step="1" // Assuming CLP doesn't typically use decimals
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

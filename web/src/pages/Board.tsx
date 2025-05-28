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
  const [isFinalizing, setIsFinalizing] = useState(false); // New state for finalization loading
  const [finalizeMessage, setFinalizeMessage] = useState<string | null>(null); // New state for success/error messages from finalize

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
    if (categories.length === 0) {
        return;
    }
    setIsLoading(true);
    setError(null);
    setFinalizeMessage(null); // Clear finalize message on new data load
    try {
      // Use the new getBoardData API endpoint
      const data = await api.getBoardData(monthId); 
      const enrichedData = data.map(line => {
        const category = categories.find(c => c.id === line.category_id);
        return {
          ...line,
          category_name: category?.name || 'Unknown Category',
          category_color: category?.color || 'bg-gray-500', // Default color
          // actual_amount and actual_id should now come directly from getBoardData
          // Ensure actual_amount is a number, defaulting to 0 if null/undefined
          actual_amount: line.actual_amount === undefined || line.actual_amount === null ? 0 : Number(line.actual_amount),
        };
      });
      setBoardLines(enrichedData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch board data');
      setBoardLines([]); // Clear board lines on error
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(() => {
    if (categories.length > 0) { 
        fetchBoardData(currentMonthId);
    }
  }, [currentMonthId, categories]); // Re-fetch when currentMonthId or categories change

  const handleFinalizeMonth = async () => {
    if (!currentMonthId) {
      setError("Month ID is not set.");
      return;
    }
    setIsFinalizing(true);
    setError(null);
    setFinalizeMessage(null);
    try {
      const response = await api.finalizeMonth(currentMonthId);
      setFinalizeMessage(response.message || "Month finalized successfully!");
      // setCurrentMonthId(response.new_month_id); // This will trigger useEffect to refetch board data for new month
      // Instead of directly setting, we might want to fetch board data for the *new* month ID explicitly
      // or rely on the user to navigate if the month ID display changes.
      // For now, setting currentMonthId will cause a re-fetch due to useEffect dependency.
      fetchBoardData(response.new_month_id); // Fetch data for the new month
      setCurrentMonthId(response.new_month_id); // Update the displayed month ID
      
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to finalize month. Check if all actuals are set or an error occurred.";
      setError(errorMessage);
      // setFinalizeMessage(errorMessage); // Show error in finalize message spot too, or rely on general error display
    } finally {
      setIsFinalizing(false);
    }
  };

  const handleActualAmountChange = async (
    budgetLineId: number, 
    actualLineId: number | undefined, 
    newActualString: string
  ) => {
    const newActual = parseFloat(newActualString);
    if (isNaN(newActual) || newActual < 0) {
      alert("Please enter a valid positive number for the actual amount.");
      fetchBoardData(currentMonthId); 
      return;
    }

    if (actualLineId === undefined) {
        setError(`Cannot update actual amount: ActualLine ID is missing for budget line ${budgetLineId}. This might mean the actual record was not created yet.`);
        console.error("ActualLine ID is undefined for budget line:", budgetLineId);
        return;
    }

    setBoardLines(prevLines =>
      prevLines.map(line =>
        line.id === budgetLineId ? { ...line, actual_amount: newActual } : line
      )
    );

    try {
      await api.updateActualLine(actualLineId, { actual: newActual });
    } catch (err) {
      alert(`Failed to update actual amount: ${err instanceof Error ? err.message : 'Unknown error'}`);
      setError(err instanceof Error ? `Failed to update: ${err.message}` : 'Failed to update actual amount.');
      fetchBoardData(currentMonthId); 
    }
  };
  
  const getRowColor = (line: api.BudgetLine): string => {
    if (line.actual_amount && Number(line.actual_amount) > 0) { // Ensure comparison with number
      return 'bg-green-700 hover:bg-green-600'; 
    }
    return 'bg-yellow-700 hover:bg-yellow-600'; 
  };

  if (isLoadingCategories) return <div className="p-4 text-white">Loading categories...</div>;

  return (
    <div className="p-4 bg-gray-900 min-h-screen text-white">
      <h1 className="text-2xl font-bold mb-6 text-center">Monthly Budget Board</h1>

      <div className={`${cardClasses} mb-6 flex flex-col sm:flex-row justify-between items-center space-y-4 sm:space-y-0`}>
        <div className="flex items-center space-x-2">
          <label htmlFor="month_id_selector_board" className="block text-sm font-medium">Month ID:</label>
          <input
            id="month_id_selector_board"
            type="number"
            value={currentMonthId}
            onChange={(e) => setCurrentMonthId(parseInt(e.target.value, 10) || 1)}
            className={`${inputClasses} w-24 !text-black`}
            min="1"
          />
        </div>
        <div>
          <button
            onClick={handleFinalizeMonth}
            disabled={isFinalizing || isLoading || boardLines.length === 0}
            className={`${buttonClasses} bg-green-600 hover:bg-green-700 disabled:bg-gray-500 disabled:cursor-not-allowed`}
          >
            {isFinalizing ? 'Finalizing...' : 'Finalize Current Month'}
          </button>
        </div>
      </div>

      {error && (
        <div className="my-4 p-3 bg-red-800 border border-red-700 text-white rounded text-center">
          <p>Error: {error}</p>
        </div>
      )}
      
      {finalizeMessage && (
        <div className={`my-4 p-3 rounded text-center ${error ? 'bg-red-800 border-red-700' : 'bg-green-800 border-green-700'} text-white`}>
          <p>{finalizeMessage}</p>
        </div>
      )}

      {isLoading && <div className="text-center py-4">Loading board data...</div>}

      {!isLoading && !error && boardLines.length === 0 && (
        <div className={`${cardClasses} text-center`}>
          <p className={textMutedClasses}>No budget lines found for Month ID: {currentMonthId}.</p>
          <p className={textMutedClasses}>You can add budget lines in the 'Manage' page for this month if it's not finalized.</p>
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
                  <td className="px-4 py-3 whitespace-nowrap text-right">{Number(line.expected).toFixed(0)}</td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <input
                      type="number"
                      defaultValue={Number(line.actual_amount) || 0}
                      onBlur={(e: ChangeEvent<HTMLInputElement>) => 
                        handleActualAmountChange(line.id, line.actual_id, e.target.value)
                      }
                      className={`${inputClasses} !text-black w-full`} // Ensure text is visible, take full cell width
                      placeholder="0"
                      min="0"
                      step="1"
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

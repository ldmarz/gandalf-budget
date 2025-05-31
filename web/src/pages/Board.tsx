import React, { useState, useEffect, ChangeEvent } from 'react';
import * as api from '../lib/api';
// import { textMutedClasses } from '../styles/commonClasses'; // Assuming this will be replaced by Tailwind
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/Card'; // Assuming these are Tailwind-styled
import { Button } from '../components/ui/Button'; // Assuming Tailwind-styled
import { Input } from '../components/ui/Input'; // Assuming Tailwind-styled
import { Loader2 } from 'lucide-react'; // For a spinner icon
import { Alert, AlertDescription, AlertTitle } from '../components/ui/Alert'; // Assuming Tailwind-styled
import { CategoryBadge } from '../components/CategoryBadge'; // Assuming Tailwind-styled
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../components/ui/Table'; // Assuming Tailwind-styled
import { formatCurrency } from '../lib/utils'; // For currency formatting

export default function BoardPage() {
  const [boardData, setBoardData] = useState<api.BoardDataPayload | null>(null);
  // Ensure budget_lines is part of BoardDataPayload or adjust type
  const [budgetLines, setBudgetLines] = useState<api.BudgetLineWithActual[]>([]);
  const [currentMonthId, setCurrentMonthId] = useState<number>(() => {
    // You might want to get this from URL params or a global state
    const params = new URLSearchParams(window.location.search);
    return parseInt(params.get('month_id') || '1', 10);
  });
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isFinalizing, setIsFinalizing] = useState(false);
  const [finalizeMessage, setFinalizeMessage] = useState<string | null>(null);

  const fetchBoardData = async (monthId: number) => {
    setIsLoading(true);
    setError(null);
    setFinalizeMessage(null);
    try {
      const data = await api.getBoardData(monthId);
      setBoardData(data);
      setBudgetLines(data.budget_lines || []); // Ensure budget_lines are set
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch board data');
      setBoardData(null);
      setBudgetLines([]);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchBoardData(currentMonthId);
    // Update URL search param when monthId changes
    const newSearch = new URLSearchParams(window.location.search);
    newSearch.set('month_id', currentMonthId.toString());
    window.history.replaceState({}, '', `${window.location.pathname}?${newSearch}`);
  }, [currentMonthId]);

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
      // Navigate to the new month or refresh data for the new month
      setCurrentMonthId(response.new_month_id);
      fetchBoardData(response.new_month_id); // fetch data for new month
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to finalize month. Check if all actuals are set or an error occurred.";
      setError(errorMessage);
    } finally {
      setIsFinalizing(false);
    }
  };

  const handleActualAmountChange = async (
    budgetLineId: number,
    actualId: number | null, // actual_id can be null if not yet set
    newActualString: string
  ) => {
    const newActual = parseFloat(newActualString);
    if (isNaN(newActual) || newActual < 0) {
      // Show a more user-friendly error, perhaps using a toast or inline message
      setError("Please enter a valid positive number for the actual amount.");
      // Re-fetch to revert optimistic update or show correct state
      fetchBoardData(currentMonthId);
      return;
    }
    setError(null); // Clear previous error

    // Optimistic UI update
    setBudgetLines(prevLines =>
      prevLines.map(line =>
        line.id === budgetLineId ? { ...line, actual_amount: newActual, actual_id: actualId || line.id } : line
      )
    );

    try {
      // If actual_id exists, it's an update, otherwise it's a create.
      // The backend API needs to handle this logic, or you need separate API calls.
      // Assuming `updateActualLine` can handle create/update based on `actual_id` or `budget_line_id`.
      // This might need adjustment based on your actual API capabilities.
      // For simplicity, let's assume `updateActualLine` handles this.
      // If it doesn't, you'd need a `createActualLine` and `updateActualLine`.
      await api.updateActualLine(budgetLineId, { actual: newActual });
    } catch (err) {
      setError(err instanceof Error ? `Failed to update: ${err.message}` : 'Failed to update actual amount.');
      // Revert optimistic update on error
      fetchBoardData(currentMonthId);
    }
  };

  if (isLoading && !boardData && !error) {
    return (
      <div className="flex flex-col justify-center items-center h-screen bg-gray-100 p-4 text-center">
        <Loader2 className="h-12 w-12 animate-spin text-blue-500 mb-4" />
        <p className="text-xl text-gray-700">Loading board data...</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-4 sm:p-6 lg:p-8 bg-gray-100 min-h-screen space-y-6">
      <header className="text-center">
        <h1 className="text-3xl font-bold text-gray-800">
          Monthly Budget Board
        </h1>
        {boardData && (
          <p className="text-xl text-gray-600">{boardData.month_name} {boardData.year}</p>
        )}
      </header>

      <Card className="shadow-lg">
        <CardContent className="p-6 flex flex-col sm:flex-row justify-between items-center space-y-4 sm:space-y-0 sm:space-x-4">
          <div className="flex items-center space-x-3">
            <label htmlFor="month_id_selector_board" className="text-sm font-medium text-gray-700">
              Select Month ID:
            </label>
            <Input
              id="month_id_selector_board"
              type="number"
              value={currentMonthId}
              onChange={(e) => setCurrentMonthId(parseInt(e.target.value, 10) || 1)}
              className="w-24 border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
              min="1"
            />
          </div>
          <Button
            onClick={handleFinalizeMonth}
            disabled={isFinalizing || isLoading || budgetLines.length === 0 }
            className="bg-green-600 hover:bg-green-700 text-white font-semibold py-2 px-4 rounded-md shadow-sm disabled:opacity-50 w-full sm:w-auto"
          >
            {isFinalizing && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isFinalizing ? 'Finalizing...' : (boardData ? 'Month Finalized' : 'Finalize Current Month')}
          </Button>
        </CardContent>
      </Card>

      {error && (
        <Alert variant="destructive" className="max-w-xl mx-auto">
          <AlertTitle className="font-semibold">Error</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      {finalizeMessage && (
         <Alert variant={error ? "destructive" : "default"} className={`max-w-xl mx-auto ${error ? 'bg-red-50 border-red-500 text-red-700' : 'bg-green-50 border-green-500 text-green-700'}`}>
          <AlertTitle className="font-semibold">{error ? 'Finalization Failed' : 'Success'}</AlertTitle>
          <AlertDescription>{finalizeMessage}</AlertDescription>
        </Alert>
      )}

      {isLoading && budgetLines.length > 0 && ( // Show spinner overlay if loading more data
        <div className="fixed inset-0 bg-white bg-opacity-75 flex justify-center items-center z-50">
             <Loader2 className="h-12 w-12 animate-spin text-blue-500" />
        </div>
      )}


      {!isLoading && budgetLines.length === 0 && !error && (
        <Card className="shadow-md">
          <CardContent className="p-6 text-center">
            <h3 className="text-xl font-semibold text-gray-700 mb-2">No Budget Lines Found</h3>
            <p className="text-gray-500">
              No budget lines found for Month ID: {currentMonthId}.
            </p>
            {!boardData?.is_finalized && (
                <p className="text-sm text-gray-500 mt-1">
                You can add budget lines on the 'Manage' page.
                </p>
            )}
          </CardContent>
        </Card>
      )}

      {!isLoading && budgetLines.length > 0 && (
        <Card className="shadow-lg overflow-hidden">
          <Table>
            <TableHeader className="bg-gray-200 ">
              <TableRow>
                <TableHead className="px-6 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">Category</TableHead>
                <TableHead className="px-6 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">Budget Line Item</TableHead>
                <TableHead className="px-6 py-3 text-right text-xs font-medium text-gray-600 uppercase tracking-wider">Expected</TableHead>
                <TableHead className="px-6 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">Actual</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="bg-white divide-y divide-gray-200">
              {budgetLines.map(line => (
                <TableRow key={line.id} className={`hover:bg-gray-50 ${line.actual_amount && Number(line.actual_amount) > 0 ? 'bg-green-50' : 'bg-yellow-50'}`}>
                  <TableCell className="px-6 py-4 whitespace-nowrap">
                    <CategoryBadge category={{ name: line.category_name, color: line.category_color || 'gray' }} />
                  </TableCell>
                  <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-gray-800">{line.label}</TableCell>
                  <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-gray-800 text-right">{formatCurrency(line.expected_amount)}</TableCell>
                  <TableCell className="px-6 py-4 whitespace-nowrap">
                    <Input
                      type="number"
                      defaultValue={Number(line.actual_amount) || ''}
                      onBlur={(e: ChangeEvent<HTMLInputElement>) => 
                        handleActualAmountChange(line.id, line.actual_id, e.target.value)
                      }
                      className="border-gray-300 rounded-md shadow-sm px-3 py-2 focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm w-full disabled:bg-gray-100"
                      placeholder="0.00"
                      min="0"
                      step="0.01"
                      disabled={boardData?.is_finalized || isFinalizing}
                    />
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Card>
      )}
    </div>
  );
}

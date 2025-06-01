import React, { useState, useEffect, ChangeEvent } from 'react';
import * as api from '../lib/api';
import { textMutedClasses } from '../styles/commonClasses';
import Card from '../components/ui/Card';
import Button from '../components/ui/Button';
import Input from '../components/ui/Input';
import LoadingSpinner from '../components/ui/LoadingSpinner';
import MessageDisplay from '../components/ui/MessageDisplay';
import CategoryBadge from '../components/CategoryBadge';

export default function BoardPage() {
  const [boardData, setBoardData] = useState<api.BoardDataPayload | null>(null);
  const [currentMonthId, setCurrentMonthId] = useState<number>(1);
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
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch board data');
      setBoardData(null);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchBoardData(currentMonthId);
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
      setCurrentMonthId(response.new_month_id);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to finalize month. Check if all actuals are set or an error occurred.";
      setError(errorMessage);
    } finally {
      setIsFinalizing(false);
    }
  };

  const handleActualAmountChange = async (
    budgetLineId: number,
    newActualString: string
  ) => {
    const newActual = parseFloat(newActualString);
    if (isNaN(newActual) || newActual < 0) {
      alert("Please enter a valid positive number for the actual amount.");
      fetchBoardData(currentMonthId);
      return;
    }

    if (boardData) {
      setBoardData({
        ...boardData,
        budget_lines: boardData.budget_lines.map(line =>
          line.id === budgetLineId ? { ...line, actual_amount: newActual } : line
        ),
      });
    }

    try {
      await api.updateActualLine(budgetLineId, { actual: newActual });
    } catch (err) {
      alert(`Failed to update actual amount: ${err instanceof Error ? err.message : 'Unknown error'}`);
      setError(err instanceof Error ? `Failed to update: ${err.message}` : 'Failed to update actual amount.');
      fetchBoardData(currentMonthId);
    }
  };

  const getRowColor = (line: api.BudgetLineWithActual): string => {
    if (line.actual_amount && Number(line.actual_amount) > 0) {
      return 'bg-green-700 hover:bg-green-600';
    }
    return 'bg-yellow-700 hover:bg-yellow-600';
  };

  if (isLoading && !boardData) return <LoadingSpinner text="Loading board data..." />;

  return (
    <div className="p-4 bg-gray-900 min-h-screen text-white">
      <h1 className="text-2xl font-bold mb-6 text-center">
        Monthly Budget Board: {boardData ? `${boardData.month_name} ${boardData.year}` : `Month ID ${currentMonthId}`}
      </h1>
      {boardData?.is_finalized && (
        <p className="text-center text-yellow-500 mb-4">(This month is finalized and read-only)</p>
      )}

      <Card className="mb-6 flex flex-col sm:flex-row justify-between items-center space-y-4 sm:space-y-0">
        <div className="flex items-center space-x-2">
          <label htmlFor="month_id_selector_board" className="block text-sm font-medium">Select Month ID:</label>
          <Input
            id="month_id_selector_board"
            type="number"
            value={currentMonthId}
            onChange={(e) => setCurrentMonthId(parseInt(e.target.value, 10) || 1)}
            className="w-24 !text-black text-sm"
            min="1"
          />
        </div>
        <div>
          <Button
            onClick={handleFinalizeMonth}
            disabled={isFinalizing || isLoading || !boardData?.budget_lines || boardData.budget_lines.length === 0 || boardData?.is_finalized}
            className={`text-sm ${boardData?.is_finalized ? 'bg-gray-500' : 'bg-green-600 hover:bg-green-700'}`}
          >
            {boardData?.is_finalized ? 'Month Finalized' : (isFinalizing ? 'Finalizing...' : 'Finalize Current Month')}
          </Button>
        </div>
      </Card>

      <MessageDisplay message={error} type="error" className="my-4 text-center" />
      <MessageDisplay message={finalizeMessage} type={error ? 'error' : 'success'} className="my-4 text-center" />

      {isLoading && <LoadingSpinner text="Loading board data..." />}

      {!isLoading && !error && (!boardData?.budget_lines || boardData.budget_lines.length === 0) && (
        <Card className="text-center">
          <p className={textMutedClasses}>No budget lines found for Month ID: {currentMonthId}.</p>
          <p className={textMutedClasses}>You can add budget lines in the 'Manage' page for this month if it's not finalized.</p>
        </Card>
      )}

      {!isLoading && boardData?.budget_lines && boardData.budget_lines.length > 0 && (
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
              {boardData?.budget_lines.map(line => (
                <tr key={line.id} className={`${getRowColor(line)} transition-colors duration-150`}>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <CategoryBadge category={{ id: line.category_id, name: line.category_name, color: line.category_color || 'bg-gray-500' }} />
                  </td>
                  <td className="px-4 py-3 whitespace-nowrap">{line.label}</td>
                  <td className="px-4 py-3 whitespace-nowrap text-right">{Number(line.expected_amount).toFixed(0)}</td>
                  <td className="px-4 py-3 whitespace-nowrap">
                    <Input
                      type="number"
                      defaultValue={Number(line.actual_amount) || 0}
                      onBlur={(e: ChangeEvent<HTMLInputElement>) => 
                        handleActualAmountChange(line.id, e.target.value)
                      }
                      className="!text-black w-full text-sm"
                      placeholder="0"
                      min="0"
                      step="1"
                      disabled={boardData?.is_finalized}
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

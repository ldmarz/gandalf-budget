import React, { useState } from 'react';
import { Button } from '../components/ui/Button'; // Assuming a Tailwind-styled Button component
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/Card'; // Assuming Tailwind-styled Card
import { DownloadCloudIcon, CheckCircleIcon, AlertTriangleIcon } from 'lucide-react'; // Icons

export default function BackupPage() {
  const [lastBackup, setLastBackup] = useState<string | null>(localStorage.getItem('lastBackupTimestamp'));
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleExportJson = async () => {
    setIsLoading(true);
    setError(null);
    setSuccessMessage(null);

    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 1000));

    try {
      // Simulate data fetching/preparation
      const dataToExport = {
        appName: "MyBudgetApp",
        version: "1.0.0",
        timestamp: new Date().toISOString(),
        userSettings: { theme: "dark", notifications: true },
        categories: [
          { id: 1, name: "Groceries", color: "bg-green-500" },
          { id: 2, name: "Utilities", color: "bg-blue-500" },
        ],
        budgetLines: [
          { id: 101, month_id: 1, category_id: 1, label: "Aldi Shopping", expected: 200.00, actual_amount: 185.50 },
          { id: 102, month_id: 1, category_id: 2, label: "Electricity Bill", expected: 75.00, actual_amount: 72.10 },
        ],
        // Add more data as needed
      };

      const jsonString = JSON.stringify(dataToExport, null, 2);
      const blob = new Blob([jsonString], { type: "application/json" });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      const timestamp = new Date();
      const formattedDate = `${timestamp.getFullYear()}${(timestamp.getMonth() + 1).toString().padStart(2, '0')}${timestamp.getDate().toString().padStart(2, '0')}`;
      const formattedTime = `${timestamp.getHours().toString().padStart(2, '0')}${timestamp.getMinutes().toString().padStart(2, '0')}`;
      link.download = `mybudgetapp_backup_${formattedDate}_${formattedTime}.json`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);

      const newLastBackup = new Date().toLocaleString();
      setLastBackup(newLastBackup);
      localStorage.setItem('lastBackupTimestamp', newLastBackup);
      setSuccessMessage("Export successful! Your download should start automatically.");

    } catch (err) {
      console.error("Export failed:", err);
      setError(err instanceof Error ? err.message : "An unknown error occurred during export.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-[calc(100vh-64px)] bg-gray-100 flex flex-col items-center justify-center p-4 sm:p-6 lg:p-8">
      {/* Assuming 64px is approx header height, adjust if navbar height is known and fixed */}
      <Card className="bg-white shadow-xl rounded-xl p-6 sm:p-8 lg:p-10 max-w-lg w-full text-center">
        <CardHeader className="border-b-0 pb-4">
          <div className="flex justify-center mb-4">
            <DownloadCloudIcon className="h-16 w-16 text-blue-600" />
          </div>
          <CardTitle className="text-3xl font-bold text-gray-800">
            Export Your Data
          </CardTitle>
          <CardDescription className="text-gray-600 mt-2 text-base">
            Download a JSON file containing all your application data. Keep it safe!
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <Button
            onClick={handleExportJson}
            disabled={isLoading}
            className="w-full py-3 px-6 bg-green-600 hover:bg-green-700 text-white rounded-lg text-lg font-semibold shadow-md hover:shadow-lg transition duration-150 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50 disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center"
          >
            {isLoading ? (
              <Loader2 className="mr-2 h-5 w-5 animate-spin" />
            ) : (
              <DownloadCloudIcon className="mr-2 h-5 w-5" />
            )}
            {isLoading ? 'Exporting...' : 'Export All Data as JSON'}
          </Button>

          {error && (
            <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-md text-red-700 flex items-center text-sm">
              <AlertTriangleIcon className="h-5 w-5 mr-2 shrink-0" />
              <span>{error}</span>
            </div>
          )}

          {successMessage && !error && (
            <div className="mt-4 p-3 bg-green-50 border border-green-200 rounded-md text-green-700 flex items-center text-sm">
              <CheckCircleIcon className="h-5 w-5 mr-2 shrink-0" />
              <span>{successMessage}</span>
            </div>
          )}

          {lastBackup && (
            <p className="text-xs text-gray-500 mt-6">
              Last backup created: {lastBackup}
            </p>
          )}
           {!lastBackup && (
            <p className="text-xs text-gray-500 mt-6">
              No backup has been created yet.
            </p>
          )}
        </CardContent>
      </Card>
       <p className="text-center text-xs text-gray-500 mt-8 max-w-md">
        Note: This is a browser-based backup. For robust data protection, consider server-side backup solutions if applicable to your application architecture.
      </p>
    </div>
  );
}

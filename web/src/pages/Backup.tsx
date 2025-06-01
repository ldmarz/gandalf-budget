import React, { useState, useEffect } from 'react';
import Button from '../components/ui/Button';
import Card from '../components/ui/Card';
import { textMutedClasses } from '../styles/commonClasses';

const LAST_BACKUP_TIMESTAMP_KEY = 'lastBackupTimestamp';

export default function BackupPage() {
  const [lastBackupDateISO, setLastBackupDateISO] = useState<string | null>(null);

  useEffect(() => {
    const timestamp = localStorage.getItem(LAST_BACKUP_TIMESTAMP_KEY);
    if (timestamp) {
      const date = new Date(parseInt(timestamp, 10));
      setLastBackupDateISO(date.toISOString());
    }
  }, []);

  const handleExportJson = () => {
    // Trigger the download
    window.location.href = '/api/v1/export/json';

    // Update localStorage
    const now = new Date();
    localStorage.setItem(LAST_BACKUP_TIMESTAMP_KEY, now.getTime().toString());
    setLastBackupDateISO(now.toISOString()); // Update UI immediately with ISO string
  };

  const calculateDaysAgo = (isoDateString: string | null): string => {
    if (!isoDateString) return 'never';

    const backupDate = new Date(isoDateString);
    const today = new Date();

    // Reset time part for accurate day difference
    backupDate.setHours(0, 0, 0, 0);
    today.setHours(0, 0, 0, 0);

    const diffTime = today.getTime() - backupDate.getTime(); // Ensure today is later or equal
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 0) return 'today';
    if (diffDays === 1) return 'yesterday';
    if (diffDays < 0) return 'in the future (check clock?)'; // Should not happen
    return `${diffDays} days ago`;
  };

  const displayLastBackupInfo = () => {
    if (!lastBackupDateISO) {
      return 'No backup performed yet.';
    }
    const dateObj = new Date(lastBackupDateISO);
    return `Last backup: ${dateObj.toLocaleDateString()} (${calculateDaysAgo(lastBackupDateISO)})`;
  };

  return (
    <div className="p-4 bg-gray-900 min-h-screen text-white">
      <h1 className="text-2xl font-bold mb-6 text-center">Backup</h1>
      <Card className="max-w-md mx-auto">
        <div className="p-6">
          <h2 className="text-xl font-semibold mb-4 text-gray-200">Export Data</h2>
          <p className={`mb-4 ${textMutedClasses}`}>
            Download all your budget data as a JSON file. This file can be used for manual backups.
          </p>
          <Button
            onClick={handleExportJson}
            className="w-full bg-blue-600 hover:bg-blue-700"
          >
            Export JSON
          </Button>
          <p className={`mt-4 text-sm ${textMutedClasses}`}>
            {displayLastBackupInfo()}
          </p>
        </div>
      </Card>
    </div>
  );
}

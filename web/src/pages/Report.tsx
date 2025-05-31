import { useState } from 'react';
import { getAnnualSnapshots, getSnapshotDetail, AnnualSnapMeta, DashboardPayload } from '../lib/api';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../components/ui/table';
import { Alert, AlertDescription, AlertTitle } from '../components/ui/alert';
import { Loader2 } from 'lucide-react';
import ReadOnlyDashboardView from '../components/ReadOnlyDashboardView';

const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A';
  try {
    const date = new Date(dateString.endsWith('Z') ? dateString : dateString + 'Z');
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZone: 'UTC'
    });
  } catch (e) {
    console.warn("Failed to parse date:", dateString, e);
    return dateString;
  }
};


export default function ReportPage() {
  const [yearInput, setYearInput] = useState<string>(new Date().getFullYear().toString());
  const [selectedYear, setSelectedYear] = useState<number | null>(null); // To store the year for which snapshots are loaded

  const [snapshots, setSnapshots] = useState<AnnualSnapMeta[]>([]);
  const [loadingSnapshots, setLoadingSnapshots] = useState<boolean>(false);
  const [errorSnapshots, setErrorSnapshots] = useState<string | null>(null);

  const [selectedSnapId, setSelectedSnapId] = useState<number | null>(null);
  const [snapshotDetail, setSnapshotDetail] = useState<DashboardPayload | null>(null);
  const [loadingDetail, setLoadingDetail] = useState<boolean>(false);
  const [errorDetail, setErrorDetail] = useState<string | null>(null);

  const handleLoadReports = async () => {
    setErrorSnapshots(null);
    setErrorDetail(null);
    setSnapshotDetail(null);
    setSelectedSnapId(null);
    setSnapshots([]);

    const yearNum = parseInt(yearInput, 10);
    if (isNaN(yearNum) || yearInput.length !== 4) {
      setErrorSnapshots("Please enter a valid four-digit year.");
      return;
    }
    setSelectedYear(yearNum);
    setLoadingSnapshots(true);
    try {
      const data = await getAnnualSnapshots(yearNum);
      setSnapshots(data);
      if (data.length === 0) {
        setErrorSnapshots(`No snapshots found for the year ${yearNum}.`);
      }
    } catch (err) {
      setErrorSnapshots(err instanceof Error ? err.message : 'Failed to load snapshots.');
      console.error("Failed to fetch annual snapshots:", err);
    } finally {
      setLoadingSnapshots(false);
    }
  };

  const handleViewSnapshot = async (snapId: number) => {
    setSelectedSnapId(snapId);
    setSnapshotDetail(null);
    setErrorDetail(null);
    setLoadingDetail(true);
    try {
      const data = await getSnapshotDetail(snapId);
      setSnapshotDetail(data);
    } catch (err) {
      setErrorDetail(err instanceof Error ? err.message : 'Failed to load snapshot detail.');
      console.error(`Failed to fetch snapshot detail for ID ${snapId}:`, err);
    } finally {
      setLoadingDetail(false);
    }
  };

  return (
    <div className="container mx-auto p-4 sm:p-6 lg:p-8 space-y-8 bg-gray-100 min-h-screen">
      <Card className="shadow-lg">
        <CardHeader>
          <CardTitle className="text-3xl font-bold text-gray-800">Annual Reports Viewer</CardTitle>
          <CardDescription className="text-gray-600">Select a year to view available monthly snapshots.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4 pt-6">
          <div className="flex flex-col sm:flex-row sm:space-x-3 sm:items-end">
            <div className="flex-grow sm:flex-grow-0 sm:w-48">
              <label htmlFor="yearInputReport" className="block text-sm font-medium text-gray-700 mb-1">Year</label>
              <Input
                id="yearInputReport"
                type="number"
                placeholder="e.g., 2023"
                value={yearInput}
                onChange={(e) => setYearInput(e.target.value)}
                className="w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                min="2000"
                max="2099"
              />
            </div>
            <Button
              onClick={handleLoadReports}
              disabled={loadingSnapshots}
              className="mt-2 sm:mt-0 w-full sm:w-auto bg-blue-600 hover:bg-blue-700 text-white font-semibold py-2 px-4 rounded-md shadow-sm disabled:opacity-50 flex items-center justify-center"
            >
              {loadingSnapshots && <Loader2 className="mr-2 h-5 w-5 animate-spin" />}
              Load Snapshots
            </Button>
          </div>
          {errorSnapshots && !loadingSnapshots && (
             <Alert variant={snapshots.length > 0 && errorSnapshots !== `No snapshots found for the year ${selectedYear}.` ? "destructive" : "info"} className="mt-4 max-w-xl mx-auto">
              <AlertTitle className="font-semibold">{snapshots.length > 0 && errorSnapshots !== `No snapshots found for the year ${selectedYear}.` ? "Error" : "Information"}</AlertTitle>
              <AlertDescription>{errorSnapshots}</AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>

      {loadingSnapshots && (
        <div className="flex flex-col justify-center items-center py-20 text-center">
          <Loader2 className="h-12 w-12 animate-spin text-blue-500 mb-4" />
          <p className="text-xl text-gray-700">Loading snapshots...</p>
        </div>
      )}

      {!loadingSnapshots && snapshots.length > 0 && selectedYear && (
        <Card className="shadow-lg">
          <CardHeader>
            <CardTitle className="text-2xl font-semibold text-gray-800">Snapshots for {selectedYear}</CardTitle>
            <CardDescription className="text-gray-600">Select a snapshot to view its details.</CardDescription>
          </CardHeader>
          <CardContent className="overflow-x-auto">
            <Table>
              <TableHeader className="bg-gray-200">
                <TableRow>
                  <TableHead className="px-6 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">Month</TableHead>
                  <TableHead className="px-6 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">Year</TableHead>
                  <TableHead className="px-6 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">Snapshot Taken At (UTC)</TableHead>
                  <TableHead className="px-6 py-3 text-right text-xs font-medium text-gray-600 uppercase tracking-wider">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody className="bg-white divide-y divide-gray-200">
                {snapshots.map((snap) => (
                  <TableRow key={snap.id} className={`hover:bg-gray-50 ${selectedSnapId === snap.id ? "bg-blue-50" : ""}`}>
                    <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{snap.month}</TableCell>
                    <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{snap.year}</TableCell>
                    <TableCell className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{formatDate(snap.snap_created_at)}</TableCell>
                    <TableCell className="px-6 py-4 whitespace-nowrap text-right">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleViewSnapshot(snap.id)}
                        disabled={loadingDetail && selectedSnapId === snap.id}
                        className="bg-indigo-600 hover:bg-indigo-700 text-white text-xs py-1 px-3 rounded-md shadow-sm disabled:opacity-50 flex items-center"
                      >
                        {loadingDetail && selectedSnapId === snap.id && <Loader2 className="mr-1.5 h-3 w-3 animate-spin" />}
                        View Detail
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}

      {loadingDetail && (
        <Card className="mt-6 shadow-md">
          <CardContent className="flex flex-col justify-center items-center py-20 text-center">
            <Loader2 className="h-12 w-12 animate-spin text-blue-500 mb-4" />
            <p className="text-xl text-gray-700">Loading snapshot detail...</p>
          </CardContent>
        </Card>
      )}
      
      {errorDetail && !loadingDetail && (
        <Alert variant="destructive" className="my-6 max-w-xl mx-auto">
          <AlertTitle className="font-semibold">Error Loading Snapshot Detail</AlertTitle>
          <AlertDescription>{errorDetail}</AlertDescription>
        </Alert>
      )}

      {snapshotDetail && !loadingDetail && selectedSnapId && (
        <Card className="mt-6 shadow-lg">
          <CardHeader>
            <CardTitle className="text-2xl font-semibold text-gray-800">
              Snapshot Detail: {snapshotDetail.month} {snapshotDetail.year}
            </CardTitle>
            <CardDescription className="text-gray-600">
              Displaying the content of the selected snapshot taken on {formatDate(snapshots.find(s=>s.id === selectedSnapId)?.snap_created_at || '')}.
            </CardDescription>
          </CardHeader>
          <CardContent className="p-0 sm:p-2 md:p-4"> {/* Adjust padding for ReadOnlyDashboardView if needed */}
            <ReadOnlyDashboardView data={snapshotDetail} />
          </CardContent>
        </Card>
      )}
    </div>
  );
}

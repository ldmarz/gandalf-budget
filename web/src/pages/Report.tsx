import { useState } from 'react';
import { getAnnualSnapshots, getSnapshotDetail, AnnualSnapMeta, DashboardPayload } from '../lib/api';
import { Button } from '../components/ui/button';
import { Input } from '../components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../components/ui/table';
import { Alert, AlertDescription, AlertTitle } from '../components/ui/alert';
import { Loader2 } from 'lucide-react';
import ReadOnlyDashboardView from '../components/ReadOnlyDashboardView'; // Added import

// Helper to format date strings
const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A';
  try {
    // Attempt to handle both ISO string with Z and without
    const date = new Date(dateString.endsWith('Z') ? dateString : dateString + 'Z');
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      timeZone: 'UTC' // Assuming created_at is UTC, display as such or convert to local
    });
  } catch (e) {
    console.warn("Failed to parse date:", dateString, e);
    return dateString; // Fallback if parsing fails
  }
};


export default function ReportPage() {
  const [yearInput, setYearInput] = useState<string>(new Date().getFullYear().toString());
  const [selectedYear, setSelectedYear] = useState<number | null>(null);
  
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
    if (isNaN(yearNum) || yearInput.length !== 4) { // Basic year format validation
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
    setSnapshotDetail(null); // Clear previous detail
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
    <div className="container mx-auto p-4 space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl">Annual Reports Viewer</CardTitle>
          <CardDescription>Select a year to view available monthly snapshots.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex space-x-2 items-center">
            <Input
              type="number"
              placeholder="Enter Year (e.g., 2023)"
              value={yearInput}
              onChange={(e) => setYearInput(e.target.value)}
              className="max-w-xs"
              min="2000" // Reasonable min year
              max="2099" // Reasonable max year
            />
            <Button onClick={handleLoadReports} disabled={loadingSnapshots}>
              {loadingSnapshots && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Load Snapshots
            </Button>
          </div>
          {errorSnapshots && !loadingSnapshots && (
             <Alert variant={snapshots.length > 0 ? "destructive" : "info"}>
              <AlertTitle>{snapshots.length > 0 ? "Error" : "Information"}</AlertTitle>
              <AlertDescription>{errorSnapshots}</AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>

      {loadingSnapshots && (
        <div className="flex justify-center items-center py-10">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="ml-2">Loading snapshots...</p>
        </div>
      )}

      {!loadingSnapshots && snapshots.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>Snapshots for {selectedYear}</CardTitle>
            <CardDescription>Select a snapshot to view its details.</CardDescription>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Month</TableHead>
                  <TableHead>Year</TableHead>
                  <TableHead>Snapshot Taken At (UTC)</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {snapshots.map((snap) => (
                  <TableRow key={snap.id} className={selectedSnapId === snap.id ? "bg-muted/50" : ""}>
                    <TableCell>{snap.month}</TableCell>
                    <TableCell>{snap.year}</TableCell>
                    <TableCell>{formatDate(snap.snap_created_at)}</TableCell>
                    <TableCell className="text-right">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleViewSnapshot(snap.id)}
                        disabled={loadingDetail && selectedSnapId === snap.id}
                      >
                        {loadingDetail && selectedSnapId === snap.id ? (
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        ) : null}
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
        <Card className="mt-6">
          <CardContent className="flex justify-center items-center py-10">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
            <p className="ml-2">Loading snapshot detail...</p>
          </CardContent>
        </Card>
      )}
      
      {errorDetail && !loadingDetail && (
        <Alert variant="destructive" className="my-4">
          <AlertTitle>Error Loading Snapshot Detail</AlertTitle>
          <AlertDescription>{errorDetail}</AlertDescription>
        </Alert>
      )}

      {snapshotDetail && !loadingDetail && (
        <Card className="mt-6">
          <CardHeader>
            {/* Title is now part of ReadOnlyDashboardView, but we can keep a general title here or remove */}
            {/* <CardTitle>Snapshot Detail: {snapshotDetail.month} {snapshotDetail.year}</CardTitle> */}
             <CardDescription>
              Displaying the content of the selected snapshot.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ReadOnlyDashboardView data={snapshotDetail} />
          </CardContent>
        </Card>
      )}
    </div>
  );
}

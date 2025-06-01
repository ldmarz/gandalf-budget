import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { getDashboardData, DashboardPayload, CategorySummary, BudgetLineDetail } from '../lib/api';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../components/ui/Card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../components/ui/Table';
import { Badge } from '../components/ui/Badge';
import { CategoryBadge } from '../components/CategoryBadge';
import { Alert, AlertDescription, AlertTitle } from '../components/ui/Alert';
import { Loader2 } from 'lucide-react';
import { formatCurrency } from '../lib/utils';

export default function Dashboard() {
  const [searchParams] = useSearchParams();
  const monthId = searchParams.get('month_id');

  const [dashboardData, setDashboardData] = useState<DashboardPayload | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!monthId) {
      setError('Month ID is required. Please select a month.');
      setLoading(false);
      return;
    }

    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await getDashboardData(monthId);
        setDashboardData(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An unknown error occurred.');
        console.error("Failed to fetch dashboard data:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [monthId]);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <p className="ml-2 text-lg">Loading dashboard...</p>
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive" className="max-w-2xl mx-auto my-4">
        <AlertTitle>Error</AlertTitle>
        <AlertDescription>{error}</AlertDescription>
      </Alert>
    );
  }

  if (!dashboardData) {
    return (
      <div className="text-center py-10">
        <p>No dashboard data available for the selected month.</p>
         {!monthId && <p>Please ensure a month is selected or a `month_id` is provided in the URL.</p>}
      </div>
    );
  }

  const {
    month,
    year,
    total_expected,
    total_actual,
    total_difference,
    category_summaries,
  } = dashboardData;

  return (
    <div className="container mx-auto p-4 space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-3xl font-bold">Dashboard for {month} {year}</CardTitle>
          <CardDescription>Overview of your budget and spending for the month.</CardDescription>
        </CardHeader>
        <CardContent className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card className="bg-blue-50 dark:bg-blue-900">
            <CardHeader>
              <CardTitle className="text-xl">Total Expected</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-semibold">{formatCurrency(total_expected)}</p>
            </CardContent>
          </Card>
          <Card className="bg-green-50 dark:bg-green-900">
            <CardHeader>
              <CardTitle className="text-xl">Total Actual</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-semibold">{formatCurrency(total_actual)}</p>
            </CardContent>
          </Card>
          <Card className={total_difference >= 0 ? "bg-yellow-50 dark:bg-yellow-900" : "bg-red-50 dark:bg-red-900"}>
            <CardHeader>
              <CardTitle className="text-xl">Total Difference</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-2xl font-semibold">{formatCurrency(total_difference)}</p>
              <Badge variant={total_difference >= 0 ? 'default' : 'destructive'} className="mt-1">
                {total_difference >= 0 ? 'Under Budget' : 'Over Budget'}
              </Badge>
            </CardContent>
          </Card>
        </CardContent>
      </Card>

      {category_summaries.length === 0 && (
        <Card>
          <CardContent>
            <p className="text-center py-4">No category summaries available for this month.</p>
          </CardContent>
        </Card>
      )}

      {category_summaries.map((summary: CategorySummary) => (
        <Card key={summary.category_id}>
          <CardHeader className="flex flex-row justify-between items-center">
            <div>
              <CategoryBadge category={{ name: summary.category_name, color: summary.category_color, id: summary.category_id }} />
              <CardDescription className="mt-1">
                Summary for {summary.category_name}
              </CardDescription>
            </div>
            <div className="text-right">
                 <p className="text-sm text-muted-foreground">Difference</p>
                 <p className={`text-lg font-semibold ${summary.difference >=0 ? 'text-green-600' : 'text-red-600'}`}>
                    {formatCurrency(summary.difference)}
                 </p>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
              <div className="p-3 bg-gray-50 dark:bg-gray-800 rounded-md">
                <p className="text-sm text-muted-foreground">Expected</p>
                <p className="text-lg font-medium">{formatCurrency(summary.total_expected)}</p>
              </div>
              <div className="p-3 bg-gray-50 dark:bg-gray-800 rounded-md">
                <p className="text-sm text-muted-foreground">Actual</p>
                <p className="text-lg font-medium">{formatCurrency(summary.total_actual)}</p>
              </div>
            </div>

            {summary.budget_lines.length > 0 ? (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Budget Line Item</TableHead>
                    <TableHead className="text-right">Expected</TableHead>
                    <TableHead className="text-right">Actual</TableHead>
                    <TableHead className="text-right">Difference</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {summary.budget_lines.map((line: BudgetLineDetail) => (
                    <TableRow key={line.budget_line_id}>
                      <TableCell>{line.label}</TableCell>
                      <TableCell className="text-right">{formatCurrency(line.expected_amount)}</TableCell>
                      <TableCell className="text-right">{formatCurrency(line.actual_amount)}</TableCell>
                      <TableCell className={`text-right font-medium ${line.difference >=0 ? 'text-green-700 dark:text-green-500' : 'text-red-700 dark:text-red-500'}`}>
                        {formatCurrency(line.difference)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            ) : (
              <p className="text-sm text-muted-foreground text-center py-3">No budget lines in this category.</p>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

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
      <div className="flex flex-col justify-center items-center h-64 p-4 text-center">
        <Loader2 className="h-12 w-12 animate-spin text-blue-500 mb-4" />
        <p className="text-xl text-gray-700">Loading dashboard...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex justify-center items-center h-64 p-4">
        <Alert variant="destructive" className="max-w-2xl mx-auto">
          <AlertTitle className="text-xl font-semibold">Error</AlertTitle>
          <AlertDescription className="text-gray-700">{error}</AlertDescription>
        </Alert>
      </div>
    );
  }

  if (!dashboardData) {
    return (
      <div className="text-center py-10 p-4">
        <p className="text-xl text-gray-700 mb-2">No dashboard data available for the selected month.</p>
         {!monthId && <p className="text-sm text-gray-500">Please ensure a month is selected or a `month_id` is provided in the URL.</p>}
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
    <div className="container mx-auto p-6 lg:p-8 space-y-8">
      <Card className="shadow-lg">
        <CardHeader>
          <CardTitle className="text-3xl font-bold text-gray-800">Dashboard for {month} {year}</CardTitle>
          <CardDescription className="text-gray-600">Overview of your budget and spending for the month.</CardDescription>
        </CardHeader>
        <CardContent className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card className="bg-blue-50 dark:bg-blue-900 border border-blue-200 dark:border-blue-700 rounded-lg">
            <CardHeader>
              <CardTitle className="text-xl text-blue-700 dark:text-blue-300">Total Expected</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-semibold text-blue-900 dark:text-blue-100">{formatCurrency(total_expected)}</p>
            </CardContent>
          </Card>
          <Card className="bg-green-50 dark:bg-green-900 border border-green-200 dark:border-green-700 rounded-lg">
            <CardHeader>
              <CardTitle className="text-xl text-green-700 dark:text-green-300">Total Actual</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-semibold text-green-900 dark:text-green-100">{formatCurrency(total_actual)}</p>
            </CardContent>
          </Card>
          <Card className={`${total_difference >= 0 ? "bg-yellow-50 dark:bg-yellow-900 border-yellow-200 dark:border-yellow-700" : "bg-red-50 dark:bg-red-900 border-red-200 dark:border-red-700"} rounded-lg`}>
            <CardHeader>
              <CardTitle className={`text-xl ${total_difference >= 0 ? "text-yellow-700 dark:text-yellow-300" : "text-red-700 dark:text-red-300"}`}>Total Difference</CardTitle>
            </CardHeader>
            <CardContent>
              <p className={`text-3xl font-semibold ${total_difference >= 0 ? "text-yellow-900 dark:text-yellow-100" : "text-red-900 dark:text-red-100"}`}>{formatCurrency(total_difference)}</p>
              <Badge variant={total_difference >= 0 ? 'default' : 'destructive'} className="mt-2 text-sm">
                {total_difference >= 0 ? 'Under Budget' : 'Over Budget'}
              </Badge>
            </CardContent>
          </Card>
        </CardContent>
      </Card>

      {category_summaries.length === 0 && (
        <Card className="shadow">
          <CardContent>
            <p className="text-center py-6 text-gray-600 text-lg">No category summaries available for this month.</p>
          </CardContent>
        </Card>
      )}

      {category_summaries.map((summary: CategorySummary) => (
        <Card key={summary.category_id} className="shadow-lg">
          <CardHeader className="flex flex-row justify-between items-start pb-4 border-b">
            <div>
              <CategoryBadge category={{ name: summary.category_name, color: summary.category_color, id: summary.category_id }} />
              <CardDescription className="mt-1 text-sm text-gray-500">
                Summary for {summary.category_name}
              </CardDescription>
            </div>
            <div className="text-right">
                 <p className="text-sm text-gray-500">Difference</p>
                 <p className={`text-xl font-semibold ${summary.difference >=0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>
                    {formatCurrency(summary.difference)}
                 </p>
            </div>
          </CardHeader>
          <CardContent className="pt-6">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 mb-6">
              <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">Expected</p>
                <p className="text-xl font-medium text-gray-800 dark:text-gray-200">{formatCurrency(summary.total_expected)}</p>
              </div>
              <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">Actual</p>
                <p className="text-xl font-medium text-gray-800 dark:text-gray-200">{formatCurrency(summary.total_actual)}</p>
              </div>
            </div>

            {summary.budget_lines.length > 0 ? (
              <Table>
                <TableHeader>
                  <TableRow className="bg-gray-50 dark:bg-gray-800">
                    <TableHead className="text-gray-700 dark:text-gray-300 font-semibold">Budget Line Item</TableHead>
                    <TableHead className="text-right text-gray-700 dark:text-gray-300 font-semibold">Expected</TableHead>
                    <TableHead className="text-right text-gray-700 dark:text-gray-300 font-semibold">Actual</TableHead>
                    <TableHead className="text-right text-gray-700 dark:text-gray-300 font-semibold">Difference</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {summary.budget_lines.map((line: BudgetLineDetail) => (
                    <TableRow key={line.budget_line_id} className="hover:bg-gray-100 dark:hover:bg-gray-700">
                      <TableCell className="text-gray-700 dark:text-gray-300">{line.label}</TableCell>
                      <TableCell className="text-right text-gray-600 dark:text-gray-400">{formatCurrency(line.expected_amount)}</TableCell>
                      <TableCell className="text-right text-gray-600 dark:text-gray-400">{formatCurrency(line.actual_amount)}</TableCell>
                      <TableCell className={`text-right font-medium ${line.difference >=0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>
                        {formatCurrency(line.difference)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            ) : (
              <p className="text-sm text-gray-500 text-center py-4">No budget lines in this category.</p>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

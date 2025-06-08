import { DashboardPayload, CategorySummary, BudgetLineDetail } from '../lib/api';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from './ui/Card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from './ui/Table';
import { Badge } from './ui/Badge';
import { CategoryBadge } from './CategoryBadge';
import { formatCurrency } from '../lib/utils';

interface ReadOnlyDashboardViewProps {
  data: DashboardPayload;
}

export default function ReadOnlyDashboardView({ data }: ReadOnlyDashboardViewProps) {
  if (!data) {
    return <p>No data provided to display.</p>;
  }

  const {
    month,
    year,
    total_expected,
    total_actual,
    total_difference,
    category_summaries,
  } = data;

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="text-3xl font-bold">Snapshot: {month} {year}</CardTitle>
          <CardDescription>Read-only view of budget and spending for this period.</CardDescription>
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

      {(!category_summaries || category_summaries.length === 0) && (
        <Card>
          <CardContent>
            <p className="text-center py-4 text-muted-foreground">No category summaries available for this snapshot.</p>
          </CardContent>
        </Card>
      )}

      {category_summaries && category_summaries.map((summary: CategorySummary) => (
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
                 <p className={`text-lg font-semibold ${summary.difference >=0 ? 'text-green-600 dark:text-green-400' : 'text-red-600 dark:text-red-400'}`}>
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

            {(!summary.budget_lines || summary.budget_lines.length === 0) ? (
              <p className="text-sm text-muted-foreground text-center py-3">No budget lines in this category for this snapshot.</p>
            ) : (
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
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

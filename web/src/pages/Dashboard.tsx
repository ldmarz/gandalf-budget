import { useState } from 'react';
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts';
import SidebarLayout from '../components/layout/SidebarLayout';
import { Avatar } from '../components/ui/Avatar';
import { Button } from '../components/ui/Button';
import { Card, CardHeader, CardTitle, CardContent, CardDescription } from '../components/ui/Card';
import { Badge } from '../components/ui/Badge';
import { Search, Cog } from 'lucide-react';

const lineData = [
  { name: 'Jan', value: 1000 },
  { name: 'Feb', value: 1200 },
  { name: 'Mar', value: 1500 },
  { name: 'Apr', value: 1700 },
  { name: 'May', value: 1600 },
  { name: 'Jun', value: 1800 },
  { name: 'Jul', value: 1900 },
  { name: 'Aug', value: 2000 },
  { name: 'Sep', value: 2200 },
  { name: 'Oct', value: 2300 },
  { name: 'Nov', value: 2400 },
  { name: 'Dec', value: 2500 },
];

const investmentData = [
  { name: 'VOO', value: 30, color: '#86efac' },
  { name: 'VTI', value: 25, color: '#60a5fa' },
  { name: 'CLN', value: 20, color: '#fcd34d' },
  { name: 'BTEX', value: 25, color: '#fda4af' },
];

export default function Dashboard() {
  const [range, setRange] = useState('Year');

  return (
    <SidebarLayout>
      <div className="flex items-start justify-between pb-6">
        <div>
          <h1 className="text-2xl font-bold">Welcome back, John</h1>
          <p className="text-sm text-gray-500">Here is an overview of your finances.</p>
        </div>
        <div className="flex items-center space-x-4">
          <Search className="w-5 h-5 text-gray-600" />
          <Cog className="w-5 h-5 text-gray-600" />
          <Avatar src="https://ui.shadcn.com/avatars/01.png" />
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <Card className="lg:col-span-2">
          <CardHeader className="flex justify-between items-center">
            <CardTitle>Account Balance</CardTitle>
            <div className="space-x-2">
              {['Day', 'Week', 'Month', 'Year'].map((label) => (
                <Button
                  key={label}
                  size="sm"
                  variant={label === range ? 'primary' : 'secondary'}
                  onClick={() => setRange(label)}
                >
                  {label}
                </Button>
              ))}
            </div>
          </CardHeader>
          <CardContent className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={lineData} margin={{ top: 20, right: 20, bottom: 20, left: 0 }}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis tickFormatter={(v) => `${v / 1000}k`} domain={[0, 5000]} />
                <Tooltip />
                <Line type="monotone" dataKey="value" stroke="#3b82f6" dot />
              </LineChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>
        <div className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="text-right">Total Balance</CardTitle>
            </CardHeader>
            <CardContent className="text-right">
              <p className="text-2xl font-semibold">$11,716.77</p>
              <Badge variant="success">+3.5%</Badge>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle className="text-right">Main Account</CardTitle>
            </CardHeader>
            <CardContent className="text-right">
              <p className="text-2xl font-semibold">$4,500.00</p>
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle className="text-right">Savings</CardTitle>
            </CardHeader>
            <CardContent className="text-right">
              <p className="text-2xl font-semibold">$7,216.77</p>
              <Badge variant="outline" className="text-blue-600">
                -1.1%
              </Badge>
            </CardContent>
          </Card>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 pt-6">
        <Card className="lg:col-span-2">
          <CardHeader className="flex justify-between items-center">
            <CardTitle>Recent Transactions</CardTitle>
            <a className="text-sm text-blue-600" href="#">
              See all
            </a>
          </CardHeader>
          <CardContent>
            <table className="w-full text-sm">
              <thead>
                <tr className="text-left">
                  <th>Name</th>
                  <th>Date</th>
                  <th>Time</th>
                  <th>Status</th>
                  <th className="text-right">Amount</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200">
                <tr>
                  <td className="flex items-center space-x-2">
                    <Avatar className="w-6 h-6" src="https://ui.shadcn.com/avatars/02.png" />
                    <span>Groceries</span>
                  </td>
                  <td>Today</td>
                  <td>15:45</td>
                  <td>
                    <Badge variant="warning">Pending</Badge>
                  </td>
                  <td className="text-right text-red-600">-$54.33</td>
                </tr>
                <tr>
                  <td className="flex items-center space-x-2">
                    <Avatar className="w-6 h-6" src="https://ui.shadcn.com/avatars/03.png" />
                    <span>Salary</span>
                  </td>
                  <td>11 Dec</td>
                  <td>09:12</td>
                  <td>
                    <Badge variant="success">Completed</Badge>
                  </td>
                  <td className="text-right text-green-600">+$2,500.00</td>
                </tr>
                <tr>
                  <td className="flex items-center space-x-2">
                    <Avatar className="w-6 h-6" src="https://ui.shadcn.com/avatars/04.png" />
                    <span>Utilities</span>
                  </td>
                  <td>10 Dec</td>
                  <td>08:00</td>
                  <td>
                    <Badge variant="success">Completed</Badge>
                  </td>
                  <td className="text-right text-red-600">-$120.00</td>
                </tr>
              </tbody>
            </table>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Investments</CardTitle>
            <CardDescription>
              +$1,203.64 <span className="text-green-600">(+27.24%)</span>
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ResponsiveContainer width="100%" height={200}>
              <PieChart>
                <Pie
                  startAngle={180}
                  endAngle={0}
                  data={investmentData}
                  dataKey="value"
                  cx="50%"
                  cy="100%"
                  innerRadius={60}
                  outerRadius={80}
                >
                  {investmentData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
              </PieChart>
            </ResponsiveContainer>
            <div className="flex justify-around text-xs mt-2">
              {investmentData.map((d) => (
                <div key={d.name} className="flex items-center space-x-1">
                  <span className="w-3 h-3 inline-block" style={{ background: d.color }}></span>
                  <span>{d.name}</span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </SidebarLayout>
  );
}

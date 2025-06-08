import { ReactNode } from 'react';
import { Home, Wallet, CreditCard, LineChart as ChartIcon, FileText, Settings, LogOut } from 'lucide-react';
import { Button } from '../ui/Button';

interface Props { children: ReactNode }

const menu = [
  { label: 'Overview', icon: Home },
  { label: 'Budgets', icon: Wallet },
  { label: 'Expenses', icon: CreditCard },
  { label: 'Investments', icon: ChartIcon },
  { label: 'Reports', icon: FileText },
  { label: 'Settings', icon: Settings },
];

export default function SidebarLayout({ children }: Props) {
  return (
    <div className="flex min-h-screen bg-gray-50">
      <aside className="w-64 bg-gray-900 text-gray-100 flex flex-col">
        <div className="p-4 text-lg font-bold">MyBudget</div>
        <nav className="flex-1 px-2 space-y-1">
          {menu.map((item) => {
            const Icon = item.icon;
            return (
              <button key={item.label} className="w-full flex items-center space-x-2 px-3 py-2 rounded hover:bg-gray-800 text-sm">
                <Icon className="w-4 h-4" />
                <span>{item.label}</span>
              </button>
            );
          })}
        </nav>
        <div className="m-3 p-3 bg-gray-800 rounded">
          <p className="text-sm mb-2">Need help?</p>
          <Button size="sm" className="w-full">Contact us</Button>
        </div>
        <Button variant="ghost" className="justify-start m-3 text-sm">
          <LogOut className="w-4 h-4 mr-2" /> Log out
        </Button>
      </aside>
      <main className="flex-1 p-6">{children}</main>
    </div>
  );
}

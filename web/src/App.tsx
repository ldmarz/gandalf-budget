import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import Board from './pages/Board';
import Report from './pages/Report';
import Manage from './pages/Manage';
import Backup from './pages/Backup';
import Navbar from './components/layout/Navbar';

function App() {
  return (
    <BrowserRouter>
      <Navbar />
      <main className="p-4 text-red-500">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/board" element={<Board />} />
          <Route path="/report" element={<Report />} />
          <Route path="/manage" element={<Manage />} />
          <Route path="/backup" element={<Backup />} />
        </Routes>
      </main>
    </BrowserRouter>
  );
}

export default App;

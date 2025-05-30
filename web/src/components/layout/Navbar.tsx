import { Link } from 'react-router-dom';

const Navbar = () => {
  return (
    <nav className="bg-gray-100 p-4">
      <ul className="flex space-x-4">
        <li>
          <Link to="/" className="text-blue-600 hover:underline">
            Dashboard
          </Link>
        </li>
        <li>
          <Link to="/board" className="text-blue-600 hover:underline">
            Board
          </Link>
        </li>
        <li>
          <Link to="/report" className="text-blue-600 hover:underline">
            Report
          </Link>
        </li>
        <li>
          <Link to="/manage" className="text-blue-600 hover:underline">
            Manage
          </Link>
        </li>
        <li>
          <Link to="/backup" className="text-blue-600 hover:underline">
            Backup
          </Link>
        </li>
      </ul>
    </nav>
  );
};

export default Navbar;

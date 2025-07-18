import { Routes as RouterRoutes, Route, BrowserRouter, Navigate } from 'react-router-dom';
import { useKindeAuth } from '@kinde-oss/kinde-auth-react';
import ProtectedRoute from './components/ProtectedRoute';
import SignIn from './pages/SignIn';
import Home from './pages/Home';
import ThemeToggle from './components/ThemeToggle';
import { Button } from './components/ui/button';

// This is a placeholder for a dashboard page component
const Dashboard = () => {
  const { logout } = useKindeAuth();

  return (
    <div className="min-h-screen flex flex-col">
      {/* Navigation Bar */}
      <nav className="border-b">
        <div className="container mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center">
            <h1 className="text-xl font-bold">RecurrentPaymentTracker</h1>
          </div>
          <div className="flex items-center gap-2">
            <ThemeToggle />
            <Button onClick={logout} variant="outline">
              Sign Out
            </Button>
          </div>
        </div>
      </nav>
      
      {/* Dashboard Content */}
      <div className="container mx-auto p-4">
        <h1 className="text-2xl font-bold mb-4">Dashboard</h1>
        <p>Your payments dashboard will be displayed here.</p>
      </div>
    </div>
  );
};

// Routes component that handles all application routing
const Routes = () => {
  const { isAuthenticated } = useKindeAuth();

  return (
    <BrowserRouter>
      <RouterRoutes>
        <Route path="/" element={<Home />} />
        <Route path="/signin" element={<SignIn />} />
        <Route 
          path="/dashboard" 
          element={
            <ProtectedRoute>
              <Dashboard />
            </ProtectedRoute>
          } 
        />
        {/* Redirect to signin if route doesn't exist */}
        <Route path="*" element={<Navigate to={isAuthenticated ? "/dashboard" : "/signin"} />} />
      </RouterRoutes>
    </BrowserRouter>
  );
};

export default Routes;
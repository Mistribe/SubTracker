import { Routes as RouterRoutes, Route, BrowserRouter, Navigate } from 'react-router-dom';
import { useKindeAuth } from '@kinde-oss/kinde-auth-react';
import ProtectedRoute from './components/ProtectedRoute';
import SignIn from './pages/SignIn';

// This is a placeholder for a home page component
const Home = () => {
  const { user } = useKindeAuth();
  
  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Welcome to Recurrent Payment Tracker</h1>
      {user ? (
        <p>Hello, {user.givenName || user.email}!</p>
      ) : (
        <p>Please log in to manage your payments.</p>
      )}
    </div>
  );
};

// This is a placeholder for a dashboard page component
const Dashboard = () => {
  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Dashboard</h1>
      <p>Your payments dashboard will be displayed here.</p>
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
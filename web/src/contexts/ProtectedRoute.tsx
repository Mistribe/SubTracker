import { ReactNode } from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "./useAuth";

interface ProtectedRouteProps {
  children: ReactNode;
}

export const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { isAuthenticated } = useAuth();
  const storedAuth = localStorage.getItem("isAuthenticated");
  
  // If not authenticated, redirect to login
  if (!isAuthenticated && storedAuth !== "true") {
    return <Navigate to="/login" replace />;
  }
  
  return <>{children}</>;
};
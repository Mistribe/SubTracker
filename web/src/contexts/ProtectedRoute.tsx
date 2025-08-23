import { useKindeAuth } from "@kinde-oss/kinde-auth-react";
import type { ReactNode } from "react";
import { Loader2 } from "lucide-react";

interface ProtectedRouteProps {
  children: ReactNode;
}

export const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { isAuthenticated, isLoading, login } = useKindeAuth();

  // Wait for auth state to initialize to avoid unnecessary redirects
  if (isLoading) {
    return (
      <div className="flex h-[60vh] items-center justify-center text-muted-foreground">
        <Loader2 className="mr-2 h-5 w-5 animate-spin" />
        Loading...
      </div>
    );
  }

  if (!isAuthenticated) {
    login();
    return null;
  }

  return <>{children}</>;
};
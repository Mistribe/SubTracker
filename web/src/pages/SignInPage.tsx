import { useEffect } from "react";
import { SignIn, useUser } from "@clerk/clerk-react";
import { useNavigate, Link } from "react-router-dom";
import { useTheme } from "@/components/theme-provider";
import { AuthLayout } from "@/layouts/AuthLayout";
import { getClerkAppearance } from "@/lib/clerk-theme";

const SignInPage = () => {
  const { isSignedIn } = useUser();
  const navigate = useNavigate();
  const { theme } = useTheme();

  // Redirect authenticated users to dashboard
  useEffect(() => {
    if (isSignedIn) {
      navigate("/dashboard", { replace: true });
    }
  }, [isSignedIn, navigate]);

  // Don't render the sign-in form if user is already authenticated
  if (isSignedIn) {
    return null;
  }

  return (
    <AuthLayout
      title="Welcome back"
      description="Sign in to your account to continue"
    >
      <div className="space-y-6">
        <SignIn
          appearance={getClerkAppearance(theme)}
          fallbackRedirectUrl={"/dashboard"}
          signUpUrl="/sign-up"
          fallback={
            <div className="flex flex-col space-y-4 w-full">
              <div className="h-10 bg-muted animate-pulse rounded" />
              <div className="h-10 bg-muted animate-pulse rounded" />
              <div className="h-10 bg-muted animate-pulse rounded" />
            </div>
          }
        />

        {/* Navigation link to sign-up page */}
        <div className="text-center text-sm">
          <span className="text-muted-foreground">Don't have an account? </span>
          <Link
            to="/sign-up"
            className="font-medium text-primary hover:text-primary/80 transition-colors"
          >
            Sign up
          </Link>
        </div>
      </div>
    </AuthLayout>
  );
};

export default SignInPage;
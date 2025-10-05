import { useEffect } from "react";
import { SignUp, useUser } from "@clerk/clerk-react";
import { useNavigate, Link } from "react-router-dom";
import { useTheme } from "@/components/theme-provider";
import { AuthLayout } from "@/layouts/AuthLayout";
import { getClerkAppearance } from "@/lib/clerk-theme";

const SignUpPage = () => {
  const { isSignedIn } = useUser();
  const navigate = useNavigate();
  const { theme } = useTheme();

  // Redirect authenticated users to dashboard
  useEffect(() => {
    if (isSignedIn) {
      navigate("/dashboard", { replace: true });
    }
  }, [isSignedIn, navigate]);

  // Don't render the sign-up form if user is already authenticated
  if (isSignedIn) {
    return null;
  }

  return (
    <AuthLayout
      title="Create your account"
      description="Sign up to start tracking your subscriptions"
    >
      <div className="space-y-6">
        <SignUp
          appearance={getClerkAppearance(theme)}
          fallbackRedirectUrl={"/dashboard"}
          signInUrl="/sign-in"
          fallback={
            <div className="flex flex-col space-y-4 w-full">
              <div className="h-10 bg-muted animate-pulse rounded" />
              <div className="h-10 bg-muted animate-pulse rounded" />
              <div className="h-10 bg-muted animate-pulse rounded" />
            </div>
          }
        />

        {/* Navigation link to sign-in page */}
        <div className="text-center text-sm">
          <span className="text-muted-foreground">Already have an account? </span>
          <Link
            to="/sign-in"
            className="font-medium text-primary hover:text-primary/80 transition-colors"
          >
            Sign in
          </Link>
        </div>
      </div>
    </AuthLayout>
  );
};

export default SignUpPage;
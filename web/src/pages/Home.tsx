import { useKindeAuth } from '@kinde-oss/kinde-auth-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { useNavigate } from 'react-router-dom';
import ThemeToggle from '@/components/ThemeToggle';

const Home = () => {
  const { user, login, register } = useKindeAuth();
  const navigate = useNavigate();

  const handleSignIn = () => {
    login();
  };

  const handleSignUp = () => {
    register();
  };

  const handleDashboard = () => {
    navigate('/dashboard');
  };

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
            {user ? (
              <Button onClick={handleDashboard} variant="default">
                Dashboard
              </Button>
            ) : (
              <div className="space-x-2">
                <Button onClick={handleSignIn} variant="outline">
                  Sign In
                </Button>
                <Button onClick={handleSignUp} variant="default">
                  Sign Up
                </Button>
              </div>
            )}
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="bg-gray-50 py-16">
        <div className="container mx-auto px-4 text-center">
          <h1 className="text-4xl font-bold mb-6">Track Your Subscriptions Effortlessly</h1>
          <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
            Never miss a payment again. Manage all your recurring subscriptions in one place and stay on top of your finances.
          </p>
          <div className="flex justify-center space-x-4">
            {user ? (
              <Button onClick={handleDashboard} size="lg">
                Go to Dashboard
              </Button>
            ) : (
              <Button onClick={handleSignUp} size="lg">
                Get Started for Free
              </Button>
            )}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold text-center mb-12">Key Features</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <Card>
              <CardHeader>
                <CardTitle>Track Subscriptions</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  Keep track of all your subscriptions in one place, never forget about a recurring payment again.
                </CardDescription>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Payment Reminders</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  Display date of the next payment and get notified before you're charged.
                </CardDescription>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Price Evolution</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  See how subscription prices change over time and identify increasing costs.
                </CardDescription>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Family Sharing</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  Share with family members to have a global view of your household's recurring expenses.
                </CardDescription>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Multi-Currency Support</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  Handle subscriptions in multiple currencies with automatic conversion to your preferred currency.
                </CardDescription>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Smart Filtering</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription>
                  Filter by labels or family members to get a clear view of specific subscription categories.
                </CardDescription>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-gray-50 py-16 mt-auto">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl font-bold mb-6">Ready to Take Control of Your Subscriptions?</h2>
          <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
            Join thousands of users who are saving money by tracking and optimizing their recurring payments.
          </p>
          {user ? (
            <Button onClick={handleDashboard} size="lg">
              Go to Dashboard
            </Button>
          ) : (
            <Button onClick={handleSignUp} size="lg">
              Sign Up Now
            </Button>
          )}
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t py-8">
        <div className="container mx-auto px-4 text-center text-gray-500">
          <p>© {new Date().getFullYear()} RecurrentPaymentTracker. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
};

export default Home;
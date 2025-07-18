import { useKindeAuth } from '@kinde-oss/kinde-auth-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { useNavigate } from 'react-router-dom';
import ThemeToggle from '@/components/ThemeToggle';
import { 
  CreditCard, 
  Bell, 
  LineChart, 
  Users, 
  Globe, 
  Filter, 
  ChevronRight,
  Calendar,
  DollarSign,
  Sparkles
} from 'lucide-react';

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
      <nav className="border-b shadow-sm animate-fadeIn">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2 animate-slideInLeft">
            <DollarSign className="h-6 w-6 text-primary" />
            <h1 className="text-xl font-bold">RecurrentPaymentTracker</h1>
            <Badge variant="secondary" className="ml-2 animate-pulse">Beta</Badge>
          </div>
          <div className="flex items-center gap-3 animate-slideInRight">
            <ThemeToggle />
            {user ? (
              <Button onClick={handleDashboard} variant="default" className="animate-fadeIn">
                <Calendar className="mr-2 h-4 w-4" />
                Dashboard
              </Button>
            ) : (
              <div className="space-x-2">
                <Button onClick={handleSignIn} variant="outline" className="animate-fadeIn">
                  Sign In
                </Button>
                <Button onClick={handleSignUp} variant="default" className="animate-fadeIn">
                  <Sparkles className="mr-2 h-4 w-4" />
                  Sign Up
                </Button>
              </div>
            )}
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="bg-gradient-to-b from-gray-50 to-gray-100 py-20 relative overflow-hidden">
        <div className="absolute inset-0 bg-grid-pattern opacity-5"></div>
        <div className="container mx-auto px-4 text-center relative z-10">
          <Badge variant="secondary" className="mb-4 animate-fadeIn">Subscription Management Made Easy</Badge>
          <h1 className="text-5xl font-bold mb-6 animate-slideInUp">
            Track Your Subscriptions <span className="text-primary">Effortlessly</span>
          </h1>
          <p className="text-xl text-gray-600 mb-10 max-w-2xl mx-auto animate-fadeIn">
            Never miss a payment again. Manage all your recurring subscriptions in one place 
            and stay on top of your finances with smart reminders and insights.
          </p>
          <div className="flex flex-col sm:flex-row justify-center gap-4 animate-bounceIn">
            {user ? (
              <Button onClick={handleDashboard} size="lg" className="group">
                <Calendar className="mr-2 h-5 w-5 group-hover:animate-bounce" />
                Go to Dashboard
                <ChevronRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
              </Button>
            ) : (
              <Button onClick={handleSignUp} size="lg" className="group">
                <Sparkles className="mr-2 h-5 w-5 group-hover:animate-spin" />
                Get Started for Free
                <ChevronRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
              </Button>
            )}
          </div>
        </div>
        <Separator className="mt-16" />
      </section>

      {/* Features Section */}
      <section className="py-20">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16 animate-fadeIn">
            <Badge variant="outline" className="mb-4">Why Choose Us</Badge>
            <h2 className="text-4xl font-bold">Powerful Features</h2>
            <Separator className="w-24 mx-auto my-6" />
            <p className="text-lg text-gray-600 max-w-2xl mx-auto">
              Everything you need to manage your subscriptions efficiently in one place
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <Card className="border-2 hover:border-primary transition-all duration-300 hover:shadow-lg group animate-fadeIn">
              <CardHeader className="flex flex-row items-center gap-4">
                <div className="bg-primary/10 p-3 rounded-full group-hover:bg-primary/20 transition-colors">
                  <CreditCard className="h-6 w-6 text-primary" />
                </div>
                <CardTitle>Track Subscriptions</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base">
                  Keep track of all your subscriptions in one place, never forget about a recurring payment again.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="border-2 hover:border-primary transition-all duration-300 hover:shadow-lg group animate-fadeIn delay-100">
              <CardHeader className="flex flex-row items-center gap-4">
                <div className="bg-primary/10 p-3 rounded-full group-hover:bg-primary/20 transition-colors">
                  <Bell className="h-6 w-6 text-primary" />
                </div>
                <CardTitle>Payment Reminders</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base">
                  Display date of the next payment and get notified before you're charged.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="border-2 hover:border-primary transition-all duration-300 hover:shadow-lg group animate-fadeIn delay-200">
              <CardHeader className="flex flex-row items-center gap-4">
                <div className="bg-primary/10 p-3 rounded-full group-hover:bg-primary/20 transition-colors">
                  <LineChart className="h-6 w-6 text-primary" />
                </div>
                <CardTitle>Price Evolution</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base">
                  See how subscription prices change over time and identify increasing costs.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="border-2 hover:border-primary transition-all duration-300 hover:shadow-lg group animate-fadeIn delay-300">
              <CardHeader className="flex flex-row items-center gap-4">
                <div className="bg-primary/10 p-3 rounded-full group-hover:bg-primary/20 transition-colors">
                  <Users className="h-6 w-6 text-primary" />
                </div>
                <CardTitle>Family Sharing</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base">
                  Share with family members to have a global view of your household's recurring expenses.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="border-2 hover:border-primary transition-all duration-300 hover:shadow-lg group animate-fadeIn delay-400">
              <CardHeader className="flex flex-row items-center gap-4">
                <div className="bg-primary/10 p-3 rounded-full group-hover:bg-primary/20 transition-colors">
                  <Globe className="h-6 w-6 text-primary" />
                </div>
                <CardTitle>Multi-Currency Support</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base">
                  Handle subscriptions in multiple currencies with automatic conversion to your preferred currency.
                </CardDescription>
              </CardContent>
            </Card>

            <Card className="border-2 hover:border-primary transition-all duration-300 hover:shadow-lg group animate-fadeIn delay-500">
              <CardHeader className="flex flex-row items-center gap-4">
                <div className="bg-primary/10 p-3 rounded-full group-hover:bg-primary/20 transition-colors">
                  <Filter className="h-6 w-6 text-primary" />
                </div>
                <CardTitle>Smart Filtering</CardTitle>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-base">
                  Filter by labels or family members to get a clear view of specific subscription categories.
                </CardDescription>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-gradient-to-r from-primary/10 to-primary/5 py-24 mt-auto relative overflow-hidden">
        {/* Decorative elements */}
        <div className="absolute top-0 left-0 w-72 h-72 bg-primary/10 rounded-full -translate-x-1/2 -translate-y-1/2 animate-pulse"></div>
        <div className="absolute bottom-0 right-0 w-96 h-96 bg-primary/10 rounded-full translate-x-1/3 translate-y-1/3 animate-pulse delay-1000"></div>
        
        <div className="container mx-auto px-4 text-center relative z-10">
          <Badge variant="secondary" className="mb-6 animate-fadeIn">Join Our Community</Badge>
          <h2 className="text-4xl font-bold mb-6 animate-slideInUp">Ready to Take Control of Your Subscriptions?</h2>
          <p className="text-xl text-gray-600 mb-10 max-w-2xl mx-auto animate-fadeIn">
            Join thousands of users who are saving money by tracking and optimizing their recurring payments.
            <span className="block mt-2 font-medium">Start managing your finances better today!</span>
          </p>
          
          <div className="flex flex-col sm:flex-row justify-center gap-4 animate-bounceIn">
            {user ? (
              <Button onClick={handleDashboard} size="lg" className="bg-primary/90 hover:bg-primary group">
                <Calendar className="mr-2 h-5 w-5" />
                Go to Dashboard
                <ChevronRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
              </Button>
            ) : (
              <>
                <Button onClick={handleSignIn} variant="outline" size="lg" className="border-primary/20 hover:border-primary/50">
                  Sign In
                </Button>
                <Button onClick={handleSignUp} size="lg" className="bg-primary/90 hover:bg-primary group">
                  <Sparkles className="mr-2 h-5 w-5" />
                  Sign Up Now
                  <ChevronRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
                </Button>
              </>
            )}
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t py-12 bg-gray-50/50">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row justify-between items-center mb-8 animate-fadeIn">
            <div className="flex items-center gap-2 mb-4 md:mb-0">
              <DollarSign className="h-6 w-6 text-primary" />
              <h3 className="text-xl font-bold">RecurrentPaymentTracker</h3>
            </div>
            <div className="flex gap-4">
              <Button variant="ghost" size="icon" className="rounded-full hover:text-primary hover:bg-primary/10">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="h-5 w-5"><path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"></path></svg>
                <span className="sr-only">Twitter</span>
              </Button>
              <Button variant="ghost" size="icon" className="rounded-full hover:text-primary hover:bg-primary/10">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="h-5 w-5"><path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"></path><rect width="4" height="12" x="2" y="9"></rect><circle cx="4" cy="4" r="2"></circle></svg>
                <span className="sr-only">LinkedIn</span>
              </Button>
              <Button variant="ghost" size="icon" className="rounded-full hover:text-primary hover:bg-primary/10">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="h-5 w-5"><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"></path><path d="M9 18c-4.51 2-5-2-7-2"></path></svg>
                <span className="sr-only">GitHub</span>
              </Button>
            </div>
          </div>
          
          <Separator className="mb-8" />
          
          <div className="flex flex-col md:flex-row justify-between items-center text-gray-500 animate-fadeIn">
            <p>© {new Date().getFullYear()} RecurrentPaymentTracker. All rights reserved.</p>
            <div className="flex gap-6 mt-4 md:mt-0">
              <a href="#" className="text-sm hover:text-primary transition-colors">Privacy Policy</a>
              <a href="#" className="text-sm hover:text-primary transition-colors">Terms of Service</a>
              <a href="#" className="text-sm hover:text-primary transition-colors">Contact</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default Home;
import { ModeToggle } from "@/components/mode-toggle";
import { Link } from "react-router-dom";
import { useAuth } from "@/contexts/useAuth";

const HomePage = () => {
  const { isAuthenticated } = useAuth();

  return (
    <div className="container mx-auto p-4">
      <header className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Recurrent Payment Tracker</h1>
        <ModeToggle />
      </header>
      
      <main>
        <section className="mb-8">
          <h2 className="text-2xl font-semibold mb-4">Welcome to Payment Tracker</h2>
          <p className="mb-4">
            Track and manage all your subscription services and recurring payments in one place.
          </p>
          <div className="flex gap-4">
            {!isAuthenticated ? (
              <Link 
                to="/login" 
                className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90"
              >
                Login
              </Link>
            ) : (
              <Link 
                to="/dashboard" 
                className="px-4 py-2 bg-secondary text-secondary-foreground rounded hover:bg-secondary/90"
              >
                Go to Dashboard
              </Link>
            )}
          </div>
        </section>
        
        <section className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="p-6 border rounded-lg">
            <h3 className="text-xl font-medium mb-2">Track Expenses</h3>
            <p>Keep track of all your recurring payments and subscriptions in one place.</p>
          </div>
          <div className="p-6 border rounded-lg">
            <h3 className="text-xl font-medium mb-2">Get Reminders</h3>
            <p>Never miss a payment with timely reminders before your subscriptions renew.</p>
          </div>
          <div className="p-6 border rounded-lg">
            <h3 className="text-xl font-medium mb-2">Analyze Spending</h3>
            <p>Visualize your spending patterns and identify opportunities to save.</p>
          </div>
        </section>
      </main>
      
      <footer className="mt-12 pt-6 border-t">
        <p className="text-center text-muted-foreground">
          © {new Date().getFullYear()} Recurrent Payment Tracker. All rights reserved.
        </p>
      </footer>
    </div>
  );
};

export default HomePage;
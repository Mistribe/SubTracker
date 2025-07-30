import {ModeToggle} from "@/components/mode-toggle";
import {Link} from "react-router-dom";
import {useKindeAuth} from "@kinde-oss/kinde-auth-react";
import {Button} from "@/components/ui/button.tsx";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card";
import {CalendarIcon, CreditCardIcon, DollarSignIcon, TagIcon, TrendingUpIcon, UsersIcon} from "lucide-react";

const HomePage = () => {
    const {login, isAuthenticated} = useKindeAuth();

    return (
        <div className="container mx-auto p-4">
            <header className="flex justify-between items-center mb-8">
                <h1 className="text-3xl font-bold">Recurrent Payment Tracker</h1>
                <ModeToggle/>
            </header>

            <main>
                {/* Hero Section */}
                <section
                    className="mb-16 py-12 px-6 bg-gradient-to-r from-primary/10 to-secondary/10 rounded-2xl text-center">
                    <h2 className="text-4xl font-bold mb-6 ">Simplify Your Subscription Management</h2>
                    <p className="text-xl mb-8 max-w-3xl mx-auto text-muted-foreground ">
                        Take control of your recurring payments, track expenses, and never miss a renewal again.
                        The smart way to manage all your subscriptions in one place.
                    </p>
                    <div className="flex gap-4 justify-center">
                        {!isAuthenticated ? (
                            <Button
                                onClick={() => login()} type="button"
                                size="lg"
                                className="text-base font-medium"
                            >
                                Get Started Now
                            </Button>
                        ) : (
                            <Link to="/dashboard">
                                <Button
                                    size="lg"
                                    className="text-base font-medium"
                                >
                                    Go to Dashboard
                                </Button>
                            </Link>
                        )}
                    </div>
                </section>

                {/* Features Section */}
                <section className="mb-16">
                    <h2 className="text-3xl font-bold mb-8 text-center">Powerful Features</h2>

                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
                        {/* Feature 1: Track subscription spending */}
                        <Card
                            className="transition-all hover:shadow-md animate-in fade-in slide-in-from-bottom duration-500">
                            <CardHeader className="pb-2">
                                <div
                                    className="p-2 bg-primary/10 rounded-full w-12 h-12 flex items-center justify-center mb-2 animate-in zoom-in duration-700 delay-100">
                                    <CreditCardIcon className="h-6 w-6 text-primary"/>
                                </div>
                                <CardTitle className="text-xl">Track Subscription Spending</CardTitle>
                                <CardDescription>Monitor how much you spend on each provider</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p className="text-muted-foreground">
                                    Get a clear overview of your subscription costs across different services and
                                    identify where your money is going each month.
                                </p>
                            </CardContent>
                        </Card>

                        {/* Feature 2: Track price changes */}
                        <Card
                            className="transition-all hover:shadow-md animate-in fade-in slide-in-from-bottom duration-500 delay-150">
                            <CardHeader className="pb-2">
                                <div
                                    className="p-2 bg-primary/10 rounded-full w-12 h-12 flex items-center justify-center mb-2 animate-in zoom-in duration-700 delay-250">
                                    <TrendingUpIcon className="h-6 w-6 text-primary"/>
                                </div>
                                <CardTitle className="text-xl">Track Price Changes</CardTitle>
                                <CardDescription>Stay informed about subscription price increases</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p className="text-muted-foreground">
                                    Never be surprised by unexpected price hikes. Get notified when your subscription
                                    costs change so you can make informed decisions.
                                </p>
                            </CardContent>
                        </Card>

                        {/* Feature 3: Share with family */}
                        <Card
                            className="transition-all hover:shadow-md animate-in fade-in slide-in-from-bottom duration-500 delay-300">
                            <CardHeader className="pb-2">
                                <div
                                    className="p-2 bg-primary/10 rounded-full w-12 h-12 flex items-center justify-center mb-2 animate-in zoom-in duration-700 delay-400">
                                    <UsersIcon className="h-6 w-6 text-primary"/>
                                </div>
                                <CardTitle className="text-xl">Share With Family</CardTitle>
                                <CardDescription>Manage subscriptions as a group</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p className="text-muted-foreground">
                                    See who is using each subscription and who is paying for it. Perfect for family
                                    plans and shared accounts.
                                </p>
                            </CardContent>
                        </Card>

                        {/* Feature 4: Tag subscriptions */}
                        <Card
                            className="transition-all hover:shadow-md animate-in fade-in slide-in-from-bottom duration-500 delay-450">
                            <CardHeader className="pb-2">
                                <div
                                    className="p-2 bg-primary/10 rounded-full w-12 h-12 flex items-center justify-center mb-2 animate-in zoom-in duration-700 delay-550">
                                    <TagIcon className="h-6 w-6 text-primary"/>
                                </div>
                                <CardTitle className="text-xl">Tag Your Subscriptions</CardTitle>
                                <CardDescription>Organize with custom labels</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p className="text-muted-foreground">
                                    Categorize your subscriptions with custom labels for better organization and
                                    filtering. Group by entertainment, productivity, or any category you choose.
                                </p>
                            </CardContent>
                        </Card>

                        {/* Feature 5: Track renewals */}
                        <Card
                            className="transition-all hover:shadow-md animate-in fade-in slide-in-from-bottom duration-500 delay-600">
                            <CardHeader className="pb-2">
                                <div
                                    className="p-2 bg-primary/10 rounded-full w-12 h-12 flex items-center justify-center mb-2 animate-in zoom-in duration-700 delay-700">
                                    <CalendarIcon className="h-6 w-6 text-primary"/>
                                </div>
                                <CardTitle className="text-xl">Track Renewals</CardTitle>
                                <CardDescription>Never miss a payment date</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p className="text-muted-foreground">
                                    Get timely notifications before your subscriptions renew so you can decide whether
                                    to continue, cancel, or modify your plan.
                                </p>
                            </CardContent>
                        </Card>

                        {/* Feature 6: Expense Analytics */}
                        <Card
                            className="transition-all hover:shadow-md animate-in fade-in slide-in-from-bottom duration-500 delay-750">
                            <CardHeader className="pb-2">
                                <div
                                    className="p-2 bg-primary/10 rounded-full w-12 h-12 flex items-center justify-center mb-2 animate-in zoom-in duration-700 delay-850">
                                    <DollarSignIcon className="h-6 w-6 text-primary"/>
                                </div>
                                <CardTitle className="text-xl">Expense Analytics</CardTitle>
                                <CardDescription>Visualize your spending patterns</CardDescription>
                            </CardHeader>
                            <CardContent>
                                <p className="text-muted-foreground">
                                    Get insights into your subscription spending with detailed charts and reports to
                                    help you optimize your budget.
                                </p>
                            </CardContent>
                        </Card>
                    </div>
                </section>

                {/* CTA Section */}
                <section className="mb-16 py-10 px-6 bg-primary/5 rounded-2xl text-center">
                    <h2 className="text-3xl font-bold mb-4">Ready to Take Control?</h2>
                    <p className="text-lg mb-6 max-w-2xl mx-auto text-muted-foreground">
                        Join thousands of users who are saving money and reducing subscription stress.
                    </p>
                    <div>
                        {!isAuthenticated ? (
                            <Button
                                onClick={() => login()} type="button"
                                size="lg"
                                className="text-base font-medium"
                            >
                                Start Managing Your Subscriptions
                            </Button>
                        ) : (
                            <Link to="/dashboard">
                                <Button
                                    size="lg"
                                    variant="secondary"
                                    className="text-base font-medium"
                                >
                                    Go to Dashboard
                                </Button>
                            </Link>
                        )}
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
import {useState} from "react";
import {useNavigate, Link} from "react-router-dom";
import {ModeToggle} from "@/components/mode-toggle";
import {useKindeAuth} from "@kinde-oss/kinde-auth-react";

interface Payment {
    id: number;
    name: string;
    amount: number;
    frequency: string;
    nextPayment: string;
}

const DashboardPage = () => {
    const navigate = useNavigate();
    const {logout} = useKindeAuth();

    // Sample data - in a real app, this would come from an API
    const [payments] = useState<Payment[]>([
        {
            id: 1,
            name: "Netflix",
            amount: 15.99,
            frequency: "Monthly",
            nextPayment: "2025-08-01"
        },
        {
            id: 2,
            name: "Spotify",
            amount: 9.99,
            frequency: "Monthly",
            nextPayment: "2025-07-25"
        },
        {
            id: 3,
            name: "Amazon Prime",
            amount: 139,
            frequency: "Yearly",
            nextPayment: "2026-01-15"
        }
    ]);

    const handleLogout = () => {
        // Use the logout function from auth context
        logout();
        // Redirect to home page
        navigate("/");
    };

    const totalMonthly = payments
        .filter(p => p.frequency === "Monthly")
        .reduce((sum, payment) => sum + payment.amount, 0);

    const totalYearly = payments
        .filter(p => p.frequency === "Yearly")
        .reduce((sum, payment) => sum + payment.amount, 0);

    return (
        <div className="container mx-auto p-4">
            <header className="flex justify-between items-center mb-8">
                <h1 className="text-3xl font-bold">Recurrent Payment Tracker</h1>
                <div className="flex items-center gap-4">
                    <Link
                        to="/profile"
                        className="text-sm hover:underline"
                    >
                        Profile
                    </Link>
                    <button
                        onClick={handleLogout}
                        className="text-sm hover:underline"
                    >
                        Logout
                    </button>
                    <ModeToggle/>
                </div>
            </header>

            <main>
                <div className="mb-8">
                    <h2 className="text-2xl font-semibold mb-4">Dashboard</h2>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                        <div className="p-6 border rounded-lg bg-card">
                            <h3 className="text-xl font-medium mb-2">Monthly Expenses</h3>
                            <p className="text-3xl font-bold">${totalMonthly.toFixed(2)}</p>
                        </div>
                        <div className="p-6 border rounded-lg bg-card">
                            <h3 className="text-xl font-medium mb-2">Yearly Expenses</h3>
                            <p className="text-3xl font-bold">${totalYearly.toFixed(2)}</p>
                        </div>
                        <div className="p-6 border rounded-lg bg-card">
                            <h3 className="text-xl font-medium mb-2">Total Annual Cost</h3>
                            <p className="text-3xl font-bold">${(totalMonthly * 12 + totalYearly).toFixed(2)}</p>
                        </div>
                    </div>
                </div>

                <div>
                    <div className="flex justify-between items-center mb-4">
                        <h3 className="text-xl font-semibold">Your Subscriptions</h3>
                        <button className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90">
                            Add New
                        </button>
                    </div>

                    <div className="overflow-x-auto">
                        <table className="w-full border-collapse">
                            <thead>
                            <tr className="bg-muted">
                                <th className="p-3 text-left">Name</th>
                                <th className="p-3 text-left">Amount</th>
                                <th className="p-3 text-left">Frequency</th>
                                <th className="p-3 text-left">Next Payment</th>
                                <th className="p-3 text-left">Actions</th>
                            </tr>
                            </thead>
                            <tbody>
                            {payments.map(payment => (
                                <tr key={payment.id} className="border-b">
                                    <td className="p-3">{payment.name}</td>
                                    <td className="p-3">${payment.amount.toFixed(2)}</td>
                                    <td className="p-3">{payment.frequency}</td>
                                    <td className="p-3">{payment.nextPayment}</td>
                                    <td className="p-3">
                                        <button className="text-blue-500 hover:underline mr-2">Edit</button>
                                        <button className="text-red-500 hover:underline">Delete</button>
                                    </td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default DashboardPage;
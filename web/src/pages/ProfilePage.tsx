import {useState} from "react";
import {useNavigate, Link} from "react-router-dom";
import {ModeToggle} from "@/components/mode-toggle";
import {useKindeAuth} from "@kinde-oss/kinde-auth-react";

const ProfilePage = () => {
    const navigate = useNavigate();
    const {logout} = useKindeAuth();

    // Sample user data - in a real app, this would come from an API or auth context
    const [user, setUser] = useState({
        name: "John Doe",
        email: "john.doe@example.com",
        joinDate: "2025-01-15"
    });

    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState({...user});

    const handleLogout = () => {
        // Use the logout function from auth context
        logout();
        // Redirect to home page
        navigate("/");
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const {name, value} = e.target;
        setFormData(prev => ({...prev, [name]: value}));
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        setUser(formData);
        setIsEditing(false);
    };

    return (
        <div className="container mx-auto p-4">
            <header className="flex justify-between items-center mb-8">
                <h1 className="text-3xl font-bold">Recurrent Payment Tracker</h1>
                <div className="flex items-center gap-4">
                    <Link
                        to="/dashboard"
                        className="text-sm hover:underline"
                    >
                        Dashboard
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

            <main className="max-w-2xl mx-auto">
                <h2 className="text-2xl font-semibold mb-6">Profile</h2>

                <div className="bg-card p-6 rounded-lg shadow-sm border">
                    {isEditing ? (
                        <form onSubmit={handleSubmit}>
                            <div className="mb-4">
                                <label htmlFor="name" className="block text-sm font-medium mb-1">
                                    Name
                                </label>
                                <input
                                    id="name"
                                    name="name"
                                    type="text"
                                    value={formData.name}
                                    onChange={handleInputChange}
                                    className="w-full p-2 border rounded"
                                    required
                                />
                            </div>

                            <div className="mb-4">
                                <label htmlFor="email" className="block text-sm font-medium mb-1">
                                    Email
                                </label>
                                <input
                                    id="email"
                                    name="email"
                                    type="email"
                                    value={formData.email}
                                    onChange={handleInputChange}
                                    className="w-full p-2 border rounded"
                                    required
                                />
                            </div>

                            <div className="flex justify-end gap-2 mt-6">
                                <button
                                    type="button"
                                    onClick={() => {
                                        setFormData({...user});
                                        setIsEditing(false);
                                    }}
                                    className="px-4 py-2 border rounded hover:bg-muted"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90"
                                >
                                    Save Changes
                                </button>
                            </div>
                        </form>
                    ) : (
                        <div>
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
                                <div>
                                    <h3 className="text-sm font-medium text-muted-foreground">Name</h3>
                                    <p className="mt-1">{user.name}</p>
                                </div>
                                <div>
                                    <h3 className="text-sm font-medium text-muted-foreground">Email</h3>
                                    <p className="mt-1">{user.email}</p>
                                </div>
                                <div>
                                    <h3 className="text-sm font-medium text-muted-foreground">Member Since</h3>
                                    <p className="mt-1">{user.joinDate}</p>
                                </div>
                            </div>

                            <div className="flex justify-end">
                                <button
                                    onClick={() => setIsEditing(true)}
                                    className="px-4 py-2 bg-primary text-primary-foreground rounded hover:bg-primary/90"
                                >
                                    Edit Profile
                                </button>
                            </div>
                        </div>
                    )}
                </div>

                <div className="mt-8 bg-card p-6 rounded-lg shadow-sm border">
                    <h3 className="text-xl font-medium mb-4">Account Settings</h3>

                    <div className="space-y-4">
                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Email Notifications</h4>
                                <p className="text-sm text-muted-foreground">Receive email notifications for upcoming
                                    payments</p>
                            </div>
                            <button className="px-3 py-1 border rounded">Enabled</button>
                        </div>

                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Change Password</h4>
                                <p className="text-sm text-muted-foreground">Update your account password</p>
                            </div>
                            <button className="px-3 py-1 border rounded">Change</button>
                        </div>

                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium text-destructive">Delete Account</h4>
                                <p className="text-sm text-muted-foreground">Permanently delete your account and all
                                    data</p>
                            </div>
                            <button className="px-3 py-1 border border-destructive text-destructive rounded">Delete
                            </button>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default ProfilePage;
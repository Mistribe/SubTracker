import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useKindeAuth } from "@kinde-oss/kinde-auth-react";

const ProfilePage = () => {
    const { user: kindeUser } = useKindeAuth();

    // Sample user data - in a real app, this would come from an API or auth context
    const [user, setUser] = useState({
        name: kindeUser?.given_name && kindeUser?.family_name 
            ? `${kindeUser.given_name} ${kindeUser.family_name}` 
            : "John Doe",
        email: kindeUser?.email || "john.doe@example.com",
        joinDate: "2025-01-15"
    });

    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState({...user});

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
        <div>
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Profile</h1>
            </div>

            <div className="max-w-2xl mx-auto">
                <div className="bg-card p-6 rounded-lg shadow-sm border">
                    {isEditing ? (
                        <form onSubmit={handleSubmit}>
                            <div className="mb-4">
                                <label htmlFor="name" className="block text-sm font-medium mb-1">
                                    Name
                                </label>
                                <Input
                                    id="name"
                                    name="name"
                                    type="text"
                                    value={formData.name}
                                    onChange={handleInputChange}
                                    required
                                />
                            </div>

                            <div className="mb-4">
                                <label htmlFor="email" className="block text-sm font-medium mb-1">
                                    Email
                                </label>
                                <Input
                                    id="email"
                                    name="email"
                                    type="email"
                                    value={formData.email}
                                    onChange={handleInputChange}
                                    required
                                />
                            </div>

                            <div className="flex justify-end gap-2 mt-6">
                                <Button
                                    type="button"
                                    variant="outline"
                                    onClick={() => {
                                        setFormData({...user});
                                        setIsEditing(false);
                                    }}
                                >
                                    Cancel
                                </Button>
                                <Button type="submit">
                                    Save Changes
                                </Button>
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
                                <Button onClick={() => setIsEditing(true)}>
                                    Edit Profile
                                </Button>
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
                                <p className="text-sm text-muted-foreground">Receive email notifications for upcoming payments</p>
                            </div>
                            <Button variant="outline" size="sm">Enabled</Button>
                        </div>

                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Change Password</h4>
                                <p className="text-sm text-muted-foreground">Update your account password</p>
                            </div>
                            <Button variant="outline" size="sm">Change</Button>
                        </div>

                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium text-destructive">Delete Account</h4>
                                <p className="text-sm text-muted-foreground">Permanently delete your account and all data</p>
                            </div>
                            <Button variant="outline" size="sm" className="border-destructive text-destructive">Delete</Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProfilePage;
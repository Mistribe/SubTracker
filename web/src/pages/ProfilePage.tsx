import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useKindeAuth } from "@kinde-oss/kinde-auth-react";
import { 
    Select, 
    SelectContent, 
    SelectItem, 
    SelectTrigger, 
    SelectValue 
} from "@/components/ui/select";
import { ModeToggle } from "@/components/mode-toggle";
import { useTheme } from "@/components/theme-provider";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";
import { PageHeader } from "@/components/ui/page-header";

const ProfilePage = () => {
    const { user: kindeUser } = useKindeAuth();
    const { theme } = useTheme();

    // Sample user data - in a real app, this would come from an API or auth context
    const [user, setUser] = useState({
        name: kindeUser?.given_name && kindeUser?.family_name 
            ? `${kindeUser.given_name} ${kindeUser.family_name}` 
            : "John Doe",
        email: kindeUser?.email || "john.doe@example.com",
        joinDate: "2025-01-15",
        preferredCurrency: "USD",
        preferredTheme: "system"
    });

    // Update theme preference when theme changes
    useEffect(() => {
        setUser(prev => ({...prev, preferredTheme: theme}));
    }, [theme]);

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
    
    // List of common currencies
    const currencies = [
        { value: "USD", label: "US Dollar ($)" },
        { value: "EUR", label: "Euro (€)" },
        { value: "GBP", label: "British Pound (£)" },
        { value: "JPY", label: "Japanese Yen (¥)" },
        { value: "CAD", label: "Canadian Dollar (C$)" },
        { value: "AUD", label: "Australian Dollar (A$)" },
        { value: "CHF", label: "Swiss Franc (Fr)" },
        { value: "CNY", label: "Chinese Yuan (¥)" },
        { value: "INR", label: "Indian Rupee (₹)" },
        { value: "BRL", label: "Brazilian Real (R$)" }
    ];

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Profile"
                description="Manage your profile and preferences"
                actionButton={
                    !isEditing && (
                        <Button onClick={() => setIsEditing(true)}>
                            Edit Profile
                        </Button>
                    )
                }
            />

            <div className="max-w-3xl mx-auto space-y-12">
                {/* User Profile Section */}
                <div className="space-y-8">
                    {/* User Avatar and Basic Info */}
                    <div className="flex flex-col items-center space-y-4 py-8 bg-primary/5 rounded-lg">
                        <Avatar className="h-24 w-24 border-4 border-background">
                            {kindeUser?.picture ? (
                                <AvatarImage src={kindeUser.picture} alt={user.name} />
                            ) : (
                                <AvatarFallback className="text-2xl">
                                    {user.name.split(' ').map(n => n[0]).join('').toUpperCase()}
                                </AvatarFallback>
                            )}
                        </Avatar>
                        <div className="text-center">
                            <h2 className="text-2xl font-bold">{user.name}</h2>
                            <p className="text-muted-foreground">{user.email}</p>
                        </div>
                    </div>

                    {/* Profile Details */}
                    {isEditing ? (
                        <form onSubmit={handleSubmit} className="space-y-6 bg-muted/20 p-6 rounded-lg">
                            <div className="space-y-4">
                                <div>
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
                                        className="w-full"
                                    />
                                </div>

                                <div>
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
                                        className="w-full"
                                    />
                                </div>
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
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div>
                                <h3 className="text-sm font-medium text-muted-foreground">Name</h3>
                                <p className="mt-1 font-medium">{user.name}</p>
                            </div>
                            <div>
                                <h3 className="text-sm font-medium text-muted-foreground">Email</h3>
                                <p className="mt-1 font-medium">{user.email}</p>
                            </div>
                            <div>
                                <h3 className="text-sm font-medium text-muted-foreground">Member Since</h3>
                                <p className="mt-1 font-medium">{user.joinDate}</p>
                            </div>
                        </div>
                    )}
                </div>

                {/* Preferences Section */}
                <div className="space-y-6">
                    <div>
                        <h2 className="text-xl font-bold">Preferences</h2>
                        <p className="text-muted-foreground">Manage your application preferences</p>
                    </div>
                    
                    <div className="space-y-6 bg-muted/20 p-6 rounded-lg">
                        {/* Currency Preference */}
                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Preferred Currency</h4>
                                <p className="text-sm text-muted-foreground">Select your preferred currency for the application</p>
                            </div>
                            <Select 
                                value={user.preferredCurrency}
                                onValueChange={(value) => setUser(prev => ({...prev, preferredCurrency: value}))}
                            >
                                <SelectTrigger className="w-[180px]">
                                    <SelectValue placeholder="Select currency" />
                                </SelectTrigger>
                                <SelectContent>
                                    {currencies.map(currency => (
                                        <SelectItem key={currency.value} value={currency.value}>
                                            {currency.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        </div>

                        <Separator />

                        {/* Theme Preference */}
                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Theme</h4>
                                <p className="text-sm text-muted-foreground">Choose between light and dark mode</p>
                            </div>
                            <ModeToggle />
                        </div>
                    </div>
                </div>

                {/* Account Settings Section */}
                <div className="space-y-6">
                    <div>
                        <h2 className="text-xl font-bold">Account Settings</h2>
                        <p className="text-muted-foreground">Manage your account settings and preferences</p>
                    </div>
                    
                    <div className="space-y-6 bg-muted/20 p-6 rounded-lg">
                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Email Notifications</h4>
                                <p className="text-sm text-muted-foreground">Receive email notifications for upcoming payments</p>
                            </div>
                            <Button variant="outline" size="sm">Enabled</Button>
                        </div>

                        <Separator />

                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Change Password</h4>
                                <p className="text-sm text-muted-foreground">Update your account password</p>
                            </div>
                            <Button variant="outline" size="sm">Change</Button>
                        </div>

                        <Separator />

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
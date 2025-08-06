import {useEffect, useState} from "react";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {useKindeAuth} from "@kinde-oss/kinde-auth-react";
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from "@/components/ui/select";
import {ModeToggle} from "@/components/mode-toggle";
import {useTheme} from "@/components/theme-provider";
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";
import {Separator} from "@/components/ui/separator";
import {PageHeader} from "@/components/ui/page-header";
import {useProfileManagement} from "@/hooks/profile/useProfileManagement";

/**
 * ProfilePage component
 * Displays and manages user profile information, including preferences like currency and theme
 */
const ProfilePage = () => {
    // Get authentication data from Kinde
    const {user: kindeUser} = useKindeAuth();
    // Get theme information from theme context
    const {theme} = useTheme();
    // Use the profile management hook to fetch and update profile data
    const {profile, isLoading, isError, updateProfile, isUpdating} = useProfileManagement();

    // Fallback user data when profile is not yet loaded or there's an error
    // This combines data from Kinde auth with default values
    const [user, setUser] = useState({
        name: kindeUser?.givenName && kindeUser?.familyName
            ? `${kindeUser.givenName} ${kindeUser.familyName}`
            : "John Doe",
        email: kindeUser?.email || "john.doe@example.com",
        joinDate: "2025-01-15",
        preferredCurrency: "USD", // Default currency, will be updated from profile when loaded
        preferredTheme: "system"
    });

    // Update local user data when profile is loaded from the backend
    useEffect(() => {
        if (profile) {
            setUser(prev => ({
                ...prev,
                preferredCurrency: profile.currency || prev.preferredCurrency
            }));
        }
    }, [profile]);

    // Update theme preference when theme changes in the theme context
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
        {value: "USD", label: "US Dollar ($)"},
        {value: "EUR", label: "Euro (€)"},
        {value: "GBP", label: "British Pound (£)"},
        {value: "JPY", label: "Japanese Yen (¥)"},
        {value: "CAD", label: "Canadian Dollar (C$)"},
        {value: "AUD", label: "Australian Dollar (A$)"},
        {value: "CHF", label: "Swiss Franc (Fr)"},
        {value: "CNY", label: "Chinese Yuan (¥)"},
        {value: "INR", label: "Indian Rupee (₹)"},
        {value: "BRL", label: "Brazilian Real (R$)"}
    ];

    return (
        <div className="container mx-auto py-6">
            <PageHeader
                title="Profile"
                description="Manage your profile and preferences"
                actionButton={
                    !isEditing && !isLoading && (
                        <Button onClick={() => setIsEditing(true)}>
                            Edit Profile
                        </Button>
                    )
                }
            />
            
            {isLoading && (
                <div className="flex justify-center my-12">
                    <div className="text-center">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
                        <p className="text-muted-foreground">Loading profile data...</p>
                    </div>
                </div>
            )}
            
            {isError && (
                <div className="bg-destructive/10 text-destructive p-4 rounded-md my-6">
                    <h3 className="font-medium">Error loading profile</h3>
                    <p>There was a problem loading your profile data. Please try refreshing the page.</p>
                </div>
            )}

            <div className="max-w-3xl mx-auto space-y-12">
                {/* User Profile Section */}
                <div className="space-y-8">
                    {/* User Avatar and Basic Info */}
                    <div className="flex flex-col items-center space-y-4 py-8 bg-primary/5 rounded-lg">
                        <Avatar className="h-24 w-24 border-4 border-background">
                            {kindeUser?.picture ? (
                                <AvatarImage src={kindeUser.picture} alt={user.name}/>
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
                        {/* Currency Preference - Changes are persisted to the backend */}
                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Preferred Currency</h4>
                                <p className="text-sm text-muted-foreground">Select your preferred currency for the
                                    application</p>
                            </div>
                            {/* 
                              * Currency selection component
                              * When the user selects a new currency:
                              * 1. Updates the local state immediately for a responsive UI
                              * 2. Calls updateProfile to persist the change to the backend
                              * 3. Disables the select while the update is in progress
                              */}
                            <Select
                                value={user.preferredCurrency}
                                onValueChange={(value) => {
                                    // Update local state for immediate UI feedback
                                    setUser(prev => ({...prev, preferredCurrency: value}));
                                    // Persist to backend using the profile management hook
                                    updateProfile(value);
                                }}
                                disabled={isUpdating}
                            >
                                <SelectTrigger className="w-[180px]">
                                    <SelectValue placeholder="Select currency"/>
                                </SelectTrigger>
                                <SelectContent>
                                    {currencies.map(currency => (
                                        <SelectItem key={currency.value} value={currency.value}>
                                            {currency.label}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                            {isUpdating && <span className="ml-2 text-sm text-muted-foreground">Saving...</span>}
                        </div>

                        <Separator/>

                        {/* Theme Preference */}
                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Theme</h4>
                                <p className="text-sm text-muted-foreground">Choose between light and dark mode</p>
                            </div>
                            <ModeToggle/>
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
                                <p className="text-sm text-muted-foreground">Receive email notifications for upcoming
                                    payments</p>
                            </div>
                            <Button variant="outline" size="sm">Enabled</Button>
                        </div>

                        <Separator/>

                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium">Change Password</h4>
                                <p className="text-sm text-muted-foreground">Update your account password</p>
                            </div>
                            <Button variant="outline" size="sm">Change</Button>
                        </div>

                        <Separator/>

                        <div className="flex justify-between items-center">
                            <div>
                                <h4 className="font-medium text-destructive">Delete Account</h4>
                                <p className="text-sm text-muted-foreground">Permanently delete your account and all
                                    data</p>
                            </div>
                            <Button variant="outline" size="sm"
                                    className="border-destructive text-destructive">Delete</Button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ProfilePage;
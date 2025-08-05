import {type ReactNode} from "react";
import {Link, useLocation, useNavigate} from "react-router-dom";
import {useKindeAuth} from "@kinde-oss/kinde-auth-react";
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuItem,
    SidebarMenuButton,
    SidebarProvider,
    SidebarTrigger,
} from "@/components/ui/sidebar";
import {ModeToggle} from "@/components/mode-toggle";
import {Avatar, AvatarFallback, AvatarImage} from "@/components/ui/avatar";
import {
    HomeIcon,
    CreditCardIcon,
    UsersIcon,
    TagIcon,
    UserIcon,
    LogOutIcon,
    PackageIcon
} from "lucide-react";

interface AppLayoutProps {
    children: ReactNode;
}

export function AppLayout({children}: AppLayoutProps) {
    const {user, logout} = useKindeAuth();
    const navigate = useNavigate();
    const location = useLocation();

    const handleLogout = () => {
        logout();
        navigate("/");
    };

    // Get user initials for avatar fallback
    const getInitials = () => {
        if (!user?.givenName && !user?.familyName) return "U";
        return `${user?.givenName?.[0] || ""}${user?.familyName ?.[0] || ""}`;
    };

    return (
        <SidebarProvider>
            <div className="flex min-h-screen w-full">
                <Sidebar>
                    <SidebarHeader className="flex items-center justify-between">
                        <Link to="/dashboard" className="flex items-center gap-2 px-2">
                            <CreditCardIcon className="h-6 w-6"/>
                            <span className="font-semibold">Payment Tracker</span>
                        </Link>
                        <SidebarTrigger/>
                    </SidebarHeader>

                    <SidebarContent>
                        <SidebarMenu>
                            <SidebarMenuItem>
                                <SidebarMenuButton
                                    asChild
                                    isActive={location.pathname === "/dashboard"}
                                    tooltip="Dashboard"
                                >
                                    <Link to="/dashboard" className="flex items-center">
                                        <HomeIcon className="mr-2 h-4 w-4"/>
                                        <span>Dashboard</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                            
                            <SidebarMenuItem>
                                <SidebarMenuButton
                                    asChild
                                    isActive={location.pathname === "/subscriptions"}
                                    tooltip="Subscriptions"
                                >
                                    <Link to="/subscriptions" className="flex items-center">
                                        <CreditCardIcon className="mr-2 h-4 w-4"/>
                                        <span>Subscriptions</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem>
                                <SidebarMenuButton
                                    asChild
                                    isActive={location.pathname === "/providers"}
                                    tooltip="Providers"
                                >
                                    <Link to="/providers" className="flex items-center">
                                        <PackageIcon className="mr-2 h-4 w-4"/>
                                        <span>Providers</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem>
                                <SidebarMenuButton
                                    asChild
                                    isActive={location.pathname === "/families"}
                                    tooltip="Families"
                                >
                                    <Link to="/families" className="flex items-center">
                                        <UsersIcon className="mr-2 h-4 w-4"/>
                                        <span>Families</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem>
                                <SidebarMenuButton
                                    asChild
                                    isActive={location.pathname === "/labels"}
                                    tooltip="Labels"
                                >
                                    <Link to="/labels" className="flex items-center">
                                        <TagIcon className="mr-2 h-4 w-4"/>
                                        <span>Labels</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>

                            <SidebarMenuItem>
                                <SidebarMenuButton
                                    asChild
                                    isActive={location.pathname === "/profile"}
                                    tooltip="Profile"
                                >
                                    <Link to="/profile" className="flex items-center">
                                        <UserIcon className="mr-2 h-4 w-4"/>
                                        <span>Profile</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        </SidebarMenu>
                    </SidebarContent>

                    <SidebarFooter>
                        <div className="flex items-center justify-between p-2">
                            <div className="flex items-center gap-2">
                                <Avatar>
                                    <AvatarImage src={user?.picture || ""} alt={user?.givenName || "User"}/>
                                    <AvatarFallback>{getInitials()}</AvatarFallback>
                                </Avatar>
                                <div className="flex flex-col">
                                    <span className="text-sm font-medium">{user?.givenName} {user?.familyName}</span>
                                    <span
                                        className="text-xs text-muted-foreground truncate max-w-[120px]">{user?.email}</span>
                                </div>
                            </div>
                            <div className="flex items-center gap-2">
                                <ModeToggle/>
                                <button
                                    onClick={handleLogout}
                                    className="text-muted-foreground hover:text-foreground"
                                    title="Logout"
                                >
                                    <LogOutIcon className="h-4 w-4"/>
                                </button>
                            </div>
                        </div>
                    </SidebarFooter>
                </Sidebar>

                <main className="flex-1 p-6 overflow-auto">
                    {children}
                </main>
            </div>
        </SidebarProvider>
    );
}
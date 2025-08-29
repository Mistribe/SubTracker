/**
 * AppSidebar Component
 *
 * A modern sidebar for the Payment Tracker application using shadcn/ui components and Tailwind CSS.
 * Features:
 * - Responsive design with collapsible sidebar
 * - Navigation menu with visual indicators for active items
 * - User profile section with avatar and account information
 * - Theme toggle and logout functionality
 * - Tooltips for better user experience
 * - Scrollable content area
 */

import {Link, useLocation,} from "react-router-dom";
import {useUser, UserButton} from "@clerk/clerk-react";
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarTrigger,
    useSidebar,
} from "@/components/ui/sidebar";
import {ScrollArea} from "@/components/ui/scroll-area";
import {Separator} from "@/components/ui/separator";
import {Tooltip, TooltipContent, TooltipProvider, TooltipTrigger} from "@/components/ui/tooltip";
import {Badge} from "@/components/ui/badge";
import {
    AlertTriangleIcon,
    CreditCardIcon,
    HomeIcon,
    PackageIcon, Settings,
    TagIcon,
    UsersIcon
} from "lucide-react";
import {envVar} from "@/lib/env";

export function AppSidebar() {
    const {user} = useUser();
    const location = useLocation();
    const {state, isMobile} = useSidebar();

    // Get environment from VITE_TARGET_ENV or default to development
    const getEnvironmentInfo = () => {
        const env = envVar('VITE_TARGET_ENV') || 'development';

        switch (env) {
            case 'development':
                return {
                    show: true,
                    text: 'alpha',
                    variant: 'secondary' as const,
                    icon: <AlertTriangleIcon className="h-3 w-3"/>
                };
            case 'staging':
                return {
                    show: true,
                    text: 'beta',
                    variant: 'outline' as const,
                    icon: <AlertTriangleIcon className="h-3 w-3"/>
                };
            case 'production':
                return {show: false};
            default:
                return {
                    show: true,
                    text: 'alpha',
                    variant: 'secondary' as const,
                    icon: <AlertTriangleIcon className="h-3 w-3"/>
                };
        }
    };

    // Navigation items
    const navItems = [
        {path: "/dashboard", icon: <HomeIcon className="h-4 w-4"/>, label: "Dashboard"},
        {path: "/subscriptions", icon: <CreditCardIcon className="h-4 w-4"/>, label: "Subscriptions"},
        {path: "/providers", icon: <PackageIcon className="h-4 w-4"/>, label: "Providers"},
        {path: "/family", icon: <UsersIcon className="h-4 w-4"/>, label: "Family"},
        {path: "/labels", icon: <TagIcon className="h-4 w-4"/>, label: "Labels"},
    ];

    // Get environment information
    const environmentInfo = getEnvironmentInfo();

    return (
        <Sidebar collapsible="icon">
            {environmentInfo.show && (
                <div className={`w-full text-center py-1 ${state === "collapsed" ? "px-2" : "px-4"}`}>
                    <Badge
                        variant={environmentInfo.variant}
                        className={`w-full flex items-center justify-center gap-1 ${state === "collapsed" ? "py-1" : ""}`}
                    >
                        {environmentInfo.icon}
                        {state === "expanded" && environmentInfo.text}
                    </Badge>
                </div>
            )}
            <SidebarHeader className="flex items-center justify-between p-4 border-b">
                <Link to="/dashboard" className="flex items-center gap-2">
                    <div className="bg-primary/10 p-2 rounded-md">
                        <CreditCardIcon className="h-5 w-5 text-primary"/>
                    </div>
                    <span className="font-semibold text-lg group-data-[collapsible=icon]:hidden">Payment Tracker</span>
                </Link>
                <SidebarTrigger className="hover:bg-muted"/>
            </SidebarHeader>

            <SidebarContent>
                <ScrollArea className="h-[calc(100vh-12rem)]">
                    <div className="px-2 py-4">
                        {state === "expanded" && (
                            <div className="mb-2 px-4">
                                <h2 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                                    Main Navigation
                                </h2>
                            </div>
                        )}
                        <SidebarMenu>
                            {navItems.map((item) => (
                                <SidebarMenuItem key={item.path}>
                                    <SidebarMenuButton
                                        asChild
                                        isActive={location.pathname === item.path}
                                        tooltip={item.label}
                                    >
                                        <Link
                                            to={item.path}
                                            className={`flex items-center gap-3 px-3 py-2 rounded-md transition-colors ${
                                                location.pathname.startsWith(item.path)
                                                    ? "bg-primary/10 text-primary font-medium"
                                                    : "hover:bg-muted"
                                            }`}
                                        >
                                            <div
                                                className={`${location.pathname === item.path ? "text-primary" : "text-muted-foreground"}`}>
                                                {item.icon}
                                            </div>
                                            <span>{item.label}</span>
                                        </Link>
                                    </SidebarMenuButton>
                                </SidebarMenuItem>
                            ))}
                        </SidebarMenu>
                    </div>
                </ScrollArea>
            </SidebarContent>

            <SidebarFooter>
                <Separator className="my-2"/>
                {state === "expanded" && (
                    <div className="flex items-center gap-3">
                        <UserButton/>
                        <div className="flex justify-between w-full">
                            <div className="flex flex-col">
                                <p className="text-sm font-medium">
                                    {user?.firstName} {user?.lastName}
                                </p>
                            </div>
                            <div>
                                <Link to="/profile" className="text-muted-foreground">
                                    <Settings
                                        className="transition-transform hover:rotate-180 duration-500"/>
                                </Link>
                            </div>
                        </div>
                    </div>
                )}
                {state === "collapsed" && (
                    <div className="flex flex-col items-center gap-3">
                        <UserButton/>
                        <TooltipProvider delayDuration={0}>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <Link to="/profile" className="text-muted-foreground ">
                                        <Settings
                                            className="h-4 w-4 transition-transform hover:rotate-180 duration-500"/>
                                    </Link>
                                </TooltipTrigger>
                                <TooltipContent side="right" hidden={isMobile}>
                                    Preference
                                </TooltipContent>
                            </Tooltip>
                        </TooltipProvider>
                    </div>
                )}
            </SidebarFooter>
        </Sidebar>
    );
}
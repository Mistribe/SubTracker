/**
 * AppSidebar Component
 *
 * A modern, sleek sidebar for the Payment Tracker application using shadcn/ui components and Tailwind CSS.
 * Features:
 * - Responsive design with smooth collapsible animations
 * - Modern navigation menu with enhanced visual indicators
 * - Polished user profile section with card design
 * - Environment badges with contextual styling
 * - Micro-interactions and hover effects
 * - Optimized for all screen sizes
 */

import { Link, useLocation } from "react-router-dom";
import { useUser, UserButton } from "@clerk/clerk-react";
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
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Badge } from "@/components/ui/badge";
import {
    AlertTriangleIcon,
    CreditCardIcon,
    HomeIcon,
    PackageIcon,
    Settings,
    TagIcon,
    UsersIcon,
    ChevronRight
} from "lucide-react";
import { envVar } from "@/lib/env";

export function AppSidebar() {
    const { user } = useUser();
    const location = useLocation();
    const { state, isMobile } = useSidebar();

    // Get environment from VITE_TARGET_ENV or default to development
    const getEnvironmentInfo = () => {
        const env = envVar('VITE_TARGET_ENV') || 'development';

        switch (env) {
            case 'development':
                return {
                    show: true,
                    text: 'Alpha',
                    variant: 'secondary' as const,
                    icon: <AlertTriangleIcon className="h-3 w-3" />,
                    className: 'bg-amber-500/10 text-amber-700 dark:text-amber-400 border-amber-500/20'
                };
            case 'staging':
                return {
                    show: true,
                    text: 'Beta',
                    variant: 'outline' as const,
                    icon: <AlertTriangleIcon className="h-3 w-3" />,
                    className: 'bg-blue-500/10 text-blue-700 dark:text-blue-400 border-blue-500/20'
                };
            case 'production':
                return { show: false };
            default:
                return {
                    show: true,
                    text: 'Alpha',
                    variant: 'secondary' as const,
                    icon: <AlertTriangleIcon className="h-3 w-3" />,
                    className: 'bg-amber-500/10 text-amber-700 dark:text-amber-400 border-amber-500/20'
                };
        }
    };

    // Navigation items
    const navItems = [
        { path: "/dashboard", icon: <HomeIcon className="h-4 w-4" />, label: "Dashboard" },
        { path: "/subscriptions", icon: <CreditCardIcon className="h-4 w-4" />, label: "Subscriptions" },
        { path: "/providers", icon: <PackageIcon className="h-4 w-4" />, label: "Providers" },
        { path: "/family", icon: <UsersIcon className="h-4 w-4" />, label: "Family" },
        { path: "/labels", icon: <TagIcon className="h-4 w-4" />, label: "Labels" },
    ];

    // Get environment information
    const environmentInfo = getEnvironmentInfo();
    const env = envVar('VITE_TARGET_ENV') || 'development';

    return (
        <Sidebar collapsible="icon" className="border-r">
            {/* Environment Badge */}
            {environmentInfo.show && (
                <div className="w-full text-center py-2 px-2 transition-all duration-300">
                    {state === "expanded" ? (
                        <Badge
                            variant={environmentInfo.variant}
                            className={`w-full flex items-center justify-center gap-1.5 py-1.5 px-3 font-medium ${environmentInfo.className} transition-all duration-300`}
                        >
                            <span className="relative flex h-2 w-2">
                                <span className={`animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${env === 'development' ? 'bg-amber-400' : 'bg-blue-400'}`}></span>
                                <span className={`relative inline-flex rounded-full h-2 w-2 ${env === 'development' ? 'bg-amber-500' : 'bg-blue-500'}`}></span>
                            </span>
                            <span className="text-xs tracking-wide">{environmentInfo.text}</span>
                        </Badge>
                    ) : (
                        <TooltipProvider delayDuration={0}>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <div className="flex justify-center">
                                        <Badge
                                            variant={environmentInfo.variant}
                                            className={`p-1.5 ${environmentInfo.className} transition-all duration-300`}
                                        >
                                            <span className="relative flex h-2 w-2">
                                                <span className={`animate-ping absolute inline-flex h-full w-full rounded-full opacity-75 ${env === 'development' ? 'bg-amber-400' : 'bg-blue-400'}`}></span>
                                                <span className={`relative inline-flex rounded-full h-2 w-2 ${env === 'development' ? 'bg-amber-500' : 'bg-blue-500'}`}></span>
                                            </span>
                                        </Badge>
                                    </div>
                                </TooltipTrigger>
                                <TooltipContent side="right" hidden={isMobile}>
                                    <p className="font-medium">{environmentInfo.text} Environment</p>
                                </TooltipContent>
                            </Tooltip>
                        </TooltipProvider>
                    )}
                </div>
            )}

            {/* Header */}
            <SidebarHeader className="border-b bg-gradient-to-br from-background to-muted/20">
                {state === "expanded" ? (
                    <div className="flex items-center justify-between px-4 py-3">
                        <Link
                            to="/dashboard"
                            className="flex items-center gap-3 group transition-all duration-300 hover:scale-105"
                        >
                            <div className="bg-gradient-to-br from-primary to-primary/80 p-2.5 rounded-xl shadow-lg group-hover:shadow-primary/25 transition-all duration-300 group-hover:rotate-3">
                                <CreditCardIcon className="h-5 w-5 text-primary-foreground" />
                            </div>
                            <span className="font-bold text-lg bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text text-transparent">
                                SubTracker
                            </span>
                        </Link>
                        <SidebarTrigger className="hover:bg-muted rounded-lg transition-all duration-200 hover:scale-110" />
                    </div>
                ) : (
                    <div className="flex flex-col items-center justify-center py-3 space-y-2">
                        <Link
                            to="/dashboard"
                            className="group transition-all duration-300 hover:scale-110"
                        >
                            <div className="bg-gradient-to-br from-primary to-primary/80 p-2.5 rounded-xl shadow-lg group-hover:shadow-primary/25 transition-all duration-300 group-hover:rotate-3">
                                <CreditCardIcon className="h-5 w-5 text-primary-foreground" />
                            </div>
                        </Link>
                        <SidebarTrigger className="hover:bg-muted rounded-lg transition-all duration-200 hover:scale-110" />
                    </div>
                )}
            </SidebarHeader>

            <SidebarContent>
                <ScrollArea className="h-[calc(100vh-12rem)]">
                    <div className="px-3 py-4 space-y-1">
                        {state === "expanded" && (
                            <div className="mb-3 px-3">
                                <h2 className="text-xs font-semibold text-muted-foreground/70 uppercase tracking-wider">
                                    Navigation
                                </h2>
                            </div>
                        )}
                        <SidebarMenu className="space-y-1">
                            {navItems.map((item) => {
                                const isActive = location.pathname === item.path;
                                return (
                                    <SidebarMenuItem key={item.path}>
                                        <SidebarMenuButton
                                            asChild
                                            isActive={isActive}
                                            tooltip={item.label}
                                        >
                                            <Link
                                                to={item.path}
                                                className={`
                                                    group flex items-center gap-3 px-3 py-2.5 rounded-lg
                                                    transition-all duration-200 relative overflow-hidden
                                                    ${isActive
                                                        ? "bg-gradient-to-r from-primary/15 to-primary/5 text-primary font-semibold shadow-sm"
                                                        : "hover:bg-muted/50 text-muted-foreground hover:text-foreground"
                                                    }
                                                `}
                                            >
                                                {/* Active indicator bar */}
                                                {isActive && (
                                                    <span className="absolute left-0 top-1/2 -translate-y-1/2 h-8 w-1 bg-primary rounded-r-full shadow-lg shadow-primary/50" />
                                                )}

                                                {/* Icon with background */}
                                                <div className={`
                                                    relative flex items-center justify-center
                                                    transition-all duration-200
                                                    ${isActive
                                                        ? "text-primary scale-110"
                                                        : "group-hover:scale-110 group-hover:text-primary"
                                                    }
                                                `}>
                                                    {item.icon}
                                                </div>

                                                {/* Label */}
                                                <span className="flex-1 transition-all duration-200">
                                                    {item.label}
                                                </span>

                                                {/* Chevron on hover */}
                                                {state === "expanded" && (
                                                    <ChevronRight className={`
                                                        h-4 w-4 transition-all duration-200
                                                        ${isActive
                                                            ? "opacity-100 translate-x-0"
                                                            : "opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0"
                                                        }
                                                    `} />
                                                )}
                                            </Link>
                                        </SidebarMenuButton>
                                    </SidebarMenuItem>
                                );
                            })}
                        </SidebarMenu>
                    </div>
                </ScrollArea>
            </SidebarContent>

            <SidebarFooter className="p-4 border-t bg-gradient-to-br from-background to-muted/10">
                {state === "expanded" ? (
                    <div className="space-y-3">
                        {/* Settings Link Card */}
                        <Link
                            to="/profile"
                            className="flex items-center gap-3 p-3 rounded-lg bg-muted/50 hover:bg-muted transition-all duration-200 group border border-transparent hover:border-primary/20 hover:shadow-sm"
                        >
                            <div className="p-2 rounded-lg bg-background group-hover:bg-primary/10 transition-all duration-200">
                                <Settings className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-all duration-200 group-hover:rotate-90" />
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className="text-sm font-medium text-foreground">Preferences</p>
                                <p className="text-xs text-muted-foreground truncate">Manage your settings</p>
                            </div>
                            <ChevronRight className="h-4 w-4 text-muted-foreground group-hover:text-primary group-hover:translate-x-1 transition-all duration-200" />
                        </Link>

                        <Separator className="my-2" />

                        {/* User Profile Card */}
                        <div className="flex items-center gap-3 p-3 rounded-lg bg-gradient-to-br from-muted/30 to-muted/10 border border-border/50">
                            <div className="relative">
                                <UserButton />
                                <span className="absolute -bottom-0.5 -right-0.5 h-3 w-3 bg-green-500 border-2 border-background rounded-full" />
                            </div>
                            <div className="flex-1 min-w-0">
                                <p className="text-sm font-semibold text-foreground truncate">
                                    {user?.firstName} {user?.lastName}
                                </p>
                                <p className="text-xs text-muted-foreground truncate">
                                    {user?.primaryEmailAddress?.emailAddress}
                                </p>
                            </div>
                        </div>
                    </div>
                ) : (
                    <div className="flex flex-col items-center gap-3 py-2">
                        <TooltipProvider delayDuration={0}>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <Link
                                        to="/profile"
                                        className="p-2.5 rounded-lg bg-muted/50 hover:bg-muted transition-all duration-200 group hover:scale-110 border border-transparent hover:border-primary/20"
                                    >
                                        <Settings className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-all duration-200 group-hover:rotate-90" />
                                    </Link>
                                </TooltipTrigger>
                                <TooltipContent side="right" hidden={isMobile}>
                                    <p className="font-medium">Preferences</p>
                                </TooltipContent>
                            </Tooltip>
                        </TooltipProvider>

                        <Separator className="w-8" />

                        <TooltipProvider delayDuration={0}>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <div className="relative">
                                        <UserButton />
                                        <span className="absolute -bottom-0.5 -right-0.5 h-2.5 w-2.5 bg-green-500 border-2 border-background rounded-full" />
                                    </div>
                                </TooltipTrigger>
                                <TooltipContent side="right" hidden={isMobile}>
                                    <p className="font-medium">{user?.firstName} {user?.lastName}</p>
                                    <p className="text-xs text-muted-foreground">{user?.primaryEmailAddress?.emailAddress}</p>
                                </TooltipContent>
                            </Tooltip>
                        </TooltipProvider>
                    </div>
                )}
            </SidebarFooter>
        </Sidebar>
    );
}
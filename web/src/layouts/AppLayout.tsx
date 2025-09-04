import {type ReactNode} from "react";
import {SidebarProvider} from "@/components/ui/sidebar";
import {AppSidebar} from "@/components/AppSidebar";

interface AppLayoutProps {
    children: ReactNode;
}

export function AppLayout({children}: AppLayoutProps) {
    return (
        <SidebarProvider>
            <div
                className="flex min-h-screen w-full bg-white relative bg-fixed bg-dot-grid dark:bg-dot-grid-dark dark:bg-black"
            >
                <AppSidebar />
                <main className="flex-1 p-6 overflow-auto">
                    {children}
                </main>
            </div>
        </SidebarProvider>
    );
}
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
                className="min-h-screen w-full bg-white relative bg-fixed bg-dot-grid dark:bg-dot-grid-dark dark:bg-black"
            >
                <AppSidebar />
                <main className="flex-1 p-6 overflow-auto">
                    {children}
                </main>
            </div>
        </SidebarProvider>
    );
}

<div className="min-h-screen w-full bg-black relative">
    {/* Dark White Dotted Grid Background */}
    <div
        className="absolute inset-0 z-0"
        style={{
            background: "#000000",
            backgroundImage: `
        radial-gradient(circle, rgba(255, 255, 255, 0.2) 1.5px, transparent 1.5px)
      `,
            backgroundSize: "30px 30px",
            backgroundPosition: "0 0",
        }}
    />
    {/* Your Content/Components */}
</div>
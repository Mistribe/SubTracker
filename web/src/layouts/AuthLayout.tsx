import { type ReactNode } from "react";
import { ModeToggle } from "@/components/mode-toggle";

interface AuthLayoutProps {
  children: ReactNode;
  title: string;
  description: string;
}

export function AuthLayout({ children, title, description }: AuthLayoutProps) {
  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-br from-background to-muted/20">
      {/* Header */}
      <header className="container mx-auto p-4 flex justify-between items-center">
        <h1 className="text-2xl font-bold">Recurrent Payment Tracker</h1>
        <ModeToggle />
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8 flex-1 flex items-center justify-center">
        <div className="w-full max-w-md space-y-6">
          {/* Title and Description */}
          <div className="text-center space-y-2">
            <h2 className="text-3xl font-bold tracking-tight">{title}</h2>
            <p className="text-muted-foreground">{description}</p>
          </div>

          {/* Auth Component Container */}
          <div className="p-6">
            {children}
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="container mx-auto px-4 py-6 border-t">
        <p className="text-center text-sm text-muted-foreground">
          © {new Date().getFullYear()} Recurrent Payment Tracker. All rights reserved.
        </p>
      </footer>
    </div>
  );
}
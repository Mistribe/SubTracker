import type { Appearance } from "@clerk/types";
import { dark } from "@clerk/themes";

/**
 * Get Clerk appearance configuration that matches the application theme
 * @param theme - The current theme ('light', 'dark', or 'system')
 * @returns Clerk appearance configuration
 */
export function getClerkAppearance(theme: string): Appearance {
    // Determine the base theme for Clerk
    const resolvedTheme = theme === 'system'
        ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
        : theme;

    return {
        baseTheme: resolvedTheme === 'dark' ? dark : undefined,
        variables: {
            colorPrimary: "hsl(var(--primary))",
            borderRadius: "0.5rem",
            fontFamily: "inherit",
        },
        elements: {
            card: {
                backgroundColor: "transparent",
                border: "none",
                boxShadow: "none",
            },
            formButtonPrimary: {
                backgroundColor: "hsl(var(--primary))",
                color: "hsl(var(--primary-foreground))",
                "&:hover": {
                    backgroundColor: "hsl(var(--primary) / 0.9)",
                },
            },
            // Hide the redundant header and subtitle
            headerTitle: {
                display: "none",
            },
            headerSubtitle: {
                display: "none",
            },
        },
    };
}

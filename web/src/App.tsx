import { ThemeProvider } from "@/components/theme-provider"

function App() {
    return (
        <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
            <div>Your content here</div>
        </ThemeProvider>
    )
}

export default App
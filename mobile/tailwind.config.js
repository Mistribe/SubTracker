/** @type {import('tailwindcss').Config} */
module.exports = {
    // NOTE: Update this to include the paths to all files that contain Nativewind classes.
    content: [
        "./App.tsx",
        "./app/**/*.{js,jsx,ts,tsx}",
        "./components/**/*.{js,jsx,ts,tsx}",
        "./lib/**/*.{js,jsx,ts,tsx}",
    ],
    presets: [require("nativewind/preset")],
    // Use class strategy so NativeWindStyleSheet.setColorScheme('dark'|'light') works
    darkMode: "class",
    theme: {
        extend: {},
    },
    plugins: [],
}
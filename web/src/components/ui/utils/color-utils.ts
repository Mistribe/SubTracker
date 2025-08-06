// Utility functions for color conversion
export const hexToArgb = (hex: string): string => {
    // Default alpha value (fully opaque)
    const alpha = "FF";

    // Remove # if present
    const cleanHex = hex.startsWith("#") ? hex.slice(1) : hex;

    // Ensure we have a valid 6-character RGB hex
    const validRgb = cleanHex.length === 3
        ? cleanHex.split('').map(c => c + c).join('') // Convert 3-char hex to 6-char
        : cleanHex.padEnd(6, '0').substring(0, 6); // Ensure 6 characters

    // Return in ARGB format
    return `#${alpha}${validRgb}`;
};

export const argbToHex = (argb: string): string => {
    // Remove # if present
    const cleanArgb = argb.startsWith("#") ? argb.slice(1) : argb;

    // Extract RGB part (last 6 characters)
    // In ARGB format, the first 2 characters are alpha, followed by 6 characters for RGB
    const rgb = cleanArgb.length >= 8 ? cleanArgb.substring(2) : cleanArgb;

    // Return in hex format
    return `#${rgb}`;
};

export const getAlphaFromArgb = (argb: string): number => {
    // Remove # if present
    const cleanArgb = argb.startsWith("#") ? argb.slice(1) : argb;

    // Extract alpha part (first 2 characters if length is 8)
    const alpha = cleanArgb.length >= 8 ? cleanArgb.substring(0, 2) : "FF";

    // Convert hex alpha to decimal (0-255) then to percentage (0-100)
    return (parseInt(alpha, 16) / 255) * 100;
};

export const setAlphaInArgb = (argb: string, alphaPercent: number): string => {
    // Remove # if present
    const cleanArgb = argb.startsWith("#") ? argb.slice(1) : argb;

    // Extract RGB part (last 6 characters or all if less than 8)
    const rgb = cleanArgb.length >= 8 ? cleanArgb.substring(2) : cleanArgb;

    // Convert alpha percentage to hex
    const alpha = Math.round((alphaPercent / 100) * 255)
        .toString(16)
        .padStart(2, "0")
        .toUpperCase();

    // Return in ARGB format
    return `#${alpha}${rgb}`;
};

// Convert ARGB string to CSS rgba format
export const argbToRgba = (argb: string | undefined): string => {
    if (!argb) return "rgba(0, 0, 0, 0)";
    // Remove # if present
    const cleanArgb = argb.startsWith("#") ? argb.slice(1) : argb;

    // Default values if format is invalid
    let r = 0, g = 0, b = 0, a = 1;

    if (cleanArgb.length >= 8) {
        // ARGB format: AARRGGBB
        a = parseInt(cleanArgb.substring(0, 2), 16) / 255;
        r = parseInt(cleanArgb.substring(2, 4), 16);
        g = parseInt(cleanArgb.substring(4, 6), 16);
        b = parseInt(cleanArgb.substring(6, 8), 16);
    } else if (cleanArgb.length >= 6) {
        // RGB format: RRGGBB
        r = parseInt(cleanArgb.substring(0, 2), 16);
        g = parseInt(cleanArgb.substring(2, 4), 16);
        b = parseInt(cleanArgb.substring(4, 6), 16);
    }

    return `rgba(${r}, ${g}, ${b}, ${a})`;
};
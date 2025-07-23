import React from "react";
import {Input} from "./input";
import {Slider} from "./slider";
import {
    argbToHex,
    getAlphaFromArgb,
    setAlphaInArgb,
    argbToRgba
} from "./utils/color-utils";

// ColorPicker component
export interface ColorPickerProps {
    color: string;
    onChange: (color: string) => void;
}

export const ColorPicker = ({color, onChange}: ColorPickerProps) => {
    // Extract hex color (without alpha) for the color input
    const hexColor = argbToHex(color);

    // Get alpha value as percentage (0-100)
    const alphaPercent = getAlphaFromArgb(color);

    // Handle color change from input
    const handleColorChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const newHexColor = e.target.value;
        const newColor = setAlphaInArgb(newHexColor, alphaPercent);
        // Preserve the current alpha when changing the color
        onChange(newColor);
    };

    // Handle alpha change from slider
    const handleAlphaChange = (value: number[]) => {
        const newAlphaPercent = value[0];
        onChange(setAlphaInArgb(color, newAlphaPercent));
    };

    return (
        <div className="flex flex-col gap-4 p-4 w-64">
            <div className="flex items-center gap-2">
                <div
                    className="w-8 h-8 rounded-md border"
                    style={{backgroundColor: argbToRgba(color)}}
                />
                <Input
                    type="color"
                    value={hexColor}
                    onChange={handleColorChange}
                    className="w-full h-8"
                />
            </div>
            <div className="space-y-2">
                <div className="flex justify-between">
                    <span className="text-sm">Opacity</span>
                    <span className="text-sm">{Math.round(alphaPercent)}%</span>
                </div>
                <Slider
                    defaultValue={[alphaPercent]}
                    max={100}
                    step={1}
                    onValueChange={handleAlphaChange}
                />
            </div>
            <Input
                value={color.toUpperCase()}
                onChange={(e) => {
                    // Validate the input is a proper ARGB format before updating
                    const validColorRegex = /^#[0-9A-Fa-f]{8}$/;
                    if (validColorRegex.test(e.target.value)) {
                        onChange(e.target.value);
                    }
                }}
                placeholder="ARGB Color (#AARRGGBB)"
                className="font-mono text-sm"
            />
        </div>
    );
};
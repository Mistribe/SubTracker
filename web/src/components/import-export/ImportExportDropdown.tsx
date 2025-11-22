import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ChevronDown, Upload, Download } from "lucide-react";

interface ImportExportDropdownProps {
    onImport: () => void;
    onExport: () => void;
    buttonText?: string;
    disabled?: boolean;
}

export function ImportExportDropdown({
    onImport,
    onExport,
    buttonText = "Import / Export",
    disabled = false,
}: ImportExportDropdownProps) {
    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant="outline" disabled={disabled}>
                    {buttonText}
                    <ChevronDown className="ml-2 h-4 w-4" />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={onImport}>
                    <Upload className="mr-2 h-4 w-4" />
                    Import from file
                </DropdownMenuItem>
                <DropdownMenuItem onClick={onExport}>
                    <Download className="mr-2 h-4 w-4" />
                    Export to file
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    );
}

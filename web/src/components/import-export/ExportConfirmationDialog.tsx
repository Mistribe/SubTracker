import { useState, useEffect } from "react";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import type { ExportFormat } from "@/services/exportService";
import { Loader2 } from "lucide-react";

interface ExportConfirmationDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onConfirm: (format: ExportFormat) => void;
    entityType: "labels" | "subscriptions" | "providers";
    isExporting?: boolean;
}

const FUNNY_MESSAGES = [
    "Convincing the hamsters to run faster...",
    "Summoning the data from the digital void...",
    "Asking the database nicely...",
    "Consulting the ancient scrolls...",
    "Doing some computational gymnastics...",
    "Bribing the servers with cookies...",
    "Untangling the data spaghetti...",
    "Waking up the sleepy bytes...",
    "Polishing the bits and bytes...",
    "Teaching the data to stand in line...",
    "Counting all the electrons...",
    "Converting coffee into data...",
    "Herding the digital cats...",
    "Negotiating with the cloud...",
    "Spinning up the magic...",
    "Reticulating splines...",
    "Gathering stardust from the servers...",
    "Asking the data gnomes for permission...",
    "Translating 1s and 0s to human...",
    "Feeding the algorithm some snacks...",
    "Performing digital acrobatics...",
    "Downloading more RAM... just kidding!",
    "Convincing bits to play nice...",
    "Calibrating the flux capacitor...",
    "Organizing the chaos into neat rows...",
];

export function ExportConfirmationDialog({
    open,
    onOpenChange,
    onConfirm,
    entityType,
    isExporting = false,
}: ExportConfirmationDialogProps) {
    const [selectedFormat, setSelectedFormat] = useState<ExportFormat>("csv");
    const [funnyMessage, setFunnyMessage] = useState("");
    const [messageIndex, setMessageIndex] = useState(0);

    // Initialize first funny message when export starts
    useEffect(() => {
        if (isExporting) {
            const randomIndex = Math.floor(Math.random() * FUNNY_MESSAGES.length);
            setMessageIndex(randomIndex);
            setFunnyMessage(FUNNY_MESSAGES[randomIndex]);
        }
    }, [isExporting]);

    // Rotate through funny messages while exporting
    useEffect(() => {
        if (!isExporting) return;

        const interval = setInterval(() => {
            setMessageIndex((prevIndex) => {
                // Get a different random message
                let newIndex;
                do {
                    newIndex = Math.floor(Math.random() * FUNNY_MESSAGES.length);
                } while (newIndex === prevIndex && FUNNY_MESSAGES.length > 1);

                setFunnyMessage(FUNNY_MESSAGES[newIndex]);
                return newIndex;
            });
        }, 2000); // Change message every 2 seconds

        return () => clearInterval(interval);
    }, [isExporting]);

    const handleConfirm = () => {
        onConfirm(selectedFormat);
    };

    const handleOpenChange = (newOpen: boolean) => {
        // Prevent closing the dialog while exporting
        if (isExporting) return;
        onOpenChange(newOpen);
    };

    return (
        <AlertDialog open={open} onOpenChange={handleOpenChange}>
            <AlertDialogContent>
                <AlertDialogHeader>
                    <AlertDialogTitle>Export {entityType}</AlertDialogTitle>
                    <AlertDialogDescription>
                        All {entityType} will be exported regardless of the current filter.
                        This action will download a file containing all your {entityType} data.
                    </AlertDialogDescription>
                </AlertDialogHeader>

                {isExporting ? (
                    <div className="flex flex-col items-center justify-center py-6 space-y-4">
                        <Loader2 className="h-8 w-8 animate-spin text-primary" />
                        <p
                            key={messageIndex}
                            className="text-sm text-muted-foreground italic animate-in fade-in-50 duration-300"
                        >
                            {funnyMessage}
                        </p>
                    </div>
                ) : (
                    <div className="py-4">
                        <Label className="text-sm font-medium mb-3 block">
                            Choose export format:
                        </Label>
                        <RadioGroup
                            value={selectedFormat}
                            onValueChange={(value) => setSelectedFormat(value as ExportFormat)}
                            className="space-y-2"
                        >
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="csv" id="csv" />
                                <Label htmlFor="csv" className="font-normal cursor-pointer">
                                    CSV - Best for spreadsheets (Excel, Google Sheets)
                                </Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="json" id="json" />
                                <Label htmlFor="json" className="font-normal cursor-pointer">
                                    JSON - Best for developers and APIs
                                </Label>
                            </div>
                            <div className="flex items-center space-x-2">
                                <RadioGroupItem value="yaml" id="yaml" />
                                <Label htmlFor="yaml" className="font-normal cursor-pointer">
                                    YAML - Human-readable configuration format
                                </Label>
                            </div>
                        </RadioGroup>
                    </div>
                )}

                <AlertDialogFooter>
                    <AlertDialogCancel disabled={isExporting}>Cancel</AlertDialogCancel>
                    <AlertDialogAction onClick={handleConfirm} disabled={isExporting}>
                        {isExporting ? "Exporting..." : "Export"}
                    </AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    );
}

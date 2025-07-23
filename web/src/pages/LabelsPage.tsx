import {useState} from "react";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle
} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {Badge} from "@/components/ui/badge";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover";
import {Slider} from "@/components/ui/slider";
import {PlusIcon, Pencil, X, Loader2} from "lucide-react";
import {useApiClient} from "@/hooks/use-api-client.ts";
import {useQuery, useMutation, useQueryClient} from "@tanstack/react-query";
import type {LabelModel} from "@/api/models";

// Utility functions for color conversion
const hexToArgb = (hex: string): string => {
    // Default alpha value (fully opaque)
    const alpha = "FF";
    
    // Remove # if present
    const cleanHex = hex.startsWith("#") ? hex.slice(1) : hex;
    
    // Return in ARGB format
    return `#${alpha}${cleanHex}`;
};

const argbToHex = (argb: string): string => {
    // Remove # if present
    const cleanArgb = argb.startsWith("#") ? argb.slice(1) : argb;
    
    // Extract RGB part (last 6 characters)
    const rgb = cleanArgb.length >= 8 ? cleanArgb.slice(2) : cleanArgb;
    
    // Return in hex format
    return `#${rgb}`;
};

const getAlphaFromArgb = (argb: string): number => {
    // Remove # if present
    const cleanArgb = argb.startsWith("#") ? argb.slice(1) : argb;
    
    // Extract alpha part (first 2 characters if length is 8)
    const alpha = cleanArgb.length >= 8 ? cleanArgb.slice(0, 2) : "FF";
    
    // Convert hex alpha to decimal (0-255) then to percentage (0-100)
    return (parseInt(alpha, 16) / 255) * 100;
};

const setAlphaInArgb = (argb: string, alphaPercent: number): string => {
    // Remove # if present
    const cleanArgb = argb.startsWith("#") ? argb.slice(1) : argb;
    
    // Extract RGB part (last 6 characters or all if less than 8)
    const rgb = cleanArgb.length >= 8 ? cleanArgb.slice(2) : cleanArgb;
    
    // Convert alpha percentage to hex
    const alpha = Math.round((alphaPercent / 100) * 255)
        .toString(16)
        .padStart(2, "0")
        .toUpperCase();
    
    // Return in ARGB format
    return `#${alpha}${rgb}`;
};

// ColorPicker component
interface ColorPickerProps {
    color: string;
    onChange: (color: string) => void;
}

const ColorPicker = ({ color, onChange }: ColorPickerProps) => {
    // Extract hex color (without alpha) for the color input
    const hexColor = argbToHex(color);
    
    // Get alpha value as percentage (0-100)
    const alphaPercent = getAlphaFromArgb(color);
    
    // Handle color change from input
    const handleColorChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const newHexColor = e.target.value;
        // Preserve the current alpha when changing the color
        onChange(setAlphaInArgb(newHexColor, alphaPercent));
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
                    style={{ backgroundColor: color }}
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
                onChange={(e) => onChange(e.target.value)}
                placeholder="ARGB Color"
                className="font-mono text-sm"
            />
        </div>
    );
};

const LabelsPage = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    // Fetch labels using TanStack Query
    const {
        data: labels = [],
        isLoading,
        error
    } = useQuery({
        queryKey: ['labels'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }
            console.log('Fetching labels...');
            console.log(apiClient);
            const result = await apiClient?.labels.get({queryParameters: {withDefault: true}});
            console.log(result);
            return result || [];
        },
        enabled: !!apiClient,
        // Optional but useful settings
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });

    // State for new label
    const [newLabel, setNewLabel] = useState("");
    const [newLabelColor, setNewLabelColor] = useState("#FF000000"); // Default ARGB color
    
    // State for editing
    const [editingId, setEditingId] = useState<number | null>(null);
    const [editingName, setEditingName] = useState("");
    const [editingColor, setEditingColor] = useState("");

    // Create label mutation
    const createLabelMutation = useMutation({
        mutationFn: async (labelData: { name: string, color: string }) => {
            return apiClient?.labels.post({
                name: labelData.name,
                color: labelData.color
            });
        },
        onSuccess: () => {
            // Invalidate and refetch
            queryClient.invalidateQueries({queryKey: ['labels']});
            setNewLabel("");
        }
    });

    // Update label mutation
    const updateLabelMutation = useMutation({
        mutationFn: async ({id, name, color}: { id: number, name: string, color: string }) => {
            return apiClient?.labels.byId(id.toString()).put({
                name: name,
                color: color,
            });
        },
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: ['labels']});
            setEditingId(null);
            setEditingName("");
        }
    });

    // Delete label mutation
    const deleteLabelMutation = useMutation({
        mutationFn: async (id: number) => {
            return apiClient?.labels.byId(id.toString()).delete();
        },
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: ['labels']});
        }
    });

    const handleAddLabel = () => {
        if (newLabel.trim()) {
            createLabelMutation.mutate({
                name: newLabel,
                color: newLabelColor
            });
        }
    };

    const startEditing = (label: LabelModel) => {
        setEditingId(label.id);
        setEditingName(label.name);
        // Convert hex color to ARGB if needed
        const color = label.color.startsWith('#') && label.color.length === 7 
            ? hexToArgb(label.color) 
            : label.color;
        setEditingColor(color);
    };

    const saveEdit = () => {
        if (editingId && editingName.trim()) {
            updateLabelMutation.mutate({
                id: editingId,
                name: editingName,
                color: editingColor
            });
        }
    };

    const cancelEdit = () => {
        setEditingId(null);
        setEditingName("");
    };

    const deleteLabel = (id: number) => {
        deleteLabelMutation.mutate(id);
    };

    // Show loading or error states
    if (isLoading) {
        return (
            <div className="flex flex-col items-center justify-center h-64">
                <Loader2 className="h-8 w-8 animate-spin text-primary mb-4"/>
                <p className="text-muted-foreground">Loading labels...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex flex-col items-center justify-center h-64">
                <p className="text-destructive mb-2">Error loading labels</p>
                <p className="text-muted-foreground text-sm">{error instanceof Error ? error.message : 'Unknown error'}</p>
            </div>
        );
    }

    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Labels</h1>
            </div>

            <Card className="mb-6">
                <CardHeader>
                    <CardTitle>Add New Label</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="flex flex-col gap-4">
                        <div className="flex gap-2">
                            <Input
                                placeholder="Enter label name"
                                value={newLabel}
                                onChange={(e) => setNewLabel(e.target.value)}
                                className="max-w-sm"
                            />
                            <Popover>
                                <PopoverTrigger asChild>
                                    <Button variant="outline" className="w-10 p-0">
                                        <div 
                                            className="w-6 h-6 rounded-md"
                                            style={{ backgroundColor: newLabelColor }}
                                        />
                                    </Button>
                                </PopoverTrigger>
                                <PopoverContent className="w-auto p-0" align="start">
                                    <ColorPicker 
                                        color={newLabelColor} 
                                        onChange={setNewLabelColor} 
                                    />
                                </PopoverContent>
                            </Popover>
                            <Button
                                onClick={handleAddLabel}
                                disabled={createLabelMutation.isPending}
                            >
                                {createLabelMutation.isPending ? (
                                    <Loader2 className="h-4 w-4 mr-2 animate-spin"/>
                                ) : (
                                    <PlusIcon className="h-4 w-4 mr-2"/>
                                )}
                                Add Label
                            </Button>
                        </div>
                        <div className="flex items-center gap-2">
                            <div 
                                className="w-4 h-4 rounded-full"
                                style={{ backgroundColor: newLabelColor }}
                            />
                            <span className="text-sm text-muted-foreground">
                                Preview with selected color
                            </span>
                        </div>
                    </div>
                </CardContent>
            </Card>

            <Card>
                <CardHeader>
                    <CardTitle>Manage Labels</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid gap-4">
                        {labels.map((label) => (
                            <div
                                key={label.id}
                                className="flex items-center justify-between p-3 border rounded-md"
                            >
                                <div className="flex items-center gap-3">
                                    <div
                                        className="w-4 h-4 rounded-full"
                                        style={{backgroundColor: label.color}}
                                    />

                                    {editingId === label.id ? (
                                        <div className="flex items-center gap-2">
                                            <Input
                                                value={editingName}
                                                onChange={(e) => setEditingName(e.target.value)}
                                                className="max-w-xs"
                                                autoFocus
                                            />
                                            <Popover>
                                                <PopoverTrigger asChild>
                                                    <Button variant="outline" className="w-10 p-0">
                                                        <div 
                                                            className="w-6 h-6 rounded-md"
                                                            style={{ backgroundColor: editingColor }}
                                                        />
                                                    </Button>
                                                </PopoverTrigger>
                                                <PopoverContent className="w-auto p-0" align="start">
                                                    <ColorPicker 
                                                        color={editingColor} 
                                                        onChange={setEditingColor} 
                                                    />
                                                </PopoverContent>
                                            </Popover>
                                        </div>
                                    ) : (
                                        <span className="font-medium">{label.name}</span>
                                    )}

                                    <Badge variant="secondary">{label.count} items</Badge>
                                </div>

                                <div className="flex items-center gap-2">
                                    {editingId === label.id ? (
                                        <>
                                            <Button
                                                variant="ghost"
                                                size="sm"
                                                onClick={saveEdit}
                                                disabled={updateLabelMutation.isPending}
                                            >
                                                {updateLabelMutation.isPending ? (
                                                    <Loader2 className="h-3 w-3 mr-2 animate-spin"/>
                                                ) : null}
                                                Save
                                            </Button>
                                            <Button variant="ghost" size="sm" onClick={cancelEdit}>
                                                Cancel
                                            </Button>
                                        </>
                                    ) : (
                                        <>
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() => startEditing(label)}
                                            >
                                                <Pencil className="h-4 w-4"/>
                                            </Button>
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                onClick={() => deleteLabel(label.id)}
                                                className="text-destructive hover:text-destructive"
                                                disabled={deleteLabelMutation.isPending && deleteLabelMutation.variables === label.id}
                                            >
                                                {deleteLabelMutation.isPending && deleteLabelMutation.variables === label.id ? (
                                                    <Loader2 className="h-4 w-4 animate-spin"/>
                                                ) : (
                                                    <X className="h-4 w-4"/>
                                                )}
                                            </Button>
                                        </>
                                    )}
                                </div>
                            </div>
                        ))}
                    </div>
                </CardContent>
            </Card>
        </div>
    );
};

export default LabelsPage;
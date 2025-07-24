import {useState} from "react";
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle
} from "@/components/ui/card";
import {Button} from "@/components/ui/button";
import {Input} from "@/components/ui/input";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover";
import {ColorPicker} from "@/components/ui/color-picker";
import {hexToArgb, argbToRgba} from "@/components/ui/utils/color-utils";
import {PlusIcon, Pencil, X, Loader2} from "lucide-react";
import {useApiClient} from "@/hooks/use-api-client.ts";
import {useQuery, useMutation, useQueryClient} from "@tanstack/react-query";
import type {LabelModel} from "@/api/models";
import Label from "@/models/label.ts";

const LabelsPage = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    const [page] = useState(1);
    const [pageSize] = useState(10);

    // Fetch labels using TanStack Query
    const {
        data: queryResponse,
        isLoading,
        error
    } = useQuery({
        queryKey: ['labels'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }
            const result = await apiClient?.labels.get({
                queryParameters: {
                    withDefault: true,
                    page: page,
                    size: pageSize
                }
            });
            if (result && result.data) {
                return {
                    labels: result.data.map((model: LabelModel) => {
                        return Label.fromModel(model);
                    }),
                    length: result.data.length,
                    total: result.total ?? 0
                }
            }
            return {labels: [], length: 0, total: 0};
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
    const [editingId, setEditingId] = useState<string | null>(null);
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
        mutationFn: async ({id, name, color}: { id: string, name: string, color: string }) => {
            return apiClient?.labels.byId(id).put({
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
        mutationFn: async (id: string) => {
            return apiClient?.labels.byId(id).delete();
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

    const startEditing = (label: Label) => {
        // Prevent editing default labels
        if (label.isDefault) {
            return;
        }
        
        setEditingId(label.id);
        setEditingName(label.name);
        // Convert hex color to ARGB if needed
        let color = label.color;

        // If it's a standard hex color (#RRGGBB)
        if (label.color.startsWith('#') && label.color.length === 7) {
            color = hexToArgb(label.color);
        }
        // If it's not in ARGB format (#AARRGGBB), make sure it is
        else if (!label.color.startsWith('#') || label.color.length !== 9) {
            // Default to fully opaque black if invalid format
            color = hexToArgb("#000000");
        }

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

    const deleteLabel = (id: string) => {
        // Find the label to check if it's a default label
        const labelToDelete = queryResponse?.labels?.find(label => label.id === id);
        
        // Prevent deleting default labels
        if (labelToDelete?.isDefault) {
            return;
        }
        
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
                <CardHeader className="py-4">
                    <CardTitle className="text-lg">Add New Label</CardTitle>
                </CardHeader>
                <CardContent className="py-2">
                    <div className="flex flex-col sm:flex-row gap-3 items-start sm:items-center">
                        <div className="flex items-center gap-2 w-full sm:w-auto">
                            <div
                                className="w-4 h-4 rounded-full flex-shrink-0"
                                style={{backgroundColor: argbToRgba(newLabelColor)}}
                            />
                            <Input
                                placeholder="Enter label name"
                                value={newLabel}
                                onChange={(e) => setNewLabel(e.target.value)}
                                className="h-9 text-sm w-full sm:w-48 md:w-64"
                            />
                        </div>
                        <div className="flex items-center gap-2">
                            <Popover>
                                <PopoverTrigger asChild>
                                    <Button variant="outline" size="sm" className="h-9">
                                        <div
                                            className="w-4 h-4 rounded-md mr-1"
                                            style={{backgroundColor: argbToRgba(newLabelColor)}}
                                        />
                                        Color
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
                                size="sm"
                                className="h-9"
                                onClick={handleAddLabel}
                                disabled={createLabelMutation.isPending}
                            >
                                {createLabelMutation.isPending ? (
                                    <Loader2 className="h-4 w-4 mr-1 animate-spin"/>
                                ) : (
                                    <PlusIcon className="h-4 w-4 mr-1"/>
                                )}
                                Add
                            </Button>
                        </div>
                    </div>
                </CardContent>
            </Card>

            <div className="mt-8">
                <h2 className="text-xl font-semibold mb-4">Labels</h2>
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-3">
                    {queryResponse?.labels?.map((label) => (
                        <div
                            key={label.id}
                            className={`flex items-center p-2 border rounded-md ${editingId === label.id ? 'bg-muted' : 'hover:bg-muted/50'}`}
                        >
                            {editingId === label.id ? (
                                <div className="w-full space-y-2">
                                    <div className="flex items-center gap-2">
                                        <div
                                            className="w-4 h-4 rounded-full flex-shrink-0"
                                            style={{backgroundColor: argbToRgba(editingColor)}}
                                        />
                                        <Input
                                            value={editingName}
                                            onChange={(e) => setEditingName(e.target.value)}
                                            className="h-7 text-sm"
                                            autoFocus
                                        />
                                    </div>
                                    <div className="flex items-center gap-1 justify-between">
                                        <Popover>
                                            <PopoverTrigger asChild>
                                                <Button variant="outline" size="sm" className="h-7 px-2">
                                                    <div
                                                        className="w-4 h-4 rounded-md mr-1"
                                                        style={{backgroundColor: argbToRgba(editingColor)}}
                                                    />
                                                    Color
                                                </Button>
                                            </PopoverTrigger>
                                            <PopoverContent className="w-auto p-0" align="start">
                                                <ColorPicker
                                                    color={editingColor}
                                                    onChange={setEditingColor}
                                                />
                                            </PopoverContent>
                                        </Popover>
                                        <div className="flex items-center gap-1">
                                            <Button
                                                variant="ghost"
                                                size="sm"
                                                className="h-7 px-2"
                                                onClick={saveEdit}
                                                disabled={updateLabelMutation.isPending}
                                            >
                                                {updateLabelMutation.isPending ? (
                                                    <Loader2 className="h-3 w-3 mr-1 animate-spin"/>
                                                ) : null}
                                                Save
                                            </Button>
                                            <Button 
                                                variant="ghost" 
                                                size="sm" 
                                                className="h-7 px-2"
                                                onClick={cancelEdit}
                                            >
                                                Cancel
                                            </Button>
                                        </div>
                                    </div>
                                </div>
                            ) : (
                                <>
                                    <div
                                        className="w-3 h-3 rounded-full mr-2 flex-shrink-0"
                                        style={{backgroundColor: argbToRgba(label.color)}}
                                    />
                                    <span className="text-sm font-medium truncate flex-grow">{label.name}</span>
                                    <div className="flex items-center gap-1 ml-1">
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            className="h-6 w-6"
                                            onClick={() => startEditing(label)}
                                            disabled={label.isDefault}
                                            title={label.isDefault ? "Default labels cannot be edited" : "Edit label"}
                                        >
                                            <Pencil className="h-3 w-3"/>
                                        </Button>
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            className="h-6 w-6 text-destructive hover:text-destructive"
                                            onClick={() => deleteLabel(label.id)}
                                            disabled={(deleteLabelMutation.isPending && deleteLabelMutation.variables === label.id) || label.isDefault}
                                            title={label.isDefault ? "Default labels cannot be deleted" : "Delete label"}
                                        >
                                            {deleteLabelMutation.isPending && deleteLabelMutation.variables === label.id ? (
                                                <Loader2 className="h-3 w-3 animate-spin"/>
                                            ) : (
                                                <X className="h-3 w-3"/>
                                            )}
                                        </Button>
                                    </div>
                                </>
                            )}
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default LabelsPage;
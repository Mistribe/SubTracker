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

    const [page, setPage] = useState(1);
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

    const startEditing = (label: Label) => {
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
                                            style={{backgroundColor: argbToRgba(newLabelColor)}}
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
                                style={{backgroundColor: argbToRgba(newLabelColor)}}
                            />
                            <span className="text-sm text-muted-foreground">
                                Preview with selected color
                            </span>
                        </div>
                    </div>
                </CardContent>
            </Card>

            <div className="mt-8">
                <h2 className="text-xl font-semibold mb-4">Labels</h2>
                <div className="grid gap-4">
                    {queryResponse?.labels?.map((label) => (
                        <div
                            key={label.id}
                            className="flex items-center justify-between p-3 border rounded-md"
                        >
                            <div className="flex items-center gap-3">
                                <div
                                    className="w-4 h-4 rounded-full"
                                    style={{backgroundColor: argbToRgba(label.color)}}
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
                                                        style={{backgroundColor: argbToRgba(editingColor)}}
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
            </div>
        </div>
    );
};

export default LabelsPage;
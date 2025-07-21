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
import {PlusIcon, Pencil, X, Loader2} from "lucide-react";
import {useApiClient} from "@/hooks/use-api-client.ts";
import {useQuery, useMutation, useQueryClient} from "@tanstack/react-query";
import type {LabelModel} from "@/api/models";

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

    const [newLabel, setNewLabel] = useState("");
    const [editingId, setEditingId] = useState<number | null>(null);
    const [editingName, setEditingName] = useState("");

    // Create label mutation
    const createLabelMutation = useMutation({
        mutationFn: async (labelData: { name: string, color: string }) => {
            return await apiClient?.labels.post({
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
        mutationFn: async ({id, name}: { id: number, name: string }) => {
            return await apiClient?.labels.byId(id.toString()).patch({
                name: name
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
            return await apiClient?.labels.byId(id.toString()).delete();
        },
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: ['labels']});
        }
    });

    const handleAddLabel = () => {
        if (newLabel.trim()) {
            const randomColor = `#${Math.floor(Math.random() * 16777215).toString(16)}`;
            createLabelMutation.mutate({
                name: newLabel,
                color: randomColor
            });
        }
    };

    const startEditing = (label: LabelModel) => {
        setEditingId(label.id);
        setEditingName(label.name);
    };

    const saveEdit = () => {
        if (editingId && editingName.trim()) {
            updateLabelMutation.mutate({
                id: editingId,
                name: editingName
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
                    <div className="flex gap-2">
                        <Input
                            placeholder="Enter label name"
                            value={newLabel}
                            onChange={(e) => setNewLabel(e.target.value)}
                            className="max-w-sm"
                        />
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
                                        <Input
                                            value={editingName}
                                            onChange={(e) => setEditingName(e.target.value)}
                                            className="max-w-xs"
                                            autoFocus
                                        />
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
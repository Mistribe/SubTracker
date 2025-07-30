import {useState} from "react";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import {Button} from "@/components/ui/button";
import {CheckIcon, Loader2, XIcon} from "lucide-react";
import Family from "@/models/family.ts";
import {useApiClient} from "@/hooks/use-api-client.ts";
import {useQuery, useQueryClient} from "@tanstack/react-query";
import type {FamilyModel, PatchFamilyModel} from "@/api/models";
import {CreateFamilyDialog} from "@/components/families/CreateFamilyDialog.tsx";
import {Input} from "@/components/ui/input";
import {AddFamilyMemberDialog} from "@/components/families/AddFamilyMemberDialog.tsx";

const FamiliesPage = () => {
    const {apiClient} = useApiClient();
    const queryClient = useQueryClient();

    const [page] = useState(1);
    const [pageSize] = useState(10);
    const [editingFamilyId, setEditingFamilyId] = useState<string | null>(null);
    const [editedName, setEditedName] = useState<string>("");
    const [isUpdating, setIsUpdating] = useState<boolean>(false);

    // Function to start editing a family
    const startEditing = (family: Family) => {
        setEditingFamilyId(family.id);
        setEditedName(family.name);
    };

    // Function to cancel editing
    const cancelEditing = () => {
        setEditingFamilyId(null);
        setEditedName("");
    };

    // Function to save changes
    const saveChanges = async (family: Family) => {
        if (!apiClient) return;

        try {
            setIsUpdating(true);

            // Only update if values have changed
            if (editedName === family.name) {
                cancelEditing();
                return;
            }

            const patchModel: Partial<PatchFamilyModel> = {
                id: family.id,
                name: editedName,
                updatedAt: new Date()
            };

            await apiClient.families.patch(patchModel);

            // Invalidate and refetch the families query
            await queryClient.invalidateQueries({queryKey: ['families']});

            // Reset editing state
            cancelEditing();
        } catch (error) {
            console.error("Failed to update family:", error);
        } finally {
            setIsUpdating(false);
        }
    };

    // Fetch families using TanStack Query
    const {
        data: queryResponse,
        isLoading,
        error
    } = useQuery({
        queryKey: ['families'],
        queryFn: async () => {
            if (!apiClient) {
                throw new Error('API client not initialized');
            }
            const result = await apiClient.families.get({
                queryParameters: {
                    page: page,
                    size: pageSize
                }
            });
            if (result && result.data) {
                return {
                    families: result.data.map((model: FamilyModel) => {
                        return Family.fromModel(model);
                    }),
                    length: result.data.length,
                    total: result.total ?? 0
                }
            }
            return {families: [], length: 0, total: 0};
        },
        enabled: !!apiClient,
        staleTime: 5 * 60 * 1000, // 5 minutes
        refetchOnWindowFocus: true,
    });

    // Loading state
    if (isLoading) {
        return (
            <div className="flex justify-center items-center h-64">
                <Loader2 className="h-8 w-8 animate-spin"/>
                <span className="ml-2">Loading families...</span>
            </div>
        );
    }

    // Error state
    if (error) {
        return (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative" role="alert">
                <strong className="font-bold">Error!</strong>
                <span
                    className="block sm:inline"> {error instanceof Error ? error.message : 'Failed to load families'}</span>
            </div>
        );
    }

    const families = queryResponse?.families || [];

    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-3xl font-bold">Families</h1>
                {families.length > 0 && families.filter(x => x.isOwner).length === 0 && <CreateFamilyDialog/>}
            </div>

            {families.length === 0 ? (
                <div className="text-center p-8 bg-gray-50 rounded-lg">
                    <p className="text-gray-500 mb-4">No families found. Create your first family to get started.</p>
                    <div className="flex justify-center">
                        <CreateFamilyDialog/>
                    </div>
                </div>
            ) : (
                <div className="grid gap-6">
                    {families.map((family) => (
                        <Card key={family.id}>
                            <CardHeader>
                                <div className="flex justify-between items-center">
                                    <div>
                                        {editingFamilyId === family.id ? (
                                            <div className="flex flex-col gap-2">
                                                <div className="flex items-center gap-2">
                                                    <Input
                                                        value={editedName}
                                                        onChange={(e) => setEditedName(e.target.value)}
                                                        className="w-full"
                                                        placeholder="Family name"
                                                    />
                                                </div>
                                            </div>
                                        ) : (
                                            <>
                                                <div className="flex items-center gap-2">
                                                    <CardTitle>{family.name}</CardTitle>
                                                </div>
                                                <CardDescription>{family.members.length} members</CardDescription>
                                            </>
                                        )}
                                    </div>
                                    <div className="flex gap-2">
                                        {editingFamilyId === family.id ? (
                                            <>
                                                <Button
                                                    variant="outline"
                                                    size="sm"
                                                    onClick={() => saveChanges(family)}
                                                    disabled={isUpdating}
                                                >
                                                    {isUpdating ? (
                                                        <Loader2 className="h-4 w-4 animate-spin mr-2"/>
                                                    ) : (
                                                        <CheckIcon className="h-4 w-4 mr-2"/>
                                                    )}
                                                    Save
                                                </Button>
                                                <Button
                                                    variant="outline"
                                                    size="sm"
                                                    onClick={cancelEditing}
                                                    disabled={isUpdating}
                                                >
                                                    <XIcon className="h-4 w-4 mr-2"/>
                                                    Cancel
                                                </Button>
                                            </>
                                        ) : (
                                            <>
                                                <Button
                                                    variant="outline"
                                                    size="sm"
                                                    onClick={() => startEditing(family)}
                                                >
                                                    Edit
                                                </Button>
                                                <AddFamilyMemberDialog familyId={family.id}/>
                                            </>
                                        )}
                                    </div>
                                </div>
                            </CardHeader>
                            <CardContent>
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            <TableHead>Name</TableHead>
                                            <TableHead>Role</TableHead>
                                            <TableHead>Email</TableHead>
                                            <TableHead className="w-[100px]">Actions</TableHead>
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {family.members.map((member) => (
                                            <TableRow key={member.id}>
                                                <TableCell className="font-medium">{member.name}</TableCell>
                                                <TableCell>{member.isKid ? 'Kid' : 'Adult'}</TableCell>
                                                <TableCell>{member.email}</TableCell>
                                                <TableCell>
                                                    <Button variant="ghost" size="sm">Edit</Button>
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
};

export default FamiliesPage;
import { useState } from "react";
import { 
  Card, 
  CardContent, 
  CardHeader, 
  CardTitle 
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { PlusIcon, Pencil, X } from "lucide-react";

interface Label {
  id: number;
  name: string;
  color: string;
  count: number;
}

const LabelsPage = () => {
  // Sample data - in a real app, this would come from an API
  const [labels, setLabels] = useState<Label[]>([
    {
      id: 1,
      name: "Subscription",
      color: "#3b82f6", // blue
      count: 12
    },
    {
      id: 2,
      name: "Utilities",
      color: "#10b981", // green
      count: 5
    },
    {
      id: 3,
      name: "Entertainment",
      color: "#8b5cf6", // purple
      count: 8
    },
    {
      id: 4,
      name: "Insurance",
      color: "#f59e0b", // amber
      count: 3
    },
    {
      id: 5,
      name: "Housing",
      color: "#ef4444", // red
      count: 2
    }
  ]);

  const [newLabel, setNewLabel] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editingName, setEditingName] = useState("");

  const handleAddLabel = () => {
    if (newLabel.trim()) {
      const randomColor = `#${Math.floor(Math.random()*16777215).toString(16)}`;
      setLabels([
        ...labels,
        {
          id: labels.length + 1,
          name: newLabel,
          color: randomColor,
          count: 0
        }
      ]);
      setNewLabel("");
    }
  };

  const startEditing = (label: Label) => {
    setEditingId(label.id);
    setEditingName(label.name);
  };

  const saveEdit = () => {
    if (editingId && editingName.trim()) {
      setLabels(labels.map(label => 
        label.id === editingId 
          ? { ...label, name: editingName } 
          : label
      ));
      setEditingId(null);
      setEditingName("");
    }
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditingName("");
  };

  const deleteLabel = (id: number) => {
    setLabels(labels.filter(label => label.id !== id));
  };

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
            <Button onClick={handleAddLabel}>
              <PlusIcon className="h-4 w-4 mr-2" />
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
                    style={{ backgroundColor: label.color }}
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
                      <Button variant="ghost" size="sm" onClick={saveEdit}>
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
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button 
                        variant="ghost" 
                        size="icon"
                        onClick={() => deleteLabel(label.id)}
                        className="text-destructive hover:text-destructive"
                      >
                        <X className="h-4 w-4" />
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
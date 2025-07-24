import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQueryClient } from "@tanstack/react-query";

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { PlusIcon, Loader2 } from "lucide-react";
import { useApiClient } from "@/hooks/use-api-client";
import type { CreateFamilyModel } from "@/api/models";

// Define the form schema
const formSchema = z.object({
  name: z.string().min(1, "Family name is required"),
  haveJointAccount: z.boolean().default(false),
});

type FormValues = z.infer<typeof formSchema>;

interface CreateFamilyDialogProps {
  onSuccess?: () => void;
}

export function CreateFamilyDialog({ onSuccess }: CreateFamilyDialogProps) {
  const { apiClient } = useApiClient();
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Initialize the form
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      haveJointAccount: false,
    },
  });

  // Handle form submission
  const onSubmit = async (values: FormValues) => {
    if (!apiClient) {
      return;
    }

    setIsSubmitting(true);

    try {
      // Create the family model
      const familyModel: CreateFamilyModel = {
        name: values.name,
        haveJointAccount: values.haveJointAccount,
      };

      // Call the API to create the family
      await apiClient.families.post(familyModel);

      // Invalidate the families query to refetch the data
      queryClient.invalidateQueries({ queryKey: ["families"] });

      // Close the dialog and reset the form
      setOpen(false);
      form.reset();

      // Call the onSuccess callback if provided
      if (onSuccess) {
        onSuccess();
      }
    } catch (error) {
      console.error("Failed to create family:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <PlusIcon className="h-4 w-4 mr-2" />
          Add Family
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create a new family</DialogTitle>
          <DialogDescription>
            Create a family to start tracking your recurring payments.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Family Name</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter family name" {...field} />
                  </FormControl>
                  <FormDescription>
                    This is the name of your family.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="haveJointAccount"
              render={({ field }) => (
                <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                  <FormControl>
                    <Checkbox
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                  <div className="space-y-1 leading-none">
                    <FormLabel>Joint Account</FormLabel>
                    <FormDescription>
                      Does your family have a joint account?
                    </FormDescription>
                  </div>
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Creating...
                  </>
                ) : (
                  "Create Family"
                )}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
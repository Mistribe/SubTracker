import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useApiClient } from "@/hooks/use-api-client";
import { UpdateProfileModel } from "@/api/models";

interface ProfileQueryOptions {
  /**
   * Whether the query should be enabled (default: true)
   */
  enabled?: boolean;
}

/**
 * Hook for managing user profile data
 * Provides functions to fetch and update profile information
 * 
 * @param options - Options for the profile query
 * @returns Object containing profile data, loading states, and update functions
 */
export const useProfileManagement = (options: ProfileQueryOptions = {}) => {
  const { enabled = true } = options;
  const { apiClient } = useApiClient();
  const queryClient = useQueryClient();

  // Query to fetch profile data from the backend
  const profileQuery = useQuery({
    queryKey: ['profile'],
    queryFn: async () => {
      if (!apiClient) {
        throw new Error('API client not initialized');
      }

      try {
        const result = await apiClient.users.profile.get();
        return result;
      } catch (error) {
        console.error('Error fetching profile data:', error);
        throw error;
      }
    },
    enabled: !!apiClient && enabled,
    // Don't refetch on window focus to avoid unnecessary API calls
    refetchOnWindowFocus: false,
  });

  // Mutation to update profile data, specifically the preferred currency
  const updateProfileMutation = useMutation({
    mutationFn: async (currency: string) => {
      if (!apiClient) {
        throw new Error('API client not initialized');
      }

      try {
        const payload: UpdateProfileModel = {
          currency,
          additionalData: {}
        };

        return await apiClient.users.profile.put(payload);
      } catch (error) {
        console.error('Error updating profile:', error);
        throw error;
      }
    },
    onSuccess: () => {
      // Invalidate and refetch profile data to ensure UI is up-to-date
      queryClient.invalidateQueries({ queryKey: ['profile'] });
    }
  });

  return {
    profile: profileQuery.data,
    isLoading: profileQuery.isLoading,
    isError: profileQuery.isError,
    error: profileQuery.error,
    updateProfile: updateProfileMutation.mutate,
    isUpdating: updateProfileMutation.isPending
  };
};
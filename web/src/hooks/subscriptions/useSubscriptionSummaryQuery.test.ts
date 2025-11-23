import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useSubscriptionSummaryQuery } from './useSubscriptionSummaryQuery';
import { useApiClient } from '@/hooks/use-api-client';
import type { Amount } from '@/models/amount';
import React from 'react';

// Mock the API client hook
vi.mock('@/hooks/use-api-client');

describe('useSubscriptionSummaryQuery - Summary Model Parsing', () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: {
          retry: false,
        },
      },
    });
    vi.clearAllMocks();
  });

  const wrapper = ({ children }: { children: React.ReactNode }) =>
    React.createElement(QueryClientProvider, { client: queryClient }, children);

  const createMockAmount = (value: number, currency: string = 'USD'): Amount => ({
    value,
    currency,
  });

  describe('Parsing new fields from API response', () => {
    it('should parse activePersonal and activeFamily fields from API response', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 10,
            activePersonal: 6,
            activeFamily: 4,
            totalMonthly: createMockAmount(100),
            totalYearly: createMockAmount(1200),
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.activeSubscriptions).toBe(10);
      expect(result.current.activePersonal).toBe(6);
      expect(result.current.activeFamily).toBe(4);
    });

    it('should parse personal spending fields from API response', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 5,
            activePersonal: 5,
            activeFamily: 0,
            totalMonthly: createMockAmount(50),
            totalYearly: createMockAmount(600),
            personalMonthly: createMockAmount(50),
            personalYearly: createMockAmount(600),
            personalLastMonth: createMockAmount(45),
            personalLastYear: createMockAmount(540),
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.personalMonthly).toEqual(createMockAmount(50));
      expect(result.current.personalYearly).toEqual(createMockAmount(600));
      expect(result.current.personalLastMonth).toEqual(createMockAmount(45));
      expect(result.current.personalLastYear).toEqual(createMockAmount(540));
    });

    it('should parse family spending fields from API response', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 3,
            activePersonal: 0,
            activeFamily: 3,
            totalMonthly: createMockAmount(30),
            totalYearly: createMockAmount(360),
            familyMonthly: createMockAmount(30),
            familyYearly: createMockAmount(360),
            familyLastMonth: createMockAmount(25),
            familyLastYear: createMockAmount(300),
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.familyMonthly).toEqual(createMockAmount(30));
      expect(result.current.familyYearly).toEqual(createMockAmount(360));
      expect(result.current.familyLastMonth).toEqual(createMockAmount(25));
      expect(result.current.familyLastYear).toEqual(createMockAmount(300));
    });

    it('should parse all breakdown fields when both personal and family exist', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 8,
            activePersonal: 5,
            activeFamily: 3,
            totalMonthly: createMockAmount(80),
            totalYearly: createMockAmount(960),
            totalLastMonth: createMockAmount(75),
            totalLastYear: createMockAmount(900),
            personalMonthly: createMockAmount(50),
            personalYearly: createMockAmount(600),
            personalLastMonth: createMockAmount(45),
            personalLastYear: createMockAmount(540),
            familyMonthly: createMockAmount(30),
            familyYearly: createMockAmount(360),
            familyLastMonth: createMockAmount(30),
            familyLastYear: createMockAmount(360),
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      // Verify all fields are parsed correctly
      expect(result.current.activeSubscriptions).toBe(8);
      expect(result.current.activePersonal).toBe(5);
      expect(result.current.activeFamily).toBe(3);
      expect(result.current.totalMonthly).toEqual(createMockAmount(80));
      expect(result.current.totalYearly).toEqual(createMockAmount(960));
      expect(result.current.personalMonthly).toEqual(createMockAmount(50));
      expect(result.current.personalYearly).toEqual(createMockAmount(600));
      expect(result.current.familyMonthly).toEqual(createMockAmount(30));
      expect(result.current.familyYearly).toEqual(createMockAmount(360));
    });
  });

  describe('Fallback to zero values when fields are missing', () => {
    it('should fallback to zero for activePersonal and activeFamily when missing', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 10,
            // activePersonal and activeFamily are missing
            totalMonthly: createMockAmount(100),
            totalYearly: createMockAmount(1200),
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.activeSubscriptions).toBe(10);
      expect(result.current.activePersonal).toBe(0);
      expect(result.current.activeFamily).toBe(0);
    });

    it('should handle missing personal spending fields gracefully', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 5,
            totalMonthly: createMockAmount(50),
            totalYearly: createMockAmount(600),
            // personal fields are missing
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.personalMonthly).toBeUndefined();
      expect(result.current.personalYearly).toBeUndefined();
      expect(result.current.personalLastMonth).toBeUndefined();
      expect(result.current.personalLastYear).toBeUndefined();
    });

    it('should handle missing family spending fields gracefully', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 5,
            totalMonthly: createMockAmount(50),
            totalYearly: createMockAmount(600),
            // family fields are missing
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.familyMonthly).toBeUndefined();
      expect(result.current.familyYearly).toBeUndefined();
      expect(result.current.familyLastMonth).toBeUndefined();
      expect(result.current.familyLastYear).toBeUndefined();
    });

    it('should maintain backward compatibility with old API responses', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            // Old API response without new fields
            active: 10,
            totalMonthly: createMockAmount(100),
            totalYearly: createMockAmount(1200),
            totalLastMonth: createMockAmount(95),
            totalLastYear: createMockAmount(1140),
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      // Old fields should work
      expect(result.current.activeSubscriptions).toBe(10);
      expect(result.current.totalMonthly).toEqual(createMockAmount(100));
      expect(result.current.totalYearly).toEqual(createMockAmount(1200));

      // New fields should have safe defaults
      expect(result.current.activePersonal).toBe(0);
      expect(result.current.activeFamily).toBe(0);
      expect(result.current.personalMonthly).toBeUndefined();
      expect(result.current.familyMonthly).toBeUndefined();
    });

    it('should handle null/undefined response gracefully', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue(null),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      // Should have safe defaults for all fields
      expect(result.current.activeSubscriptions).toBe(0);
      expect(result.current.activePersonal).toBe(0);
      expect(result.current.activeFamily).toBe(0);
    });

    it('should handle empty response object', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({}),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      // Should have safe defaults for all fields
      expect(result.current.activeSubscriptions).toBe(0);
      expect(result.current.activePersonal).toBe(0);
      expect(result.current.activeFamily).toBe(0);
      expect(result.current.topProviders).toEqual([]);
      expect(result.current.topLabels).toEqual([]);
      expect(result.current.upcomingRenewals).toEqual([]);
    });
  });

  describe('Edge cases', () => {
    it('should handle zero values correctly', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 0,
            activePersonal: 0,
            activeFamily: 0,
            totalMonthly: createMockAmount(0),
            totalYearly: createMockAmount(0),
            personalMonthly: createMockAmount(0),
            personalYearly: createMockAmount(0),
            familyMonthly: createMockAmount(0),
            familyYearly: createMockAmount(0),
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.activeSubscriptions).toBe(0);
      expect(result.current.activePersonal).toBe(0);
      expect(result.current.activeFamily).toBe(0);
      expect(result.current.personalMonthly).toEqual(createMockAmount(0));
      expect(result.current.familyMonthly).toEqual(createMockAmount(0));
    });

    it('should handle only personal subscriptions', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 5,
            activePersonal: 5,
            activeFamily: 0,
            totalMonthly: createMockAmount(50),
            totalYearly: createMockAmount(600),
            personalMonthly: createMockAmount(50),
            personalYearly: createMockAmount(600),
            familyMonthly: createMockAmount(0),
            familyYearly: createMockAmount(0),
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.activePersonal).toBe(5);
      expect(result.current.activeFamily).toBe(0);
      expect(result.current.personalMonthly).toEqual(createMockAmount(50));
      expect(result.current.familyMonthly).toEqual(createMockAmount(0));
    });

    it('should handle only family subscriptions', async () => {
      const mockApiClient = {
        subscriptions: {
          subscriptionsSummaryGet: vi.fn().mockResolvedValue({
            active: 3,
            activePersonal: 0,
            activeFamily: 3,
            totalMonthly: createMockAmount(30),
            totalYearly: createMockAmount(360),
            personalMonthly: createMockAmount(0),
            personalYearly: createMockAmount(0),
            familyMonthly: createMockAmount(30),
            familyYearly: createMockAmount(360),
            topProviders: [],
            topLabels: [],
            upcomingRenewals: [],
          }),
        },
      };

      vi.mocked(useApiClient).mockReturnValue({ apiClient: mockApiClient as any });

      const { result } = renderHook(() => useSubscriptionSummaryQuery(), { wrapper });

      await waitFor(() => expect(result.current.isLoading).toBe(false));

      expect(result.current.activePersonal).toBe(0);
      expect(result.current.activeFamily).toBe(3);
      expect(result.current.personalMonthly).toEqual(createMockAmount(0));
      expect(result.current.familyMonthly).toEqual(createMockAmount(30));
    });
  });
});

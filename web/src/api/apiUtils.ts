import { ApiClient } from './apiClient';

/**
 * API instance for use throughout the application
 */
export const api = new ApiClient();

/**
 * Utility function to handle API errors in components
 * @param error The error to handle
 * @param fallbackMessage A fallback message if the error doesn't have one
 * @returns An error message
 */
export const handleApiError = (error: unknown, fallbackMessage = 'An error occurred'): string => {
  if (error instanceof Error) {
    return error.message;
  }
  return fallbackMessage;
};

/**
 * Format a date string for display
 * @param dateString ISO date string
 * @returns Formatted date string
 */
export const formatDate = (dateString: string): string => {
  try {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    }).format(date);
  } catch (error) {
    console.error('Error formatting date:', error);
    return dateString;
  }
};

/**
 * Format a currency value for display
 * @param amount The amount to format
 * @param currency The currency code (e.g., 'USD', 'EUR')
 * @returns Formatted currency string
 */
export const formatCurrency = (amount: number, currency: string): string => {
  try {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(amount);
  } catch (error) {
    console.error('Error formatting currency:', error);
    return `${amount} ${currency}`;
  }
};

/**
 * Calculate the monthly cost of a subscription based on its payments
 * @param price The total price
 * @param months The number of months
 * @returns The monthly cost
 */
export const calculateMonthlyCost = (price: number, months: number): number => {
  if (!months || months <= 0) return price;
  return price / months;
};

/**
 * Check if a subscription is active based on its latest payment
 * @param endDate The end date of the latest payment
 * @returns Whether the subscription is active
 */
export const isSubscriptionActive = (endDate: string): boolean => {
  try {
    const end = new Date(endDate);
    const now = new Date();
    return end >= now;
  } catch (error) {
    console.error('Error checking subscription status:', error);
    return false;
  }
};

/**
 * Get the time remaining until a subscription expires
 * @param endDate The end date of the subscription
 * @returns A string representing the time remaining
 */
export const getTimeRemaining = (endDate: string): string => {
  try {
    const end = new Date(endDate);
    const now = new Date();
    
    if (end < now) {
      return 'Expired';
    }
    
    const diffMs = end.getTime() - now.getTime();
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    
    if (diffDays > 30) {
      const diffMonths = Math.floor(diffDays / 30);
      return `${diffMonths} month${diffMonths !== 1 ? 's' : ''} left`;
    }
    
    return `${diffDays} day${diffDays !== 1 ? 's' : ''} left`;
  } catch (error) {
    console.error('Error calculating time remaining:', error);
    return 'Unknown';
  }
};
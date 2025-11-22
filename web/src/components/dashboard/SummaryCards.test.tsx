import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import * as fc from 'fast-check';
import SummaryCards from './SummaryCards';
import { zeroAmount, type Amount } from '@/models/amount.ts';

// Helper function to compute percentage change
const computeChange = (current: Amount, previous: Amount): number | null => {
  const prevVal = previous?.value;
  const currVal = current?.value;
  if (prevVal == null || currVal == null || prevVal === 0) return null;
  const delta = currVal - prevVal;
  return (delta / prevVal) * 100;
};

// Mock the Money component
vi.mock('@/components/ui/money', () => ({
  default: ({ amount, className }: { amount: Amount; className?: string }) => (
    <span className={className} data-testid="money">
      {amount.currency} {amount.value.toFixed(2)}
    </span>
  ),
}));

describe('SummaryCards', () => {
  const createAmount = (value: number, currency: string = 'USD'): Amount => ({
    value,
    currency,
  });

  describe('Unit Tests - Card Rendering', () => {
    it('should render all three values in Monthly Expenses card', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          personalLastMonth={createAmount(50)}
          personalLastYear={createAmount(600)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          familyLastMonth={createAmount(40)}
          familyLastYear={createAmount(500)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Check Monthly Expenses card
      const monthlyCard = screen.getByTestId('summary-card-monthly');
      expect(monthlyCard).toBeInTheDocument();
      expect(screen.getByTestId('monthly-personal')).toBeInTheDocument();
      expect(screen.getByTestId('monthly-family')).toBeInTheDocument();
      expect(screen.getByTestId('monthly-total')).toBeInTheDocument();
    });

    it('should render all three values in Yearly Expenses card', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          personalLastMonth={createAmount(50)}
          personalLastYear={createAmount(600)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          familyLastMonth={createAmount(40)}
          familyLastYear={createAmount(500)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Check Yearly Expenses card
      const yearlyCard = screen.getByTestId('summary-card-yearly');
      expect(yearlyCard).toBeInTheDocument();
      expect(screen.getByTestId('yearly-personal')).toBeInTheDocument();
      expect(screen.getByTestId('yearly-family')).toBeInTheDocument();
      expect(screen.getByTestId('yearly-total')).toBeInTheDocument();
    });

    it('should render all three values in Active Subscriptions card', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Check Active Subscriptions card
      const activeCard = screen.getByTestId('summary-card-active');
      expect(activeCard).toBeInTheDocument();
      expect(screen.getByTestId('active-personal')).toBeInTheDocument();
      expect(screen.getByTestId('active-family')).toBeInTheDocument();
      expect(screen.getByTestId('active-total')).toBeInTheDocument();
    });

    it('should handle edge case: no personal subscriptions', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(40)}
          totalYearly={createAmount(480)}
          totalLastMonth={createAmount(40)}
          totalLastYear={createAmount(500)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          familyLastMonth={createAmount(40)}
          familyLastYear={createAmount(500)}
          activeSubscriptionsCount={4}
          activePersonal={0}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Personal rows should not be rendered when no personal data
      expect(screen.queryByTestId('monthly-personal')).not.toBeInTheDocument();
      expect(screen.queryByTestId('yearly-personal')).not.toBeInTheDocument();
      
      // Family and total should still be rendered
      expect(screen.getByTestId('monthly-family')).toBeInTheDocument();
      expect(screen.getByTestId('monthly-total')).toBeInTheDocument();
      expect(screen.getByTestId('active-family')).toBeInTheDocument();
      expect(screen.getByTestId('active-total')).toBeInTheDocument();
    });

    it('should handle edge case: no family subscriptions', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(60)}
          totalYearly={createAmount(720)}
          totalLastMonth={createAmount(50)}
          totalLastYear={createAmount(600)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          personalLastMonth={createAmount(50)}
          personalLastYear={createAmount(600)}
          activeSubscriptionsCount={6}
          activePersonal={6}
          activeFamily={0}
          isLoading={false}
        />
      );

      // Family rows should not be rendered when no family data
      expect(screen.queryByTestId('monthly-family')).not.toBeInTheDocument();
      expect(screen.queryByTestId('yearly-family')).not.toBeInTheDocument();
      
      // Personal and total should still be rendered
      expect(screen.getByTestId('monthly-personal')).toBeInTheDocument();
      expect(screen.getByTestId('monthly-total')).toBeInTheDocument();
      expect(screen.getByTestId('active-personal')).toBeInTheDocument();
      expect(screen.getByTestId('active-total')).toBeInTheDocument();
    });

    it('should handle edge case: no subscriptions at all', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(0)}
          totalYearly={createAmount(0)}
          totalLastMonth={createAmount(0)}
          totalLastYear={createAmount(0)}
          activeSubscriptionsCount={0}
          activePersonal={0}
          activeFamily={0}
          isLoading={false}
        />
      );

      // Only total rows should be rendered with zero values
      expect(screen.getByTestId('monthly-total')).toBeInTheDocument();
      expect(screen.getByTestId('yearly-total')).toBeInTheDocument();
      expect(screen.getByTestId('active-total')).toBeInTheDocument();
      
      // Personal and family should not be rendered
      expect(screen.queryByTestId('monthly-personal')).not.toBeInTheDocument();
      expect(screen.queryByTestId('monthly-family')).not.toBeInTheDocument();
    });

    it('should display skeleton loaders during loading', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={true}
        />
      );

      // All three cards should show skeleton loaders
      const monthlyCard = screen.getByTestId('summary-card-monthly');
      const yearlyCard = screen.getByTestId('summary-card-yearly');
      const activeCard = screen.getByTestId('summary-card-active');

      // Each card should have skeleton elements
      expect(monthlyCard.querySelector('.animate-pulse')).toBeInTheDocument();
      expect(yearlyCard.querySelector('.animate-pulse')).toBeInTheDocument();
      expect(activeCard.querySelector('.animate-pulse')).toBeInTheDocument();

      // Breakdown rows should not be visible during loading
      expect(screen.queryByTestId('monthly-personal')).not.toBeInTheDocument();
      expect(screen.queryByTestId('monthly-family')).not.toBeInTheDocument();
      expect(screen.queryByTestId('monthly-total')).not.toBeInTheDocument();
    });
  });

  describe('Unit Tests - Trend Indicator Logic', () => {
    it('should calculate percentage change with specific values', () => {
      // Test increase: from 100 to 120 = +20%
      const increase = computeChange(createAmount(120), createAmount(100));
      expect(increase).toBe(20);

      // Test decrease: from 100 to 80 = -20%
      const decrease = computeChange(createAmount(80), createAmount(100));
      expect(decrease).toBe(-20);

      // Test no change: from 100 to 100 = 0%
      const noChange = computeChange(createAmount(100), createAmount(100));
      expect(noChange).toBe(0);

      // Test large increase: from 50 to 150 = +200%
      const largeIncrease = computeChange(createAmount(150), createAmount(50));
      expect(largeIncrease).toBe(200);

      // Test small change: from 100 to 101 = +1%
      const smallIncrease = computeChange(createAmount(101), createAmount(100));
      expect(smallIncrease).toBe(1);
    });

    it('should return null when historical data is missing', () => {
      // Previous amount is zero
      const zeroResult = computeChange(createAmount(100), createAmount(0));
      expect(zeroResult).toBeNull();

      // Previous amount is null
      const nullResult = computeChange(createAmount(100), { value: null as any, currency: 'USD' });
      expect(nullResult).toBeNull();

      // Current amount is null
      const currentNullResult = computeChange({ value: null as any, currency: 'USD' }, createAmount(100));
      expect(currentNullResult).toBeNull();

      // Both are null
      const bothNullResult = computeChange(
        { value: null as any, currency: 'USD' },
        { value: null as any, currency: 'USD' }
      );
      expect(bothNullResult).toBeNull();
    });

    it('should display neutral indicator when historical data is unavailable', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(0)} // Zero previous value
          totalLastYear={createAmount(0)} // Zero previous value
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // The component should render without errors
      expect(screen.getByTestId('monthly-total')).toBeInTheDocument();
      expect(screen.getByTestId('yearly-total')).toBeInTheDocument();

      // Neutral indicators (—) should be present when change is null
      const monthlyTotal = screen.getByTestId('monthly-total');
      expect(monthlyTotal.textContent).toContain('—');
    });

    it('should apply correct color coding for positive changes (increases)', () => {
      const { container } = render(
        <SummaryCards
          totalMonthly={createAmount(120)}
          totalYearly={createAmount(1440)}
          totalLastMonth={createAmount(100)}
          totalLastYear={createAmount(1200)}
          personalMonthly={createAmount(72)}
          personalYearly={createAmount(864)}
          personalLastMonth={createAmount(60)}
          personalLastYear={createAmount(720)}
          familyMonthly={createAmount(48)}
          familyYearly={createAmount(576)}
          familyLastMonth={createAmount(40)}
          familyLastYear={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Positive changes (increases in spending) should have red styling
      const redBadges = container.querySelectorAll('.bg-red-100, .dark\\:bg-red-500\\/10');
      expect(redBadges.length).toBeGreaterThan(0);

      // Check that percentage is displayed
      const monthlyTotal = screen.getByTestId('monthly-total');
      expect(monthlyTotal.textContent).toContain('20.0%');
    });

    it('should apply correct color coding for negative changes (decreases)', () => {
      const { container } = render(
        <SummaryCards
          totalMonthly={createAmount(80)}
          totalYearly={createAmount(960)}
          totalLastMonth={createAmount(100)}
          totalLastYear={createAmount(1200)}
          personalMonthly={createAmount(48)}
          personalYearly={createAmount(576)}
          personalLastMonth={createAmount(60)}
          personalLastYear={createAmount(720)}
          familyMonthly={createAmount(32)}
          familyYearly={createAmount(384)}
          familyLastMonth={createAmount(40)}
          familyLastYear={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Negative changes (decreases in spending) should have green styling
      const greenBadges = container.querySelectorAll('.bg-green-100, .dark\\:bg-green-500\\/10');
      expect(greenBadges.length).toBeGreaterThan(0);

      // Check that percentage is displayed
      const monthlyTotal = screen.getByTestId('monthly-total');
      expect(monthlyTotal.textContent).toContain('20.0%');
    });

    it('should apply neutral color coding for zero change', () => {
      const { container } = render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(100)}
          totalLastYear={createAmount(1200)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          personalLastMonth={createAmount(60)}
          personalLastYear={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          familyLastMonth={createAmount(40)}
          familyLastYear={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Zero change should have muted styling
      const mutedBadges = container.querySelectorAll('.bg-muted');
      expect(mutedBadges.length).toBeGreaterThan(0);

      // Check that 0.0% is displayed
      const monthlyTotal = screen.getByTestId('monthly-total');
      expect(monthlyTotal.textContent).toContain('0.0%');
    });

    it('should calculate trends for personal monthly expenses', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(70)}
          personalYearly={createAmount(840)}
          personalLastMonth={createAmount(60)}
          personalLastYear={createAmount(720)}
          familyMonthly={createAmount(30)}
          familyYearly={createAmount(360)}
          familyLastMonth={createAmount(30)}
          familyLastYear={createAmount(380)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      const personalMonthly = screen.getByTestId('monthly-personal');
      // Personal increased from 60 to 70 = +16.67%
      expect(personalMonthly.textContent).toContain('16.7%');
    });

    it('should calculate trends for family monthly expenses', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          personalLastMonth={createAmount(60)}
          personalLastYear={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          familyLastMonth={createAmount(30)}
          familyLastYear={createAmount(380)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      const familyMonthly = screen.getByTestId('monthly-family');
      // Family increased from 30 to 40 = +33.33%
      expect(familyMonthly.textContent).toContain('33.3%');
    });

    it('should calculate trends for personal yearly expenses', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(840)}
          personalLastMonth={createAmount(60)}
          personalLastYear={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(360)}
          familyLastMonth={createAmount(30)}
          familyLastYear={createAmount(380)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      const personalYearly = screen.getByTestId('yearly-personal');
      // Personal increased from 720 to 840 = +16.67%
      expect(personalYearly.textContent).toContain('16.7%');
    });

    it('should calculate trends for family yearly expenses', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          personalLastMonth={createAmount(60)}
          personalLastYear={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          familyLastMonth={createAmount(40)}
          familyLastYear={createAmount(400)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      const familyYearly = screen.getByTestId('yearly-family');
      // Family increased from 400 to 480 = +20%
      expect(familyYearly.textContent).toContain('20.0%');
    });

    it('should calculate trends for total monthly expenses', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(110)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(100)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      const totalMonthly = screen.getByTestId('monthly-total');
      // Total increased from 100 to 110 = +10%
      expect(totalMonthly.textContent).toContain('10.0%');
    });

    it('should calculate trends for total yearly expenses', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1320)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1200)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      const totalYearly = screen.getByTestId('yearly-total');
      // Total increased from 1200 to 1320 = +10%
      expect(totalYearly.textContent).toContain('10.0%');
    });
  });

  describe('Responsive Design Tests', () => {
    it('should use responsive grid layout', () => {
      const { container } = render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      const gridContainer = screen.getByTestId('summary-cards');
      
      // Verify responsive grid classes are applied
      expect(gridContainer).toHaveClass('grid');
      expect(gridContainer).toHaveClass('grid-cols-1'); // Mobile: single column
      expect(gridContainer).toHaveClass('md:grid-cols-3'); // Desktop: three columns
      expect(gridContainer).toHaveClass('gap-6'); // Consistent spacing
    });

    it('should render all three cards in the grid', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // All three cards should be present
      expect(screen.getByTestId('summary-card-monthly')).toBeInTheDocument();
      expect(screen.getByTestId('summary-card-yearly')).toBeInTheDocument();
      expect(screen.getByTestId('summary-card-active')).toBeInTheDocument();
    });

    it('should stack breakdown values vertically within each card', () => {
      const { container } = render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          personalLastMonth={createAmount(50)}
          personalLastYear={createAmount(600)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          familyLastMonth={createAmount(40)}
          familyLastYear={createAmount(500)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Check that breakdown rows use flex layout for proper stacking
      const monthlyPersonal = screen.getByTestId('monthly-personal');
      expect(monthlyPersonal).toHaveClass('flex');
      expect(monthlyPersonal).toHaveClass('items-center');
      expect(monthlyPersonal).toHaveClass('justify-between');
      
      // Verify vertical spacing between rows
      expect(monthlyPersonal).toHaveClass('py-2');
      expect(monthlyPersonal).toHaveClass('border-b');
    });

    it('should maintain readability with proper spacing and borders', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          personalLastMonth={createAmount(50)}
          personalLastYear={createAmount(600)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          familyLastMonth={createAmount(40)}
          familyLastYear={createAmount(500)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Check Monthly card breakdown rows
      const monthlyPersonal = screen.getByTestId('monthly-personal');
      const monthlyFamily = screen.getByTestId('monthly-family');
      const monthlyTotal = screen.getByTestId('monthly-total');

      // All rows should have borders except the last one
      expect(monthlyPersonal).toHaveClass('border-b');
      expect(monthlyFamily).toHaveClass('border-b');
      expect(monthlyTotal).toHaveClass('last:border-b-0');

      // All rows should have consistent padding
      expect(monthlyPersonal).toHaveClass('py-2');
      expect(monthlyFamily).toHaveClass('py-2');
      expect(monthlyTotal).toHaveClass('py-2');
    });

    it('should use appropriate text sizes for mobile readability', () => {
      render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      const monthlyPersonal = screen.getByTestId('monthly-personal');
      
      // Labels should use readable text size
      const label = monthlyPersonal.querySelector('.text-sm');
      expect(label).toBeInTheDocument();
      
      // Amounts should be prominent
      const amount = monthlyPersonal.querySelector('.text-lg');
      expect(amount).toBeInTheDocument();
    });

    it('should handle long currency values without breaking layout', () => {
      // Test with large numbers that might wrap
      render(
        <SummaryCards
          totalMonthly={createAmount(999999.99)}
          totalYearly={createAmount(11999999.88)}
          totalLastMonth={createAmount(888888.88)}
          totalLastYear={createAmount(10666666.66)}
          personalMonthly={createAmount(599999.99)}
          personalYearly={createAmount(7199999.88)}
          familyMonthly={createAmount(400000.00)}
          familyYearly={createAmount(4800000.00)}
          activeSubscriptionsCount={999}
          activePersonal={599}
          activeFamily={400}
          isLoading={false}
        />
      );

      // All cards should still be rendered
      expect(screen.getByTestId('summary-card-monthly')).toBeInTheDocument();
      expect(screen.getByTestId('summary-card-yearly')).toBeInTheDocument();
      expect(screen.getByTestId('summary-card-active')).toBeInTheDocument();

      // Breakdown rows should maintain flex layout
      const monthlyTotal = screen.getByTestId('monthly-total');
      expect(monthlyTotal).toHaveClass('flex');
      expect(monthlyTotal).toHaveClass('justify-between');
    });

    it('should maintain card structure with minimal content', () => {
      // Test with only total values (no personal/family breakdown)
      render(
        <SummaryCards
          totalMonthly={createAmount(10)}
          totalYearly={createAmount(120)}
          totalLastMonth={createAmount(10)}
          totalLastYear={createAmount(120)}
          activeSubscriptionsCount={1}
          activePersonal={0}
          activeFamily={0}
          isLoading={false}
        />
      );

      // Cards should still render properly with minimal content
      expect(screen.getByTestId('summary-card-monthly')).toBeInTheDocument();
      expect(screen.getByTestId('summary-card-yearly')).toBeInTheDocument();
      expect(screen.getByTestId('summary-card-active')).toBeInTheDocument();

      // Only total rows should be visible
      expect(screen.getByTestId('monthly-total')).toBeInTheDocument();
      expect(screen.getByTestId('yearly-total')).toBeInTheDocument();
      expect(screen.getByTestId('active-total')).toBeInTheDocument();
    });

    it('should apply consistent spacing in card content areas', () => {
      const { container } = render(
        <SummaryCards
          totalMonthly={createAmount(100)}
          totalYearly={createAmount(1200)}
          totalLastMonth={createAmount(90)}
          totalLastYear={createAmount(1100)}
          personalMonthly={createAmount(60)}
          personalYearly={createAmount(720)}
          familyMonthly={createAmount(40)}
          familyYearly={createAmount(480)}
          activeSubscriptionsCount={10}
          activePersonal={6}
          activeFamily={4}
          isLoading={false}
        />
      );

      // Check that card content uses space-y for vertical spacing
      const monthlyCard = screen.getByTestId('summary-card-monthly');
      const contentArea = monthlyCard.querySelector('.space-y-1');
      expect(contentArea).toBeInTheDocument();
    });
  });

  describe('Property-Based Tests', () => {
    /**
     * Feature: dashboard-spending-breakdown, Property 5: Percentage change calculation
     * Validates: Requirements 3.1, 3.2, 3.3, 3.4, 3.5, 3.6
     * 
     * For any current amount and previous amount (where previous is non-zero),
     * the percentage change should equal ((current - previous) / previous) * 100
     */
    it('should correctly calculate percentage change for any valid amounts', () => {
      fc.assert(
        fc.property(
          // Generate current value (can be any number including negative)
          fc.double({ min: -1000000, max: 1000000, noNaN: true }),
          // Generate previous value (non-zero to avoid division by zero)
          // Filter out very small numbers that cause floating point precision issues
          fc.double({ min: -1000000, max: 1000000, noNaN: true })
            .filter(v => Math.abs(v) > 1e-10),
          fc.constantFrom('USD', 'EUR', 'GBP', 'JPY'),
          (currentValue, previousValue, currency) => {
            const current: Amount = { value: currentValue, currency };
            const previous: Amount = { value: previousValue, currency };

            const result = computeChange(current, previous);

            // The result should not be null when previous is non-zero
            expect(result).not.toBeNull();

            if (result !== null) {
              // Calculate expected percentage change
              const expected = ((currentValue - previousValue) / previousValue) * 100;

              // Both should be NaN or both should be valid numbers
              if (isNaN(expected)) {
                expect(isNaN(result)).toBe(true);
              } else {
                // Allow for small floating point errors
                expect(Math.abs(result - expected)).toBeLessThan(0.0001);
              }
            }
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should return null when previous amount is zero', () => {
      fc.assert(
        fc.property(
          fc.double({ min: -1000000, max: 1000000, noNaN: true }),
          fc.constantFrom('USD', 'EUR', 'GBP', 'JPY'),
          (currentValue, currency) => {
            const current: Amount = { value: currentValue, currency };
            const previous: Amount = { value: 0, currency };

            const result = computeChange(current, previous);

            // Should return null when previous is zero (avoid division by zero)
            expect(result).toBeNull();
          }
        ),
        { numRuns: 100 }
      );
    });

    it('should return null when previous amount is null or undefined', () => {
      fc.assert(
        fc.property(
          fc.double({ min: -1000000, max: 1000000, noNaN: true }),
          fc.constantFrom('USD', 'EUR', 'GBP', 'JPY'),
          (currentValue, currency) => {
            const current: Amount = { value: currentValue, currency };
            const previous: Amount = { value: null as any, currency };

            const result = computeChange(current, previous);

            // Should return null when previous value is null
            expect(result).toBeNull();
          }
        ),
        { numRuns: 100 }
      );
    });
  });
});

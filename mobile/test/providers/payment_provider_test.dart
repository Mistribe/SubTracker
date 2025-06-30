import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/models/payment.dart';
import 'package:mobile/providers/payment_provider.dart';

void main() {
  group('PaymentProvider', () {
    late PaymentProvider provider;

    setUp(() {
      provider = PaymentProvider();
    });

    test('activePaymentsCount returns correct count', () {
      // Initially there should be 0 payments
      expect(provider.activePaymentsCount, 0);

      // Add a payment
      provider.addPayment('Netflix', 15.99, false);
      expect(provider.activePaymentsCount, 1);

      // Add another payment
      provider.addPayment('Spotify', 9.99, false);
      expect(provider.activePaymentsCount, 2);

      // Remove a payment
      final payments = provider.payments;
      provider.removePayment(payments[0].id);
      expect(provider.activePaymentsCount, 1);
    });

    test('totalAmountSpent calculates correctly', () {
      // Create payment dates
      final threeMonthsAgo = DateTime.now().subtract(const Duration(days: 90));
      final twoYearsAgo = DateTime.now().subtract(const Duration(days: 730));

      // Initially total spent should be 0
      expect(provider.totalAmountSpent, 0);

      // Add a monthly payment that started 3 months ago
      provider.addPayment('Netflix', 15.99, false, paymentDate: threeMonthsAgo);
      // Should have paid 3 times (once per month)
      expect(provider.totalAmountSpent, closeTo(15.99 * 3, 0.01));

      // Add an annual payment that started 2 years ago
      provider.addPayment('Amazon Prime', 119.99, true, paymentDate: twoYearsAgo);
      // Should have paid 2 times (once per year)
      // Total should be (15.99 * 3) + (119.99 * 2)
      expect(provider.totalAmountSpent, closeTo((15.99 * 3) + (119.99 * 2), 0.01));

      // Add a new payment (should not contribute to total spent yet)
      provider.addPayment('New Service', 9.99, false);
      // Total should still be (15.99 * 3) + (119.99 * 2)
      expect(provider.totalAmountSpent, closeTo((15.99 * 3) + (119.99 * 2), 0.01));
    });
  });
}
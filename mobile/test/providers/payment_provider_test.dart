import 'package:flutter_test/flutter_test.dart';
import 'package:hive/hive.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:mobile/models/payment.dart';
import 'package:mobile/providers/payment_provider.dart';
import 'package:mobile/repositories/payment_repository.dart';
import 'dart:io';

// Mock repository for testing
class MockPaymentRepository implements PaymentRepository {
  final List<Payment> _payments = [];

  @override
  Future<void> initialize() async {
    // No initialization needed for mock
  }

  @override
  List<Payment> getAllPayments() {
    return List.from(_payments); // Return a copy to avoid unintended modifications
  }

  @override
  Future<void> addPayment(Payment payment) async {
    _payments.add(payment);
  }

  @override
  Future<void> updatePayment(Payment payment) async {
    final index = _payments.indexWhere((p) => p.id == payment.id);
    if (index >= 0) {
      _payments[index] = payment;
    }
  }

  @override
  Future<void> deletePayment(String id) async {
    _payments.removeWhere((payment) => payment.id == id);
  }

  @override
  Future<void> clearAll() async {
    _payments.clear();
  }
}

void main() {
  setUpAll(() async {
    // Set up a temporary directory for Hive
    final tempDir = await Directory.systemTemp.createTemp('hive_test');
    Hive.init(tempDir.path);

    // Register adapters
    if (!Hive.isAdapterRegistered(0)) {
      Hive.registerAdapter(PaymentAdapter());
    }
    if (!Hive.isAdapterRegistered(1)) {
      Hive.registerAdapter(PriceChangeAdapter());
    }
  });

  tearDownAll(() async {
    // Clean up after all tests
    await Hive.close();
  });

  group('PaymentProvider', () {
    late PaymentProvider provider;
    late MockPaymentRepository mockRepository;

    setUp(() async {
      mockRepository = MockPaymentRepository();
      provider = PaymentProvider(paymentRepository: mockRepository);

      // Wait for the provider to initialize
      await Future.delayed(Duration.zero);
    });

    test('activePaymentsCount returns correct count', () async {
      // Initially there should be 0 payments
      expect(provider.activePaymentsCount, 0);

      // Add a payment
      await provider.addPayment('Netflix', 15.99, false);
      expect(provider.activePaymentsCount, 1);

      // Add another payment
      await provider.addPayment('Spotify', 9.99, false);
      expect(provider.activePaymentsCount, 2);

      // Remove a payment
      final payments = provider.payments;
      await provider.removePayment(payments[0].id);

      // Verify that the payment was removed
      expect(provider.activePaymentsCount, 1);

      // Verify that the correct payment remains
      expect(provider.payments.length, 1);
      expect(provider.payments[0].name, 'Spotify');
    });

    test('stopPayment and reactivatePayment work correctly', () async {
      // Add a payment
      await provider.addPayment('Netflix', 15.99, false);
      final payments = provider.payments;
      final paymentId = payments[0].id;

      // Initially the payment should be active
      expect(payments[0].isStopped, false);
      expect(provider.activePaymentsCount, 1);

      // Stop the payment with default stop date (last payment date)
      await provider.stopPayment(paymentId);

      // Verify that the payment is stopped
      final stoppedPayments = provider.payments;
      expect(stoppedPayments[0].isStopped, true);
      expect(stoppedPayments[0].reactivationDate, null);
      expect(stoppedPayments[0].stopDate, isNotNull); // Should have a stop date

      // Verify that the payment is not counted in active payments
      expect(provider.activePaymentsCount, 0);

      // Reactivate the payment with a future date
      final reactivationDate = DateTime.now().add(const Duration(days: 30));
      await provider.reactivatePayment(paymentId, reactivationDate);

      // Verify that the payment is still stopped but has a reactivation date
      final reactivatedPayments = provider.payments;
      expect(reactivatedPayments[0].isStopped, true);
      expect(reactivatedPayments[0].reactivationDate, reactivationDate);
      expect(reactivatedPayments[0].stopDate, isNotNull); // Should still have a stop date

      // Verify that the payment is still not counted in active payments
      expect(provider.activePaymentsCount, 0);
    });

    test('stopPayment with custom stop date works correctly', () async {
      // Add a payment
      await provider.addPayment('Netflix', 15.99, false);
      final payments = provider.payments;
      final paymentId = payments[0].id;

      // Create a custom stop date
      final customStopDate = DateTime.now().subtract(const Duration(days: 15));

      // Stop the payment with custom stop date
      await provider.stopPayment(paymentId, stopDate: customStopDate);

      // Verify that the payment is stopped with the custom stop date
      final stoppedPayments = provider.payments;
      expect(stoppedPayments[0].isStopped, true);
      expect(stoppedPayments[0].stopDate, customStopDate);
    });

    test('totalAmountSpent calculates correctly', () async {
      // Create payment dates
      final threeMonthsAgo = DateTime.now().subtract(const Duration(days: 90));
      final twoYearsAgo = DateTime.now().subtract(const Duration(days: 730));

      // Initially total spent should be 0
      expect(provider.totalAmountSpent, 0);

      // Add a monthly payment that started 3 months ago
      await provider.addPayment('Netflix', 15.99, false, paymentDate: threeMonthsAgo);
      // Should have paid 3 times (once per month)
      expect(provider.totalAmountSpent, closeTo(15.99 * 3, 0.01));

      // Get the current total
      final firstTotal = provider.totalAmountSpent;

      // Add an annual payment that started 2 years ago
      await provider.addPayment('Amazon Prime', 119.99, true, paymentDate: twoYearsAgo);
      // Should have paid 2 times (once per year)
      // Total should be firstTotal + (119.99 * 2)
      final expectedTotal = firstTotal + (119.99 * 2);
      expect(provider.totalAmountSpent, closeTo(expectedTotal, 0.01));

      // Add a new payment (should not contribute to total spent yet)
      await provider.addPayment('New Service', 9.99, false);
      // Total should still be the same as before
      expect(provider.totalAmountSpent, closeTo(expectedTotal, 0.01));
    });

    test('totalAmountSpent excludes stopped payments', () async {
      // Create payment dates
      final threeMonthsAgo = DateTime.now().subtract(const Duration(days: 90));

      // Add a monthly payment that started 3 months ago
      await provider.addPayment('Netflix', 15.99, false, paymentDate: threeMonthsAgo);
      // Should have paid 3 times (once per month)
      final initialTotal = 15.99 * 3;
      expect(provider.totalAmountSpent, closeTo(initialTotal, 0.01));

      // Stop the payment
      final payments = provider.payments;
      final paymentId = payments[0].id;
      await provider.stopPayment(paymentId);

      // Total amount spent should be 0 since the payment is stopped
      expect(provider.totalAmountSpent, 0);

      // Reactivate the payment with a future date
      final reactivationDate = DateTime.now().add(const Duration(days: 30));
      await provider.reactivatePayment(paymentId, reactivationDate);

      // Total amount spent should still be 0 since the payment is still stopped
      expect(provider.totalAmountSpent, 0);

      // Add another active payment
      await provider.addPayment('Spotify', 9.99, false, paymentDate: threeMonthsAgo);

      // Total amount spent should only include the active payment
      expect(provider.totalAmountSpent, closeTo(9.99 * 3, 0.01));
    });
  });
}

import 'package:flutter_test/flutter_test.dart';
import 'package:hive/hive.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:mobile/models/payment.dart';
import 'dart:io';
import 'package:path/path.dart' as path;

void main() {
  setUp(() async {
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

  tearDown(() async {
    // Clean up after each test
    await Hive.close();
  });
  group('Payment', () {
    test('totalAmountSpent for monthly payment', () {
      // Create a payment that started 3 months ago
      final threeMonthsAgo = DateTime.now().subtract(const Duration(days: 90));
      final payment = Payment(
        id: '1',
        name: 'Netflix',
        price: 15.99,
        isAnnual: false,
        paymentDate: threeMonthsAgo,
      );

      // Should have paid 3 times (once per month)
      expect(payment.totalAmountSpent, closeTo(15.99 * 3, 0.01));
    });

    test('totalAmountSpent for annual payment', () {
      // Create a payment that started 2 years ago
      final twoYearsAgo = DateTime.now().subtract(const Duration(days: 730));
      final payment = Payment(
        id: '2',
        name: 'Amazon Prime',
        price: 119.99,
        isAnnual: true,
        paymentDate: twoYearsAgo,
      );

      // Should have paid 2 times (once per year)
      expect(payment.totalAmountSpent, closeTo(119.99 * 2, 0.01));
    });

    test('totalAmountSpent for new payment', () {
      // Create a payment that started today
      final payment = Payment(
        id: '3',
        name: 'New Service',
        price: 9.99,
        isAnnual: false,
        paymentDate: DateTime.now(),
      );

      // Should have paid 0 times (just started)
      expect(payment.totalAmountSpent, 0);
    });

    test('formattedTotalAmountSpent', () {
      // Create a payment that started 2 months ago
      final twoMonthsAgo = DateTime.now().subtract(const Duration(days: 60));
      final payment = Payment(
        id: '4',
        name: 'Spotify',
        price: 9.99,
        isAnnual: false,
        paymentDate: twoMonthsAgo,
      );

      // Should have paid 2 times and format correctly
      expect(payment.formattedTotalAmountSpent, '\$${(9.99 * 2).toStringAsFixed(2)}');
    });
  });
}

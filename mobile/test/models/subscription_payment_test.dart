import 'dart:io';
import 'package:flutter_test/flutter_test.dart';
import 'package:hive/hive.dart';
import 'package:subscription_tracker/models/subscription_payment.dart';

void main() {
  group('SubscriptionPayment', () {
    setUpAll(() async {
      // Set up a temporary directory for Hive
      final tempDir = await Directory.systemTemp.createTemp('hive_test');
      Hive.init(tempDir.path);

      // Register adapters
      if (!Hive.isAdapterRegistered(1)) {
        Hive.registerAdapter(SubscriptionPaymentAdapter());
      }
    });

    tearDownAll(() async {
      // Clean up after all tests
      await Hive.close();
    });

    // Helper function to create a payment with default values
    SubscriptionPayment createPayment({
      String id = '1',
      double price = 9.99,
      DateTime? startDate,
      DateTime? endDate,
      int months = 1,
      String currency = 'USD',
    }) {
      return SubscriptionPayment(
        id: id,
        price: price,
        startDate: startDate ?? DateTime(2023, 1, 1),
        endDate: endDate,
        months: months,
        currency: currency,
      );
    }

    test('constructor creates payment with correct values', () {
      final startDate = DateTime(2023, 1, 1);
      final endDate = DateTime(2023, 12, 31);
      
      final payment = SubscriptionPayment(
        id: '1',
        price: 9.99,
        startDate: startDate,
        endDate: endDate,
        months: 1,
        currency: 'USD',
      );

      expect(payment.id, equals('1'));
      expect(payment.price, equals(9.99));
      expect(payment.startDate, equals(startDate));
      expect(payment.endDate, equals(endDate));
      expect(payment.months, equals(1));
      expect(payment.currency, equals('USD'));
    });

    group('isStarted', () {
      test('returns true if startDate is in the past', () {
        final pastDate = DateTime.now().subtract(const Duration(days: 10));
        final payment = createPayment(startDate: pastDate);
        
        expect(payment.isStarted, isTrue);
      });

      test('returns false if startDate is in the future', () {
        final futureDate = DateTime.now().add(const Duration(days: 10));
        final payment = createPayment(startDate: futureDate);
        
        expect(payment.isStarted, isFalse);
      });
    });

    group('isActive', () {
      test('returns true if endDate is null', () {
        final payment = createPayment(endDate: null);
        
        expect(payment.isActive, isTrue);
      });

      test('returns true if endDate is in the future', () {
        final futureDate = DateTime.now().add(const Duration(days: 10));
        final payment = createPayment(endDate: futureDate);
        
        expect(payment.isActive, isTrue);
      });

      test('returns false if endDate is in the past', () {
        final pastDate = DateTime.now().subtract(const Duration(days: 10));
        final payment = createPayment(endDate: pastDate);
        
        expect(payment.isActive, isFalse);
      });
    });

    group('totalMonthElapsed', () {
      test('calculates months between startDate and endDate', () {
        final startDate = DateTime(2023, 1, 1);
        final endDate = DateTime(2023, 6, 1); // 5 months difference
        final payment = createPayment(startDate: startDate, endDate: endDate);
        
        expect(payment.totalMonthElapsed, equals(5));
      });

      test('calculates months between startDate and now if endDate is null', () {
        // This test is time-dependent, so we need to be careful
        final now = DateTime.now();
        final startDate = DateTime(now.year, now.month - 3, now.day); // 3 months ago
        final payment = createPayment(startDate: startDate, endDate: null);
        
        expect(payment.totalMonthElapsed, equals(3));
      });

      test('handles year boundaries correctly', () {
        final startDate = DateTime(2022, 11, 1);
        final endDate = DateTime(2023, 2, 1); // 3 months difference across year boundary
        final payment = createPayment(startDate: startDate, endDate: endDate);
        
        expect(payment.totalMonthElapsed, equals(3));
      });
    });

    test('monthlyCost calculates correctly', () {
      final payment = createPayment(price: 12.0, months: 3);
      
      expect(payment.monthlyCost, equals(4.0)); // 12 / 3 = 4
    });

    test('totalAmountSpent calculates correctly', () {
      final startDate = DateTime(2023, 1, 1);
      final endDate = DateTime(2023, 4, 1); // 3 months
      final payment = createPayment(
        startDate: startDate,
        endDate: endDate,
        price: 30.0,
        months: 3, // Monthly cost is 10.0
      );
      
      expect(payment.totalAmountSpent, equals(30.0)); // 3 months * 10.0 = 30.0
    });

    group('nextPaymentDate', () {
      test('returns startDate if it is in the future', () {
        final futureDate = DateTime.now().add(const Duration(days: 10));
        final payment = createPayment(startDate: futureDate);
        
        expect(payment.nextPaymentDate, equals(futureDate));
      });

      test('returns far future date if payment is not active', () {
        final pastDate = DateTime.now().subtract(const Duration(days: 10));
        final payment = createPayment(
          startDate: DateTime(2023, 1, 1),
          endDate: pastDate,
        );
        
        expect(payment.nextPaymentDate, equals(DateTime(9999, 12, 31)));
      });

      test('calculates next payment date for monthly subscription', () {
        // Use a fixed date for testing to avoid time-dependent issues
        final now = DateTime(2023, 5, 15);
        final startDate = DateTime(2023, 1, 1);
        final payment = createPayment(startDate: startDate, months: 1);
        
        // Mock DateTime.now() for this test
        final originalNow = DateTime.now;
        DateTime.now = () => now;
        
        try {
          // Next payment should be on June 1, 2023
          expect(payment.nextPaymentDate, equals(DateTime(2023, 6, 1)));
        } finally {
          // Restore original DateTime.now
          DateTime.now = originalNow;
        }
      });

      test('calculates next payment date for quarterly subscription', () {
        // Use a fixed date for testing
        final now = DateTime(2023, 5, 15);
        final startDate = DateTime(2023, 1, 1);
        final payment = createPayment(startDate: startDate, months: 3);
        
        // Mock DateTime.now() for this test
        final originalNow = DateTime.now;
        DateTime.now = () => now;
        
        try {
          // Next payment should be on August 1, 2023
          expect(payment.nextPaymentDate, equals(DateTime(2023, 8, 1)));
        } finally {
          // Restore original DateTime.now
          DateTime.now = originalNow;
        }
      });

      test('calculates next payment date for annual subscription', () {
        // Use a fixed date for testing
        final now = DateTime(2023, 5, 15);
        final startDate = DateTime(2023, 1, 1);
        final payment = createPayment(startDate: startDate, months: 12);
        
        // Mock DateTime.now() for this test
        final originalNow = DateTime.now;
        DateTime.now = () => now;
        
        try {
          // Next payment should be on January 1, 2024
          expect(payment.nextPaymentDate, equals(DateTime(2024, 1, 1)));
        } finally {
          // Restore original DateTime.now
          DateTime.now = originalNow;
        }
      });
    });

    test('lastOccurrencePaid calculates correctly', () {
      // Use a fixed date for testing
      final now = DateTime(2023, 5, 15);
      final startDate = DateTime(2023, 1, 1);
      final payment = createPayment(startDate: startDate, months: 3);
      
      // Mock DateTime.now() for this test
      final originalNow = DateTime.now;
      DateTime.now = () => now;
      
      try {
        // Next payment is August 1, 2023, so last payment was May 1, 2023
        expect(payment.lastOccurrencePaid, equals(DateTime(2023, 5, 1)));
      } finally {
        // Restore original DateTime.now
        DateTime.now = originalNow;
      }
    });

    test('copyWith creates a new payment with updated values', () {
      final startDate = DateTime(2023, 1, 1);
      final endDate = DateTime(2023, 12, 31);
      final payment = createPayment(
        id: '1',
        price: 9.99,
        startDate: startDate,
        endDate: endDate,
        months: 1,
        currency: 'USD',
      );

      final newStartDate = DateTime(2023, 2, 1);
      final newEndDate = DateTime(2024, 1, 31);
      final updatedPayment = payment.copyWith(
        price: 19.99,
        startDate: newStartDate,
        endDate: newEndDate,
        months: 3,
        currency: 'EUR',
      );

      // Original payment should be unchanged
      expect(payment.id, equals('1'));
      expect(payment.price, equals(9.99));
      expect(payment.startDate, equals(startDate));
      expect(payment.endDate, equals(endDate));
      expect(payment.months, equals(1));
      expect(payment.currency, equals('USD'));

      // Updated payment should have new values
      expect(updatedPayment.id, equals('1')); // ID should remain the same
      expect(updatedPayment.price, equals(19.99));
      expect(updatedPayment.startDate, equals(newStartDate));
      expect(updatedPayment.endDate, equals(newEndDate));
      expect(updatedPayment.months, equals(3));
      expect(updatedPayment.currency, equals('EUR'));
    });

    test('toJson converts payment to JSON correctly', () {
      final startDate = DateTime(2023, 1, 1);
      final endDate = DateTime(2023, 12, 31);
      final payment = createPayment(
        id: '1',
        price: 9.99,
        startDate: startDate,
        endDate: endDate,
        months: 1,
        currency: 'USD',
      );

      final json = payment.toJson();

      expect(json['id'], equals('1'));
      expect(json['price'], equals(9.99));
      expect(json['startDate'], equals(startDate.millisecondsSinceEpoch));
      expect(json['endDate'], equals(endDate.millisecondsSinceEpoch));
      expect(json['months'], equals(1));
      expect(json['currency'], equals('USD'));
    });

    test('fromJson creates payment from JSON correctly', () {
      final startDate = DateTime(2023, 1, 1);
      final endDate = DateTime(2023, 12, 31);
      final json = {
        'id': '1',
        'price': 9.99,
        'startDate': startDate.millisecondsSinceEpoch,
        'endDate': endDate.millisecondsSinceEpoch,
        'months': 1,
        'currency': 'USD',
      };

      final payment = SubscriptionPayment.fromJson(json);

      expect(payment.id, equals('1'));
      expect(payment.price, equals(9.99));
      expect(payment.startDate, equals(startDate));
      expect(payment.endDate, equals(endDate));
      expect(payment.months, equals(1));
      expect(payment.currency, equals('USD'));
    });

    test('fromJson handles null endDate', () {
      final startDate = DateTime(2023, 1, 1);
      final json = {
        'id': '1',
        'price': 9.99,
        'startDate': startDate.millisecondsSinceEpoch,
        'endDate': null,
        'months': 1,
        'currency': 'USD',
      };

      final payment = SubscriptionPayment.fromJson(json);

      expect(payment.endDate, isNull);
    });

    test('fromJson uses USD as default currency for backward compatibility', () {
      final startDate = DateTime(2023, 1, 1);
      final json = {
        'id': '1',
        'price': 9.99,
        'startDate': startDate.millisecondsSinceEpoch,
        'endDate': null,
        'months': 1,
        // No currency field
      };

      final payment = SubscriptionPayment.fromJson(json);

      expect(payment.currency, equals('USD'));
    });
  });
}
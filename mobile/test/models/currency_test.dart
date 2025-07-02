import 'package:flutter_test/flutter_test.dart';
import 'package:subscription_tracker/models/currency.dart';

void main() {
  group('Currency', () {
    test('fromCode returns correct currency for valid code', () {
      expect(Currency.fromCode('USD'), equals(Currency.USD));
      expect(Currency.fromCode('EUR'), equals(Currency.EUR));
      expect(Currency.fromCode('GBP'), equals(Currency.GBP));
      expect(Currency.fromCode('JPY'), equals(Currency.JPY));
      expect(Currency.fromCode('CAD'), equals(Currency.CAD));
      expect(Currency.fromCode('AUD'), equals(Currency.AUD));
      expect(Currency.fromCode('CHF'), equals(Currency.CHF));
      expect(Currency.fromCode('CNY'), equals(Currency.CNY));
      expect(Currency.fromCode('INR'), equals(Currency.INR));
    });

    test('fromCode returns USD for invalid code', () {
      expect(Currency.fromCode('INVALID'), equals(Currency.USD));
      expect(Currency.fromCode(''), equals(Currency.USD));
    });

    test('codes returns all currency codes', () {
      final codes = Currency.codes;
      expect(codes, contains('USD'));
      expect(codes, contains('EUR'));
      expect(codes, contains('GBP'));
      expect(codes, contains('JPY'));
      expect(codes, contains('CAD'));
      expect(codes, contains('AUD'));
      expect(codes, contains('CHF'));
      expect(codes, contains('CNY'));
      expect(codes, contains('INR'));
      expect(codes.length, equals(9)); // Total number of currencies
    });

    test('formatAmount formats USD correctly', () {
      final currency = Currency.USD;
      expect(currency.formatAmount(10.5), equals('\$10.50'));
      expect(currency.formatAmount(0), equals('\$0.00'));
      expect(currency.formatAmount(1000), equals('\$1000.00'));
    });

    test('formatAmount formats EUR correctly', () {
      final currency = Currency.EUR;
      expect(currency.formatAmount(10.5), equals('10.50 €'));
      expect(currency.formatAmount(0), equals('0.00 €'));
      expect(currency.formatAmount(1000), equals('1000.00 €'));
    });

    test('formatAmount formats other currencies correctly', () {
      final gbp = Currency.GBP;
      expect(gbp.formatAmount(10.5), equals('£10.50'));
      
      final jpy = Currency.JPY;
      expect(jpy.formatAmount(10.5), equals('¥10.50'));
    });

    test('defaultCode is USD', () {
      expect(Currency.defaultCode, equals('USD'));
    });
  });
}
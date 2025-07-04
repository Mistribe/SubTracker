import 'package:flutter_test/flutter_test.dart';
import 'package:subscription_tracker/models/currency.dart';
import 'package:subscription_tracker/services/currency_converter.dart';

void main() {
  group('CurrencyConverter', () {
    test('convert returns same amount when currencies are the same', () {
      final amount = 100.0;
      final result = CurrencyConverter.convert(amount, 'USD', 'USD');
      expect(result, equals(amount));
    });

    test('convert correctly converts from USD to EUR', () {
      final amount = 100.0;
      final usdToEurRate = Currency.EUR.exchangeRate; // EUR rate relative to USD
      final expected = amount * usdToEurRate;
      final result = CurrencyConverter.convert(amount, 'USD', 'EUR');
      expect(result, equals(expected));
    });

    test('convert correctly converts from EUR to USD', () {
      final amount = 100.0;
      final eurToUsdRate = 1 / Currency.EUR.exchangeRate; // USD rate relative to EUR
      final expected = amount * eurToUsdRate;
      final result = CurrencyConverter.convert(amount, 'EUR', 'USD');
      expect(result, closeTo(expected, 0.001)); // Using closeTo to handle floating point precision
    });

    test('convert correctly converts between non-USD currencies', () {
      final amount = 100.0;
      // Convert from EUR to GBP
      // First convert EUR to USD, then USD to GBP
      final eurToUsdRate = 1 / Currency.EUR.exchangeRate;
      final usdToGbpRate = Currency.GBP.exchangeRate;
      final expected = amount * eurToUsdRate * usdToGbpRate;
      final result = CurrencyConverter.convert(amount, 'EUR', 'GBP');
      expect(result, closeTo(expected, 0.001));
    });

    test('getCurrencySymbol returns correct symbol for USD', () {
      final symbol = CurrencyConverter.getCurrencySymbol('USD');
      expect(symbol, equals('\$'));
    });

    test('getCurrencySymbol returns correct symbol for EUR', () {
      final symbol = CurrencyConverter.getCurrencySymbol('EUR');
      expect(symbol, equals('€'));
    });

    test('formatAmountWithCurrency formats USD correctly', () {
      final amount = 100.0;
      final formatted = CurrencyConverter.formatAmountWithCurrency(amount, 'USD');
      expect(formatted, equals('\$100.00'));
    });

    test('formatAmountWithCurrency formats EUR correctly', () {
      final amount = 100.0;
      final formatted = CurrencyConverter.formatAmountWithCurrency(amount, 'EUR');
      expect(formatted, equals('100.00 €'));
    });
  });
}
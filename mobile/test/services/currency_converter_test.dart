import 'package:flutter_test/flutter_test.dart';
import 'package:mobile/services/currency_converter.dart';

void main() {
  group('CurrencyConverter', () {
    test('should convert USD to EUR correctly', () {
      // Arrange
      const double amountInUSD = 100.0;
      const String fromCurrency = 'USD';
      const String toCurrency = 'EUR';

      // Act
      final result = CurrencyConverter.convert(amountInUSD, fromCurrency, toCurrency);

      // Assert
      expect(result, equals(85.0)); // Based on our fake exchange rate
    });

    test('should convert EUR to USD correctly', () {
      // Arrange
      const double amountInEUR = 85.0;
      const String fromCurrency = 'EUR';
      const String toCurrency = 'USD';

      // Act
      final result = CurrencyConverter.convert(amountInEUR, fromCurrency, toCurrency);

      // Assert
      expect(result, equals(100.0)); // Based on our fake exchange rate
    });

    test('should return the same amount when currencies are the same', () {
      // Arrange
      const double amount = 100.0;
      const String currency = 'USD';

      // Act
      final result = CurrencyConverter.convert(amount, currency, currency);

      // Assert
      expect(result, equals(amount));
    });

    test('should format amount with correct currency symbol', () {
      // Arrange
      const double amount = 100.0;

      // Act & Assert
      expect(CurrencyConverter.formatAmountWithCurrency(amount, 'USD'), equals('\$ 100.00'));
      expect(CurrencyConverter.formatAmountWithCurrency(amount, 'EUR'), equals('€ 100.00'));
      expect(CurrencyConverter.formatAmountWithCurrency(amount, 'GBP'), equals('£ 100.00'));
    });
  });
}

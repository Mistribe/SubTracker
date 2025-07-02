import '../models/currency.dart';

class CurrencyConverter {
  // Convert an amount from one currency to another
  static double convert(double amount, String fromCurrencyCode, String toCurrencyCode) {
    // If currencies are the same, no conversion needed
    if (fromCurrencyCode == toCurrencyCode) {
      return amount;
    }

    final fromCurrency = Currency.fromCode(fromCurrencyCode);
    final toCurrency = Currency.fromCode(toCurrencyCode);

    // Convert from source currency to USD (base currency)
    double amountInUSD = amount / fromCurrency.exchangeRate;

    // Convert from USD to target currency
    return amountInUSD * toCurrency.exchangeRate;
  }

  // Get the currency symbol for display
  static String getCurrencySymbol(String currencyCode) {
    return Currency.fromCode(currencyCode).symbol;
  }

  // Format an amount with the appropriate currency symbol
  static String formatAmountWithCurrency(double amount, String currencyCode) {
    return Currency.fromCode(currencyCode).formatAmount(amount);
  }
}

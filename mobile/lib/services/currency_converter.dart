class CurrencyConverter {
  // Fake exchange rates relative to USD
  static const Map<String, double> _exchangeRates = {
    'USD': 1.0,
    'EUR': 0.85,
    'GBP': 0.75,
    'JPY': 110.0,
    'CAD': 1.25,
    'AUD': 1.35,
    'CHF': 0.92,
    'CNY': 6.45,
    'INR': 74.5,
  };

  // Convert an amount from one currency to another
  static double convert(double amount, String fromCurrency, String toCurrency) {
    // If currencies are the same, no conversion needed
    if (fromCurrency == toCurrency) {
      return amount;
    }

    // Convert from source currency to USD (base currency)
    double amountInUSD = amount / _exchangeRates[fromCurrency]!;
    
    // Convert from USD to target currency
    return amountInUSD * _exchangeRates[toCurrency]!;
  }

  // Get the currency symbol for display
  static String getCurrencySymbol(String currencyCode) {
    switch (currencyCode) {
      case 'USD':
        return '\$';
      case 'EUR':
        return '€';
      case 'GBP':
        return '£';
      case 'JPY':
        return '¥';
      case 'CAD':
        return 'C\$';
      case 'AUD':
        return 'A\$';
      case 'CHF':
        return 'CHF';
      case 'CNY':
        return '¥';
      case 'INR':
        return '₹';
      default:
        return currencyCode;
    }
  }

  // Format an amount with the appropriate currency symbol
  static String formatAmountWithCurrency(double amount, String currencyCode) {
    return '${getCurrencySymbol(currencyCode)} ${amount.toStringAsFixed(2)}';
  }
}
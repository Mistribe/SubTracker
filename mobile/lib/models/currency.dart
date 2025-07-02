enum Currency {
  USD(code: 'USD', symbol: '\$', exchangeRate: 1.0),
  EUR(code: 'EUR', symbol: '€', exchangeRate: 0.85),
  GBP(code: 'GBP', symbol: '£', exchangeRate: 0.75),
  JPY(code: 'JPY', symbol: '¥', exchangeRate: 110.0),
  CAD(code: 'CAD', symbol: 'C\$', exchangeRate: 1.25),
  AUD(code: 'AUD', symbol: 'A\$', exchangeRate: 1.35),
  CHF(code: 'CHF', symbol: 'CHF', exchangeRate: 0.92),
  CNY(code: 'CNY', symbol: '¥', exchangeRate: 6.45),
  INR(code: 'INR', symbol: '₹', exchangeRate: 74.5);

  static const String defaultCode = 'USD';

  const Currency({
    required this.code,
    required this.symbol,
    required this.exchangeRate,
  });

  final String code;
  final String symbol;
  final double exchangeRate;

  // Get a Currency from its code
  static Currency fromCode(String code) {
    return Currency.values.firstWhere(
      (currency) => currency.code == code,
      orElse: () => Currency.USD, // Default to USD if not found
    );
  }

  // Get all currency codes
  static List<String> get codes => Currency.values.map((c) => c.code).toList();

  // Format an amount with the currency symbol
  String formatAmount(double amount) {
    return '$symbol ${amount.toStringAsFixed(2)}';
  }
}

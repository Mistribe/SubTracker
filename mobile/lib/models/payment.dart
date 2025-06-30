class PriceChange {
  final double price;
  final DateTime endDate;

  PriceChange({
    required this.price,
    required this.endDate,
  });

  Map<String, dynamic> toJson() {
    return {
      'price': price,
      'endDate': endDate.millisecondsSinceEpoch,
    };
  }

  factory PriceChange.fromJson(Map<String, dynamic> json) {
    return PriceChange(
      price: json['price'],
      endDate: DateTime.fromMillisecondsSinceEpoch(json['endDate']),
    );
  }
}

class Payment {
  final String id;
  final String name;
  final double price;
  final bool isAnnual;
  final DateTime paymentDate;
  final List<PriceChange> priceHistory;

  Payment({
    required this.id,
    required this.name,
    required this.price,
    required this.isAnnual,
    DateTime? paymentDate,
    List<PriceChange>? priceHistory,
  }) : 
    paymentDate = paymentDate ?? DateTime.now(),
    priceHistory = priceHistory ?? [];

  // Get price at a specific date, considering price history
  double getPriceAtDate(DateTime date) {
    // If the date is before the payment start date, return the initial price
    if (date.isBefore(paymentDate)) {
      return price;
    }

    // Sort price history by date (newest first)
    final sortedHistory = List<PriceChange>.from(priceHistory)
      ..sort((a, b) => b.endDate.compareTo(a.endDate));

    // Find the most recent price change before or at the given date
    for (var priceChange in sortedHistory) {
      if (!date.isBefore(priceChange.endDate)) {
        return priceChange.price;
      }
    }

    // If no applicable price change found, return the current price
    return price;
  }

  // Calculate the monthly cost based on current price
  double get monthlyCost => isAnnual ? price / 12 : price;

  // Calculate the monthly cost at a specific date
  double getMonthlyCostAtDate(DateTime date) {
    final priceAtDate = getPriceAtDate(date);
    return isAnnual ? priceAtDate / 12 : priceAtDate;
  }

  // Calculate the next payment date
  DateTime get nextPaymentDate {
    final now = DateTime.now();
    if (isAnnual) {
      // For annual payments, add 1 year from the payment date
      DateTime nextDate = DateTime(paymentDate.year + 1, paymentDate.month, paymentDate.day);

      // If the next date is in the past, keep adding years until it's in the future
      while (nextDate.isBefore(now)) {
        nextDate = DateTime(nextDate.year + 1, nextDate.month, nextDate.day);
      }
      return nextDate;
    } else {
      final currentDate = DateTime.now();
      // For monthly payments, add 1 month from the payment date
      DateTime nextDate = DateTime(currentDate.year, currentDate.month + 1, paymentDate.day);

      return nextDate;
    }
  }

  // Format the next payment date as a string
  String get formattedNextPaymentDate {
    final date = nextPaymentDate;
    return '${date.month}/${date.day}/${date.year}';
  }

  // Calculate the total amount spent since the payment started, considering price history
  double get totalAmountSpent {
    final now = DateTime.now();
    final startDate = paymentDate;
    double total = 0.0;

    if (isAnnual) {
      // For annual payments, calculate each year's payment separately
      DateTime currentDate = DateTime(startDate.year, startDate.month, startDate.day);

      while (currentDate.isBefore(now)) {
        // Get the price at this payment date
        final priceAtDate = getPriceAtDate(currentDate);
        total += priceAtDate;

        // Move to next year
        currentDate = DateTime(currentDate.year + 1, currentDate.month, currentDate.day);
      }
    } else {
      // For monthly payments, calculate each month's payment separately
      DateTime currentDate = DateTime(startDate.year, startDate.month, startDate.day);

      while (currentDate.isBefore(now)) {
        // Get the price at this payment date
        final priceAtDate = getPriceAtDate(currentDate);
        total += priceAtDate;

        // Move to next month
        currentDate = DateTime(currentDate.year, currentDate.month + 1, currentDate.day);
      }
    }

    return total;
  }

  // Format the total amount spent as a string
  String get formattedTotalAmountSpent {
    return '\$${totalAmountSpent.toStringAsFixed(2)}';
  }

  // Create a copy of this payment with updated fields
  Payment copyWith({
    String? id,
    String? name,
    double? price,
    bool? isAnnual,
    DateTime? paymentDate,
    List<PriceChange>? priceHistory,
  }) {
    return Payment(
      id: id ?? this.id,
      name: name ?? this.name,
      price: price ?? this.price,
      isAnnual: isAnnual ?? this.isAnnual,
      paymentDate: paymentDate ?? this.paymentDate,
      priceHistory: priceHistory ?? this.priceHistory,
    );
  }

  // Convert to and from JSON for potential future persistence
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'price': price,
      'isAnnual': isAnnual,
      'paymentDate': paymentDate.millisecondsSinceEpoch,
      'priceHistory': priceHistory.map((change) => change.toJson()).toList(),
    };
  }

  factory Payment.fromJson(Map<String, dynamic> json) {
    return Payment(
      id: json['id'],
      name: json['name'],
      price: json['price'],
      isAnnual: json['isAnnual'],
      paymentDate: DateTime.fromMillisecondsSinceEpoch(json['paymentDate']),
      priceHistory: json['priceHistory'] != null
          ? (json['priceHistory'] as List)
              .map((item) => PriceChange.fromJson(item))
              .toList()
          : [],
    );
  }
}

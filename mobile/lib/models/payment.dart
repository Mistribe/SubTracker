import 'package:hive/hive.dart';

part 'payment.g.dart';

@HiveType(typeId: 1)
class PriceChange {
  @HiveField(0)
  final double price;

  @HiveField(1)
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

@HiveType(typeId: 2)
class PaymentHistory {
  @HiveField(0)
  final double price;

  @HiveField(1)
  final DateTime startDate;

  @HiveField(2)
  final DateTime endDate;

  @HiveField(3)
  final bool isActive;

  PaymentHistory({
    required this.price,
    required this.startDate,
    required this.endDate,
    this.isActive = true,
  });

  Map<String, dynamic> toJson() {
    return {
      'price': price,
      'startDate': startDate.millisecondsSinceEpoch,
      'endDate': endDate.millisecondsSinceEpoch,
      'isActive': isActive,
    };
  }

  factory PaymentHistory.fromJson(Map<String, dynamic> json) {
    return PaymentHistory(
      price: json['price'],
      startDate: DateTime.fromMillisecondsSinceEpoch(json['startDate']),
      endDate: DateTime.fromMillisecondsSinceEpoch(json['endDate']),
      isActive: json['isActive'] ?? true,
    );
  }
}

@HiveType(typeId: 0)
class Payment extends HiveObject {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String name;

  @HiveField(2)
  final double price;

  @HiveField(3)
  final bool isAnnual;

  @HiveField(4)
  final DateTime paymentDate;

  @HiveField(5)
  final List<PriceChange> priceHistory;

  @HiveField(6)
  final bool isStopped;

  @HiveField(7)
  final DateTime? reactivationDate;

  @HiveField(8)
  final DateTime? stopDate;

  @HiveField(9)
  final List<PaymentHistory> paymentHistory;

  Payment({
    required this.id,
    required this.name,
    required this.price,
    required this.isAnnual,
    DateTime? paymentDate,
    List<PriceChange>? priceHistory,
    List<PaymentHistory>? paymentHistory,
    this.isStopped = false,
    this.reactivationDate,
    this.stopDate,
  }) : 
    paymentDate = paymentDate ?? DateTime.now(),
    priceHistory = priceHistory ?? [],
    paymentHistory = paymentHistory ?? [];

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
  double get monthlyCost {
    // If payment is stopped and not yet reactivated, return 0
    if (isStopped && (reactivationDate == null || reactivationDate!.isAfter(DateTime.now()))) {
      return 0;
    }
    return isAnnual ? price / 12 : price;
  }

  // Calculate the monthly cost at a specific date
  double getMonthlyCostAtDate(DateTime date) {
    // If payment is stopped and the date is before reactivation date, return 0
    if (isStopped && (reactivationDate == null || date.isBefore(reactivationDate!))) {
      return 0;
    }
    final priceAtDate = getPriceAtDate(date);
    return isAnnual ? priceAtDate / 12 : priceAtDate;
  }

  // Calculate the next payment date
  DateTime get nextPaymentDate {
    final now = DateTime.now();

    // If payment is stopped, the next payment date is the reactivation date or far future if none
    if (isStopped) {
      return reactivationDate ?? DateTime(9999, 12, 31); // Far future date if no reactivation date
    }

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

    // If the payment starts today or in the future, return 0
    if (!startDate.isBefore(DateTime(now.year, now.month, now.day))) {
      return 0.0;
    }

    // If payment is stopped and never activated, return 0
    if (isStopped && reactivationDate == null) {
      return 0.0;
    }

    if (isAnnual) {
      // For annual payments, calculate each year's payment separately
      DateTime currentDate = DateTime(startDate.year, startDate.month, startDate.day);

      while (currentDate.isBefore(now)) {
        // Skip if payment is stopped at this date and not yet reactivated
        if (isStopped && (reactivationDate == null || currentDate.isBefore(reactivationDate!))) {
          // Move to next year
          currentDate = DateTime(currentDate.year + 1, currentDate.month, currentDate.day);
          continue;
        }

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
        // Skip if payment is stopped at this date and not yet reactivated
        if (isStopped && (reactivationDate == null || currentDate.isBefore(reactivationDate!))) {
          // Move to next month
          currentDate = DateTime(currentDate.year, currentDate.month + 1, currentDate.day);
          continue;
        }

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
    List<PaymentHistory>? paymentHistory,
    bool? isStopped,
    DateTime? reactivationDate,
    DateTime? stopDate,
  }) {
    return Payment(
      id: id ?? this.id,
      name: name ?? this.name,
      price: price ?? this.price,
      isAnnual: isAnnual ?? this.isAnnual,
      paymentDate: paymentDate ?? this.paymentDate,
      priceHistory: priceHistory ?? this.priceHistory,
      paymentHistory: paymentHistory ?? this.paymentHistory,
      isStopped: isStopped ?? this.isStopped,
      reactivationDate: reactivationDate ?? this.reactivationDate,
      stopDate: stopDate ?? this.stopDate,
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
      'paymentHistory': paymentHistory.map((history) => history.toJson()).toList(),
      'isStopped': isStopped,
      'reactivationDate': reactivationDate?.millisecondsSinceEpoch,
      'stopDate': stopDate?.millisecondsSinceEpoch,
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
      paymentHistory: json['paymentHistory'] != null
          ? (json['paymentHistory'] as List)
              .map((item) => PaymentHistory.fromJson(item))
              .toList()
          : [],
      isStopped: json['isStopped'] ?? false,
      reactivationDate: json['reactivationDate'] != null
          ? DateTime.fromMillisecondsSinceEpoch(json['reactivationDate'])
          : null,
      stopDate: json['stopDate'] != null
          ? DateTime.fromMillisecondsSinceEpoch(json['stopDate'])
          : null,
    );
  }
}

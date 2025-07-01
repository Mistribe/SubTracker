import 'package:hive/hive.dart';

part 'subscription_payment.g.dart';

@HiveType(typeId: 1)
class SubscriptionPayment {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final double price;

  @HiveField(2)
  final DateTime startDate;

  @HiveField(3)
  final DateTime? endDate;

  // The number of month cover by the price before the next payment.
  // If 1 this is a monthly subscription
  // If 3, the subscription is renew every 3 months and the price cover 3 months. To have the monthly price we should make price / 3
  // If 12, the subscription is annually.
  // todo rename in frequency ?
  @HiveField(4)
  final int months;

  @HiveField(5)
  final String currency;

  SubscriptionPayment({
    required this.id,
    required this.price,
    required this.startDate,
    required this.endDate,
    required this.months,
    this.currency = 'USD',
  });

  bool get isActive {
    if (endDate == null) {
      return true;
    }
    return endDate!.isAfter(DateTime.now());
  }

  int get totalMonthElapsed {
    final until = endDate == null ? DateTime.now() : endDate!;
    int yearDifference = until.year - startDate.year;
    int monthDifference = until.month - startDate.month;

    return yearDifference * 12 + monthDifference;
  }

  double get totalAmountSpent {
    return totalMonthElapsed * monthlyCost;
  }

  double get monthlyCost {
    return price / months;
  }

  DateTime get lastOccurrencePaid {
    final year = (months / 12).truncate();
    final month = (months % 12);

    final nextPaymentDate = this.nextPaymentDate;
    return DateTime(
      nextPaymentDate.year - year,
      nextPaymentDate.month - month,
      nextPaymentDate.day,
    );
  }

  DateTime get nextPaymentDate {
    final now = DateTime.now();

    // If payment is stopped, the next payment date is the reactivation date or far future if none
    if (isActive == false) {
      return DateTime(9999, 12, 31); // Far future date if no reactivation date
    }

    if (startDate.isAfter(now)) {
      return startDate;
    }

    final year = (months / 12).truncate();
    final month = (months % 12);
    DateTime nextDate = DateTime(
      now.year + year,
      now.month + month,
      startDate.day,
    );

    return nextDate;
  }

  SubscriptionPayment copyWith({
    String? id,
    double? price,
    DateTime? startDate,
    DateTime? endDate,
    int? months,
    String? currency,
  }) {
    return SubscriptionPayment(
      id: id ?? this.id,
      price: price ?? this.price,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
      months: months ?? this.months,
      currency: currency ?? this.currency,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'price': price,
      'startDate': startDate.millisecondsSinceEpoch,
      'endDate': endDate?.millisecondsSinceEpoch,
      'months': months,
      'currency': currency,
    };
  }

  factory SubscriptionPayment.fromJson(Map<String, dynamic> json) {
    return SubscriptionPayment(
      id: json['id'],
      price: json['price'],
      startDate: DateTime.fromMillisecondsSinceEpoch(json['startDate']),
      endDate: json['endDate'] != null
          ? DateTime.fromMillisecondsSinceEpoch(json['endDate'])
          : null,
      months: json['months'],
      currency: json['currency'] ?? 'USD', // Default to USD for backward compatibility
    );
  }
}

import 'package:hive/hive.dart';
import 'package:subscription_tracker/models/subscription_state.dart';
import 'currency.dart';

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

  @HiveField(6)
  final DateTime createdAt;

  @HiveField(7)
  final DateTime updatedAt;

  @HiveField(8)
  final String eTag;

  SubscriptionPayment({
    required this.id,
    required this.price,
    required this.startDate,
    required this.endDate,
    required this.months,
    required this.currency,
    this.eTag = '',
    DateTime? createdAt,
    DateTime? updatedAt,
  }) : createdAt = createdAt ?? DateTime.now(),
       updatedAt = updatedAt ?? DateTime.now();

  SubscriptionState get state {
    if (startDate.isAfter(DateTime.now())) {
      return SubscriptionState.notStarted;
    }
    if (endDate == null || endDate!.isAfter(DateTime.now())) {
      return SubscriptionState.active;
    }
    return SubscriptionState.ended;
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
    if (state == SubscriptionState.ended) {
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
    String? eTag,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return SubscriptionPayment(
      id: id ?? this.id,
      price: price ?? this.price,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
      months: months ?? this.months,
      currency: currency ?? this.currency,
      eTag: eTag ?? this.eTag,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'price': price,
      'start_date': startDate.toIso8601String(),
      'end_date': endDate?.toIso8601String(),
      'months': months,
      'currency': currency,
      'etag': eTag,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  factory SubscriptionPayment.fromJson(Map<String, dynamic> json) {
    return SubscriptionPayment(
      id: json['id'],
      price: (json['price'] is int)
          ? (json['price'] as int).toDouble()
          : json['price'],
      startDate: DateTime.parse(json['start_date']),
      endDate: json['end_date'] != null
          ? DateTime.parse(json['end_date'])
          : null,
      months: json['months'],
      currency: json['currency'] ?? Currency.USD.code,
      eTag: json['etag'],
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : null,
      updatedAt: json['updated_at'] != null
          ? DateTime.parse(json['updated_at'])
          : null,
    );
  }
}

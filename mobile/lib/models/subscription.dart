import 'package:hive/hive.dart';
import 'package:subscription_tracker/models/subscription_payment.dart';

part 'subscription.g.dart';

@HiveType(typeId: 0)
class Subscription extends HiveObject {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String name;

  @HiveField(3)
  final List<SubscriptionPayment> subscriptionPayments;

  Subscription({
    required this.id,
    required this.name,
    List<SubscriptionPayment>? subscriptionPayments,
  }) : subscriptionPayments = subscriptionPayments ?? [];

  SubscriptionPayment getLastPaymentDetail() {
    final sortedHistory = List<SubscriptionPayment>.from(subscriptionPayments)
      ..sort((a, b) => b.startDate.compareTo(a.startDate));

    return sortedHistory[0];
  }

  SubscriptionPayment getPaymentDetailAtDate(DateTime date) {
    if (subscriptionPayments.isEmpty) {
      throw StateError('No payment details available');
    }
    // Sort payment detail by start date (newest first)
    final sortedDetails = List<SubscriptionPayment>.from(subscriptionPayments)
      ..sort((a, b) => b.startDate.compareTo(a.startDate));

    // Find the payment detail entry that covers the given date
    for (var detail in sortedDetails) {
      if (date.isAfter(detail.startDate) &&
              (detail.endDate != null && date.isBefore(detail.endDate!) ||
                  detail.endDate == null) ||
          date.isAtSameMomentAs(detail.startDate)) {
        return detail;
      }
    }
    throw StateError('No payment detail found for the given date');
  }

  // Get price at a specific date, considering payment history
  double getPriceAtDate(DateTime date) {
    return getPaymentDetailAtDate(date).price;
  }

  // Get if the subscription is active
  bool get isActive {
    final detail = getLastPaymentDetail();
    return detail.isActive;
  }

  // Calculate the monthly cost based on current price
  double get monthlyCost {
    var currentDetail = getLastPaymentDetail();

    if (!currentDetail.isActive) {
      return 0;
    }
    return currentDetail.monthlyCost;
  }

  double get annualCost {
    return monthlyCost * 12;
  }

  // Calculate the monthly cost at a specific date
  double getMonthlyCostAtDate(DateTime date) {
    final detail = getPaymentDetailAtDate(date);
    // If payment is stopped and the date is before reactivation date, return 0
    if (!detail.isActive) {
      return 0;
    }
    return detail.monthlyCost;
  }

  // Calculate the next payment date
  DateTime get nextPaymentDate {
    final currentDetail = getLastPaymentDetail();
    return currentDetail.nextPaymentDate;
  }

  DateTime get lastPaymentDate {
    final currentDetail = getLastPaymentDetail();
    return currentDetail.lastOccurrencePaid;
  }

  SubscriptionPayment? findDetailById(String paymentDetailId) {
    final idx = subscriptionPayments.indexWhere(
      (detail) => detail.id == paymentDetailId,
    );
    if (idx < 0) {
      return null;
    }

    return subscriptionPayments[idx];
  }

  void setPaymentDetail(SubscriptionPayment detail) {
    final idx = subscriptionPayments.indexWhere(
      (detail) => detail.id == detail.id,
    );
    if (idx < 0) {
      throw StateError('No payment details available');
    }

    subscriptionPayments[idx] = detail;
  }

  void setEndDateToCurrentPaymentDetail(DateTime effectiveDate) {
    final currentDetail = getLastPaymentDetail();
    setPaymentDetail(currentDetail.copyWith(endDate: effectiveDate));
  }

  void addPaymentDetail(SubscriptionPayment paymentDetail) {
    subscriptionPayments.add(paymentDetail);
  }

  // Format the next payment date as a string
  String get formattedNextPaymentDate {
    final date = nextPaymentDate;
    return '${date.month}/${date.day}/${date.year}';
  }

  // Calculate the total amount spent since the payment started, considering price history
  double get totalAmountSpent {
    final now = DateTime.now();
    double total = 0.0;

    // If the payment starts today or in the future, return 0
    if (subscriptionPayments.isEmpty ||
        subscriptionPayments.length == 1 &&
            getLastPaymentDetail().startDate.isAfter(
              DateTime(now.year, now.month, now.day),
            )) {
      return 0.0;
    }

    for (var detail in subscriptionPayments) {
      total += detail.monthlyCost * detail.totalMonthElapsed;
    }

    return total;
  }

  // Format the total amount spent as a string
  String get formattedTotalAmountSpent {
    return '\$${totalAmountSpent.toStringAsFixed(2)}';
  }

  // Create a copy of this payment with updated fields
  Subscription copyWith({
    String? id,
    String? name,
    List<SubscriptionPayment>? subscriptionPayments,
  }) {
    return Subscription(
      id: id ?? this.id,
      name: name ?? this.name,
      subscriptionPayments: subscriptionPayments ?? this.subscriptionPayments,
    );
  }

  // Convert to and from JSON for potential future persistence
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'paymentDetails': subscriptionPayments
          .map((detail) => detail.toJson())
          .toList(),
    };
  }

  factory Subscription.fromJson(Map<String, dynamic> json) {
    return Subscription(
      id: json['id'],
      name: json['name'],
      subscriptionPayments: json['paymentDetails'] != null
          ? (json['paymentDetails'] as List)
                .map((item) => SubscriptionPayment.fromJson(item))
                .toList()
          : [],
    );
  }

  void endCurrentPaymentDetail(DateTime endDate) {
    final currentDetail = getLastPaymentDetail();
    setPaymentDetail(currentDetail.copyWith(endDate: endDate));
  }
}

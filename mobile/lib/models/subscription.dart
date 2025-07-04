import 'package:hive/hive.dart';
import 'package:subscription_tracker/models/subscription_payment.dart';
import 'package:subscription_tracker/models/label.dart';
import 'package:subscription_tracker/models/currency.dart';
import 'package:subscription_tracker/models/family_member.dart';
import 'package:subscription_tracker/models/subscription_state.dart';

part 'subscription.g.dart';

@HiveType(typeId: 0)
class Subscription extends HiveObject {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String name;

  @HiveField(3)
  final List<SubscriptionPayment> subscriptionPayments;

  @HiveField(4)
  final List<Label> labels;

  @HiveField(5)
  final List<FamilyMember> userFamilyMembers;

  @HiveField(6)
  final FamilyMember? payerFamilyMember;

  Subscription({
    required this.id,
    required this.name,
    List<SubscriptionPayment>? subscriptionPayments,
    List<Label>? labels,
    List<FamilyMember>? userFamilyMembers,
    this.payerFamilyMember,
  }) : subscriptionPayments = subscriptionPayments ?? [],
       labels = labels ?? [],
       userFamilyMembers = userFamilyMembers ?? [];

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

  // Get the subscription state
  SubscriptionState get state {
    final detail = getLastPaymentDetail();
    return detail.state;
  }

  // Calculate the monthly cost based on current price
  double get monthlyCost {
    var currentDetail = getLastPaymentDetail();

    if (currentDetail.state == SubscriptionState.ended) {
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
    if (detail.state != SubscriptionState.active) {
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
    final idx = subscriptionPayments.indexWhere((d) => d.id == detail.id);
    if (idx < 0) {
      throw StateError('No payment details available');
    }

    subscriptionPayments[idx] = detail;
    subscriptionPayments.sort((a, b) => a.startDate.compareTo(b.startDate));
  }

  void setEndDateToCurrentPaymentDetail(DateTime effectiveDate) {
    final currentDetail = getLastPaymentDetail();
    setPaymentDetail(currentDetail.copyWith(endDate: effectiveDate));
  }

  void addPaymentDetail(SubscriptionPayment paymentDetail) {
    subscriptionPayments.add(paymentDetail);
    subscriptionPayments.sort((a, b) => a.startDate.compareTo(b.startDate));
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
    final currencyCode = subscriptionPayments.isNotEmpty
        ? subscriptionPayments.last.currency
        : Currency.USD.code;
    final currency = Currency.fromCode(currencyCode);
    return currency.formatAmount(totalAmountSpent);
  }

  // Create a copy of this payment with updated fields
  Subscription copyWith({
    String? id,
    String? name,
    List<SubscriptionPayment>? subscriptionPayments,
    List<Label>? labels,
    List<FamilyMember>? userFamilyMembers,
    FamilyMember? payerFamilyMember,
  }) {
    return Subscription(
      id: id ?? this.id,
      name: name ?? this.name,
      subscriptionPayments: subscriptionPayments ?? this.subscriptionPayments,
      labels: labels ?? this.labels,
      userFamilyMembers: userFamilyMembers ?? this.userFamilyMembers,
      payerFamilyMember: payerFamilyMember ?? this.payerFamilyMember,
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
      'labels': labels.map((label) => label.toJson()).toList(),
      'userFamilyMembers': userFamilyMembers
          .map((member) => member.toJson())
          .toList(),
      'payerFamilyMember': payerFamilyMember?.toJson(),
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
      labels: json['labels'] != null
          ? (json['labels'] as List)
                .map((item) => Label.fromJson(item))
                .toList()
          : [],
      userFamilyMembers: json['userFamilyMembers'] != null
          ? (json['userFamilyMembers'] as List)
                .map((item) => FamilyMember.fromJson(item))
                .toList()
          : [],
      payerFamilyMember: json['payerFamilyMember'] != null
          ? FamilyMember.fromJson(json['payerFamilyMember'])
          : null,
    );
  }

  void endCurrentPaymentDetail(DateTime endDate) {
    final currentDetail = getLastPaymentDetail();
    setPaymentDetail(currentDetail.copyWith(endDate: endDate));
  }

  // Remove a payment detail by ID
  void removePaymentDetail(String paymentDetailId) {
    final index = subscriptionPayments.indexWhere(
      (detail) => detail.id == paymentDetailId,
    );

    if (index < 0) {
      throw StateError('Payment detail not found');
    }

    // Remove the payment detail
    subscriptionPayments.removeAt(index);
  }
}

import 'package:hive/hive.dart';

part 'payment.g.dart';

@HiveType(typeId: 1)
class PaymentDetail {
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

  PaymentDetail({
    required this.id,
    required this.price,
    required this.startDate,
    required this.endDate,
    required this.months,
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

    final year = (months / 12).truncate();
    final month = (months % 12);
    DateTime nextDate = DateTime(
      now.year + year,
      now.month + month,
      startDate.day,
    );

    return nextDate;
  }

  PaymentDetail copyWith({
    String? id,
    double? price,
    DateTime? startDate,
    DateTime? endDate,
    int? months,
  }) {
    return PaymentDetail(
      id: id ?? this.id,
      price: price ?? this.price,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
      months: months ?? this.months,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'price': price,
      'startDate': startDate.millisecondsSinceEpoch,
      'endDate': endDate?.millisecondsSinceEpoch,
      'months': months,
    };
  }

  factory PaymentDetail.fromJson(Map<String, dynamic> json) {
    return PaymentDetail(
      id: json['id'],
      price: json['price'],
      startDate: DateTime.fromMillisecondsSinceEpoch(json['startDate']),
      endDate: json['endDate'] != null
          ? DateTime.fromMillisecondsSinceEpoch(json['endDate'])
          : null,
      months: json['months'],
    );
  }
}

@HiveType(typeId: 0)
class Payment extends HiveObject {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String name;

  @HiveField(3)
  final List<PaymentDetail> paymentDetails;

  Payment({
    required this.id,
    required this.name,
    List<PaymentDetail>? paymentDetails,
  }) : paymentDetails = paymentDetails ?? [];

  PaymentDetail getLastPaymentDetail() {
    final sortedHistory = List<PaymentDetail>.from(paymentDetails)
      ..sort((a, b) => b.startDate.compareTo(a.startDate));

    return sortedHistory[0];
  }

  PaymentDetail getPaymentDetailAtDate(DateTime date) {
    if (paymentDetails.isEmpty) {
      throw StateError('No payment details available');
    }
    // Sort payment detail by start date (newest first)
    final sortedDetails = List<PaymentDetail>.from(paymentDetails)
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

  PaymentDetail? findDetailById(String paymentDetailId) {
    final idx = paymentDetails.indexWhere(
      (detail) => detail.id == paymentDetailId,
    );
    if (idx < 0) {
      return null;
    }

    return paymentDetails[idx];
  }

  void setPaymentDetail(PaymentDetail detail) {
    final idx = paymentDetails.indexWhere((detail) => detail.id == detail.id);
    if (idx < 0) {
      throw StateError('No payment details available');
    }

    paymentDetails[idx] = detail;
  }

  void setEndDateToCurrentPaymentDetail(DateTime effectiveDate) {
    final currentDetail = getLastPaymentDetail();
    setPaymentDetail(currentDetail.copyWith(endDate: effectiveDate));
  }

  void addPaymentDetail(PaymentDetail paymentDetail) {
    paymentDetails.add(paymentDetail);
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
    if (paymentDetails.isEmpty ||
        paymentDetails.length == 1 &&
            getLastPaymentDetail().startDate.isAfter(
              DateTime(now.year, now.month, now.day),
            )) {
      return 0.0;
    }

    for (var detail in paymentDetails) {
      total += detail.monthlyCost * detail.totalMonthElapsed;
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
    List<PaymentDetail>? paymentDetails,
  }) {
    return Payment(
      id: id ?? this.id,
      name: name ?? this.name,
      paymentDetails: paymentDetails ?? this.paymentDetails,
    );
  }

  // Convert to and from JSON for potential future persistence
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'paymentDetails': paymentDetails
          .map((detail) => detail.toJson())
          .toList(),
    };
  }

  factory Payment.fromJson(Map<String, dynamic> json) {
    return Payment(
      id: json['id'],
      name: json['name'],
      paymentDetails: json['paymentDetails'] != null
          ? (json['paymentDetails'] as List)
                .map((item) => PaymentDetail.fromJson(item))
                .toList()
          : [],
    );
  }

  void endCurrentPaymentDetail(DateTime endDate) {
    final currentDetail = getLastPaymentDetail();
    setPaymentDetail(currentDetail.copyWith(endDate: endDate));
  }
}

class Payment {
  final String id;
  final String name;
  final double price;
  final bool isAnnual;
  final DateTime paymentDate;

  Payment({
    required this.id,
    required this.name,
    required this.price,
    required this.isAnnual,
    DateTime? paymentDate,
  }) : this.paymentDate = paymentDate ?? DateTime.now();

  // Calculate the monthly cost
  double get monthlyCost => isAnnual ? price / 12 : price;

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
      // For monthly payments, add 1 month from the payment date
      DateTime nextDate = DateTime(paymentDate.year, paymentDate.month + 1, paymentDate.day);

      // If the next date is in the past, keep adding months until it's in the future
      while (nextDate.isBefore(now)) {
        nextDate = DateTime(nextDate.year, nextDate.month + 1, nextDate.day);
      }
      return nextDate;
    }
  }

  // Format the next payment date as a string
  String get formattedNextPaymentDate {
    final date = nextPaymentDate;
    return '${date.month}/${date.day}/${date.year}';
  }

  // Calculate the total amount spent since the payment started
  double get totalAmountSpent {
    final now = DateTime.now();
    final startDate = paymentDate;

    if (isAnnual) {
      // Calculate years passed (including partial years)
      int yearsPassed = now.year - startDate.year;
      if (now.month < startDate.month || 
          (now.month == startDate.month && now.day < startDate.day)) {
        yearsPassed--;
      }
      return yearsPassed > 0 ? price * yearsPassed : 0;
    } else {
      // Calculate months passed
      int monthsPassed = (now.year - startDate.year) * 12 + now.month - startDate.month;
      if (now.day < startDate.day) {
        monthsPassed--;
      }
      return monthsPassed > 0 ? price * monthsPassed : 0;
    }
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
  }) {
    return Payment(
      id: id ?? this.id,
      name: name ?? this.name,
      price: price ?? this.price,
      isAnnual: isAnnual ?? this.isAnnual,
      paymentDate: paymentDate ?? this.paymentDate,
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
    };
  }

  factory Payment.fromJson(Map<String, dynamic> json) {
    return Payment(
      id: json['id'],
      name: json['name'],
      price: json['price'],
      isAnnual: json['isAnnual'],
      paymentDate: DateTime.fromMillisecondsSinceEpoch(json['paymentDate']),
    );
  }
}

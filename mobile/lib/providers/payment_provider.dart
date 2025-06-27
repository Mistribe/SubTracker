import 'dart:math';
import 'package:flutter/foundation.dart';
import '../models/payment.dart';

class PaymentProvider with ChangeNotifier {
  final List<Payment> _payments = [];

  // Getter for the payments list
  List<Payment> get payments => List.unmodifiable(_payments);

  // Calculate total monthly cost
  double get totalMonthlyCost {
    return _payments.fold(0, (sum, payment) => sum + payment.monthlyCost);
  }

  // Add a new payment
  void addPayment(String name, double price, bool isAnnual, {DateTime? paymentDate}) {
    final payment = Payment(
      id: _generateId(),
      name: name,
      price: price,
      isAnnual: isAnnual,
      paymentDate: paymentDate,
    );

    _payments.add(payment);
    notifyListeners();
  }

  // Remove a payment
  void removePayment(String id) {
    _payments.removeWhere((payment) => payment.id == id);
    notifyListeners();
  }

  // Update an existing payment
  void updatePayment(Payment updatedPayment) {
    final index = _payments.indexWhere((payment) => payment.id == updatedPayment.id);

    if (index >= 0) {
      _payments[index] = updatedPayment;
      notifyListeners();
    }
  }

  // Generate a unique ID for a new payment
  String _generateId() {
    return DateTime.now().millisecondsSinceEpoch.toString() + 
           Random().nextInt(10000).toString();
  }
}

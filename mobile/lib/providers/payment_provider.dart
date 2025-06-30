import 'dart:math';
import 'package:flutter/foundation.dart';
import '../models/payment.dart';
import '../repositories/payment_repository.dart';

class PaymentProvider with ChangeNotifier {
  final PaymentRepository paymentRepository;
  List<Payment> _payments = [];

  PaymentProvider({required this.paymentRepository}) {
    // Load payments from repository
    _loadPayments();
  }

  // Load payments from repository
  Future<void> _loadPayments() async {
    _payments = paymentRepository.getAllPayments();
    notifyListeners();
  }

  // Getter for the payments list
  List<Payment> get payments => List.unmodifiable(_payments);

  // Calculate total monthly cost
  double get totalMonthlyCost {
    return _payments.fold(0, (sum, payment) => sum + payment.monthlyCost);
  }

  // Get the count of active payments
  int get activePaymentsCount {
    return _payments.length;
  }

  // Calculate total amount spent across all payments
  double get totalAmountSpent {
    return _payments.fold(0, (sum, payment) => sum + payment.totalAmountSpent);
  }

  // Add a new payment
  Future<void> addPayment(String name, double price, bool isAnnual, {DateTime? paymentDate}) async {
    final payment = Payment(
      id: _generateId(),
      name: name,
      price: price,
      isAnnual: isAnnual,
      paymentDate: paymentDate,
    );

    // Add to local list
    _payments.add(payment);

    // Persist to storage
    await paymentRepository.addPayment(payment);

    notifyListeners();
  }

  // Remove a payment
  Future<void> removePayment(String id) async {
    _payments.removeWhere((payment) => payment.id == id);

    // Remove from storage
    await paymentRepository.deletePayment(id);

    notifyListeners();
  }

  // Update an existing payment
  Future<void> updatePayment(Payment updatedPayment) async {
    final index = _payments.indexWhere((payment) => payment.id == updatedPayment.id);

    if (index >= 0) {
      final oldPayment = _payments[index];

      // Check if price has changed
      if (oldPayment.price != updatedPayment.price) {
        // Create a copy of the price history
        final newPriceHistory = List<PriceChange>.from(updatedPayment.priceHistory);

        // Add the old price to the history with the current date
        // This records when the price was changed
        newPriceHistory.add(PriceChange(
          price: oldPayment.price,
          endDate: DateTime.now(),
        ));

        // Create a new payment with updated price history
        updatedPayment = updatedPayment.copyWith(
          priceHistory: newPriceHistory,
        );
      }

      // Update local list
      _payments[index] = updatedPayment;

      // Persist to storage
      await paymentRepository.updatePayment(updatedPayment);

      notifyListeners();
    }
  }

  // Add a price change at a specific date
  Future<void> addPriceChange(String paymentId, double newPrice, DateTime effectiveDate) async {
    final index = _payments.indexWhere((payment) => payment.id == paymentId);

    if (index >= 0) {
      final payment = _payments[index];

      // Create a copy of the price history
      final newPriceHistory = List<PriceChange>.from(payment.priceHistory);

      // Add the old price to the history with the effective date
      // This records when the price was changed
      newPriceHistory.add(PriceChange(
        price: payment.price,
        endDate: effectiveDate,
      ));

      // Sort price history by date (newest first)
      newPriceHistory.sort((a, b) => b.endDate.compareTo(a.endDate));

      // Create a new payment with updated price history and new current price
      final updatedPayment = payment.copyWith(
        price: newPrice,
        priceHistory: newPriceHistory,
      );

      // Update local list
      _payments[index] = updatedPayment;

      // Persist to storage
      await paymentRepository.updatePayment(updatedPayment);

      notifyListeners();
    }
  }

  // Update a price change at a specific index
  Future<void> updatePriceChange(String paymentId, int priceChangeIndex, double newPrice, DateTime newDate) async {
    final index = _payments.indexWhere((payment) => payment.id == paymentId);

    if (index >= 0) {
      final payment = _payments[index];

      // Make sure the index is valid
      if (priceChangeIndex >= 0 && priceChangeIndex < payment.priceHistory.length) {
        // Create a copy of the price history
        final newPriceHistory = List<PriceChange>.from(payment.priceHistory);

        // Update the price change at the specified index
        newPriceHistory[priceChangeIndex] = PriceChange(
          price: newPrice,
          endDate: newDate,
        );

        // Sort price history by date (newest first)
        newPriceHistory.sort((a, b) => b.endDate.compareTo(a.endDate));

        // Create a new payment with updated price history
        final updatedPayment = payment.copyWith(
          priceHistory: newPriceHistory,
        );

        // Update local list
        _payments[index] = updatedPayment;

        // Persist to storage
        await paymentRepository.updatePayment(updatedPayment);

        notifyListeners();
      }
    }
  }

  // Generate a unique ID for a new payment
  String _generateId() {
    return DateTime.now().millisecondsSinceEpoch.toString() + 
           Random().nextInt(10000).toString();
  }
}

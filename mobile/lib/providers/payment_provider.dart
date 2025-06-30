import 'dart:math';
import 'package:uuid/uuid.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';
import '../models/payment.dart';
import '../repositories/payment_repository.dart';

var uuid = Uuid();

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

  // Get the count of active payments (excluding stopped payments)
  int get activePaymentsCount {
    return _payments.where((payment) => payment.isActive).length;
  }

  // Calculate total amount spent across all payments
  double get totalAmountSpent {
    return _payments.fold(0, (sum, payment) => sum + payment.totalAmountSpent);
  }

  // Add a new payment
  Future<void> addPayment(
    String name,
    double price,
    int months,
    DateTime startDate, {
    DateTime? endDate,
  }) async {
    // Create initial payment history entry
    final initialDetail = [
      PaymentDetail(
        id: _generateId(),
        price: price,
        startDate: startDate,
        endDate: endDate,
        // Far future date
        months: months,
      ),
    ];

    final payment = Payment(
      id: _generateId(),
      name: name,
      paymentDetails: initialDetail,
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
  Future<void> updatePayment(String id, String name) async {
    final index = _payments.indexWhere((payment) => payment.id == id);

    if (index >= 0) {
      final payment = _payments[index];

      _payments[index] = payment.copyWith(name: name);

      // Persist to storage
      await paymentRepository.updatePayment(payment);

      notifyListeners();
    }
  }

  // Add a price change at a specific date
  Future<void> addPaymentDetailEntry(
    String paymentId,
    double newPrice,
    DateTime effectiveDate, {
    DateTime? endDate,
    int? months,
  }) async {
    final index = _payments.indexWhere((payment) => payment.id == paymentId);

    if (index >= 0) {
      final payment = _payments[index];

      final previousDetail = payment.getLastPaymentDetail();
      // Update the current payment detail with end date
      payment.setEndDateToCurrentPaymentDetail(effectiveDate);
      // Add new payment detail
      payment.addPaymentDetail(
        PaymentDetail(
          id: _generateId(),
          price: newPrice,
          startDate: effectiveDate,
          endDate: endDate,
          months: months ?? previousDetail.months,
        ),
      );

      // Persist to storage
      await paymentRepository.updatePayment(payment);

      notifyListeners();
    }
  }

  // Update a payment history entry
  Future<void> updatePaymentDetailEntry(
    String paymentId,
    String paymentDetailId,
    double newPrice,
    int months,
    DateTime startDate, {
    DateTime? endDate,
  }) async {
    final index = _payments.indexWhere((payment) => payment.id == paymentId);

    if (index >= 0) {
      final payment = _payments[index];

      payment.setPaymentDetail(
        PaymentDetail(
          id: paymentDetailId,
          price: newPrice,
          startDate: startDate,
          endDate: endDate,
          months: months,
        ),
      );

      // Persist to storage
      await paymentRepository.updatePayment(payment);

      notifyListeners();
    }
  }

  // Stop a payment
  Future<void> stopPayment(String paymentId, {DateTime? stopDate}) async {
    final index = _payments.indexWhere((payment) => payment.id == paymentId);

    if (index >= 0) {
      final payment = _payments[index];

      // If stopDate is not provided, use the last payment date
      final effectiveStopDate = stopDate ?? payment.lastPaymentDate;
      payment.endCurrentPaymentDetail(effectiveStopDate);

      // Persist to storage
      await paymentRepository.updatePayment(payment);

      notifyListeners();
    }
  }

  // Reactivate a payment at a specific date
  Future<void> reactivatePayment(
    String paymentId,
    DateTime reactivationDate,
  ) async {
    final index = _payments.indexWhere((payment) => payment.id == paymentId);

    if (index >= 0) {
      final payment = _payments[index];

      final previousDetail = payment.getLastPaymentDetail();

      payment.addPaymentDetail(
        PaymentDetail(
          id: _generateId(),
          price: previousDetail.price,
          startDate: reactivationDate,
          endDate: null,
          months: previousDetail.months,
        ),
      );

      // Persist to storage
      await paymentRepository.updatePayment(payment);

      notifyListeners();
    }
  }

  // Generate a unique ID for a new payment
  String _generateId() {
    return uuid.v7();
  }
}

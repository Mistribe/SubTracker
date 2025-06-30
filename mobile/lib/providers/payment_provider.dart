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

  // Get the count of active payments (excluding stopped payments)
  int get activePaymentsCount {
    return _payments.where((payment) => !payment.isStopped).length;
  }

  // Calculate total amount spent across all payments
  double get totalAmountSpent {
    return _payments.fold(0, (sum, payment) => sum + payment.totalAmountSpent);
  }

  // Add a new payment
  Future<void> addPayment(String name, double price, bool isAnnual, {DateTime? paymentDate}) async {
    final actualPaymentDate = paymentDate ?? DateTime.now();

    // Create initial payment history entry
    final initialHistory = [
      PaymentHistory(
        price: price,
        startDate: actualPaymentDate,
        endDate: DateTime(9999, 12, 31), // Far future date
        isActive: true,
      )
    ];

    final payment = Payment(
      id: _generateId(),
      name: name,
      price: price,
      isAnnual: isAnnual,
      paymentDate: actualPaymentDate,
      paymentHistory: initialHistory,
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
      final now = DateTime.now();

      // Check if price has changed
      if (oldPayment.price != updatedPayment.price) {
        // Create a copy of the price history
        final newPriceHistory = List<PriceChange>.from(updatedPayment.priceHistory);

        // Add the old price to the history with the current date
        // This records when the price was changed
        newPriceHistory.add(PriceChange(
          price: oldPayment.price,
          endDate: now,
        ));

        // Create a copy of the payment history
        final newPaymentHistory = List<PaymentHistory>.from(updatedPayment.paymentHistory);

        // Find the latest active payment history entry
        final latestHistoryIndex = newPaymentHistory.indexWhere((history) => 
          history.isActive && history.endDate.isAfter(now));

        if (latestHistoryIndex >= 0) {
          // Update the end date of the latest history entry
          final latestHistory = newPaymentHistory[latestHistoryIndex];
          newPaymentHistory[latestHistoryIndex] = PaymentHistory(
            price: latestHistory.price,
            startDate: latestHistory.startDate,
            endDate: now,
            isActive: latestHistory.isActive,
          );

          // Add a new history entry with the new price
          newPaymentHistory.add(PaymentHistory(
            price: updatedPayment.price,
            startDate: now,
            endDate: DateTime(9999, 12, 31), // Far future date
            isActive: true,
          ));
        }

        // Create a new payment with updated price history and payment history
        updatedPayment = updatedPayment.copyWith(
          priceHistory: newPriceHistory,
          paymentHistory: newPaymentHistory,
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

      // Create a copy of the payment history
      final newPaymentHistory = List<PaymentHistory>.from(payment.paymentHistory);

      // Find the payment history entry that covers the effective date
      final historyIndex = newPaymentHistory.indexWhere((history) => 
        history.startDate.isBefore(effectiveDate) && history.endDate.isAfter(effectiveDate));

      if (historyIndex >= 0) {
        // Update the end date of the found history entry
        final history = newPaymentHistory[historyIndex];
        newPaymentHistory[historyIndex] = PaymentHistory(
          price: history.price,
          startDate: history.startDate,
          endDate: effectiveDate,
          isActive: history.isActive,
        );

        // Add a new history entry with the new price
        newPaymentHistory.add(PaymentHistory(
          price: newPrice,
          startDate: effectiveDate,
          endDate: DateTime(9999, 12, 31), // Far future date
          isActive: history.isActive, // Maintain the active status
        ));

        // Sort payment history by start date (newest first)
        newPaymentHistory.sort((a, b) => b.startDate.compareTo(a.startDate));
      }

      // Create a new payment with updated price history, payment history, and new current price
      final updatedPayment = payment.copyWith(
        price: newPrice,
        priceHistory: newPriceHistory,
        paymentHistory: newPaymentHistory,
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
        // Get the old price change
        final oldPriceChange = payment.priceHistory[priceChangeIndex];
        final oldDate = oldPriceChange.endDate;

        // Create a copy of the price history
        final newPriceHistory = List<PriceChange>.from(payment.priceHistory);

        // Update the price change at the specified index
        newPriceHistory[priceChangeIndex] = PriceChange(
          price: newPrice,
          endDate: newDate,
        );

        // Sort price history by date (newest first)
        newPriceHistory.sort((a, b) => b.endDate.compareTo(a.endDate));

        // Create a copy of the payment history
        final newPaymentHistory = List<PaymentHistory>.from(payment.paymentHistory);

        // Find payment history entries that need to be updated
        // We need to handle two cases:
        // 1. If the date is moved earlier, we need to update the entry that contains the new date
        // 2. If the date is moved later, we need to update the entry that contains the old date

        if (newDate.isBefore(oldDate)) {
          // Case 1: Date moved earlier
          final historyIndex = newPaymentHistory.indexWhere((history) => 
            history.startDate.isBefore(newDate) && history.endDate.isAfter(newDate));

          if (historyIndex >= 0) {
            // Update the end date of the found history entry
            final history = newPaymentHistory[historyIndex];
            newPaymentHistory[historyIndex] = PaymentHistory(
              price: history.price,
              startDate: history.startDate,
              endDate: newDate,
              isActive: history.isActive,
            );

            // Add a new history entry with the new price
            newPaymentHistory.add(PaymentHistory(
              price: newPrice,
              startDate: newDate,
              endDate: history.endDate,
              isActive: history.isActive,
            ));
          }
        } else if (newDate.isAfter(oldDate)) {
          // Case 2: Date moved later
          final historyIndex = newPaymentHistory.indexWhere((history) => 
            history.startDate.isBefore(oldDate) && history.endDate.isAfter(oldDate));

          if (historyIndex >= 0) {
            // Update the end date of the found history entry
            final history = newPaymentHistory[historyIndex];

            // If there's another history entry between oldDate and newDate, we need to update that one too
            final nextHistoryIndex = newPaymentHistory.indexWhere((h) => 
              h.startDate.isAfter(oldDate) && h.startDate.isBefore(newDate));

            if (nextHistoryIndex >= 0) {
              // Update the next history entry to extend from oldDate to its current end date
              final nextHistory = newPaymentHistory[nextHistoryIndex];
              newPaymentHistory[nextHistoryIndex] = PaymentHistory(
                price: nextHistory.price,
                startDate: oldDate,
                endDate: nextHistory.endDate,
                isActive: nextHistory.isActive,
              );
            } else {
              // No history entry between oldDate and newDate, so extend the current one
              newPaymentHistory[historyIndex] = PaymentHistory(
                price: history.price,
                startDate: history.startDate,
                endDate: newDate,
                isActive: history.isActive,
              );
            }
          }
        }

        // Sort payment history by start date (newest first)
        newPaymentHistory.sort((a, b) => b.startDate.compareTo(a.startDate));

        // Create a new payment with updated price history and payment history
        final updatedPayment = payment.copyWith(
          priceHistory: newPriceHistory,
          paymentHistory: newPaymentHistory,
        );

        // Update local list
        _payments[index] = updatedPayment;

        // Persist to storage
        await paymentRepository.updatePayment(updatedPayment);

        notifyListeners();
      }
    }
  }

  // Stop a payment
  Future<void> stopPayment(String paymentId, {DateTime? stopDate}) async {
    final index = _payments.indexWhere((payment) => payment.id == paymentId);

    if (index >= 0) {
      final payment = _payments[index];

      // If stopDate is not provided, use the last payment date
      final effectiveStopDate = stopDate ?? _calculateLastPaymentDate(payment);

      // Create a copy of the payment history
      final newPaymentHistory = List<PaymentHistory>.from(payment.paymentHistory);

      // Find the latest active payment history entry
      final latestHistoryIndex = newPaymentHistory.indexWhere((history) => 
        history.isActive && history.endDate.isAfter(effectiveStopDate));

      if (latestHistoryIndex >= 0) {
        // Update the end date of the latest history entry
        final latestHistory = newPaymentHistory[latestHistoryIndex];
        newPaymentHistory[latestHistoryIndex] = PaymentHistory(
          price: latestHistory.price,
          startDate: latestHistory.startDate,
          endDate: effectiveStopDate,
          isActive: latestHistory.isActive,
        );

        // Add a new history entry with isActive=false
        newPaymentHistory.add(PaymentHistory(
          price: payment.price,
          startDate: effectiveStopDate,
          endDate: DateTime(9999, 12, 31), // Far future date
          isActive: false, // Inactive because payment is stopped
        ));

        // Sort payment history by start date (newest first)
        newPaymentHistory.sort((a, b) => b.startDate.compareTo(a.startDate));
      }

      // Create a new payment with stopped status and updated payment history
      final updatedPayment = payment.copyWith(
        isStopped: true,
        reactivationDate: null,
        stopDate: effectiveStopDate,
        paymentHistory: newPaymentHistory,
      );

      // Update local list
      _payments[index] = updatedPayment;

      // Persist to storage
      await paymentRepository.updatePayment(updatedPayment);

      notifyListeners();
    }
  }

  // Calculate the last payment date based on payment frequency and start date
  DateTime _calculateLastPaymentDate(Payment payment) {
    final now = DateTime.now();
    final startDate = payment.paymentDate;

    if (payment.isAnnual) {
      // For annual payments, find the most recent payment date before now
      DateTime lastDate = DateTime(startDate.year, startDate.month, startDate.day);
      while (true) {
        DateTime nextDate = DateTime(lastDate.year + 1, lastDate.month, lastDate.day);
        if (nextDate.isAfter(now)) {
          break;
        }
        lastDate = nextDate;
      }
      return lastDate;
    } else {
      // For monthly payments, find the most recent payment date before now
      DateTime lastDate = DateTime(startDate.year, startDate.month, startDate.day);
      while (true) {
        DateTime nextDate = DateTime(lastDate.year, lastDate.month + 1, lastDate.day);
        if (nextDate.isAfter(now)) {
          break;
        }
        lastDate = nextDate;
      }
      return lastDate;
    }
  }

  // Reactivate a payment at a specific date
  Future<void> reactivatePayment(String paymentId, DateTime reactivationDate) async {
    final index = _payments.indexWhere((payment) => payment.id == paymentId);

    if (index >= 0) {
      final payment = _payments[index];

      // Create a copy of the payment history
      final newPaymentHistory = List<PaymentHistory>.from(payment.paymentHistory);

      // Find the latest inactive payment history entry
      final latestInactiveIndex = newPaymentHistory.indexWhere((history) => 
        !history.isActive && history.endDate.isAfter(reactivationDate));

      if (latestInactiveIndex >= 0) {
        // Update the end date of the latest inactive history entry
        final latestInactiveHistory = newPaymentHistory[latestInactiveIndex];
        newPaymentHistory[latestInactiveIndex] = PaymentHistory(
          price: latestInactiveHistory.price,
          startDate: latestInactiveHistory.startDate,
          endDate: reactivationDate,
          isActive: false,
        );

        // Add a new history entry with isActive=true
        newPaymentHistory.add(PaymentHistory(
          price: payment.price,
          startDate: reactivationDate,
          endDate: DateTime(9999, 12, 31), // Far future date
          isActive: true, // Active because payment is reactivated
        ));

        // Sort payment history by start date (newest first)
        newPaymentHistory.sort((a, b) => b.startDate.compareTo(a.startDate));
      }

      // Create a new payment with reactivated status and updated payment history
      final updatedPayment = payment.copyWith(
        isStopped: true, // Still stopped until reactivation date
        reactivationDate: reactivationDate,
        paymentHistory: newPaymentHistory,
      );

      // Update local list
      _payments[index] = updatedPayment;

      // Persist to storage
      await paymentRepository.updatePayment(updatedPayment);

      notifyListeners();
    }
  }

  // Generate a unique ID for a new payment
  String _generateId() {
    return DateTime.now().millisecondsSinceEpoch.toString() + 
           Random().nextInt(10000).toString();
  }
}

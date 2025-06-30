import 'package:hive/hive.dart';
import 'package:hive_flutter/hive_flutter.dart';
import '../models/payment.dart';

/// Repository for handling payment data persistence
/// 
/// This class abstracts the data access layer, making it easier to switch
/// between different storage mechanisms in the future (e.g., local storage to API).
class PaymentRepository {
  static const String _boxName = 'payments';
  late Box<Payment> _box;
  
  /// Initialize the repository
  /// 
  /// This method must be called before using any other methods in this class.
  Future<void> initialize() async {
    // Open the Hive box for payments
    _box = await Hive.openBox<Payment>(_boxName);
  }
  
  /// Get all payments
  List<Payment> getAllPayments() {
    return _box.values.toList();
  }
  
  /// Add a new payment
  Future<void> addPayment(Payment payment) async {
    await _box.put(payment.id, payment);
  }
  
  /// Update an existing payment
  Future<void> updatePayment(Payment payment) async {
    await _box.put(payment.id, payment);
  }
  
  /// Delete a payment
  Future<void> deletePayment(String id) async {
    await _box.delete(id);
  }
  
  /// Clear all payments
  Future<void> clearAll() async {
    await _box.clear();
  }
}
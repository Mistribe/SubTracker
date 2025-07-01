import 'package:hive/hive.dart';
import 'package:hive_flutter/hive_flutter.dart';
import '../models/subscription.dart';

/// Repository for handling subscription data persistence
///
/// This class abstracts the data access layer, making it easier to switch
/// between different storage mechanisms in the future (e.g., local storage to API).
class SubscriptionRepository {
  static const String _boxName = 'subscriptions';
  late Box<Subscription> _box;

  /// Initialize the repository
  ///
  /// This method must be called before using any other methods in this class.
  Future<void> initialize() async {
    // Open the Hive box for subscriptions
    _box = await Hive.openBox<Subscription>(_boxName);
  }

  /// Get all subscriptions
  List<Subscription> getAll() {
    return _box.values.toList();
  }

  /// Add a new subscription
  Future<void> add(Subscription subscription) async {
    await _box.put(subscription.id, subscription);
  }

  /// Update an existing subscription
  Future<void> update(Subscription subscription) async {
    await _box.put(subscription.id, subscription);
  }

  /// Delete a subscription
  Future<void> delete(String id) async {
    await _box.delete(id);
  }

  /// Clear all subscriptions
  Future<void> clearAll() async {
    await _box.clear();
  }
}

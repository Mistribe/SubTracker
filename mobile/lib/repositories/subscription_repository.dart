import 'package:hive/hive.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:subscription_tracker/models/subscription_payment.dart';
import '../models/subscription.dart';
import '../providers/sync_provider.dart';

/// Repository for handling subscription data persistence
///
/// This class abstracts the data access layer, making it easier to switch
/// between different storage mechanisms in the future (e.g., local storage to API).
class SubscriptionRepository {
  static const String _boxName = 'subscriptions';
  late Box<Subscription> _box;
  SyncProvider? _syncProvider;

  /// Initialize the repository
  ///
  /// This method must be called before using any other methods in this class.
  Future<void> initialize() async {
    // Open the Hive box for subscriptions
    _box = await Hive.openBox<Subscription>(_boxName);
  }

  /// Set the sync provider
  void setSyncProvider(SyncProvider syncProvider) {
    _syncProvider = syncProvider;
  }

  /// Get all subscriptions
  List<Subscription> getAll() {
    return _box.values.toList();
  }

  /// Add a new subscription
  Future<void> add(Subscription subscription) async {
    // Save to local storage
    await _box.put(subscription.id, subscription);

    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueCreateSubscription(subscription);
    }
  }

  /// Update an existing subscription
  Future<void> update(Subscription subscription, {bool withSync = true}) async {
    // Save to local storage
    await _box.put(subscription.id, subscription);

    // Queue for sync if provider is available
    if (_syncProvider != null && withSync) {
      await _syncProvider!.queueUpdateSubscription(subscription);
    }
  }

  Future<void> createPayment(String id, SubscriptionPayment payment) async {
    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueCreateSubscriptionPayment(payment, id);
    }
  }

  Future<void> updatePayment(String id, SubscriptionPayment payment) async {
    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueUpdateSubscriptionPayment(payment, id);
    }
  }

  Future<void> deletePayment(String id, String paymentId) async {
    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueDeleteSubscriptionPayment(paymentId, id);
    }
  }

  /// Delete a subscription
  Future<void> delete(String id, {bool withSync = true}) async {
    // Delete from local storage
    await _box.delete(id);

    // Queue for sync if provider is available
    if (_syncProvider != null && withSync) {
      await _syncProvider!.queueDeleteSubscription(id);
    }
  }

  /// Clear all subscriptions
  Future<void> clearAll() async {
    // Get all IDs before clearing
    final ids = _box.keys.cast<String>().toList();

    // Clear local storage
    await _box.clear();

    // Queue deletes for sync if provider is available
    if (_syncProvider != null) {
      for (final id in ids) {
        await _syncProvider!.queueDeleteSubscription(id);
      }
    }
  }
}

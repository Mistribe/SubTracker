import 'package:flutter/material.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/subscription.dart';
import '../services/api_service.dart';
import '../services/sync_service.dart';
import '../repositories/subscription_repository.dart';

/// Provider for the sync service
class SyncProvider extends ChangeNotifier {
  late final ApiService _apiService;
  late final SyncService _syncService;
  bool _isInitialized = false;
  bool _isSyncing = false;
  DateTime? _lastSyncTime;
  bool _hasPendingOperations = false;

  SyncProvider({required SubscriptionRepository subscriptionRepository}) {
    _initialize(subscriptionRepository);
  }

  /// Initialize the provider
  Future<void> _initialize(
    SubscriptionRepository subscriptionRepository,
  ) async {
    if (_isInitialized) return;

    // Initialize API service
    _apiService = ApiService(baseUrl: 'http://localhost:5042/api');

    // Initialize shared preferences
    final prefs = await SharedPreferences.getInstance();

    // Initialize sync service
    _syncService = SyncService(
      apiService: _apiService,
      subscriptionRepository: subscriptionRepository,
      connectivity: Connectivity(),
      prefs: prefs,
    );

    // Initialize the sync service
    await _syncService.initialize();

    // Update state
    _lastSyncTime = _syncService.getLastSyncTime();
    _hasPendingOperations = _syncService.hasPendingOperations();
    _isInitialized = true;

    notifyListeners();
  }

  /// Get whether the provider is initialized
  bool get isInitialized => _isInitialized;

  /// Get whether a sync is in progress
  bool get isSyncing => _isSyncing;

  /// Get the last sync time
  DateTime? get lastSyncTime => _lastSyncTime;

  /// Get whether there are pending operations
  bool get hasPendingOperations => _hasPendingOperations;

  /// Trigger a manual sync
  Future<void> sync() async {
    if (_isSyncing || !_isInitialized) return;

    _isSyncing = true;
    notifyListeners();

    try {
      await _syncService.sync();
      _lastSyncTime = _syncService.getLastSyncTime();
      _hasPendingOperations = _syncService.hasPendingOperations();
    } finally {
      _isSyncing = false;
      notifyListeners();
    }
  }

  /// Queue a create operation
  Future<void> queueCreate(Subscription subscription) async {
    if (!_isInitialized) return;

    await _syncService.queueCreate(subscription);
    _hasPendingOperations = _syncService.hasPendingOperations();
    notifyListeners();
  }

  /// Queue an update operation
  Future<void> queueUpdate(Subscription subscription) async {
    if (!_isInitialized) return;

    await _syncService.queueUpdate(subscription);
    _hasPendingOperations = _syncService.hasPendingOperations();
    notifyListeners();
  }

  /// Queue a delete operation
  Future<void> queueDelete(String id) async {
    if (!_isInitialized) return;

    await _syncService.queueDelete(id);
    _hasPendingOperations = _syncService.hasPendingOperations();
    notifyListeners();
  }

  @override
  void dispose() {
    _apiService.dispose();
    _syncService.dispose();
    super.dispose();
  }
}

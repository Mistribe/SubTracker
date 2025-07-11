import 'package:flutter/material.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/subscription.dart';
import '../models/label.dart';
import '../models/family_member.dart';
import '../services/api_service.dart';
import '../services/sync_service.dart';
import '../repositories/subscription_repository.dart';
import '../repositories/label_repository.dart';
import '../repositories/family_member_repository.dart';

/// Provider for the sync service
class SyncProvider extends ChangeNotifier {
  late final ApiService _apiService;
  late final SyncService _syncService;
  bool _isInitialized = false;
  bool _isSyncing = false;
  DateTime? _lastSyncTime;
  bool _hasPendingOperations = false;
  bool _hasSyncHistory = false;

  SyncProvider({
    required SubscriptionRepository subscriptionRepository,
    required LabelRepository labelRepository,
    required FamilyMemberRepository familyMemberRepository,
  }) {
    _initialize(
      subscriptionRepository,
      labelRepository,
      familyMemberRepository,
    );
  }

  /// Initialize the provider
  Future<void> _initialize(
    SubscriptionRepository subscriptionRepository,
    LabelRepository labelRepository,
    FamilyMemberRepository familyMemberRepository,
  ) async {
    if (_isInitialized) return;

    // Initialize API service
    _apiService = ApiService(baseUrl: 'http://localhost:5042/api');

    // Initialize shared preferences
    final prefs = SharedPreferencesAsync();

    // Initialize sync service
    _syncService = SyncService(
      apiService: _apiService,
      subscriptionRepository: subscriptionRepository,
      labelRepository: labelRepository,
      familyMemberRepository: familyMemberRepository,
      connectivity: Connectivity(),
      prefs: prefs,
    );

    // Initialize the sync service
    await _syncService.initialize();

    // Update state
    _lastSyncTime = await _syncService.getLastSyncTime();
    _hasPendingOperations = await _syncService.hasPendingOperations();
    _hasSyncHistory = await _syncService.hasSyncHistory();
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

  /// Get whether there are operations in the sync history
  bool get hasSyncHistory => _hasSyncHistory;

  /// Trigger a manual sync
  Future<void> sync() async {
    if (_isSyncing || !_isInitialized) return;

    _isSyncing = true;
    notifyListeners();

    try {
      await _syncService.sync();
      _lastSyncTime = await _syncService.getLastSyncTime();
      _hasPendingOperations = await _syncService.hasPendingOperations();
      _hasSyncHistory = await _syncService.hasSyncHistory();
    } finally {
      _isSyncing = false;
      notifyListeners();
    }
  }

  /// Queue a create operation
  Future<void> queueCreateSubscription(Subscription subscription) async {
    if (!_isInitialized) return;

    await _syncService.queueCreate(subscription);
    _hasPendingOperations = await _syncService.hasPendingOperations();
    _hasSyncHistory = await _syncService.hasSyncHistory();
    notifyListeners();
  }

  Future<void> queueCreateLabel(Label label) async {
    if (!_isInitialized) return;

    await _syncService.queueCreate(label);
    _hasPendingOperations = await _syncService.hasPendingOperations();
    _hasSyncHistory = await _syncService.hasSyncHistory();
    notifyListeners();
  }

  Future<void> queueCreateFamilyMember(FamilyMember familyMember) async {
    if (!_isInitialized) return;

    await _syncService.queueCreate(familyMember);
    _hasPendingOperations = await _syncService.hasPendingOperations();
    _hasSyncHistory = await _syncService.hasSyncHistory();
    notifyListeners();
  }

  /// Queue an update operation
  Future<void> queueUpdateSubscription(Subscription subscription) async {
    if (!_isInitialized) return;

    await _syncService.queueUpdate(subscription);
    _hasPendingOperations = await _syncService.hasPendingOperations();
    _hasSyncHistory = await _syncService.hasSyncHistory();
    notifyListeners();
  }

  Future<void> queueUpdateLabel(Label label) async {
    if (!_isInitialized) return;

    await _syncService.queueUpdate(label);
    _hasPendingOperations = await _syncService.hasPendingOperations();
    _hasSyncHistory = await _syncService.hasSyncHistory();
    notifyListeners();
  }

  Future<void> queueUpdateFamilyMember(FamilyMember familyMember) async {
    if (!_isInitialized) return;

    await _syncService.queueUpdate(familyMember);
    _hasPendingOperations = await _syncService.hasPendingOperations();
    _hasSyncHistory = await _syncService.hasSyncHistory();
    notifyListeners();
  }

  Future<void> _queueDelete(String id, SyncDataType dataType) async {
    if (!_isInitialized) return;

    await _syncService.queueDelete(id, dataType);
    _hasPendingOperations = await _syncService.hasPendingOperations();
    _hasSyncHistory = await _syncService.hasSyncHistory();
    notifyListeners();
  }

  /// Queue a delete operation
  Future<void> queueDeleteSubscription(String id) async {
    await _queueDelete(id, SyncDataType.subscription);
  }

  Future<void> queueDeleteLabel(String id) async {
    await _queueDelete(id, SyncDataType.label);
  }

  Future<void> queueDeleteFamilyMember(String id) async {
    await _queueDelete(id, SyncDataType.familyMember);
  }

  @override
  void dispose() {
    _apiService.dispose();
    _syncService.dispose();
    super.dispose();
  }
}

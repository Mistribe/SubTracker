import 'dart:async';
import 'dart:convert';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:subscription_tracker/models/family.dart';
import '../models/subscription.dart';
import '../models/label.dart';
import '../models/family_member.dart';
import '../models/subscription_payment.dart';
import '../repositories/subscription_repository.dart';
import '../repositories/label_repository.dart';
import '../repositories/family_repository.dart';
import 'api_service.dart';

/// Type of sync operation
enum SyncOperationType { create, update, delete }

enum SyncDataType {
  label,
  subscription,
  familyMember,
  subscriptionPayment,
  family,
}

/// Pending sync operation
class PendingSyncOperation {
  final String id;
  final SyncOperationType type;
  final SyncDataType dataType;
  final Map<String, dynamic>? data;
  final DateTime timestamp;

  PendingSyncOperation({
    required this.id,
    required this.type,
    required this.dataType,
    this.data,
    DateTime? timestamp,
  }) : timestamp = timestamp ?? DateTime.now();

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'type': type.toString().split('.').last,
      'data_type': dataType.toString(),
      'data': data,
      'timestamp': timestamp.toIso8601String(),
    };
  }

  factory PendingSyncOperation.fromJson(Map<String, dynamic> json) {
    return PendingSyncOperation(
      id: json['id'],
      type: SyncOperationType.values.firstWhere(
        (e) => e.toString().split('.').last == json['type'],
      ),
      dataType: SyncDataType.values.firstWhere(
        (e) => e.toString() == json['data_type'],
      ),
      data: json['data'],
      timestamp: DateTime.parse(json['timestamp']),
    );
  }
}

/// Service for synchronizing data between local storage and backend
class SyncService {
  final ApiService _apiService;
  final SubscriptionRepository _subscriptionRepository;
  final LabelRepository _labelRepository;
  final FamilyRepository _familyRepository;
  final Connectivity _connectivity;
  final SharedPreferencesAsync _prefs;

  bool _isInitialized = false;
  bool _isSyncing = false;
  bool _isOnline = false;

  StreamSubscription<List<ConnectivityResult>>? _connectivitySubscription;
  Timer? _syncTimer;

  static const String _pendingOperationsKey = 'pending_sync_operations';
  static const String _syncHistoryKey = 'sync_history_operations';
  static const String _lastSyncTimeKey = 'last_sync_time';
  static const Duration _syncInterval = Duration(minutes: 15);

  SyncService({
    required ApiService apiService,
    required SubscriptionRepository subscriptionRepository,
    required LabelRepository labelRepository,
    required FamilyRepository familyMemberRepository,
    required Connectivity connectivity,
    required SharedPreferencesAsync prefs,
  }) : _apiService = apiService,
       _subscriptionRepository = subscriptionRepository,
       _labelRepository = labelRepository,
       _familyRepository = familyMemberRepository,
       _connectivity = connectivity,
       _prefs = prefs;

  /// Initialize the sync service
  Future<void> initialize() async {
    if (_isInitialized) return;

    // Check initial connectivity
    final connectivityResult = await _connectivity.checkConnectivity();
    _isOnline = connectivityResult != ConnectivityResult.none;

    // Listen for connectivity changes
    _connectivitySubscription = _connectivity.onConnectivityChanged.listen((
      List<ConnectivityResult> results,
    ) {
      final wasOnline = _isOnline;
      final result = results.isNotEmpty
          ? results.first
          : ConnectivityResult.none;
      _isOnline = result != ConnectivityResult.none;

      // If we just came online, trigger a sync
      if (!wasOnline && _isOnline) {
        sync();
      }
    });

    // Set up periodic sync
    _syncTimer = Timer.periodic(_syncInterval, (_) {
      if (_isOnline) {
        sync();
      }
    });

    _isInitialized = true;

    // Check if we need to replay operations from history
    await _replayOperationsIfNeeded();

    // Perform initial sync if online
    if (_isOnline) {
      sync();
    }
  }

  /// Replay operations from history if needed
  Future<void> _replayOperationsIfNeeded() async {
    // Get pending operations
    final pendingOperations = await _getPendingOperations();

    // If there are already pending operations, no need to replay
    if (pendingOperations.isNotEmpty) {
      return;
    }

    // Get sync history
    final syncHistory = await _getSyncHistory();

    // If there's no history, nothing to replay
    if (syncHistory.isEmpty) {
      return;
    }

    // Add all operations from history to pending operations
    // This will allow them to be processed on the next sync
    await _savePendingOperations(syncHistory);

    print('Replayed ${syncHistory.length} operations from sync history');
  }

  /// Get pending operations from shared preferences
  Future<List<PendingSyncOperation>> _getPendingOperations() async {
    final pendingOperationsJson =
        await _prefs.getStringList(_pendingOperationsKey) ?? [];
    return pendingOperationsJson
        .map((json) => PendingSyncOperation.fromJson(jsonDecode(json)))
        .toList();
  }

  /// Save pending operations to shared preferences
  Future<void> _savePendingOperations(
    List<PendingSyncOperation> operations,
  ) async {
    final operationsJson = operations
        .map((op) => jsonEncode(op.toJson()))
        .toList();
    await _prefs.setStringList(_pendingOperationsKey, operationsJson);
  }

  /// Get sync history operations from shared preferences
  Future<List<PendingSyncOperation>> _getSyncHistory() async {
    final syncHistoryJson = await _prefs.getStringList(_syncHistoryKey) ?? [];
    return syncHistoryJson
        .map((json) => PendingSyncOperation.fromJson(jsonDecode(json)))
        .toList();
  }

  /// Save sync history operations to shared preferences
  Future<void> _saveSyncHistory(List<PendingSyncOperation> operations) async {
    final operationsJson = operations
        .map((op) => jsonEncode(op.toJson()))
        .toList();
    await _prefs.setStringList(_syncHistoryKey, operationsJson);
  }

  /// Add operations to sync history
  Future<void> _addToSyncHistory(List<PendingSyncOperation> operations) async {
    if (operations.isEmpty) return;

    final history = await _getSyncHistory();

    // Add all operations to history
    history.addAll(operations);

    // Save the updated history
    await _saveSyncHistory(history);
  }

  /// Clear the sync history
  Future<void> _clearSyncHistory() async {
    await _prefs.setStringList(_syncHistoryKey, []);
  }

  /// Add a pending operation
  Future<void> _addPendingOperation(PendingSyncOperation operation) async {
    final operations = await _getPendingOperations();

    // Remove any existing operations for the same ID
    operations.removeWhere(
      (op) => op.id == operation.id && op.type == operation.type,
    );

    // Add the new operation
    operations.add(operation);

    // Save the updated list
    await _savePendingOperations(operations);
  }

  SyncDataType getDataTypeFromEntity(dynamic entity) {
    if (entity is Subscription) {
      return SyncDataType.subscription;
    } else if (entity is Label) {
      return SyncDataType.label;
    } else if (entity is SubscriptionPayment) {
      return SyncDataType.subscriptionPayment;
    } else if (entity is Family) {
      return SyncDataType.family;
    } else {
      return SyncDataType.familyMember;
    }
  }

  /// Queue a create operation for a subscription
  Future<void> queueCreate(dynamic entity, {String? subscriptionId}) async {
    if (entity is Subscription ||
        entity is Label ||
        entity is Family ||
        entity is FamilyMember ||
        entity is SubscriptionPayment) {
      final data = entity.toJson();

      // For SubscriptionPayment, we need to store the subscription ID
      if (entity is SubscriptionPayment) {
        if (subscriptionId == null) {
          throw ArgumentError(
            'Subscription ID is required for SubscriptionPayment operations',
          );
        }
        data['subscriptionId'] = subscriptionId;
      }

      await _addPendingOperation(
        PendingSyncOperation(
          id: entity.id,
          type: SyncOperationType.create,
          dataType: getDataTypeFromEntity(entity),
          data: data,
        ),
      );

      // Try to sync immediately if online
      if (_isOnline) {
        sync();
      }
    } else {
      throw ArgumentError(
        'Entity must be a Subscription, Label, FamilyMember, or SubscriptionPayment',
      );
    }
  }

  /// Queue an update operation for a subscription
  Future<void> queueUpdate(dynamic entity, {String? subscriptionId}) async {
    if (entity is Subscription ||
        entity is Label ||
        entity is FamilyMember ||
        entity is SubscriptionPayment) {
      final data = entity.toJson();

      // For SubscriptionPayment, we need to store the subscription ID
      if (entity is SubscriptionPayment) {
        if (subscriptionId == null) {
          throw ArgumentError(
            'Subscription ID is required for SubscriptionPayment operations',
          );
        }
        data['subscriptionId'] = subscriptionId;
      }

      await _addPendingOperation(
        PendingSyncOperation(
          id: entity.id,
          type: SyncOperationType.update,
          dataType: getDataTypeFromEntity(entity),
          data: data,
        ),
      );

      // Try to sync immediately if online
      if (_isOnline) {
        sync();
      }
    } else {
      throw ArgumentError(
        'Entity must be a Subscription, Label, FamilyMember, or SubscriptionPayment',
      );
    }
  }

  /// Queue a delete operation
  Future<void> queueDelete(
    String id,
    SyncDataType dataType, {
    String? subscriptionId,
  }) async {
    // For SubscriptionPayment, we need the subscription ID
    Map<String, dynamic>? data;
    if (dataType == SyncDataType.subscriptionPayment) {
      if (subscriptionId == null) {
        throw ArgumentError(
          'Subscription ID is required for SubscriptionPayment delete operations',
        );
      }
      data = {'subscriptionId': subscriptionId};
    }

    await _addPendingOperation(
      PendingSyncOperation(
        id: id,
        type: SyncOperationType.delete,
        dataType: dataType,
        data: data,
      ),
    );

    // Try to sync immediately if online
    if (_isOnline) {
      sync();
    }
  }

  /// Process pending operations and return successful operations
  Future<List<PendingSyncOperation>> _processPendingOperations() async {
    final pendingOperations = await _getPendingOperations();
    print('Processing ${pendingOperations.length} pending operations');

    final successfulOperations = <PendingSyncOperation>[];

    for (final operation in pendingOperations) {
      try {
        switch (operation.type) {
          case SyncOperationType.create:
            if (operation.data != null) {
              switch (operation.dataType) {
                case SyncDataType.subscription:
                  final subscription = Subscription.fromJson(operation.data!);
                  await _apiService.createSubscription(subscription);
                  successfulOperations.add(operation);
                  break;
                case SyncDataType.label:
                  final label = Label.fromJson(operation.data!);
                  await _apiService.createLabel(label);
                  successfulOperations.add(operation);
                  break;
                case SyncDataType.family:
                  final family = Family.fromJson(operation.data!);
                  await _apiService.createFamily(family);
                  successfulOperations.add(operation);
                  break;
                case SyncDataType.familyMember:
                  final familyMember = FamilyMember.fromJson(operation.data!);
                  await _apiService.createFamilyMember(familyMember);
                  successfulOperations.add(operation);
                  break;
                case SyncDataType.subscriptionPayment:
                  final payment = SubscriptionPayment.fromJson(operation.data!);
                  // We need the subscription ID for this operation
                  final subscriptionId = operation.data!['subscriptionId'];
                  if (subscriptionId == null) {
                    throw Exception(
                      'Subscription ID is required for subscription payment operations',
                    );
                  }
                  await _apiService.createSubscriptionPayment(
                    subscriptionId,
                    payment,
                  );
                  successfulOperations.add(operation);
                  break;
              }
            }
            break;
          case SyncOperationType.update:
            if (operation.data != null) {
              // Determine the entity type based on the data
              switch (operation.dataType) {
                case SyncDataType.subscription:
                  // This is a Subscription
                  final subscription = Subscription.fromJson(operation.data!);
                  await _apiService.updateSubscription(subscription);
                  successfulOperations.add(operation);
                  break;
                case SyncDataType.label:
                  // This is a Label
                  final label = Label.fromJson(operation.data!);
                  await _apiService.updateLabel(label);
                  successfulOperations.add(operation);
                  break;
                case SyncDataType.family:
                  final family = Family.fromJson(operation.data!);
                  await _apiService.updateFamily(family);
                  successfulOperations.add(operation);
                  break;
                case SyncDataType.familyMember:
                  // This is a FamilyMember
                  final familyMember = FamilyMember.fromJson(operation.data!);
                  await _apiService.updateFamilyMember(familyMember);
                  successfulOperations.add(operation);
                  break;
                case SyncDataType.subscriptionPayment:
                  // This is a SubscriptionPayment
                  final payment = SubscriptionPayment.fromJson(operation.data!);
                  // We need the subscription ID for this operation
                  final subscriptionId = operation.data!['subscriptionId'];
                  if (subscriptionId == null) {
                    throw Exception(
                      'Subscription ID is required for subscription payment operations',
                    );
                  }
                  await _apiService.updateSubscriptionPayment(
                    subscriptionId,
                    payment,
                  );
                  successfulOperations.add(operation);
                  break;
              }
            }
            break;
          case SyncOperationType.delete:
            // We need to determine the entity type for delete operations
            // This could be done by adding a 'type' field to the operation
            // or by checking if the ID exists in each repository
            switch (operation.dataType) {
              case SyncDataType.subscription:
                await _apiService.deleteSubscription(operation.id);
                successfulOperations.add(operation);
                break;
              case SyncDataType.label:
                await _apiService.deleteLabel(operation.id);
                successfulOperations.add(operation);
                break;
              case SyncDataType.family:
                await _apiService.deleteFamily(operation.id);
                successfulOperations.add(operation);
                break;
              case SyncDataType.familyMember:
                await _apiService.deleteFamilyMember(operation.id);
                successfulOperations.add(operation);
                break;
              case SyncDataType.subscriptionPayment:
                // For subscription payment, we need both the subscription ID and the payment ID
                // The subscription ID should be stored in the metadata
                final subscriptionId = operation.data?['subscriptionId'];
                if (subscriptionId == null) {
                  throw Exception(
                    'Subscription ID is required for subscription payment operations',
                  );
                }
                await _apiService.deleteSubscriptionPayment(
                  subscriptionId,
                  operation.id,
                );
                successfulOperations.add(operation);
                break;
            }
        }
      } catch (e) {
        print(
          'Error processing operation ${operation.type} for ${operation.id}: $e',
        );
        // We'll retry this operation on the next sync
      }
    }

    // Add successful operations to sync history
    if (successfulOperations.isNotEmpty) {
      await _addToSyncHistory(successfulOperations);

      // Remove successful operations from the pending list
      final remainingOperations = pendingOperations
          .where((op) => !successfulOperations.contains(op))
          .toList();
      await _savePendingOperations(remainingOperations);
    }

    return successfulOperations;
  }

  /// Fetch latest data from backend
  Future<Map<String, dynamic>> _fetchRemoteData() async {
    final remoteSubscriptions = await _apiService.getSubscriptions();
    final remoteLabels = await _apiService.getLabels(withDefault: true);
    final remoteFamilies = await _apiService.getFamilies();

    return {
      'subscriptions': remoteSubscriptions,
      'labels': remoteLabels,
      'families': remoteFamilies,
    };
  }

  /// Update local storage with remote data
  Future<void> _updateLocalStorage(Map<String, dynamic> remoteData) async {
    final remoteSubscriptions =
        remoteData['subscriptions'] as List<Subscription>;
    final remoteLabels = remoteData['labels'] as List<Label>;
    final remoteFamilies = remoteData['families'] as List<Family>;

    final localSubscriptions = _subscriptionRepository.getAll();
    final localLabels = _labelRepository.getAll();
    final localFamilies = _familyRepository.getAll();

    // Update local storage with remote subscriptions
    for (final remoteSubscription in remoteSubscriptions) {
      final localSubscription = localSubscriptions.firstWhere(
        (s) => s.id == remoteSubscription.id,
        orElse: () => Subscription.empty(),
      );

      // If the remote subscription is newer, update local
      if (localSubscription.updatedAt.isBefore(remoteSubscription.updatedAt)) {
        await _subscriptionRepository.update(
          remoteSubscription,
          withSync: false,
        );
      }
    }

    // Remove local labels that don't exist in remote data
    for (final localSubscription in localSubscriptions) {
      if (!remoteSubscriptions.any((r) => r.id == localSubscription.id)) {
        await _subscriptionRepository.delete(
          localSubscription.id,
          withSync: false,
        );
      }
    }

    // Update local storage with remote labels
    for (final remoteLabel in remoteLabels) {
      final localLabel = localLabels.firstWhere(
        (l) => l.id == remoteLabel.id,
        orElse: () => Label.empty(),
      );

      // If the remote label is newer, update local
      if (localLabel.updatedAt.isBefore(remoteLabel.updatedAt)) {
        await _labelRepository.update(remoteLabel, withSync: false);
      }
    }

    // Remove local labels that don't exist in remote data
    for (final localLabel in localLabels) {
      if (!remoteLabels.any((r) => r.id == localLabel.id)) {
        await _labelRepository.delete(localLabel.id, withSync: false);
      }
    }

    // Update local storage with remote family members
    for (final remoteFamily in remoteFamilies) {
      final localMember = localFamilies.firstWhere(
        (m) => m.id == remoteFamily.id,
        orElse: () => Family.empty(),
      );

      // If the remote member is newer, update local
      if (localMember.updatedAt.isBefore(remoteFamily.updatedAt)) {
        await _familyRepository.update(
          remoteFamily.id,
          remoteFamily.name,
          haveJointAccount: remoteFamily.haveJointAccount,
          members: remoteFamily.members,
          updatedAt: remoteFamily.updatedAt,
          withSync: false,
        );
      }
    }

    // Remove local labels that don't exist in remote data
    for (final localMember in localFamilies) {
      if (!remoteFamilies.any((r) => r.id == localMember.id)) {
        await _familyRepository.delete(localMember.id, withSync: false);
      }
    }
  }

  /// Update last sync time and clean up
  Future<void> _updateSyncTimeAndCleanup() async {
    // Update last sync time
    await _prefs.setString(_lastSyncTimeKey, DateTime.now().toIso8601String());

    // If there are no pending operations left, we can clear the sync history
    // as all operations have been successfully processed
    if ((await _getPendingOperations()).isEmpty) {
      await _clearSyncHistory();
    }
  }

  /// Synchronize data with the backend
  Future<void> sync() async {
    // Prevent multiple syncs from running simultaneously
    if (_isSyncing || !_isOnline) return;

    _isSyncing = true;

    try {
      // Process pending operations
      await _processPendingOperations();

      // Fetch latest data from backend
      final remoteData = await _fetchRemoteData();

      // Update local storage with remote data
      await _updateLocalStorage(remoteData);

      // Update last sync time and clean up
      await _updateSyncTimeAndCleanup();
    } catch (e) {
      print('Sync error: $e');
    } finally {
      _isSyncing = false;
    }
  }

  /// Get the last sync time
  Future<DateTime?> getLastSyncTime() async {
    final lastSyncTimeStr = await _prefs.getString(_lastSyncTimeKey);
    if (lastSyncTimeStr == null) return null;

    try {
      return DateTime.parse(lastSyncTimeStr);
    } catch (_) {
      return null;
    }
  }

  /// Check if there are pending operations
  Future<bool> hasPendingOperations() async {
    return (await _getPendingOperations()).isNotEmpty;
  }

  /// Check if there are operations in the sync history
  Future<bool> hasSyncHistory() async {
    return (await _getSyncHistory()).isNotEmpty;
  }

  /// Get the count of pending operations
  Future<int> getPendingOperationsCount() async {
    return (await _getPendingOperations()).length;
  }

  /// Dispose resources
  void dispose() {
    _connectivitySubscription?.cancel();
    _syncTimer?.cancel();
  }
}

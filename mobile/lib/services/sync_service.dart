import 'dart:async';
import 'dart:convert';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/subscription.dart';
import '../models/label.dart';
import '../models/family_member.dart';
import '../repositories/subscription_repository.dart';
import '../repositories/label_repository.dart';
import '../repositories/family_member_repository.dart';
import 'api_service.dart';

/// Type of sync operation
enum SyncOperationType { create, update, delete }

enum SyncDataType { label, subscription, familyMember }

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
  final FamilyMemberRepository _familyMemberRepository;
  final Connectivity _connectivity;
  final SharedPreferencesAsync _prefs;

  bool _isInitialized = false;
  bool _isSyncing = false;
  bool _isOnline = false;

  StreamSubscription<ConnectivityResult>? _connectivitySubscription;
  Timer? _syncTimer;

  static const String _pendingOperationsKey = 'pending_sync_operations';
  static const String _syncHistoryKey = 'sync_history_operations';
  static const String _lastSyncTimeKey = 'last_sync_time';
  static const Duration _syncInterval = Duration(minutes: 15);

  SyncService({
    required ApiService apiService,
    required SubscriptionRepository subscriptionRepository,
    required LabelRepository labelRepository,
    required FamilyMemberRepository familyMemberRepository,
    required Connectivity connectivity,
    required SharedPreferencesAsync prefs,
  }) : _apiService = apiService,
       _subscriptionRepository = subscriptionRepository,
       _labelRepository = labelRepository,
       _familyMemberRepository = familyMemberRepository,
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
      result,
    ) {
      final wasOnline = _isOnline;
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
    } else {
      return SyncDataType.familyMember;
    }
  }

  /// Queue a create operation for a subscription
  Future<void> queueCreate(dynamic entity) async {
    if (entity is Subscription || entity is Label || entity is FamilyMember) {
      await _addPendingOperation(
        PendingSyncOperation(
          id: entity.id,
          type: SyncOperationType.create,
          dataType: getDataTypeFromEntity(entity),
          data: entity.toJson(),
        ),
      );

      // Try to sync immediately if online
      if (_isOnline) {
        sync();
      }
    } else {
      throw ArgumentError(
        'Entity must be a Subscription, Label, or FamilyMember',
      );
    }
  }

  /// Queue an update operation for a subscription
  Future<void> queueUpdate(dynamic entity) async {
    if (entity is Subscription || entity is Label || entity is FamilyMember) {
      await _addPendingOperation(
        PendingSyncOperation(
          id: entity.id,
          type: SyncOperationType.update,
          dataType: getDataTypeFromEntity(entity),
          data: entity.toJson(),
        ),
      );

      // Try to sync immediately if online
      if (_isOnline) {
        sync();
      }
    } else {
      throw ArgumentError(
        'Entity must be a Subscription, Label, or FamilyMember',
      );
    }
  }

  /// Queue a delete operation
  Future<void> queueDelete(String id, SyncDataType dataType) async {
    await _addPendingOperation(
      PendingSyncOperation(
        id: id,
        type: SyncOperationType.delete,
        dataType: dataType,
      ),
    );

    // Try to sync immediately if online
    if (_isOnline) {
      sync();
    }
  }

  /// Synchronize data with the backend
  Future<void> sync() async {
    // Prevent multiple syncs from running simultaneously
    if (_isSyncing || !_isOnline) return;

    _isSyncing = true;

    try {
      // Process pending operations
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
                  case SyncDataType.familyMember:
                    final familyMember = FamilyMember.fromJson(operation.data!);
                    await _apiService.createFamilyMember(familyMember);
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
                  case SyncDataType.familyMember:
                    // This is a FamilyMember
                    final familyMember = FamilyMember.fromJson(operation.data!);
                    await _apiService.updateFamilyMember(familyMember);
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
                case SyncDataType.familyMember:
                  await _apiService.deleteFamilyMember(operation.id);
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

      // Fetch latest data from backend
      final remoteSubscriptions = await _apiService.getSubscriptions();
      final remoteLabels = await _apiService.getLabels();
      final remoteFamilyMembers = await _apiService.getFamilyMembers();

      final localSubscriptions = _subscriptionRepository.getAll();
      final localLabels = _labelRepository.getAll();
      final localFamilyMembers = _familyMemberRepository.getAll();

      // Update local storage with remote data
      for (final remoteSubscription in remoteSubscriptions) {
        final localSubscription = localSubscriptions.firstWhere(
          (s) => s.id == remoteSubscription.id,
          orElse: () => remoteSubscription,
        );

        // If the remote subscription is newer, update local
        if (localSubscription != remoteSubscription) {
          await _subscriptionRepository.update(remoteSubscription);
        }
      }

      // Update local storage with remote labels
      for (final remoteLabel in remoteLabels) {
        final localLabel = localLabels.firstWhere(
          (l) => l.id == remoteLabel.id,
          orElse: () => remoteLabel,
        );

        // If the remote label is newer, update local
        if (localLabel != remoteLabel) {
          await _labelRepository.update(remoteLabel);
        }
      }

      // Update local storage with remote family members
      for (final remoteMember in remoteFamilyMembers) {
        final localMember = localFamilyMembers.firstWhere(
          (m) => m.id == remoteMember.id,
          orElse: () => remoteMember,
        );

        // If the remote member is newer, update local
        if (localMember != remoteMember) {
          await _familyMemberRepository.update(remoteMember);
        }
      }

      // Update last sync time
      for (final remoteSubscription in remoteSubscriptions) {
        final localSubscription = localSubscriptions.firstWhere(
          (s) => s.id == remoteSubscription.id,
          orElse: () => remoteSubscription,
        );

        // If the remote subscription is newer, update local
        if (localSubscription != remoteSubscription) {
          await _subscriptionRepository.update(remoteSubscription);
        }
      }

      // Update last sync time
      await _prefs.setString(
        _lastSyncTimeKey,
        DateTime.now().toIso8601String(),
      );

      // If there are no pending operations left, we can clear the sync history
      // as all operations have been successfully processed
      if ((await _getPendingOperations()).isEmpty) {
        await _clearSyncHistory();
      }
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

  /// Dispose resources
  void dispose() {
    _connectivitySubscription?.cancel();
    _syncTimer?.cancel();
  }
}

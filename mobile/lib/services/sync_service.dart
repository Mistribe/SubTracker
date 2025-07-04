import 'dart:async';
import 'dart:convert';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/subscription.dart';
import '../repositories/subscription_repository.dart';
import 'api_service.dart';

/// Type of sync operation
enum SyncOperationType {
  create,
  update,
  delete,
}

/// Pending sync operation
class PendingSyncOperation {
  final String id;
  final SyncOperationType type;
  final Map<String, dynamic>? data;
  final DateTime timestamp;

  PendingSyncOperation({
    required this.id,
    required this.type,
    this.data,
    DateTime? timestamp,
  }) : timestamp = timestamp ?? DateTime.now();

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'type': type.toString().split('.').last,
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
      data: json['data'],
      timestamp: DateTime.parse(json['timestamp']),
    );
  }
}

/// Service for synchronizing data between local storage and backend
class SyncService {
  final ApiService _apiService;
  final SubscriptionRepository _subscriptionRepository;
  final Connectivity _connectivity;
  final SharedPreferences _prefs;
  
  bool _isInitialized = false;
  bool _isSyncing = false;
  bool _isOnline = false;
  
  StreamSubscription<ConnectivityResult>? _connectivitySubscription;
  Timer? _syncTimer;
  
  static const String _pendingOperationsKey = 'pending_sync_operations';
  static const String _lastSyncTimeKey = 'last_sync_time';
  static const Duration _syncInterval = Duration(minutes: 15);

  SyncService({
    required ApiService apiService,
    required SubscriptionRepository subscriptionRepository,
    required Connectivity connectivity,
    required SharedPreferences prefs,
  })  : _apiService = apiService,
        _subscriptionRepository = subscriptionRepository,
        _connectivity = connectivity,
        _prefs = prefs;

  /// Initialize the sync service
  Future<void> initialize() async {
    if (_isInitialized) return;
    
    // Check initial connectivity
    final connectivityResult = await _connectivity.checkConnectivity();
    _isOnline = connectivityResult != ConnectivityResult.none;
    
    // Listen for connectivity changes
    _connectivitySubscription = _connectivity.onConnectivityChanged.listen((result) {
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
    
    // Perform initial sync if online
    if (_isOnline) {
      sync();
    }
  }

  /// Get pending operations from shared preferences
  List<PendingSyncOperation> _getPendingOperations() {
    final pendingOperationsJson = _prefs.getStringList(_pendingOperationsKey) ?? [];
    return pendingOperationsJson
        .map((json) => PendingSyncOperation.fromJson(jsonDecode(json)))
        .toList();
  }

  /// Save pending operations to shared preferences
  Future<void> _savePendingOperations(List<PendingSyncOperation> operations) async {
    final operationsJson = operations
        .map((op) => jsonEncode(op.toJson()))
        .toList();
    await _prefs.setStringList(_pendingOperationsKey, operationsJson);
  }

  /// Add a pending operation
  Future<void> _addPendingOperation(PendingSyncOperation operation) async {
    final operations = _getPendingOperations();
    
    // Remove any existing operations for the same ID
    operations.removeWhere((op) => op.id == operation.id && op.type == operation.type);
    
    // Add the new operation
    operations.add(operation);
    
    // Save the updated list
    await _savePendingOperations(operations);
  }

  /// Queue a create operation
  Future<void> queueCreate(Subscription subscription) async {
    await _addPendingOperation(
      PendingSyncOperation(
        id: subscription.id,
        type: SyncOperationType.create,
        data: subscription.toJson(),
      ),
    );
    
    // Try to sync immediately if online
    if (_isOnline) {
      sync();
    }
  }

  /// Queue an update operation
  Future<void> queueUpdate(Subscription subscription) async {
    await _addPendingOperation(
      PendingSyncOperation(
        id: subscription.id,
        type: SyncOperationType.update,
        data: subscription.toJson(),
      ),
    );
    
    // Try to sync immediately if online
    if (_isOnline) {
      sync();
    }
  }

  /// Queue a delete operation
  Future<void> queueDelete(String id) async {
    await _addPendingOperation(
      PendingSyncOperation(
        id: id,
        type: SyncOperationType.delete,
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
      final pendingOperations = _getPendingOperations();
      final successfulOperations = <PendingSyncOperation>[];
      
      for (final operation in pendingOperations) {
        try {
          switch (operation.type) {
            case SyncOperationType.create:
              if (operation.data != null) {
                final subscription = Subscription.fromJson(operation.data!);
                await _apiService.createSubscription(subscription);
                successfulOperations.add(operation);
              }
              break;
            case SyncOperationType.update:
              if (operation.data != null) {
                final subscription = Subscription.fromJson(operation.data!);
                await _apiService.updateSubscription(subscription);
                successfulOperations.add(operation);
              }
              break;
            case SyncOperationType.delete:
              await _apiService.deleteSubscription(operation.id);
              successfulOperations.add(operation);
              break;
          }
        } catch (e) {
          print('Error processing operation ${operation.type} for ${operation.id}: $e');
          // We'll retry this operation on the next sync
        }
      }
      
      // Remove successful operations from the pending list
      if (successfulOperations.isNotEmpty) {
        final remainingOperations = pendingOperations
            .where((op) => !successfulOperations.contains(op))
            .toList();
        await _savePendingOperations(remainingOperations);
      }
      
      // Fetch latest data from backend
      final remoteSubscriptions = await _apiService.getSubscriptions();
      final localSubscriptions = _subscriptionRepository.getAll();
      
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
      
      // Update last sync time
      await _prefs.setString(_lastSyncTimeKey, DateTime.now().toIso8601String());
    } catch (e) {
      print('Sync error: $e');
    } finally {
      _isSyncing = false;
    }
  }

  /// Get the last sync time
  DateTime? getLastSyncTime() {
    final lastSyncTimeStr = _prefs.getString(_lastSyncTimeKey);
    if (lastSyncTimeStr == null) return null;
    
    try {
      return DateTime.parse(lastSyncTimeStr);
    } catch (_) {
      return null;
    }
  }

  /// Check if there are pending operations
  bool hasPendingOperations() {
    return _getPendingOperations().isNotEmpty;
  }

  /// Dispose resources
  void dispose() {
    _connectivitySubscription?.cancel();
    _syncTimer?.cancel();
  }
}
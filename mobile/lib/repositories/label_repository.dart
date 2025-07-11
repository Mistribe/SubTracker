import 'package:hive/hive.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:uuid/uuid.dart';
import '../models/label.dart';
import '../models/subscription.dart';
import '../providers/sync_provider.dart';

/// Repository for handling label data persistence
///
/// This class abstracts the data access layer for labels, making it easier to switch
/// between different storage mechanisms in the future.
class LabelRepository {
  static const String _boxName = 'labels';
  late Box<Label> _box;
  final _uuid = Uuid();
  SyncProvider? _syncProvider;

  // Default labels
  static const List<Map<String, dynamic>> _defaultLabels = [
    {'name': 'Music', 'color': '#81C784'},
    {'name': 'Internet', 'color': '#64B5F6'},
    {'name': 'Mobile', 'color': '#FFD54F'},
    {'name': 'Utilities', 'color': '#9575CD'},
    {'name': 'Streaming', 'color': '#F06292'},
    {'name': 'Gaming', 'color': '#4DB6AC'},
    {'name': 'Software', 'color': '#7986CB'},
  ];

  /// Initialize the repository
  ///
  /// This method must be called before using any other methods in this class.
  Future<void> initialize() async {
    // Open the Hive box for labels
    _box = await Hive.openBox<Label>(_boxName);

    // Create default labels if the box is empty
    if (_box.isEmpty) {
      await _createDefaultLabels();
    }
  }

  /// Set the sync provider
  void setSyncProvider(SyncProvider syncProvider) {
    _syncProvider = syncProvider;
  }

  /// Create default labels
  Future<void> _createDefaultLabels() async {
    for (var labelData in _defaultLabels) {
      final label = Label(
        id: _uuid.v7(),
        name: labelData['name'],
        isDefault: true,
        color: labelData['color'],
      );
      await _box.put(label.id, label);
    }
  }

  /// Get all labels
  List<Label> getAll() {
    return _box.values.toList();
  }

  /// Get default labels
  List<Label> getDefaultLabels() {
    return _box.values.where((label) => label.isDefault).toList();
  }

  /// Get custom labels
  List<Label> getCustomLabels() {
    return _box.values.where((label) => !label.isDefault).toList();
  }

  /// Get a label by ID
  Label? get(String id) {
    return _box.get(id);
  }

  /// Add a new label
  Future<Label> add(String name, String color) async {
    final label = Label(
      id: _uuid.v7(),
      name: name,
      isDefault: false,
      color: color,
    );
    // Save to local storage
    await _box.put(label.id, label);

    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueCreate(label);
    }

    return label;
  }

  /// Update an existing label
  Future<void> update(Label label) async {
    // Save to local storage
    await _box.put(label.id, label);

    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueUpdate(label);
    }
  }

  /// Delete a label
  Future<void> delete(String id) async {
    // Delete from local storage
    await _box.delete(id);

    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueDelete(id);
    }
  }

  /// Clear all custom labels
  Future<void> clearCustomLabels() async {
    final customLabels = getCustomLabels();
    for (var label in customLabels) {
      await delete(label.id); // Use delete method to ensure sync
    }
  }

  /// Clear all labels (including default)
  Future<void> clearAll() async {
    // Get all IDs before clearing
    final ids = _box.keys.cast<String>().toList();

    // Clear local storage
    await _box.clear();

    // Queue deletes for sync if provider is available
    if (_syncProvider != null) {
      for (final id in ids) {
        await _syncProvider!.queueDelete(id);
      }
    }
  }
}

import 'package:hive/hive.dart';
import 'package:uuid/uuid.dart';
import '../models/family_member.dart';
import '../providers/sync_provider.dart';

class FamilyMemberRepository {
  static const String boxName = 'family_members';
  final Box<FamilyMember> _box;
  final Uuid _uuid = Uuid();
  SyncProvider? _syncProvider;

  FamilyMemberRepository(this._box);

  /// Set the sync provider
  void setSyncProvider(SyncProvider syncProvider) {
    _syncProvider = syncProvider;
  }

  // Initialize the repository
  static Future<FamilyMemberRepository> initialize() async {
    final box = await Hive.openBox<FamilyMember>(boxName);
    return FamilyMemberRepository(box);
  }

  // Get all family members
  List<FamilyMember> getAll() {
    return _box.values.toList();
  }

  // Get a family member by ID
  FamilyMember? get(String id) {
    return _box.get(id);
  }

  // Add a new family member
  Future<FamilyMember> add(String name, {bool isKid = false}) async {
    final familyMember = FamilyMember(id: _uuid.v7(), name: name, isKid: isKid);
    // Save to local storage
    await _box.put(familyMember.id, familyMember);

    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueCreateFamilyMember(familyMember);
    }

    return familyMember;
  }

  // Update an existing family member
  Future<void> update(FamilyMember familyMember) async {
    // Save to local storage
    await _box.put(familyMember.id, familyMember);

    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueUpdateFamilyMember(familyMember);
    }
  }

  // Delete a family member
  Future<void> delete(String id) async {
    // Delete from local storage
    await _box.delete(id);

    // Queue for sync if provider is available
    if (_syncProvider != null) {
      await _syncProvider!.queueDeleteFamilyMember(id);
    }
  }

  // Clear all family members
  Future<void> clearAll() async {
    // Get all IDs before clearing
    final ids = _box.keys.cast<String>().toList();

    // Clear local storage
    await _box.clear();

    // Queue deletes for sync if provider is available
    if (_syncProvider != null) {
      for (final id in ids) {
        await _syncProvider!.queueDeleteFamilyMember(id);
      }
    }
  }

  // Close the box when done
  Future<void> close() async {
    await _box.close();
  }
}

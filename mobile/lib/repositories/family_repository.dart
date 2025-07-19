import 'package:hive/hive.dart';
import 'package:subscription_tracker/models/family.dart';
import 'package:uuid/uuid.dart';
import '../models/family_member.dart';
import '../providers/sync_provider.dart';

class FamilyRepository {
  static const String boxName = 'families';
  late Box<Family> _box;
  final Uuid _uuid = Uuid();
  SyncProvider? _syncProvider;

  /// Set the sync provider
  void setSyncProvider(SyncProvider syncProvider) {
    _syncProvider = syncProvider;
  }

  // Initialize the repository
  Future<void> initialize() async {
    _box = await Hive.openBox<Family>(boxName);
  }

  // Get all family members
  List<Family> getAll() {
    return _box.values.toList();
  }

  // Get a family member by ID
  Family? get(String id) {
    return _box.get(id);
  }

  // Check if an owned family exists
  bool hasOwnedFamily() {
    return _box.values.any((family) => family.isOwner);
  }

  Future<Family?> create(String name, {bool haveJointAccount = true}) async {
    // Prevent creating a family if the user already owns one
    if (hasOwnedFamily()) {
      return null;
    }

    final family = Family(
      id: _uuid.v7(),
      name: name,
      isOwner: true,
      familyMembers: [],
      haveJointAccount: haveJointAccount,
    );

    await _box.put(family.id, family);

    if (_syncProvider != null) {
      await _syncProvider!.queueCreateFamily(family);
    }

    return family;
  }

  Future<void> update(
    String id,
    String name, {
    bool haveJointAccount = true,
    List<FamilyMember>? members,
    DateTime? updatedAt,
    bool withSync = true,
  }) async {
    final family = get(id);
    if (family == null) {
      return;
    }

    final updatedFamily = family.copyWith(
      name: name,
      members: members,
      updatedAt: updatedAt ?? DateTime.now(),
      haveJointAccount: haveJointAccount,
    );

    await _box.put(id, updatedFamily);

    if (_syncProvider != null && withSync) {
      await _syncProvider!.queueUpdateFamily(updatedFamily);
    }
  }

  Future<void> delete(String id, {bool withSync = true}) async {
    final family = get(id);
    if (family == null) {
      return;
    }

    await _box.delete(id);

    if (_syncProvider != null && withSync) {
      await _syncProvider!.queueDeleteFamily(id);
    }
  }

  // Add a new family member
  Future<FamilyMember?> addMember(
    String familyId,
    String name, {
    String? email,
    bool isKid = false,
  }) async {
    final family = get(familyId);
    if (family == null) {
      return null;
    }
    final familyMember = FamilyMember(
      id: _uuid.v7(),
      familyId: familyId,
      name: name,
      email: email,
      isKid: isKid,
    );

    final updatedFamily = family.copyWith(
      members: [...family.members, familyMember],
    );

    await _box.put(family.id, updatedFamily);

    if (_syncProvider != null) {
      await _syncProvider!.queueCreateFamilyMember(familyMember);
    }

    return familyMember;
  }

  // Update an existing family member
  Future<void> updateMember(
    String familyId,
    String memberId,
    String name, {
    String? email,
    bool isKid = false,
    bool withSync = true,
  }) async {
    final family = get(familyId);
    if (family == null) {
      return;
    }

    final familyMember = family.members.firstWhere(
      (member) => member.id == memberId,
      orElse: () => FamilyMember.empty(),
    );

    if (familyMember == FamilyMember.empty()) {
      return;
    }

    final updatedFamilyMember = familyMember.copyWith(
      name: name,
      email: email,
      isKid: isKid,
    );

    final updatedMembers = family.members
        .map(
          (member) => member.id == updatedFamilyMember.id
              ? updatedFamilyMember
              : member,
        )
        .toList();

    final updatedFamily = family.copyWith(members: updatedMembers);

    // Save to local storage
    await _box.put(familyId, updatedFamily);

    // Queue for sync if provider is available
    if (_syncProvider != null && withSync) {
      await _syncProvider!.queueUpdateFamilyMember(updatedFamilyMember);
    }
  }

  // Delete a family member
  Future<void> deleteMember(
    String familyId,
    String memberId, {
    bool withSync = true,
  }) async {
    final family = get(familyId);
    if (family == null) {
      return;
    }

    final updatedFamily = family.copyWith(
      members: family.members.where((member) => member.id != memberId).toList(),
    );

    await _box.put(familyId, updatedFamily);

    // Queue for sync if provider is available
    if (_syncProvider != null && withSync) {
      await _syncProvider!.queueDeleteFamilyMember(familyId, memberId);
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
        await _syncProvider!.queueDeleteFamily(id);
      }
    }
  }

  // Close the box when done
  Future<void> close() async {
    await _box.close();
  }
}

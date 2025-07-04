import 'package:hive/hive.dart';
import 'package:uuid/uuid.dart';
import '../models/family_member.dart';

class FamilyMemberRepository {
  static const String boxName = 'family_members';
  final Box<FamilyMember> _box;
  final Uuid _uuid = Uuid();

  FamilyMemberRepository(this._box);

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
    final familyMember = FamilyMember(
      id: _uuid.v7(),
      name: name,
      isKid: isKid,
    );
    await _box.put(familyMember.id, familyMember);
    return familyMember;
  }

  // Update an existing family member
  Future<void> update(FamilyMember familyMember) async {
    await _box.put(familyMember.id, familyMember);
  }

  // Delete a family member
  Future<void> delete(String id) async {
    await _box.delete(id);
  }

  // Close the box when done
  Future<void> close() async {
    await _box.close();
  }
}

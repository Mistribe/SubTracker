import 'package:flutter/foundation.dart';
import '../models/family.dart';
import '../models/family_member.dart';
import '../repositories/family_repository.dart';

class FamilyProvider with ChangeNotifier {
  final FamilyRepository familyRepository;

  String? _selectedFamilyId;

  FamilyProvider({required this.familyRepository});

  String? get selectedFamilyId => _selectedFamilyId;

  void setSelectedFamilyId(String? id) {
    _selectedFamilyId = id;
    notifyListeners();
  }

  bool get haveFamilySelected => _selectedFamilyId != null;

  bool get hasFamilies => familyRepository.getAll().isNotEmpty;

  bool get hasOwnedFamily => familyRepository.hasOwnedFamily();

  Family? get selectedFamily => _selectedFamilyId != null
      ? familyRepository.get(_selectedFamilyId!)
      : null;

  bool get canEditSelectedFamily => selectedFamily?.isOwner ?? false;

  bool get hasFamilyMembers => _selectedFamilyId == null
      ? false
      : (familyRepository.get(_selectedFamilyId!)?.members.isEmpty ?? false);

  // Getter for the family members list
  List<Family> get families => List.unmodifiable(familyRepository.getAll());

  List<FamilyMember> get familyMembers => _selectedFamilyId == null
      ? []
      : (familyRepository.get(_selectedFamilyId!)?.members ?? []);

  // Create a new family
  Future<Family?> createFamily(
    String name, {
    bool haveJointAccount = true,
  }) async {
    final family = await familyRepository.create(
      name,
      haveJointAccount: haveJointAccount,
    );
    if (family != null) {
      _selectedFamilyId = family.id;
      notifyListeners();
    }
    return family;
  }

  // Update an existing family
  Future<void> updateFamily(
    String id,
    String name, {
    bool haveJointAccount = true,
    List<FamilyMember>? members,
  }) async {
    await familyRepository.update(
      id,
      name,
      haveJointAccount: haveJointAccount,
      members: members,
    );
    notifyListeners();
  }

  // Delete a family
  Future<void> deleteFamily(String id) async {
    await familyRepository.delete(id);
    if (_selectedFamilyId == id) {
      _selectedFamilyId = null;
    }
    notifyListeners();
  }

  // Add a new family member
  Future<FamilyMember?> addFamilyMember(
    String familyId,
    String name, {
    bool isKid = false,
  }) async {
    final familyMember = await familyRepository.addMember(
      familyId,
      name,
      isKid: isKid,
    );
    if (familyMember == null) {
      // todo handle error
      return null;
    }
    notifyListeners();
    return familyMember;
  }

  // Update an existing family member
  Future<void> updateFamilyMember(
    String familyId,
    String memberId,
    String name, {
    bool isKid = false,
    String? email,
  }) async {
    final familyMember = familyRepository.get(memberId);
    if (familyMember == null) {
      return;
    }

    await familyRepository.updateMember(
      familyId,
      memberId,
      name,
      email: email,
      isKid: isKid,
    );

    notifyListeners();
  }

  // Remove a family member
  Future<void> removeFamilyMember(String familyId, String memberId) async {
    // Remove from storage
    await familyRepository.deleteMember(familyId, memberId);

    notifyListeners();
  }

  // Get a family member by ID
  FamilyMember? getFamilyMemberById(String? id) {
    if (id == null) return null;
    if (_selectedFamilyId == null) {
      return null;
    }

    final family = familyRepository.get(_selectedFamilyId!);
    if (family == null) return null;

    return family.members.firstWhere(
      (member) => member.id == id,
      orElse: () => FamilyMember.empty(),
    );
  }
}

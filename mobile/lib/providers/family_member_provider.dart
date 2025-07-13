import 'package:flutter/foundation.dart';
import '../models/family_member.dart';
import '../repositories/family_member_repository.dart';

class FamilyMemberProvider with ChangeNotifier {
  final FamilyMemberRepository familyMemberRepository;

  FamilyMemberProvider({required this.familyMemberRepository});

  // Getter for the family members list
  List<FamilyMember> get familyMembers =>
      List.unmodifiable(familyMemberRepository.getAll());

  // Check if there are any family members
  bool get hasFamilyMembers => familyMemberRepository.getAll().isNotEmpty;

  // Add a new family member
  Future<FamilyMember> addFamilyMember(
    String name, {
    bool isKid = false,
  }) async {
    final familyMember = await familyMemberRepository.add(name, isKid: isKid);
    notifyListeners();
    return familyMember;
  }

  // Update an existing family member
  Future<void> updateFamilyMember(String id, String name, {bool? isKid}) async {
    final familyMember = familyMemberRepository.get(id);
    if (familyMember == null) {
      return;
    }
    final updatedFamilyMember = familyMember.copyWith(name: name, isKid: isKid);

    await familyMemberRepository.update(updatedFamilyMember);

    notifyListeners();
  }

  // Remove a family member
  Future<void> removeFamilyMember(String id) async {
    // Remove from storage
    await familyMemberRepository.delete(id);

    notifyListeners();
  }

  // Get a family member by ID
  FamilyMember? getFamilyMemberById(String? id) {
    if (id == null) return null;

    return familyMemberRepository.get(id);
  }
}

import 'package:flutter/foundation.dart';
import '../models/family_member.dart';
import '../repositories/family_member_repository.dart';

class FamilyMemberProvider with ChangeNotifier {
  final FamilyMemberRepository familyMemberRepository;
  List<FamilyMember> _familyMembers = [];

  FamilyMemberProvider({
    required this.familyMemberRepository,
  }) {
    // Load family members from repository
    _loadFamilyMembers();
  }

  // Load family members from repository
  Future<void> _loadFamilyMembers() async {
    _familyMembers = familyMemberRepository.getAll();
    notifyListeners();
  }

  // Getter for the family members list
  List<FamilyMember> get familyMembers => List.unmodifiable(_familyMembers);

  // Check if there are any family members
  bool get hasFamilyMembers => _familyMembers.isNotEmpty;

  // Add a new family member
  Future<void> addFamilyMember(String name) async {
    final familyMember = await familyMemberRepository.add(name);
    _familyMembers.add(familyMember);
    notifyListeners();
  }

  // Update an existing family member
  Future<void> updateFamilyMember(String id, String name) async {
    final index = _familyMembers.indexWhere(
      (familyMember) => familyMember.id == id,
    );

    if (index >= 0) {
      final familyMember = _familyMembers[index];
      final updatedFamilyMember = familyMember.copyWith(name: name);
      
      await familyMemberRepository.update(updatedFamilyMember);
      
      _familyMembers[index] = updatedFamilyMember;
      notifyListeners();
    }
  }

  // Remove a family member
  Future<void> removeFamilyMember(String id) async {
    _familyMembers.removeWhere((familyMember) => familyMember.id == id);
    
    // Remove from storage
    await familyMemberRepository.delete(id);
    
    notifyListeners();
  }

  // Get a family member by ID
  FamilyMember? getFamilyMemberById(String? id) {
    if (id == null) return null;
    
    final index = _familyMembers.indexWhere(
      (familyMember) => familyMember.id == id,
    );
    
    if (index >= 0) {
      return _familyMembers[index];
    }
    
    return null;
  }
}
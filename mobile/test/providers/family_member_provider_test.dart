import 'package:flutter_test/flutter_test.dart';
import 'package:subscription_tracker/models/family_member.dart';
import 'package:subscription_tracker/providers/family_provider.dart';
import 'package:subscription_tracker/repositories/family_repository.dart';

// Mock implementation of FamilyMemberRepository
class MockFamilyMemberRepository implements FamilyRepository {
  final List<FamilyMember> _familyMembers = [];

  @override
  List<FamilyMember> getAll() {
    return List.from(_familyMembers);
  }

  @override
  FamilyMember? get(String id) {
    try {
      return _familyMembers.firstWhere((member) => member.id == id);
    } catch (e) {
      return null;
    }
  }

  @override
  Future<FamilyMember> add(String name, {bool isKid = false}) async {
    final familyMember = FamilyMember(
      id: 'mock-id-${_familyMembers.length + 1}',
      name: name,
      isKid: isKid,
    );
    _familyMembers.add(familyMember);
    return familyMember;
  }

  @override
  Future<void> update(FamilyMember familyMember) async {
    final index = _familyMembers.indexWhere(
      (member) => member.id == familyMember.id,
    );

    if (index >= 0) {
      _familyMembers[index] = familyMember;
    }
  }

  @override
  Future<void> delete(String id) async {
    _familyMembers.removeWhere((member) => member.id == id);
  }

  // Unused methods from the interface
  @override
  Future<void> close() async {}

  @override
  dynamic noSuchMethod(Invocation invocation) => super.noSuchMethod(invocation);
}

void main() {
  group('FamilyMemberProvider', () {
    late MockFamilyMemberRepository mockRepository;
    late FamilyProvider provider;

    setUp(() {
      mockRepository = MockFamilyMemberRepository();
      provider = FamilyProvider(familyRepository: mockRepository);

      // Wait for the provider to initialize
      Future.delayed(Duration.zero);
    });

    test('initializes with empty list', () {
      expect(provider.families, isEmpty);
      expect(provider.hasFamilyMembers, isFalse);
    });

    test('addFamilyMember adds a family member', () async {
      await provider.addFamilyMember('John Doe');

      expect(provider.families.length, equals(1));
      expect(provider.families[0].name, equals('John Doe'));
      expect(provider.families[0].isKid, isFalse);
      expect(provider.hasFamilyMembers, isTrue);
    });

    test('addFamilyMember adds a kid', () async {
      await provider.addFamilyMember('Child', isKid: true);

      expect(provider.families.length, equals(1));
      expect(provider.families[0].name, equals('Child'));
      expect(provider.families[0].isKid, isTrue);
    });

    test('updateFamilyMember updates a family member', () async {
      // Add a family member first
      await provider.addFamilyMember('John Doe');
      final id = provider.families[0].id;

      // Update the family member
      await provider.updateFamilyMember(id, 'Jane Doe', isKid: true);

      expect(provider.families.length, equals(1));
      expect(provider.families[0].id, equals(id));
      expect(provider.families[0].name, equals('Jane Doe'));
      expect(provider.families[0].isKid, isTrue);
    });

    test('updateFamilyMember does nothing for non-existent ID', () async {
      // Add a family member first
      await provider.addFamilyMember('John Doe');

      // Try to update a non-existent family member
      await provider.updateFamilyMember('non-existent-id', 'Jane Doe');

      expect(provider.families.length, equals(1));
      expect(provider.families[0].name, equals('John Doe'));
    });

    test('removeFamilyMember removes a family member', () async {
      // Add a family member first
      await provider.addFamilyMember('John Doe');
      final id = provider.families[0].id;

      // Remove the family member
      await provider.removeFamilyMember(id);

      expect(provider.families, isEmpty);
      expect(provider.hasFamilyMembers, isFalse);
    });

    test('getFamilyMemberById returns the correct family member', () async {
      // Add two family members
      await provider.addFamilyMember('John Doe');
      await provider.addFamilyMember('Jane Doe');

      final id1 = provider.families[0].id;
      final id2 = provider.families[1].id;

      final member1 = provider.getFamilyMemberById(id1);
      final member2 = provider.getFamilyMemberById(id2);

      expect(member1?.name, equals('John Doe'));
      expect(member2?.name, equals('Jane Doe'));
    });

    test('getFamilyMemberById returns null for non-existent ID', () {
      final member = provider.getFamilyMemberById('non-existent-id');
      expect(member, isNull);
    });

    test('getFamilyMemberById returns null for null ID', () {
      final member = provider.getFamilyMemberById(null);
      expect(member, isNull);
    });
  });
}

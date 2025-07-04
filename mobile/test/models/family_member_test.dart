import 'dart:io';
import 'package:flutter_test/flutter_test.dart';
import 'package:hive/hive.dart';
import 'package:subscription_tracker/models/family_member.dart';

void main() {
  group('FamilyMember', () {
    setUpAll(() async {
      // Set up a temporary directory for Hive
      final tempDir = await Directory.systemTemp.createTemp('hive_test');
      Hive.init(tempDir.path);

      // Register adapters
      if (!Hive.isAdapterRegistered(4)) {
        Hive.registerAdapter(FamilyMemberAdapter());
      }
    });

    tearDownAll(() async {
      // Clean up after all tests
      await Hive.close();
    });

    test('constructor creates family member with correct values', () {
      final familyMember = FamilyMember(
        id: '1',
        name: 'John Doe',
        isKid: true,
      );

      expect(familyMember.id, equals('1'));
      expect(familyMember.name, equals('John Doe'));
      expect(familyMember.isKid, equals(true));
    });

    test('constructor uses default value for isKid', () {
      final familyMember = FamilyMember(
        id: '1',
        name: 'John Doe',
      );

      expect(familyMember.isKid, equals(false));
    });

    test('copyWith creates a new family member with updated values', () {
      final familyMember = FamilyMember(
        id: '1',
        name: 'John Doe',
        isKid: false,
      );

      final updatedFamilyMember = familyMember.copyWith(
        name: 'Jane Doe',
        isKid: true,
      );

      // Original family member should be unchanged
      expect(familyMember.id, equals('1'));
      expect(familyMember.name, equals('John Doe'));
      expect(familyMember.isKid, equals(false));

      // Updated family member should have new values
      expect(updatedFamilyMember.id, equals('1')); // ID should remain the same
      expect(updatedFamilyMember.name, equals('Jane Doe'));
      expect(updatedFamilyMember.isKid, equals(true));
    });

    test('toJson converts family member to JSON correctly', () {
      final familyMember = FamilyMember(
        id: '1',
        name: 'John Doe',
        isKid: true,
      );

      final json = familyMember.toJson();

      expect(json['id'], equals('1'));
      expect(json['name'], equals('John Doe'));
      expect(json['isKid'], equals(true));
    });

    test('fromJson creates family member from JSON correctly', () {
      final json = {
        'id': '1',
        'name': 'John Doe',
        'isKid': true,
      };

      final familyMember = FamilyMember.fromJson(json);

      expect(familyMember.id, equals('1'));
      expect(familyMember.name, equals('John Doe'));
      expect(familyMember.isKid, equals(true));
    });

    test('fromJson handles missing isKid field', () {
      final json = {
        'id': '1',
        'name': 'John Doe',
      };

      final familyMember = FamilyMember.fromJson(json);

      expect(familyMember.isKid, equals(false));
    });
  });
}
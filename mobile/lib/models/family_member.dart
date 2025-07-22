import 'package:hive/hive.dart';
import 'package:uuid/uuid.dart';

part 'family_member.g.dart';

@HiveType(typeId: 4)
class FamilyMember extends HiveObject {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String name;

  @HiveField(2)
  final String? email;

  @HiveField(3)
  final bool isKid;

  @HiveField(4)
  final String familyId;

  @HiveField(5)
  final DateTime createdAt;

  @HiveField(6)
  final DateTime updatedAt;

  @HiveField(7)
  final String eTag;

  FamilyMember({
    required this.id,
    required this.name,
    required this.familyId,
    this.isKid = false,
    this.email,
    this.eTag = '',
    DateTime? createdAt,
    DateTime? updatedAt,
  }) : createdAt = createdAt ?? DateTime.now(),
       updatedAt = updatedAt ?? DateTime.now();

  // Create a copy of this family member with updated fields
  FamilyMember copyWith({
    String? id,
    String? name,
    String? email,
    bool? isKid,
    String? eTag,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return FamilyMember(
      id: id ?? this.id,
      familyId: familyId,
      name: name ?? this.name,
      eTag: eTag ?? this.eTag,
      isKid: isKid ?? this.isKid,
      email: email ?? this.email,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? DateTime.now(),
    );
  }

  // Convert to and from JSON for potential future persistence
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'family_id': familyId,
      'name': name,
      'is_kid': isKid,
      'email': email,
      'etag': eTag,
      'created_at': createdAt.toUtc().toIso8601String(),
      'updated_at': updatedAt.toUtc().toIso8601String(),
    };
  }

  factory FamilyMember.fromJson(Map<String, dynamic> json) {
    return FamilyMember(
      id: json['id'],
      familyId: json['family_id'],
      name: json['name'],
      isKid: json['is_kid'] ?? false,
      email: json['email'],
      eTag: json['etag'],
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : null,
      updatedAt: json['updated_at'] != null
          ? DateTime.parse(json['updated_at'])
          : null,
    );
  }

  factory FamilyMember.empty() {
    const uuid = Uuid();
    return FamilyMember(
      id: uuid.v7(),
      familyId: '',
      name: '',
      isKid: false,
      eTag: '',
      email: null,
      createdAt: DateTime.fromMillisecondsSinceEpoch(0),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(0),
    );
  }
}

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
  final bool isKid;

  @HiveField(3)
  final DateTime createdAt;

  @HiveField(4)
  final DateTime updatedAt;

  FamilyMember({
    required this.id,
    required this.name,
    this.isKid = false,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) : createdAt = createdAt ?? DateTime.now(),
       updatedAt = updatedAt ?? DateTime.now();

  // Create a copy of this family member with updated fields
  FamilyMember copyWith({
    String? id,
    String? name,
    bool? isKid,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return FamilyMember(
      id: id ?? this.id,
      name: name ?? this.name,
      isKid: isKid ?? this.isKid,
      createdAt: createdAt ?? this.createdAt,
      updatedAt:
          updatedAt ?? DateTime.now(), // Always update the updatedAt field
    );
  }

  // Convert to and from JSON for potential future persistence
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'isKid': isKid,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }

  factory FamilyMember.fromJson(Map<String, dynamic> json) {
    return FamilyMember(
      id: json['id'],
      name: json['name'],
      isKid: json['isKid'] ?? false,
      createdAt: json['createdAt'] != null
          ? DateTime.parse(json['createdAt'])
          : null,
      updatedAt: json['updatedAt'] != null
          ? DateTime.parse(json['updatedAt'])
          : null,
    );
  }

  factory FamilyMember.empty() {
    const uuid = Uuid();
    return FamilyMember(
      id: uuid.v7(),
      name: '',
      isKid: false,
      createdAt: DateTime.fromMillisecondsSinceEpoch(0),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(0),
    );
  }
}

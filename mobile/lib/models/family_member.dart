import 'package:hive/hive.dart';

part 'family_member.g.dart';

@HiveType(typeId: 4)
class FamilyMember extends HiveObject {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String name;

  FamilyMember({
    required this.id,
    required this.name,
  });

  // Create a copy of this family member with updated fields
  FamilyMember copyWith({
    String? id,
    String? name,
  }) {
    return FamilyMember(
      id: id ?? this.id,
      name: name ?? this.name,
    );
  }

  // Convert to and from JSON for potential future persistence
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
    };
  }

  factory FamilyMember.fromJson(Map<String, dynamic> json) {
    return FamilyMember(
      id: json['id'],
      name: json['name'],
    );
  }
}
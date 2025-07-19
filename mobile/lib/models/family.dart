import 'package:hive/hive.dart';
import 'package:subscription_tracker/models/family_member.dart';
import 'package:uuid/uuid.dart';

part 'family.g.dart';

@HiveType(typeId: 6)
class Family extends HiveObject {
  @HiveField(0)
  final String id;
  @HiveField(1)
  final String name;
  @HiveField(2)
  final List<FamilyMember> members;
  @HiveField(3)
  final bool haveJointAccount;
  @HiveField(4)
  final bool isOwner;
  @HiveField(5)
  final DateTime createdAt;
  @HiveField(6)
  final DateTime updatedAt;
  @HiveField(7)
  final String eTag;

  Family({
    required this.id,
    required this.name,
    required this.isOwner,
    required this.haveJointAccount,
    this.eTag = "",
    List<FamilyMember>? familyMembers,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) : members = familyMembers ?? [],
       createdAt = createdAt ?? DateTime.now(),
       updatedAt = updatedAt ?? DateTime.now();

  Family copyWith({
    String? id,
    String? name,
    List<FamilyMember>? members,
    bool? haveJointAccount,
    String? eTag,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Family(
      id: id ?? this.id,
      name: name ?? this.name,
      isOwner: isOwner,
      eTag: eTag ?? this.eTag,
      familyMembers: members ?? this.members,
      haveJointAccount: haveJointAccount ?? this.haveJointAccount,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'is_owner': isOwner,
      'members': members.map((member) => member.toJson()).toList(),
      'have_joint_account': haveJointAccount,
      'etag': eTag,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  factory Family.fromJson(Map<String, dynamic> json) {
    return Family(
      id: json['id'],
      name: json['name'],
      isOwner: json['is_owner'],
      haveJointAccount: json['have_joint_account'],
      familyMembers: (json['members'] as List?)
          ?.map((member) => FamilyMember.fromJson(member))
          .toList(),
      eTag: json['etag'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }

  factory Family.empty() {
    const uuid = Uuid();
    return Family(
      id: uuid.v7(),
      name: '',
      isOwner: false,
      familyMembers: [],
      haveJointAccount: true,
      eTag: '',
      createdAt: DateTime.fromMillisecondsSinceEpoch(0),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(0),
    );
  }
}

import 'package:hive/hive.dart';
import 'package:uuid/uuid.dart';

part 'label.g.dart';

@HiveType(typeId: 3)
class Label extends HiveObject {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String name;

  @HiveField(2)
  final bool isDefault;

  @HiveField(3)
  final String color;

  @HiveField(4)
  final DateTime createdAt;

  @HiveField(5)
  final DateTime updatedAt;

  @HiveField(6)
  final String eTag;

  Label({
    required this.id,
    required this.name,
    required this.color,
    this.isDefault = false,
    this.eTag = '',
    DateTime? createdAt,
    DateTime? updatedAt,
  }) : createdAt = createdAt ?? DateTime.now(),
       updatedAt = updatedAt ?? DateTime.now();

  // Create a copy of this label with updated fields
  Label copyWith({
    String? id,
    String? name,
    bool? isDefault,
    String? color,
    String? eTag,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Label(
      id: id ?? this.id,
      name: name ?? this.name,
      isDefault: isDefault ?? this.isDefault,
      color: color ?? this.color,
      eTag: eTag ?? this.eTag,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? DateTime.now(),
    );
  }

  // Convert to and from JSON for potential future persistence
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'is_default': isDefault,
      'etag': eTag,
      'color': color,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  factory Label.fromJson(Map<String, dynamic> json) {
    return Label(
      id: json['id'],
      name: json['name'],
      isDefault: json['is_default'] ?? false,
      color: json['color'],
      eTag: json['etag'],
      createdAt: json['created_at'] != null
          ? DateTime.parse(json['created_at'])
          : null,
      updatedAt: json['updated_at'] != null
          ? DateTime.parse(json['updated_at'])
          : null,
    );
  }

  factory Label.empty() {
    const uuid = Uuid();
    return Label(
      id: uuid.v7(),
      name: '',
      isDefault: false,
      eTag: '',
      color: '#000000',
      createdAt: DateTime.fromMillisecondsSinceEpoch(0),
      updatedAt: DateTime.fromMillisecondsSinceEpoch(0),
    );
  }
}

import 'package:hive/hive.dart';

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

  Label({
    required this.id,
    required this.name,
    this.isDefault = false,
    required this.color,
  });

  // Create a copy of this label with updated fields
  Label copyWith({String? id, String? name, bool? isDefault, String? color}) {
    return Label(
      id: id ?? this.id,
      name: name ?? this.name,
      isDefault: isDefault ?? this.isDefault,
      color: color ?? this.color,
    );
  }

  // Convert to and from JSON for potential future persistence
  Map<String, dynamic> toJson() {
    return {'id': id, 'name': name, 'isDefault': isDefault, 'color': color};
  }

  factory Label.fromJson(Map<String, dynamic> json) {
    return Label(
      id: json['id'],
      name: json['name'],
      isDefault: json['isDefault'] ?? false,
      color: json['color'],
    );
  }
}

// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'family.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class FamilyAdapter extends TypeAdapter<Family> {
  @override
  final int typeId = 6;

  @override
  Family read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return Family(
      id: fields[0] as String,
      name: fields[1] as String,
      isOwner: fields[4] as bool,
      haveJointAccount: fields[3] as bool,
      createdAt: fields[5] as DateTime?,
      updatedAt: fields[6] as DateTime?,
    );
  }

  @override
  void write(BinaryWriter writer, Family obj) {
    writer
      ..writeByte(7)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.name)
      ..writeByte(2)
      ..write(obj.members)
      ..writeByte(3)
      ..write(obj.haveJointAccount)
      ..writeByte(4)
      ..write(obj.isOwner)
      ..writeByte(5)
      ..write(obj.createdAt)
      ..writeByte(6)
      ..write(obj.updatedAt);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is FamilyAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}

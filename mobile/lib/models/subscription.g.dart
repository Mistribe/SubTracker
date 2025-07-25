// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'subscription.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class SubscriptionAdapter extends TypeAdapter<Subscription> {
  @override
  final int typeId = 0;

  @override
  Subscription read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return Subscription(
      id: fields[0] as String,
      name: fields[1] as String,
      subscriptionPayments: (fields[3] as List?)?.cast<SubscriptionPayment>(),
      labelIds: (fields[4] as List?)?.cast<String>(),
      userFamilyMemberIds: (fields[5] as List?)?.cast<String>(),
      payerId: fields[6] as String?,
      payedByJointAccount: fields[7] as bool,
      eTag: fields[10] as String,
      createdAt: fields[8] as DateTime?,
      updatedAt: fields[9] as DateTime?,
    );
  }

  @override
  void write(BinaryWriter writer, Subscription obj) {
    writer
      ..writeByte(10)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.name)
      ..writeByte(3)
      ..write(obj.payments)
      ..writeByte(4)
      ..write(obj.labelIds)
      ..writeByte(5)
      ..write(obj.familyMemberIds)
      ..writeByte(6)
      ..write(obj.payerId)
      ..writeByte(7)
      ..write(obj.payedByJointAccount)
      ..writeByte(8)
      ..write(obj.createdAt)
      ..writeByte(9)
      ..write(obj.updatedAt)
      ..writeByte(10)
      ..write(obj.eTag);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SubscriptionAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}

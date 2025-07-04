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
      labels: (fields[4] as List?)?.cast<Label>(),
      userFamilyMembers: (fields[5] as List?)?.cast<FamilyMember>(),
      payerFamilyMember: fields[6] as FamilyMember?,
    );
  }

  @override
  void write(BinaryWriter writer, Subscription obj) {
    writer
      ..writeByte(6)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.name)
      ..writeByte(3)
      ..write(obj.subscriptionPayments)
      ..writeByte(4)
      ..write(obj.labels)
      ..writeByte(5)
      ..write(obj.userFamilyMembers)
      ..writeByte(6)
      ..write(obj.payerFamilyMember);
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

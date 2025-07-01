// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'subscription_payment.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class SubscriptionPaymentAdapter extends TypeAdapter<SubscriptionPayment> {
  @override
  final int typeId = 1;

  @override
  SubscriptionPayment read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return SubscriptionPayment(
      id: fields[0] as String,
      price: fields[1] as double,
      startDate: fields[2] as DateTime,
      endDate: fields[3] as DateTime?,
      months: fields[4] as int,
    );
  }

  @override
  void write(BinaryWriter writer, SubscriptionPayment obj) {
    writer
      ..writeByte(5)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.price)
      ..writeByte(2)
      ..write(obj.startDate)
      ..writeByte(3)
      ..write(obj.endDate)
      ..writeByte(4)
      ..write(obj.months);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SubscriptionPaymentAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}

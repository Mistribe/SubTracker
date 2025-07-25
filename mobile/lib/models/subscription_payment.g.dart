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
      currency: fields[5] as String,
      eTag: fields[8] as String,
      freeTrialMonths: fields[9] as int,
      createdAt: fields[6] as DateTime?,
      updatedAt: fields[7] as DateTime?,
    );
  }

  @override
  void write(BinaryWriter writer, SubscriptionPayment obj) {
    writer
      ..writeByte(10)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.price)
      ..writeByte(2)
      ..write(obj.startDate)
      ..writeByte(3)
      ..write(obj.endDate)
      ..writeByte(4)
      ..write(obj.months)
      ..writeByte(5)
      ..write(obj.currency)
      ..writeByte(6)
      ..write(obj.createdAt)
      ..writeByte(7)
      ..write(obj.updatedAt)
      ..writeByte(8)
      ..write(obj.eTag)
      ..writeByte(9)
      ..write(obj.freeTrialMonths);
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

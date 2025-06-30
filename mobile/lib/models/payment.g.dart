// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'payment.dart';

// **************************************************************************
// TypeAdapterGenerator
// **************************************************************************

class PaymentAdapter extends TypeAdapter<Payment> {
  @override
  final int typeId = 0;

  @override
  Payment read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return Payment(
      id: fields[0] as String,
      name: fields[1] as String,
      price: fields[2] as double,
      isAnnual: fields[3] as bool,
      paymentDate: fields[4] as DateTime,
      priceHistory: (fields[5] as List?)?.cast<PriceChange>(),
    );
  }

  @override
  void write(BinaryWriter writer, Payment obj) {
    writer
      ..writeByte(6)
      ..writeByte(0)
      ..write(obj.id)
      ..writeByte(1)
      ..write(obj.name)
      ..writeByte(2)
      ..write(obj.price)
      ..writeByte(3)
      ..write(obj.isAnnual)
      ..writeByte(4)
      ..write(obj.paymentDate)
      ..writeByte(5)
      ..write(obj.priceHistory);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is PaymentAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}

class PriceChangeAdapter extends TypeAdapter<PriceChange> {
  @override
  final int typeId = 1;

  @override
  PriceChange read(BinaryReader reader) {
    final numOfFields = reader.readByte();
    final fields = <int, dynamic>{
      for (int i = 0; i < numOfFields; i++) reader.readByte(): reader.read(),
    };
    return PriceChange(
      price: fields[0] as double,
      endDate: fields[1] as DateTime,
    );
  }

  @override
  void write(BinaryWriter writer, PriceChange obj) {
    writer
      ..writeByte(2)
      ..writeByte(0)
      ..write(obj.price)
      ..writeByte(1)
      ..write(obj.endDate);
  }

  @override
  int get hashCode => typeId.hashCode;

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is PriceChangeAdapter &&
          runtimeType == other.runtimeType &&
          typeId == other.typeId;
}
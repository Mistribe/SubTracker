import 'dart:io';
import 'package:flutter_test/flutter_test.dart';
import 'package:hive/hive.dart';
import 'package:subscription_tracker/models/label.dart';

void main() {
  group('Label', () {
    setUpAll(() async {
      // Set up a temporary directory for Hive
      final tempDir = await Directory.systemTemp.createTemp('hive_test');
      Hive.init(tempDir.path);

      // Register adapters
      if (!Hive.isAdapterRegistered(3)) {
        Hive.registerAdapter(LabelAdapter());
      }
    });

    tearDownAll(() async {
      // Clean up after all tests
      await Hive.close();
    });

    test('constructor creates a label with correct values', () {
      final label = Label(
        id: '1',
        name: 'Entertainment',
        isDefault: true,
        color: '#FF0000',
      );

      expect(label.id, equals('1'));
      expect(label.name, equals('Entertainment'));
      expect(label.isDefault, isTrue);
      expect(label.color, equals('#FF0000'));
    });

    test('constructor creates a label with default isDefault value', () {
      final label = Label(
        id: '1',
        name: 'Entertainment',
        color: '#FF0000',
      );

      expect(label.isDefault, isFalse);
    });

    test('copyWith creates a new label with updated values', () {
      final label = Label(
        id: '1',
        name: 'Entertainment',
        isDefault: true,
        color: '#FF0000',
      );

      final updatedLabel = label.copyWith(
        name: 'Movies',
        isDefault: false,
        color: '#00FF00',
      );

      // Original label should be unchanged
      expect(label.id, equals('1'));
      expect(label.name, equals('Entertainment'));
      expect(label.isDefault, isTrue);
      expect(label.color, equals('#FF0000'));

      // Updated label should have new values
      expect(updatedLabel.id, equals('1')); // ID should remain the same
      expect(updatedLabel.name, equals('Movies'));
      expect(updatedLabel.isDefault, isFalse);
      expect(updatedLabel.color, equals('#00FF00'));
    });

    test('copyWith with no parameters returns a label with the same values', () {
      final label = Label(
        id: '1',
        name: 'Entertainment',
        isDefault: true,
        color: '#FF0000',
      );

      final updatedLabel = label.copyWith();

      expect(updatedLabel.id, equals(label.id));
      expect(updatedLabel.name, equals(label.name));
      expect(updatedLabel.isDefault, equals(label.isDefault));
      expect(updatedLabel.color, equals(label.color));
    });

    test('toJson converts label to JSON correctly', () {
      final label = Label(
        id: '1',
        name: 'Entertainment',
        isDefault: true,
        color: '#FF0000',
      );

      final json = label.toJson();

      expect(json['id'], equals('1'));
      expect(json['name'], equals('Entertainment'));
      expect(json['isDefault'], isTrue);
      expect(json['color'], equals('#FF0000'));
    });

    test('fromJson creates label from JSON correctly', () {
      final json = {
        'id': '1',
        'name': 'Entertainment',
        'isDefault': true,
        'color': '#FF0000',
      };

      final label = Label.fromJson(json);

      expect(label.id, equals('1'));
      expect(label.name, equals('Entertainment'));
      expect(label.isDefault, isTrue);
      expect(label.color, equals('#FF0000'));
    });

    test('fromJson handles missing isDefault', () {
      final json = {
        'id': '1',
        'name': 'Entertainment',
        'color': '#FF0000',
      };

      final label = Label.fromJson(json);

      expect(label.isDefault, isFalse);
    });
  });
}
import 'dart:io';
import 'package:flutter_test/flutter_test.dart';
import 'package:hive/hive.dart';
import 'package:subscription_tracker/models/settings.dart';
import 'package:subscription_tracker/models/currency.dart';

void main() {
  group('Settings', () {
    setUpAll(() async {
      // Set up a temporary directory for Hive
      final tempDir = await Directory.systemTemp.createTemp('hive_test');
      Hive.init(tempDir.path);

      // Register adapters
      if (!Hive.isAdapterRegistered(2)) {
        Hive.registerAdapter(SettingsAdapter());
      }
    });

    tearDownAll(() async {
      // Clean up after all tests
      await Hive.close();
    });

    test('constructor creates settings with correct values', () {
      final settings = Settings(
        isDarkMode: true,
        useSystemTheme: false,
        defaultCurrency: 'EUR',
      );

      expect(settings.isDarkMode, isTrue);
      expect(settings.useSystemTheme, isFalse);
      expect(settings.defaultCurrency, equals('EUR'));
    });

    test('constructor creates settings with default values', () {
      final settings = Settings();

      expect(settings.isDarkMode, isFalse);
      expect(settings.useSystemTheme, isTrue);
      expect(settings.defaultCurrency, equals(Currency.defaultCode));
    });

    test('defaultCurrencyEnum returns correct Currency enum', () {
      final settings = Settings(defaultCurrency: 'EUR');
      expect(settings.defaultCurrencyEnum, equals(Currency.EUR));

      final settings2 = Settings(defaultCurrency: 'USD');
      expect(settings2.defaultCurrencyEnum, equals(Currency.USD));

      final settings3 = Settings(defaultCurrency: 'GBP');
      expect(settings3.defaultCurrencyEnum, equals(Currency.GBP));
    });

    test('copyWith creates new settings with updated values', () {
      final settings = Settings(
        isDarkMode: false,
        useSystemTheme: true,
        defaultCurrency: 'USD',
      );

      final updatedSettings = settings.copyWith(
        isDarkMode: true,
        useSystemTheme: false,
        defaultCurrency: 'EUR',
      );

      // Original settings should be unchanged
      expect(settings.isDarkMode, isFalse);
      expect(settings.useSystemTheme, isTrue);
      expect(settings.defaultCurrency, equals('USD'));

      // Updated settings should have new values
      expect(updatedSettings.isDarkMode, isTrue);
      expect(updatedSettings.useSystemTheme, isFalse);
      expect(updatedSettings.defaultCurrency, equals('EUR'));
    });

    test('copyWith with defaultCurrencyEnum updates defaultCurrency', () {
      final settings = Settings(defaultCurrency: 'USD');
      
      final updatedSettings = settings.copyWith(
        defaultCurrencyEnum: Currency.EUR,
      );

      expect(updatedSettings.defaultCurrency, equals('EUR'));
    });

    test('copyWith with no parameters returns settings with the same values', () {
      final settings = Settings(
        isDarkMode: true,
        useSystemTheme: false,
        defaultCurrency: 'EUR',
      );

      final updatedSettings = settings.copyWith();

      expect(updatedSettings.isDarkMode, equals(settings.isDarkMode));
      expect(updatedSettings.useSystemTheme, equals(settings.useSystemTheme));
      expect(updatedSettings.defaultCurrency, equals(settings.defaultCurrency));
    });

    test('copyWith prioritizes defaultCurrencyEnum over defaultCurrency', () {
      final settings = Settings(defaultCurrency: 'USD');
      
      final updatedSettings = settings.copyWith(
        defaultCurrency: 'GBP',
        defaultCurrencyEnum: Currency.EUR,
      );

      expect(updatedSettings.defaultCurrency, equals('EUR'));
    });
  });
}
import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:subscription_tracker/models/settings.dart';
import 'package:subscription_tracker/providers/theme_provider.dart';
import 'package:subscription_tracker/repositories/settings_repository.dart';

// Mock implementation of SettingsRepository
class MockSettingsRepository implements SettingsRepository {
  Settings _settings = Settings();
  
  @override
  Settings getSettings() {
    return _settings;
  }
  
  @override
  Future<void> saveSettings(Settings settings) async {
    _settings = settings;
  }
  
  @override
  Future<void> updateThemeSettings({bool? isDarkMode, bool? useSystemTheme}) async {
    _settings = _settings.copyWith(
      isDarkMode: isDarkMode,
      useSystemTheme: useSystemTheme,
    );
  }
  
  @override
  Future<void> updateDefaultCurrency(String currency) async {
    _settings = _settings.copyWith(
      defaultCurrency: currency,
    );
  }
  
  @override
  Future<void> initialize() async {
    // No-op for mock
  }
  
  // For testing purposes
  void setSettings(Settings settings) {
    _settings = settings;
  }
}

void main() {
  group('ThemeProvider', () {
    late MockSettingsRepository mockRepository;
    
    setUp(() {
      mockRepository = MockSettingsRepository();
    });
    
    test('initializes with system theme when no repository is provided', () {
      final provider = ThemeProvider();
      
      expect(provider.useSystemTheme, isTrue);
      // We can't test the actual isDarkMode value here as it depends on the system theme
    });
    
    test('initializes with settings from repository', () {
      mockRepository.setSettings(Settings(
        isDarkMode: true,
        useSystemTheme: false,
      ));
      
      final provider = ThemeProvider(settingsRepository: mockRepository);
      
      // Wait for the provider to initialize
      Future.delayed(Duration.zero);
      
      expect(provider.isDarkMode, isTrue);
      expect(provider.useSystemTheme, isFalse);
    });
    
    test('toggleTheme switches between light and dark mode', () async {
      mockRepository.setSettings(Settings(
        isDarkMode: false,
        useSystemTheme: true,
      ));
      
      final provider = ThemeProvider(settingsRepository: mockRepository);
      
      // Wait for the provider to initialize
      await Future.delayed(Duration.zero);
      
      await provider.toggleTheme();
      
      expect(provider.isDarkMode, isTrue);
      expect(provider.useSystemTheme, isFalse);
      
      // Check that settings were updated in repository
      final settings = mockRepository.getSettings();
      expect(settings.isDarkMode, isTrue);
      expect(settings.useSystemTheme, isFalse);
    });
    
    test('setSystemTheme sets theme to follow system', () async {
      mockRepository.setSettings(Settings(
        isDarkMode: true,
        useSystemTheme: false,
      ));
      
      final provider = ThemeProvider(settingsRepository: mockRepository);
      
      // Wait for the provider to initialize
      await Future.delayed(Duration.zero);
      
      await provider.setSystemTheme();
      
      expect(provider.useSystemTheme, isTrue);
      
      // Check that settings were updated in repository
      final settings = mockRepository.getSettings();
      expect(settings.useSystemTheme, isTrue);
    });
    
    test('themeMode returns correct value based on settings', () {
      // Test with dark mode and not using system theme
      mockRepository.setSettings(Settings(
        isDarkMode: true,
        useSystemTheme: false,
      ));
      
      var provider = ThemeProvider(settingsRepository: mockRepository);
      
      // Wait for the provider to initialize
      Future.delayed(Duration.zero);
      
      expect(provider.themeMode, equals(ThemeMode.dark));
      
      // Test with light mode and not using system theme
      mockRepository.setSettings(Settings(
        isDarkMode: false,
        useSystemTheme: false,
      ));
      
      provider = ThemeProvider(settingsRepository: mockRepository);
      
      // Wait for the provider to initialize
      Future.delayed(Duration.zero);
      
      expect(provider.themeMode, equals(ThemeMode.light));
      
      // Test with using system theme
      mockRepository.setSettings(Settings(
        useSystemTheme: true,
      ));
      
      provider = ThemeProvider(settingsRepository: mockRepository);
      
      // Wait for the provider to initialize
      Future.delayed(Duration.zero);
      
      expect(provider.themeMode, equals(ThemeMode.system));
    });
    
    test('lightTheme returns a ThemeData with light brightness', () {
      final provider = ThemeProvider();
      final theme = provider.lightTheme;
      
      expect(theme, isA<ThemeData>());
      expect(theme.brightness, equals(Brightness.light));
    });
    
    test('darkTheme returns a ThemeData with dark brightness', () {
      final provider = ThemeProvider();
      final theme = provider.darkTheme;
      
      expect(theme, isA<ThemeData>());
      expect(theme.brightness, equals(Brightness.dark));
    });
  });
}
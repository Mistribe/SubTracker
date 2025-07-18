import 'package:hive/hive.dart';
import 'package:hive_flutter/hive_flutter.dart';
import '../models/settings.dart';

/// Repository for handling application settings persistence
class SettingsRepository {
  static const String _boxName = 'settings';
  static const String _settingsKey = 'app_settings';
  late Box<Settings> _box;

  /// Initialize the repository
  ///
  /// This method must be called before using any other methods in this class.
  Future<void> initialize() async {
    // Open the Hive box for settings
    _box = await Hive.openBox<Settings>(_boxName);

    // Initialize with default settings if none exist
    if (!_box.containsKey(_settingsKey)) {
      await _box.put(_settingsKey, Settings());
    }
  }

  /// Get the current settings
  Settings getSettings() {
    return _box.get(_settingsKey) ?? Settings();
  }

  /// Save updated settings
  Future<void> saveSettings(Settings settings) async {
    await _box.put(_settingsKey, settings);
  }

  /// Update theme settings
  Future<void> updateThemeSettings({bool? isDarkMode, bool? useSystemTheme}) async {
    final currentSettings = getSettings();
    final updatedSettings = currentSettings.copyWith(
      isDarkMode: isDarkMode,
      useSystemTheme: useSystemTheme,
    );
    await saveSettings(updatedSettings);
  }

  /// Update default currency
  Future<void> updateDefaultCurrency(String currency) async {
    final currentSettings = getSettings();
    final updatedSettings = currentSettings.copyWith(
      defaultCurrency: currency,
    );
    await saveSettings(updatedSettings);
  }

  /// Update selected family ID
  Future<void> updateSelectedFamilyId(String? familyId) async {
    final currentSettings = getSettings();
    final updatedSettings = currentSettings.copyWith(
      selectedFamilyId: familyId,
    );
    await saveSettings(updatedSettings);
  }
}

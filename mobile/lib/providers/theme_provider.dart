import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import '../repositories/settings_repository.dart';

class ThemeProvider with ChangeNotifier {
  bool _isDarkMode = false;
  bool _useSystemTheme = true;
  final SettingsRepository? _settingsRepository;

  ThemeProvider({SettingsRepository? settingsRepository}) 
      : _settingsRepository = settingsRepository {
    // Load settings if repository is provided
    if (_settingsRepository != null) {
      _loadSettings();
    } else {
      // Initialize with system theme if no repository
      _initializeWithSystemTheme();
    }
  }

  Future<void> _loadSettings() async {
    final settings = _settingsRepository!.getSettings();
    _isDarkMode = settings.isDarkMode;
    _useSystemTheme = settings.useSystemTheme;

    // If using system theme, update dark mode based on system
    if (_useSystemTheme) {
      _initializeWithSystemTheme();
    }

    notifyListeners();
  }

  void _initializeWithSystemTheme() {
    var brightness = SchedulerBinding.instance.platformDispatcher.platformBrightness;
    _isDarkMode = brightness == Brightness.dark;
  }

  bool get isDarkMode => _isDarkMode;
  bool get useSystemTheme => _useSystemTheme;

  ThemeMode get themeMode {
    if (_useSystemTheme) {
      return ThemeMode.system;
    }
    return _isDarkMode ? ThemeMode.dark : ThemeMode.light;
  }

  Future<void> toggleTheme() async {
    _useSystemTheme = false;
    _isDarkMode = !_isDarkMode;

    // Persist settings if repository is available
    if (_settingsRepository != null) {
      await _settingsRepository.updateThemeSettings(
        isDarkMode: _isDarkMode,
        useSystemTheme: _useSystemTheme,
      );
    }

    notifyListeners();
  }

  Future<void> setSystemTheme() async {
    _useSystemTheme = true;
    _initializeWithSystemTheme();

    // Persist settings if repository is available
    if (_settingsRepository != null) {
      await _settingsRepository.updateThemeSettings(
        isDarkMode: _isDarkMode,
        useSystemTheme: _useSystemTheme,
      );
    }

    notifyListeners();
  }

  // Light theme
  ThemeData get lightTheme {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: Colors.blue,
        brightness: Brightness.light,
      ),
      appBarTheme: const AppBarTheme(
        centerTitle: true,
        elevation: 0,
      ),
      // cardTheme property removed due to type incompatibility
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        filled: true,
      ),
    );
  }

  // Dark theme
  ThemeData get darkTheme {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: Colors.blue,
        brightness: Brightness.dark,
      ),
      appBarTheme: const AppBarTheme(
        centerTitle: true,
        elevation: 0,
      ),
      // cardTheme property removed due to type incompatibility
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        filled: true,
      ),
    );
  }
}
